package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetJobList() {
	// reset automation flag vars
	ResetFlags()
}

func TestJobListCmdWithAuthenticationFlags(t *testing.T) {
	responseBody := `[{"request_id": "10ed3681-0f48-4564-a3fd-c7dfcb0c87c1", "agent": "execute", "action": "tarball", "status": "complete", "created_at": "2016-06-24T11:52:06.834057Z", "user": {"name": "user123"}}]`
	testServer := TestServer(200, responseBody, map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	want := `+--------------------------------------+----------+---------+---------+-------------------+-----------------------------+
|              REQUEST ID              |  STATUS  | ACTION  |  AGENT  |       USER        |         CREATED AT          |
+--------------------------------------+----------+---------+---------+-------------------+-----------------------------+
| 10ed3681-0f48-4564-a3fd-c7dfcb0c87c1 | complete | tarball | execute | map[name:user123] | 2016-06-24T11:52:06.834057Z |
+--------------------------------------+----------+---------+---------+-------------------+-----------------------------+`

	// reset stuff
	resetJobList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job list --auth-url=%s --user-id=%s --project-id=%s --password=%s", "some_test_url", "miau", "bup", "123456789"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response body doesn't match. \n \n %s", diffString)
	}
}

func TestJobListCmdWithNoEnvEndpointsAndTokenSet(t *testing.T) {
	resetJobList()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra job list")
	resetJobList()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra job list")
	resetJobList()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra job list")
}

func TestJobListCmdResultTable(t *testing.T) {
	// set test server
	responseBody := `[{"request_id": "10ed3681-0f48-4564-a3fd-c7dfcb0c87c1", "agent": "execute", "action": "tarball", "status": "complete", "created_at": "2016-06-24T11:52:06.834057Z", "user": {"name": "user123"}}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+--------------------------------------+----------+---------+---------+-------------------+-----------------------------+
|              REQUEST ID              |  STATUS  | ACTION  |  AGENT  |       USER        |         CREATED AT          |
+--------------------------------------+----------+---------+---------+-------------------+-----------------------------+
| 10ed3681-0f48-4564-a3fd-c7dfcb0c87c1 | complete | tarball | execute | map[name:user123] | 2016-06-24T11:52:06.834057Z |
+--------------------------------------+----------+---------+---------+-------------------+-----------------------------+`

	// reset stuff
	resetJobList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "http://somewhere.com", server.URL, "token123"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
		return
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response body doesn't match. \n \n %s", diffString)
	}
}

func TestJobListCmdWithResultJSON(t *testing.T) {
	// set test server
	responseBody := `[{"request_id": "f1b18c11-5838-44d2-8651-66aa4083bd19", "agent": "chef", "action": "zero", "status": "failed", "created_at": "2016-04-07T15:47:02.260715Z", "user_id": "u-fa35bbc5f"}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetJobList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --json", "http://somewhere.com", server.URL, "token123"))

	eq, err := JsonListDiff(responseBody, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !eq {
		t.Error("Json response body and print out Json do not match.")
	}
}

func TestJobListCmdWithPaginationResultTable(t *testing.T) {
	// set test server
	server := jobPaginationServer()
	defer server.Close()
	want := `+------------+--------+--------+-------+-------------------+-----------------------------+
| REQUEST ID | STATUS | ACTION | AGENT |       USER        |         CREATED AT          |
+------------+--------+--------+-------+-------------------+-----------------------------+
| 1          | failed | zero   | chef  | map[name:user123] | 2016-04-07T15:47:02.260715Z |
| 2          | failed | zero   | chef  | map[name:user123] | 2016-04-07T15:47:12.260715Z |
| 3          | failed | zero   | chef  | map[name:user123] | 2016-04-07T15:47:22.260715Z |
+------------+--------+--------+-------+-------------------+-----------------------------+`

	resetJobList()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", "http://somewhere.com", server.URL, "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response body doesn't match. \n \n %s", diffString)
	}
}

func TestJobListCmdWithPaginationResultJSON(t *testing.T) {
	// set test server
	server := jobPaginationServer()
	defer server.Close()

	resetJobList()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --json", "http://somewhere.com", server.URL, "token123"))

	responseBody := `[{"request_id": "1", "agent": "chef", "action": "zero", "status": "failed", "created_at": "2016-04-07T15:47:02.260715Z", "user": {"name": "user123"}},
{"request_id": "2", "agent": "chef", "action": "zero", "status": "failed", "created_at": "2016-04-07T15:47:12.260715Z", "user": {"name": "user123"}},
{"request_id": "3", "agent": "chef", "action": "zero", "status": "failed", "created_at": "2016-04-07T15:47:22.260715Z", "user": {"name": "user123"}}]`

	eq, err := JsonListDiff(responseBody, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !eq {
		t.Error("Json response body and print out Json do not match.")
	}
}

func jobPaginationServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		page := r.URL.Query().Get("page")
		if page == "1" {
			w.Header().Set("Pagination-Page", "1")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"request_id": "1", "agent": "chef", "action": "zero", "status": "failed", "created_at": "2016-04-07T15:47:02.260715Z", "user": {"name": "user123"}}]`)
		} else if page == "2" {
			w.Header().Set("Pagination-Page", "2")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"request_id": "2", "agent": "chef", "action": "zero", "status": "failed", "created_at": "2016-04-07T15:47:12.260715Z", "user": {"name": "user123"}}]`)
		} else if page == "3" {
			w.Header().Set("Pagination-Page", "3")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"request_id": "3", "agent": "chef", "action": "zero", "status": "failed", "created_at": "2016-04-07T15:47:22.260715Z", "user": {"name": "user123"}}]`)
		}
	}))
	return server
}
