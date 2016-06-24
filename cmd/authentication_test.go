package cmd

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
)

func resetAuthenticate() {
	// reset automation flag vars
	ResetFlags()
	// mock interface
	AuthenticationV3 = newMockAuthenticationV3
}

func TestAuthenticationUserIdOrNameRequired(t *testing.T) {
	// reset params
	resetAuthenticate()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --project-id=%s --password=%s", "http://some_test_url", "bup", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestAuthenticationProjectIdOrNameRequired(t *testing.T) {
	// reset params
	resetAuthenticate()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --password=%s", "http://some_test_url", "bup", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestAuthenticationPasswordfromStdInputWhenEmpty(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***
export OS_TOKEN=test_token_id`

	// keep backup of the real stdout
	oldStdout := os.Stdout

	// write passowrd
	_, err := pipeToStdin("password!!!\n")
	if err != nil {
		t.Error(err.Error())
		return
	}

	// pipe std out
	_, w, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
	}
	os.Stdout = w

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s", "http://some_test_url", "miau", "bup", "123456789"))

	// flush, restore close
	os.Stdout = oldStdout
	flushStdin()
	w.Close()

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationWithAllFlags(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --username=%s --password=%s --project-id=%s --project-name=%s --user-domain-name=%s --user-domain-id=%s --project-domain-name=%s --project-domain-id=%s --region=%s", "http://some_test_url", "userid", "username", "passwrod", "projectid", "projectname", "userdomainname", "userdomainid", "projectdomainid", "projectdomainname", "staging"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationResultString(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s", "http://some_test_url", "miau", "bup", "123456789"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationResultJSON(t *testing.T) {
	want := `{"ARC_SERVICE_ENDPOINT": "https://arc.staging.***REMOVED***","LYRA_SERVICE_ENDPOINT": "https://lyra.staging.***REMOVED***","OS_TOKEN": "test_token_id"}`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s --json", "http://some_test_url", "miau", "bup", "123456789"))

	wantSource := map[string]string{}
	err := json.Unmarshal([]byte(want), &wantSource)
	if err != nil {
		t.Error(err.Error())
		return
	}

	response := map[string]string{}
	err = json.Unmarshal([]byte(resulter.Output), &response)
	if err != nil {
		t.Error(err.Error())
		return
	}

	eq := reflect.DeepEqual(wantSource, response)
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}

//
// Test region
//

func TestAuthenticationNotGivenARegion(t *testing.T) {
	// should return first entry from endpoints
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s", "http://some_test_url", "miau", "bup", "123456789"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationGivenARegion(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s --region=%s", "http://some_test_url", "miau", "bup", "123456789", "staging"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationGivenAWrongRegion(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=
export ARC_SERVICE_ENDPOINT=
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s --region=%s", "http://some_test_url", "miau", "bup", "123456789", "wrong"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

//
// Mock authentication interface
//

type MockV3 struct {
	AuthOpts   LyraAuthOps
	client     *gophercloud.ServiceClient
	TestServer *httptest.Server
}

func newMockAuthenticationV3(authOpts LyraAuthOps) Authentication {
	return &MockV3{AuthOpts: authOpts}
}

func (a *MockV3) CheckAuthenticationParams() error {
	return checkAuthenticationParams(&a.AuthOpts)
}

func (a *MockV3) GetToken() (*tokens.Token, error) {
	token := tokens.Token{ID: "test_token_id"}
	return &token, nil
}

func (a *MockV3) GetServicePublicEndpoint(serviceType string) (string, error) {
	if a.TestServer != nil {
		return a.TestServer.URL, nil
	} else {

		// get entry from catalog
		serviceEntry, err := getServiceEntry(serviceType, &catalog1)
		if err != nil {
			return "", err
		}

		// get endpoint
		endpoint, err := getServicePublicEndpoint(a.AuthOpts.Region, serviceEntry)
		if err != nil {
			return "", err
		}

		return endpoint, nil
	}
	return "", nil
}

var catalog1 = tokens.ServiceCatalog{
	Entries: []tokens.CatalogEntry{
		{ID: "s-8be070817", Name: "Arc", Type: "arc", Endpoints: []tokens.Endpoint{
			{ID: "e-884f431c9", Region: "staging", Interface: "public", URL: "https://arc.staging.***REMOVED***"},
		}},
		{ID: "s-d5e793744", Name: "Lyra", Type: "automation", Endpoints: []tokens.Endpoint{
			{ID: "e-54b8d28fc", Region: "staging", Interface: "public", URL: "https://lyra.staging.***REMOVED***"},
		}},
	},
}
