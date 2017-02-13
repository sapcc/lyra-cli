package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func TestNodeTagListCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra node tag list")
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra node tag list")
	ResetFlags()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra node tag list")
}

func TestNodeTagListCmdWithUserAuthenticationFlags(t *testing.T) {
	testServer := TestServer(200, "{}", map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node tag list --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-id=%s", "some_test_url", "kuak", "bup", "123456789", "node_id"))
	if resulter.Error != nil {
		t.Error(fmt.Sprint(`Command expected to not get an error: `, resulter.Error))
	}
}

func TestNodeTagListCmdRightParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		method := r.Method
		path := r.URL
		if !strings.Contains(method, "GET") {
			diffString := StringDiff(method, "GET")
			t.Error(fmt.Sprintf("Command API method doesn't match. \n \n %s", diffString))
		}
		if !strings.Contains(path.String(), "agents/123456789/tags") {
			diffString := StringDiff(method, "agents/123456789/tags")
			t.Error(fmt.Sprintf("Command API path doesn't match. \n \n %s", diffString))
		}
	}))
	defer server.Close()
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra node tag list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=%s", "https://somewhere.com", server.URL, "token123", "123456789"))
}

func TestNodeTagListSuccessTable(t *testing.T) {
	// set test server
	server := TestServer(200, `{"test1":"test1","test2":"test2","test3":"test 3","test4":"test 4"}`, map[string]string{})
	defer server.Close()
	want := `+-------+--------+
|  KEY  | VALUE  |
+-------+--------+
| test1 | test1  |
| test2 | test2  |
| test3 | test 3 |
| test4 | test 4 |
+-------+--------+`

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node tag list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=123456789", "http://somewhere.com", server.URL, "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestNodeTagListSuccessJson(t *testing.T) {
	// set test server
	want := `{"test1":"test1","test2":"test2","test3":"test 3","test4":"test 4"}`
	server := TestServer(200, want, map[string]string{})
	defer server.Close()

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node tag list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=123456789 --json", "http://somewhere.com", server.URL, "token123"))

	eq, err := JsonDiff(want, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}
