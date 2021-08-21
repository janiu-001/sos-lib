package usercase

import (
	"bytes"
	"fmt"
	"github.com/janiu-001/sos-lib/graphql_aws"
	"github.com/janiu-001/sos-lib/salesforce"
	"github.com/janiu-001/sos-lib/util"
	"github.com/simpleforce/simpleforce"
)

func SyncAll() {

	//lambda.Start(HandleRequest)
	ret, err := graphql_aws.GetOrganizationsList()
	if err != nil {
		fmt.Println("SyncAll error", err)
		return
	}
	for _, v := range ret.ListOrganizations.OrganizationItem {
		if v.SyncSFObject == false {
			continue
		}
		err = SyncSfObjects(v, "")
		fmt.Println("SyncAll error", err)
	}

}

type SyncOneOrganizationDTO struct {
	OrganizationID string
	ObjectName     string
	Condition      string
}

func SyncWithConditions(req SyncOneOrganizationDTO) {

	//1. get the organization
	resp, err := graphql_aws.GetOrganizationByOrganizationId(req.OrganizationID)
	if err != nil {
		return
	}
	organization := resp.Organization
	//2. get the sf object
	listObjects, err := graphql_aws.GetSfObjectsByOrganizationIdOrName(req.OrganizationID, req.ObjectName)
	if err != nil {
		return
	}
	organization.SfObjects = listObjects.SFObject

	//3. syn
	err = SyncSfObjects(organization, req.Condition)

}

func SyncSfObjects(o graphql_aws.Organization, condition string) (err error) {
	//1.create for each organization
	client := &salesforce.SimpleForceClient{
		SfURL:      o.SfLoginUrl,
		SfUser:     o.SfUsername,
		SfPassword: o.SfPassword,
		SfToken:    o.SfSecretToken,
	}
	err = client.CreateSimpleForceClient()
	if err != nil {
		return
	}

	//2.sync all the object
	for _, v := range o.SfObjects.SFObjectItem {
		err = DoSync(&v, client, condition, &o)
		if err != nil {
			return
		}
	}
	return
}

func DoSync(sf *graphql_aws.SFObject, client *salesforce.SimpleForceClient, condition string, o *graphql_aws.Organization) (err error) {
	// 1. recover from panic if use routine
	if r := recover(); r != nil {
		fmt.Println("Sync Recovered in f", r)
	}

	// 2. create the sosql scripts
	var temp string
	//updates := map[string]SfField{}
	var buffer bytes.Buffer
	fieldsMap := make(map[string]string)
	for _, field := range sf.Fields {
		buffer.WriteString(field.ApiName)
		buffer.WriteString(",")
		fieldsMap[field.ApiName] = field.Name
	}
	temp = buffer.String()
	apiNames := temp[0 : len(temp)-1]

	totalSize := 0
	Offset := 0
	Limit := 1
	allIdMap := make(map[string]string)
	for {
		//3.query records from sales force
		querySOQL := fmt.Sprintf(graphql_aws.SelectString, apiNames, sf.SfObjectName, Limit, Offset)
		if condition != "" {
			querySOQL = fmt.Sprintf("querySOQL %s", querySOQL, condition)
		}
		var result *simpleforce.QueryResult
		result, err = client.QuerySaleForce(querySOQL)
		if err != nil {
			break
		}
		totalSize = result.TotalSize + totalSize
		if result.TotalSize < 1 {
			break
		}

		Offset = Offset + Limit + 1
		//records := make([]map[string]interface{}, 1)
		//_ = util.ConvertStructData(result.Records, &records)
		var idList []string
		idList, err = SyncDynamodb(sf.SfObjectName, result.Records, fieldsMap, o.ID)
		if err != nil {
			return
		}

		for _, v := range idList {
			allIdMap[v] = v
		}
		if totalSize == 1 {
			break
		}

	}

	//3 delete the object
	err = DeleteDynamodb(sf.SfObjectName, o.ID, &allIdMap)
	return
}

func SyncDynamodb(objectName string, records []simpleforce.SObject, fieldsMap map[string]string, orgId string) (idList []string, err error) {
	updates := make(map[string]interface{})
	// every object should add the organizationID
	updates["organizationID"] = orgId
	var id interface{}
	for _, record := range records {
		// 1. get the updates data
		for k, v := range record {

			name := fieldsMap[k]
			if name == "" {
				continue
			}
			if name == "id" {
				id = v
				idList = append(idList, v.(string))
			}
			if v != nil {
				updates[name] = v
			}
		}

		// 2.query the salesforce
		var condition string
		condition, err = graphql_aws.GetGeneralTableScripts(objectName, "query")
		if err != nil {
			return
		}
		var ret interface{}
		graphqlClient := graphql_aws.NewGraphqlClient()
		ret, err = graphqlClient.GraphqlOperation(condition, map[string]interface{}{"id": id})
		if err != nil {
			return
		}
		//3 .check if the object is exist
		retMap := make(map[string]interface{})
		_ = util.ConvertStructData(ret, &retMap)
		queryTable := fmt.Sprintf("get%s", objectName)
		action := ""
		if retMap[queryTable] == nil {
			action = "create"

		} else {
			action = "update"
		}

		//4. generate the create or update the sql
		condition, err = graphql_aws.GetGeneralTableScripts(objectName, action)
		if err != nil {
			return
		}

		//5. sync
		_, err = graphqlClient.GraphqlOperation(condition, map[string]interface{}{"input": updates})
	}

	return
}

type IDList struct {
	ID string `json:"id"`
}

func DeleteDynamodb(objectName, organizationId string, idMap *map[string]string) (err error) {
	condition, err := graphql_aws.GetGeneralTableScripts(objectName, "list")
	if err != nil {
		return
	}
	var ret interface{}
	graphqlClient := graphql_aws.NewGraphqlClient()
	value := map[string]interface{}{
		"organizationID": map[string]interface{}{"eq": organizationId},
	}
	ret, err = graphqlClient.GraphqlOperation(condition, map[string]interface{}{"filter": value})
	if err != nil {
		return
	}
	retMap := make(map[string]interface{})
	_ = util.ConvertStructData(ret, &retMap)
	queryTable := fmt.Sprintf("list%ss", objectName)
	if retMap[queryTable] == nil {
	}

	listValue := retMap[queryTable]
	items := map[string]interface{}{}
	_ = util.ConvertStructData(listValue, &items)
	itemsValue := items["items"]
	var idList []IDList
	_ = util.ConvertStructData(itemsValue, &idList)
	var deleteIdList []string
	for _, v := range *idMap {
		for _, value := range idList {
			if v == value.ID {
				continue
			}
		}
		deleteIdList = append(deleteIdList, v)
	}
	for _, v := range deleteIdList {

		condition, err = graphql_aws.GetGeneralTableScripts(objectName, "delete")
		if err != nil {
			return
		}
		value := map[string]interface{}{"id": v}
		_, err = graphqlClient.GraphqlOperation(condition, map[string]interface{}{"input": value})
	}
	return
}
