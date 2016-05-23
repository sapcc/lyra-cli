package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/rackspace/gophercloud"
	"reflect"
	"strings"
	"testing"
)

func resetAuthenticate() {
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	authenticateCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(authenticateCmd)
	// mock interface
	AuthenticationV3 = newMockAuthenticationV3
}

func TestAuthenticationUserIdOrNameRequired(t *testing.T) {}

func TestAuthenticationPasswordfromStdInputWhenEmpty(t *testing.T) {}

func TestAuthenticationWithAllFlags(t *testing.T) {}

func TestAuthenticationResultString(t *testing.T) {
	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --identity-endpoint=%s --user-id=%s --project-id=%s --password=%s", "http://some_test_url", "miau", "bup", "123456789"))
	want := `export LYRA_SERVICE_ENDPOINT=test_public_endpoint
export ARC_SERVICE_ENDPOINT=test_public_endpoint
export OS_TOKEN=test_token_id`

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationResultJSON(t *testing.T) {
	want := `{"ARC_SERVICE_ENDPOINT": "test_public_endpoint","LYRA_SERVICE_ENDPOINT": "test_public_endpoint","OS_TOKEN": "test_token_id"}`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --identity-endpoint=%s --user-id=%s --project-id=%s --password=%s --json", "http://some_test_url", "miau", "bup", "123456789"))

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
// Mock authentication interface
//

type MockV3 struct {
	AuthOpts LyraAuthOps
	client   *gophercloud.ServiceClient
}

func newMockAuthenticationV3(authOpts LyraAuthOps) Authentication {
	return &MockV3{AuthOpts: authOpts}
}

func (a *MockV3) GetToken() (string, error) {
	return "test_token_id", nil
}

func (a *MockV3) GetServicePublicEndpoint(serviceId string) (string, error) {
	return "test_public_endpoint", nil
}

func (a *MockV3) GetServiceId(serviceType string) (string, error) {
	return "Test_service_id", nil
}
