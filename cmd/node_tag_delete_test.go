package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func TestNodeTagDeleteCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra node tag delete")
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra node tag delete")
	ResetFlags()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra node tag delete")
}

func TestNodeTagDeleteCmdWithUserAuthenticationFlags(t *testing.T) {
	testServer := TestServer(200, "", map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node tag delete --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-id=%s", "some_test_url", "kuak", "bup", "123456789", "node_id"))
	if resulter.Error != nil {
		t.Error(fmt.Sprint(`Command expected to not get an error: `, resulter.Error))
	}
}

func TestNodeTagDeleteCmdRightParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		method := r.Method
		path := r.URL
		if !strings.Contains(method, "DELETE") {
			diffString := StringDiff(method, "DELETE")
			t.Error(fmt.Sprintf("Command API method doesn't match. \n \n %s", diffString))
		}
		if !strings.Contains(path.String(), "agents/123456789/tags/test123") {
			diffString := StringDiff(method, "agents/123456789/tags/test123")
			t.Error(fmt.Sprintf("Command API path doesn't match. \n \n %s", diffString))
		}
	}))
	defer server.Close()
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra node tag delete --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=%s test123", "https://somewhere.com", server.URL, "token123", "123456789"))
}
