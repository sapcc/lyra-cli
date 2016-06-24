package goOpenstackAuth

import (
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
)

//
// Mock authentication interface
//

type MockV3 struct {
	Options     AuthV3Options
	tokenResult *tokens.CreateResult
}

func NewMockAuthenticationV3(authOpts AuthV3Options) Authentication {
	return &MockV3{Options: authOpts}
}

func (a *MockV3) GetToken() (*tokens.Token, error) {
	token := tokens.Token{ID: "test_token_id"}
	return &token, nil
}

func (a *MockV3) GetServiceEndpoint(serviceType, region, serviceInterface string) (string, error) {
	// get entry from catalog
	serviceEntry, err := getServiceEntry(serviceType, &Catalog1)
	if err != nil {
		return "", err
	}

	// get endpoint
	endpoint, err := getServiceEndpoint(region, serviceInterface, serviceEntry)
	if err != nil {
		return "", err
	}

	return endpoint, nil
}

func (a *MockV3) GetProject() (*Project, error) {
	return extractProject(CommonResult1)
}

var Catalog1 = tokens.ServiceCatalog{
	Entries: []tokens.CatalogEntry{
		{ID: "s-8be070817", Name: "Arc", Type: "arc", Endpoints: []tokens.Endpoint{
			{ID: "e-904f431c9", Region: "staging", Interface: "internal", URL: "https://arc.staging.***REMOVED***/internal"},
			{ID: "e-904f431c9", Region: "staging", Interface: "admin", URL: "https://arc.staging.***REMOVED***/admin"},
			{ID: "e-884f431c9", Region: "staging", Interface: "public", URL: "https://arc.staging.***REMOVED***/public"},
			{ID: "e-904f431c9", Region: "production", Interface: "internal", URL: "https://arc.production.***REMOVED***/internal"},
			{ID: "e-904f431c9", Region: "production", Interface: "admin", URL: "https://arc.production.***REMOVED***/admin"},
			{ID: "e-884f431c9", Region: "production", Interface: "public", URL: "https://arc.production.***REMOVED***/public"},
		}},
		{ID: "s-d5e793744", Name: "Lyra", Type: "automation", Endpoints: []tokens.Endpoint{
			{ID: "e-54b8d28fc", Region: "staging", Interface: "public", URL: "https://lyra.staging.***REMOVED***"},
		}},
	},
}

var CommonResult1 = map[string]interface{}{"token": map[string]interface{}{"project": map[string]string{"id": "p-9597d2775", "domain_id": "o-monsoon2", "name": "Arc_Development"}}}
