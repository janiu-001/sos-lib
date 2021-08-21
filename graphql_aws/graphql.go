package graphql_aws

/* Amplify Params - DO NOT EDIT
	API_PAX_GRAPHQLAPIENDPOINTOUTPUT
	API_PAX_GRAPHQLAPIIDOUTPUT
	API_PAX_GRAPHQLAPIKEYOUTPUT
	ENV
	REGION
Amplify Params - DO NOT EDIT */

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"os"
)

var (
	graphqlClient *graphql.Client
	URL           string
	ApiKey        string
)

func init() {
	URL = os.Getenv("API_PAX_GRAPHQLAPIENDPOINTOUTPUT")
	ApiKey = os.Getenv("API_PAX_GRAPHQLAPIKEYOUTPUT")
	graphqlClient = graphql.NewClient(URL)
}

// GraphqlClient Graphql client Related
type GraphqlClient struct {
}

func NewGraphqlClient() *GraphqlClient {
	return &GraphqlClient{}
}

func (g *GraphqlClient) GraphqlOperation(condition string, value map[string]interface{}) (resp interface{}, err error) {
	graphqlRequest := graphql.NewRequest(condition)
	graphqlRequest.Header.Set("x-api-key", ApiKey)
	for k, v := range value {
		graphqlRequest.Var(k, v)
	}
	if err = graphqlClient.Run(context.Background(), graphqlRequest, &resp); err != nil {
		fmt.Println(err)
		return
	}
	return
}
