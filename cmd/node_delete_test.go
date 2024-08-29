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

func TestNodeDeleteCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra node delete")
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra node delete")
	ResetFlags()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra node delete")
}

func TestNodeDeleteCmdWithUserAuthenticationFlags(t *testing.T) {
	testServer := TestServer(204, "Node deleted successfully", map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node delete --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-id=%s", "some_test_url", "kuak", "bup", "123456789", "node_id"))
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
}

func TestNodeDeleteCmdMissingId(t *testing.T) {
	// set test server
	server := TestServer(204, "", map[string]string{})
	defer server.Close()

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node delete --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "https://somewhere.com", server.URL, "token123"))

	errorMsg := locales.ErrorMessages("node-id-missing")
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Errorf("Command error doesn't match. \n \n %s", diffString)
	}
}

func TestNodeDeleteCmdSuccess(t *testing.T) {
	// set test server
	server := TestServer(204, "", map[string]string{})
	defer server.Close()

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node delete --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=%s", "https://somewhere.com", server.URL, "token123", "123456789"))

	if !strings.Contains(resulter.ErrorOutput, "Node") || !strings.Contains(resulter.ErrorOutput, "123456789") || !strings.Contains(resulter.ErrorOutput, "deleted") {
		diffString := StringDiff(resulter.ErrorOutput, "deleted")
		t.Errorf("Command error doesn't match. \n \n %s", diffString)
	}
}

func TestNodeDeleteCmdRightParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		method := r.Method
		path := r.URL
		if !strings.Contains(method, "DELETE") {
			diffString := StringDiff(method, "DELETE")
			t.Errorf("Command API method doesn't match. \n \n %s", diffString)
		}
		if !strings.Contains(path.String(), "agents/123456789") {
			diffString := StringDiff(method, "agents/123456789")
			t.Errorf("Command API path doesn't match. \n \n %s", diffString)
		}
	}))
	defer server.Close()
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra node delete --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=%s", "https://somewhere.com", server.URL, "token123", "123456789"))
}
