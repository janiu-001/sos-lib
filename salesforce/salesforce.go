package salesforce

import (
	"errors"
	"fmt"
	"github.com/simpleforce/simpleforce"
)

// SimpleForceClient SimpleForce client Related
type SimpleForceClient struct {
	SfURL      string
	SfUser     string
	SfPassword string
	SfToken    string
	c          *simpleforce.Client
}

// CreateSimpleForceClient create the sale force clinet
func (s *SimpleForceClient) CreateSimpleForceClient() (err error) {
	s.c = simpleforce.NewClient(s.SfURL, simpleforce.DefaultClientID, simpleforce.DefaultAPIVersion)
	if s.c == nil {
		errors.New("failed to new client")
		// handle the error
		return
	}
	fmt.Println(*s)
	err = s.c.LoginPassword(s.SfUser, s.SfPassword, s.SfToken)
	if err != nil {
		// handle the error
		return
	}
	return
}
func (s *SimpleForceClient) QuerySaleForce(condition string) (ret *simpleforce.QueryResult, err error) {
	if s.c == nil {
		err = errors.New("sale force client ins null")
		return
	}

	ret = &simpleforce.QueryResult{}
	ret, err = s.c.Query(condition)
	return
}
