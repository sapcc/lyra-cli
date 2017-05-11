package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func TestNodeListCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra node list")
	ResetFlags()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra node list")
	ResetFlags()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra node list")
}

func TestNodeListCmdWithUserAuthenticationFlags(t *testing.T) {
	testServer := TestServer(200, "[]", map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node list --auth-url=%s --user-id=%s --project-id=%s --password=%s", "some_test_url", "grrrrr", "bup", "123456789"))
	if resulter.Error != nil {
		t.Error(fmt.Sprint(`Command expected to not get an error: `, resulter.Error))
	}
}

func TestNodeListCmdSuccess(t *testing.T) {
	// set test server
	responseBody := `[{"agent_id":"aa50283d-81d3-40f0-8bbd-42fe1751bff0","display_name":"test1","project":"abcdefghijklmnopqrstuwxyz1234567","organization":"abcdefghyjklmnopqrstuwxyz9876543","tags":{"name":"arturo"},"created_at":"2017-02-06T15:19:17.503575Z","updated_at":"2017-02-08T12:40:17.891757Z","updated_with":"b6825bae-1c0d-4805-a87b-4bd17bd05279","updated_by":"linux"}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+--------------------------------------+--------------+----------------------------------+----------------------------------+-----------------------------+-----------------------------+------------+--------------------------------------+
|               AGENT ID               | DISPLAY NAME |           ORGANIZATION           |             PROJECT              |         CREATED AT          |         UPDATED AT          | UPDATED BY |             UPDATED WITH             |
+--------------------------------------+--------------+----------------------------------+----------------------------------+-----------------------------+-----------------------------+------------+--------------------------------------+
| aa50283d-81d3-40f0-8bbd-42fe1751bff0 | test1        | abcdefghyjklmnopqrstuwxyz9876543 | abcdefghijklmnopqrstuwxyz1234567 | 2017-02-06T15:19:17.503575Z | 2017-02-08T12:40:17.891757Z | linux      | b6825bae-1c0d-4805-a87b-4bd17bd05279 |
+--------------------------------------+--------------+----------------------------------+----------------------------------+-----------------------------+-----------------------------+------------+--------------------------------------+`

	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "https://somewhere.com", server.URL, "token123"))
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
	}
}

func TestNodeListCmdWithPaginationResultTable(t *testing.T) {
	// set test server
	server := nodePaginationServer()
	defer server.Close()
	want := `+----------+--------------+----------------------------------+----------------------------------+-----------------------------+-----------------------------+------------+--------------------------------------+
| AGENT ID | DISPLAY NAME |           ORGANIZATION           |             PROJECT              |         CREATED AT          |         UPDATED AT          | UPDATED BY |             UPDATED WITH             |
+----------+--------------+----------------------------------+----------------------------------+-----------------------------+-----------------------------+------------+--------------------------------------+
| 1        | test1        | abcdefghyjklmnopqrstuwxyz9876543 | abcdefghijklmnopqrstuwxyz1234567 | 2017-02-06T15:19:17.503575Z | 2017-02-08T12:40:17.891757Z | linux      | b6825bae-1c0d-4805-a87b-4bd17bd05279 |
| 2        | test2        | abcdefghyjklmnopqrstuwxyz9876543 | abcdefghijklmnopqrstuwxyz1234567 | 2017-02-06T15:19:17.503575Z | 2017-02-08T12:40:17.891757Z | linux      | b6825bae-1c0d-4805-a87b-4bd17bd05279 |
| 3        | test3        | abcdefghyjklmnopqrstuwxyz9876543 | abcdefghijklmnopqrstuwxyz1234567 | 2017-02-06T15:19:17.503575Z | 2017-02-08T12:40:17.891757Z | linux      | b6825bae-1c0d-4805-a87b-4bd17bd05279 |
+----------+--------------+----------------------------------+----------------------------------+-----------------------------+-----------------------------+------------+--------------------------------------+`

	ResetFlags()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "http://somewhere.com", server.URL, "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestNodeListCmdWithPaginationResultJson(t *testing.T) {
	// set test server
	server := nodePaginationServer()
	defer server.Close()
	responseBody := `[{"agent_id":"1","display_name":"test1","project":"abcdefghijklmnopqrstuwxyz1234567","organization":"abcdefghyjklmnopqrstuwxyz9876543","tags":{"name":"arturo"},"created_at":"2017-02-06T15:19:17.503575Z","updated_at":"2017-02-08T12:40:17.891757Z","updated_with":"b6825bae-1c0d-4805-a87b-4bd17bd05279","updated_by":"linux"},
	{"agent_id":"2","display_name":"test2","project":"abcdefghijklmnopqrstuwxyz1234567","organization":"abcdefghyjklmnopqrstuwxyz9876543","tags":{"name":"arturo"},"created_at":"2017-02-06T15:19:17.503575Z","updated_at":"2017-02-08T12:40:17.891757Z","updated_with":"b6825bae-1c0d-4805-a87b-4bd17bd05279","updated_by":"linux"},
	{"agent_id":"3","display_name":"test3","project":"abcdefghijklmnopqrstuwxyz1234567","organization":"abcdefghyjklmnopqrstuwxyz9876543","tags":{"name":"arturo"},"created_at":"2017-02-06T15:19:17.503575Z","updated_at":"2017-02-08T12:40:17.891757Z","updated_with":"b6825bae-1c0d-4805-a87b-4bd17bd05279","updated_by":"linux"}]`

	ResetFlags()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra node list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --json", "http://somewhere.com", server.URL, "token123"))

	eq, err := JsonListDiff(responseBody, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}

func TestNodeListCmdWithFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL

		if !strings.Contains(path.Query()["q"][0], "selector_filter") {
			diffString := StringDiff(path.Query()["q"][0], "selector_filter")
			t.Error(fmt.Sprintf("Command API selector doesn't match. \n \n %s", diffString))
		}
	}))
	defer server.Close()
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra node list --selector=%s --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "selector_filter", "https://somewhere.com", server.URL, "token123"))
}

func TestNodeListCmdRightParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		method := r.Method
		path := r.URL
		if !strings.Contains(method, "GET") {
			diffString := StringDiff(method, "GET")
			t.Error(fmt.Sprintf("Command API method doesn't match. \n \n %s", diffString))
		}
		if !strings.Contains(path.String(), "agents") {
			diffString := StringDiff(method, "agents")
			t.Error(fmt.Sprintf("Command API path doesn't match. \n \n %s", diffString))
		}
	}))
	defer server.Close()
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra node list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "https://somewhere.com", server.URL, "token123"))
}

func nodePaginationServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		page := r.URL.Query().Get("page")
		if page == "1" {
			w.Header().Set("Pagination-Elements", "3")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"agent_id":"1","display_name":"test1","project":"abcdefghijklmnopqrstuwxyz1234567","organization":"abcdefghyjklmnopqrstuwxyz9876543","tags":{"name":"arturo"},"created_at":"2017-02-06T15:19:17.503575Z","updated_at":"2017-02-08T12:40:17.891757Z","updated_with":"b6825bae-1c0d-4805-a87b-4bd17bd05279","updated_by":"linux"}]`)
		} else if page == "2" {
			w.Header().Set("Pagination-Elements", "3")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"agent_id":"2","display_name":"test2","project":"abcdefghijklmnopqrstuwxyz1234567","organization":"abcdefghyjklmnopqrstuwxyz9876543","tags":{"name":"arturo"},"created_at":"2017-02-06T15:19:17.503575Z","updated_at":"2017-02-08T12:40:17.891757Z","updated_with":"b6825bae-1c0d-4805-a87b-4bd17bd05279","updated_by":"linux"}]`)
		} else if page == "3" {
			w.Header().Set("Pagination-Elements", "3")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"agent_id":"3","display_name":"test3","project":"abcdefghijklmnopqrstuwxyz1234567","organization":"abcdefghyjklmnopqrstuwxyz9876543","tags":{"name":"arturo"},"created_at":"2017-02-06T15:19:17.503575Z","updated_at":"2017-02-08T12:40:17.891757Z","updated_with":"b6825bae-1c0d-4805-a87b-4bd17bd05279","updated_by":"linux"}]`)
		}
	}))
	return server
}
