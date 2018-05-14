package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func TestNodeFactListCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra node fact list")
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra node fact list")
	ResetFlags()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra node fact list")
}

func TestNodeFactListCmdWithUserAuthenticationFlags(t *testing.T) {
	testServer := TestServer(200, "{}", map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node fact list --auth-url=%s --user-id=%s --project-id=%s --password=%s --node-id=%s", "some_test_url", "grrrrr", "bup", "123456789", "node id"))
	if resulter.Error != nil {
		t.Error(fmt.Sprint(`Command expected to not get an error: `, resulter.Error))
	}
}

func TestNodeFactListCmdRightParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		method := r.Method
		path := r.URL
		if !strings.Contains(method, "GET") {
			diffString := StringDiff(method, "GET")
			t.Error(fmt.Sprintf("Command API method doesn't match. \n \n %s", diffString))
		}
		if !strings.Contains(path.String(), "agents/node_id/facts") {
			diffString := StringDiff(method, "agents/node_id/facts")
			t.Error(fmt.Sprintf("Command API path doesn't match. \n \n %s", diffString))
		}
	}))
	defer server.Close()
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra node fact list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=%s", "https://somewhere.com", server.URL, "token123", "node_id"))
}

var nodeTagListResponseBody = `{
"arc_version": "20170209.2 (5807b68), go1.7.4",
"default_gateway": "10.0.0.1",
"default_interface": "eth0",
"domain": "novalocal",
"fqdn": "rel7.novalocal",
"hostname": "rel7",
"identity": "aa50283d-81d3-40f0-8bbd-42fe1751bff0",
"init_package": "systemd",
"ipaddress": "10.0.218.62",
"macaddress": "fa:16:3e:1d:fa:35",
"memory_available": 1.742552e+09,
"memory_total": 1.883996e+09,
"memory_used": 3.19296e+08,
"memory_used_percent": 8,
"online": true,
"organization": "abcdefghyjklmnopqrstuwxyz9876543",
"os": "linux",
"platform": "redhat",
"platform_family": "rhel",
"platform_version": "7.3",
"project": "abcdefghijklmnopqrstuwxyz1234567"
}`

func TestNodeFactListSuccessTable(t *testing.T) {
	// set test server
	server := TestServer(200, nodeTagListResponseBody, map[string]string{})
	defer server.Close()
	want := `+---------------------+--------------------------------------+
|         KEY         |                VALUE                 |
+---------------------+--------------------------------------+
| arc_version         | 20170209.2                           |
|                     | (5807b68), go1.7.4                   |
| default_gateway     | 10.0.0.1                             |
| default_interface   | eth0                                 |
| domain              | novalocal                            |
| fqdn                | rel7.novalocal                       |
| hostname            | rel7                                 |
| identity            | aa50283d-81d3-40f0-8bbd-42fe1751bff0 |
| init_package        | systemd                              |
| ipaddress           | 10.0.218.62                          |
| macaddress          | fa:16:3e:1d:fa:35                    |
| memory_available    | 1.742552e+09                         |
| memory_total        | 1.883996e+09                         |
| memory_used         | 3.19296e+08                          |
| memory_used_percent | 8                                    |
| online              | true                                 |
| organization        | abcdefghyjklmnopqrstuwxyz9876543     |
| os                  | linux                                |
| platform            | redhat                               |
| platform_family     | rhel                                 |
| platform_version    | 7.3                                  |
| project             | abcdefghijklmnopqrstuwxyz1234567     |
+---------------------+--------------------------------------+`

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node fact list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=123456789", "http://somewhere.com", server.URL, "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestNodeFactListSuccessJSON(t *testing.T) {
	// set test server
	server := TestServer(200, nodeTagListResponseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node fact list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --node-id=123456789 --json", "http://somewhere.com", server.URL, "token123"))

	eq, err := JsonDiff(nodeTagListResponseBody, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !eq {
		t.Error("Json response body and print out Json do not match.")
	}
}
