package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/locales"
)

func resetAutomationDeleteFlagVars() {
	// reset flags
	ResetFlags()
}

func TestAutomationDeleteCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	resetAutomationDeleteFlagVars()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra automation delete")
	resetAutomationDeleteFlagVars()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra automation delete")
	resetAutomationDeleteFlagVars()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra automation delete")
}

func TestAutomationDeleteCmdWithUserAuthenticationFlags(t *testing.T) {
	testServer := TestServer(204, "", map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)

	// reset stuff
	resetAutomationDeleteFlagVars()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation delete --auth-url=%s --user-id=%s --project-id=%s --password=%s --automation-id=%s", "some_test_url", "miau", "bup", "123456789", "automation_id"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
}

func TestAutomationDeleteCmdMissingId(t *testing.T) {
	// set test server
	server := TestServer(200, "", map[string]string{})
	defer server.Close()

	// reset stuff
	resetAutomationDeleteFlagVars()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation delete --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "https://somewhere.com", "token123"))

	errorMsg := locales.ErrorMessages("automation-id-missing")
	if !strings.Contains(resulter.Output, errorMsg) {
		diffString := StringDiff(resulter.Output, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationDeleteCmdSuccess(t *testing.T) {
	// set test server
	server := TestServer(204, "", map[string]string{})
	defer server.Close()

	// reset stuff
	resetAutomationDeleteFlagVars()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation delete --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --automation-id=%s", server.URL, "https://somewhere.com", "token123", "123456789"))
	if !strings.Contains(resulter.Output, "123456789") || !strings.Contains(resulter.Output, "deleted") {
		diffString := StringDiff(resulter.Output, "deleted")
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationDeleteCmdRightParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		method := r.Method
		path := r.URL
		if !strings.Contains(method, "DELETE") {
			diffString := StringDiff(method, "DELETE")
			t.Error(fmt.Sprintf("Command API method doesn't match. \n \n %s", diffString))
		}
		if !strings.Contains(path.String(), "automations/123456789") {
			diffString := StringDiff(method, "automations/123456789")
			t.Error(fmt.Sprintf("Command API path doesn't match. \n \n %s", diffString))
		}
	}))
	defer server.Close()
	// reset stuff
	resetAutomationDeleteFlagVars()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra automation delete --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --automation-id=%s", server.URL, "https://somewhere.com", "token123", "123456789"))
}
