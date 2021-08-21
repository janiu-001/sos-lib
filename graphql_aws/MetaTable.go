package graphql_aws

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/janiu-001/sos-lib/util"
)

type MetaTableField struct {
	Order  int    `json:"order"`
	Name   string `json:"name"`
	NameJp string `json:"name_jp"`
}
type MetaTable struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	NameJp string           `json:"name_jp"`
	Fields []MetaTableField `json:"fields"`
}

type MetaTableItems struct {
	MetaTableItems []MetaTable `json:"items"`
}

type MetaTables struct {
	ListMetaTables MetaTableItems `json:"listMetaTables"`
}

func GetMetaTableByName(tableName string) (resp *MetaTables, err error) {
	value := map[string]interface{}{"name": map[string]interface{}{"eq": tableName}}

	graphqlClient := NewGraphqlClient()
	ret, err := graphqlClient.GraphqlOperation(ListMetaTables, map[string]interface{}{"filter": value})
	//ret, err := GraphqlOperation(condition, map[string]interface{}{"limit": 10})
	err = util.ConvertStructData(ret, &resp)

	return
}

func GetGeneralTableScripts(objectName string, action string) (resp string, err error) {
	//1. query the graphql meta table to generate the create/update scripts
	ret := &MetaTables{}
	ret, err = GetMetaTableByName(objectName)
	if err != nil {
		return
	}
	if len(ret.ListMetaTables.MetaTableItems) == 0 {
		err = errors.New("no table exits")
		return
	}
	fields := ret.ListMetaTables.MetaTableItems[0].Fields
	if len(fields) == 0 {
		err = errors.New("no fields exits")
		return
	}

	prefix := ""
	suffix := ""
	switch action {
	case "create":
		prefix = fmt.Sprintf(CreateTablePrefix, objectName, objectName, objectName)
		suffix = CreateTableSuffix
	case "update":
		prefix = fmt.Sprintf(UpdateTablePrefix, objectName, objectName, objectName)
		suffix = UpdageTableSuffix
	case "query":
		resp = fmt.Sprintf(QueryTableSuffix, objectName, objectName)
		return
	case "delete":
		resp = fmt.Sprintf(DeleteTableSuffix, objectName, objectName, objectName, objectName)
		return
	case "list":
		resp = fmt.Sprintf(ListPrefix, objectName, objectName, objectName)

		return
	default:
		err = errors.New("invalid action")
		return
	}
	var buffer bytes.Buffer
	for _, v := range fields {
		buffer.WriteString(" ")
		buffer.WriteString(v.Name)
	}
	resp = fmt.Sprintf("%s %s %s", prefix, buffer.String(), suffix)
	return
}
