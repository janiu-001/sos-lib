package graphql_aws

import (
	"github.com/janiu-001/sos-lib/util"
)

type SfField struct {
	Order   int    `json:"order"`
	Name    string `json:"name"`
	NameJp  string `json:"name_jp"`
	ApiName string `json:"apiName"`
}

type SFObject struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	SfObjectName string    `json:"sfObjectName"`
	Fields       []SfField `json:"fields"`
}

type SFObjectItem struct {
	SFObjectItem []SFObject `json:"items"`
}

type SFObjects struct {
	SFObjectItems []SFObjectItem `json:"listSFObjects"`
}

type GEtSFObjectsModel struct {
	SFObject SFObjectItem `json:"listSFObjects"`
}

type Organization struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	SfLoginUrl    string       `json:"sfLoginUrl"`
	SfUsername    string       `json:"sfUsername"`
	SfPassword    string       `json:"sfPassword"`
	SfSecretToken string       `json:"sfSecretToken"`
	SyncSFObject  bool         `json:"syncSFObject"`
	SfObjects     SFObjectItem `json:"sfObjects"`
}

type OrganizationItem struct {
	OrganizationItem []Organization `json:"items"`
}

type OrganizationListModel struct {
	ListOrganizations OrganizationItem `json:"listOrganizations"`
}

type GetOrganizationModel struct {
	Organization Organization `json:"getOrganization"`
}

func GetOrganizationByOrganizationId(organizationId string) (resp GetOrganizationModel, err error) {
	value := map[string]interface{}{"id": organizationId}

	graphqlClient := NewGraphqlClient()

	ret, err := graphqlClient.GraphqlOperation(GetOrganization, value)
	err = util.ConvertStructData(ret, &resp)
	return
}

func GetSfObjectsByOrganizationIdOrName(organizationId string, objetName string) (resp GEtSFObjectsModel, err error) {
	if organizationId == "" {
		return
	}
	value := map[string]interface{}{}
	value["organizationID"] = organizationId
	if objetName != "" {
		value["sfObjectName"] = objetName
	}
	graphqlClient := NewGraphqlClient()
	ret, err := graphqlClient.GraphqlOperation(ListSFObjects, value)
	err = util.ConvertStructData(ret, &resp)
	return
}

func GetOrganizationsList() (resp OrganizationListModel, err error) {

	graphqlClient := NewGraphqlClient()
	ret, err := graphqlClient.GraphqlOperation(ListOrganizations, nil)
	err = util.ConvertStructData(ret, &resp)
	return
}
