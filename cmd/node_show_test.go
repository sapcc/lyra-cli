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

func TestNodeShowCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra node show")
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra node show")
	ResetFlags()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra node show")
}

func TestNodeShowCmdWithUserAuthenticationFlags(t *testing.T) {
	testServer := TestServer(200, "{}", map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node show --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-id=%s", "some_test_url", "kuak", "bup", "123456789", "node_id"))
	if resulter.Error != nil {
		t.Error(fmt.Sprint(`Command expected to not get an error: `, resulter.Error))
	}
}

func TestNodeShowCmdMissingId(t *testing.T) {
	// set test server
	server := TestServer(200, "{}", map[string]string{})
	defer server.Close()

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "https://somewhere.com", server.URL, "token123"))

	errorMsg := locales.ErrorMessages("node-id-missing")
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
	}
}

func TestNodeShowCmdResultTable(t *testing.T) {
	// set test server
	responseBody := `{"agent_id":"99da3762-921b-47fa-b02d-90b44fa63eba","project":"p-e03c46eab","organization":"o-monsoon2","created_at":"2016-12-09T14:40:56.1885Z","updated_at":"2017-02-08T14:36:01.762113Z","updated_with":"391c410f-bdad-443a-95af-9f106ce38779","updated_by":"linux"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+--------------+--------------------------------------+
|     KEY      |                VALUE                 |
+--------------+--------------------------------------+
| agent_id     | 99da3762-921b-47fa-b02d-90b44fa63eba |
| created_at   | 2016-12-09T14:40:56.1885Z            |
| organization | o-monsoon2                           |
| project      | p-e03c46eab                          |
| updated_at   | 2017-02-08T14:36:01.762113Z          |
| updated_by   | linux                                |
| updated_with | 391c410f-bdad-443a-95af-9f106ce38779 |
+--------------+--------------------------------------+`

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=123456789", "http://somewhere.com", server.URL, "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestNodeShowCmdResultJson(t *testing.T) {
	// set test server
	responseBody := `{"agent_id":"99da3762-921b-47fa-b02d-90b44fa63eba","project":"p-e03c46eab","organization":"o-monsoon2","created_at":"2016-12-09T14:40:56.1885Z","updated_at":"2017-02-08T14:36:01.762113Z","updated_with":"391c410f-bdad-443a-95af-9f106ce38779","updated_by":"linux"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=123456789 --json", "http://somewhere.com", server.URL, "token123"))
	eq, err := JsonDiff(responseBody, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}

func TestNodeShowCmdRightParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		method := r.Method
		path := r.URL
		if !strings.Contains(method, "GET") {
			diffString := StringDiff(method, "GET")
			t.Error(fmt.Sprintf("Command API method doesn't match. \n \n %s", diffString))
		}
		if !strings.Contains(path.String(), "agents/123456789") {
			diffString := StringDiff(method, "agents/123456789")
			t.Error(fmt.Sprintf("Command API path doesn't match. \n \n %s", diffString))
		}
	}))
	defer server.Close()
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra node show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=%s", "https://somewhere.com", server.URL, "token123", "123456789"))
}
