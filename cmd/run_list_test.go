package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func resetRunList() {
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	RunCmd.ResetCommands()
	RunListCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(RunCmd)
	RunCmd.AddCommand(RunListCmd)
}

func TestRunListCmdWithNoEnvEndpointsAndTokenSet(t *testing.T) {
	resetRunList()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra run list")
	resetRunList()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra run list")
	resetRunList()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra run list")
}

func TestRunListCmdWithEndpointTokenFlag(t *testing.T) {
	// set test server
	responseBody := `[{"name":"bup"}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetRunList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra run list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "token123"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
}

func TestRunListCmdResultTable(t *testing.T) {
	// set test server
	responseBody := `[{
  "id": "30",
  "log": "Selecting nodes with selector @identity='0128e993-c709-4ce1-bccf-e06eb10900a0'\nSelected nodes:\n0128e993-c709-4ce1-bccf-e06eb10900a0 mo-85b92ea6f",
  "created_at": "2016-04-07T21:42:07.416Z",
  "updated_at": "2016-04-07T21:42:14.294Z",
  "repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc",
  "state": "executing",
  "jobs": [
    "b843bbe9-fa95-4a0b-9329-aed05d1de8b8"
  ],
  "owner": "u-fa35bbc5f",
  "automation_id": "6",
  "automation_name": "Chef_test",
  "selector": "@identity='0128e993-c709-4ce1-bccf-e06eb10900a0'"
}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+----+---------------+-----------------+-----------+-------------+--------------------------+
| ID | AUTOMATION ID | AUTOMATION NAME |   STATE   |    OWNER    |        CREATED AT        |
+----+---------------+-----------------+-----------+-------------+--------------------------+
| 30 | 6             | Chef_test       | executing | u-fa35bbc5f | 2016-04-07T21:42:07.416Z |
+----+---------------+-----------------+-----------+-------------+--------------------------+`

	// reset stuff
	resetRunList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra run list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "http://somewhere.com", "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestRunListCmdWithResultJSON(t *testing.T) {
	// set test server
	responseBody := `[{
  "id": "30",
  "log": "Selecting nodes with selector @identity='0128e993-c709-4ce1-bccf-e06eb10900a0'\nSelected nodes:\n0128e993-c709-4ce1-bccf-e06eb10900a0 mo-85b92ea6f",
  "created_at": "2016-04-07T21:42:07.416Z",
  "updated_at": "2016-04-07T21:42:14.294Z",
  "repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc",
  "state": "executing",
  "jobs": [
    "b843bbe9-fa95-4a0b-9329-aed05d1de8b8"
  ],
  "owner": "u-fa35bbc5f",
  "automation_id": "6",
  "automation_name": "Chef_test",
  "selector": "@identity='0128e993-c709-4ce1-bccf-e06eb10900a0'"
}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetRunList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra run list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --json", server.URL, "http://somewhere.com", "token123"))

	source := []map[string]interface{}{}
	err := json.Unmarshal([]byte(responseBody), &source)
	if err != nil {
		t.Error(err.Error())
		return
	}

	response := []map[string]interface{}{}
	err = json.Unmarshal([]byte(resulter.Output), &response)
	if err != nil {
		t.Error(err.Error())
		return
	}

	eq := reflect.DeepEqual(source, response)
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}

func TestRunListCmdResultTableExtraCustomColumns(t *testing.T) {}

func TestRunListCmdWithPaginationResultTable(t *testing.T) {
	// set test server
	server := runPaginationServer()
	defer server.Close()

	want := `+----+---------------+-----------------+-----------+-------------+--------------------------+
| ID | AUTOMATION ID | AUTOMATION NAME |   STATE   |    OWNER    |        CREATED AT        |
+----+---------------+-----------------+-----------+-------------+--------------------------+
| 1  | 6             | Chef_test       | executing | u-fa35bbc5f | 2016-04-07T21:42:07.416Z |
| 2  | 6             | Chef_test       | failed    | u-fa35bbc5f | 2016-04-07T21:42:17.416Z |
| 3  | 6             | Chef_test       | completed | u-fa35bbc5f | 2016-04-07T21:42:27.416Z |
+----+---------------+-----------------+-----------+-------------+--------------------------+`

	resetRunList()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra run list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "http://somewhere.com", "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestRunListCmdWithPaginationResultJSON(t *testing.T) {
	// set test server
	server := runPaginationServer()
	defer server.Close()

	resetRunList()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra run list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --json", server.URL, "http://somewhere.com", "token123"))

	responseBody := `[{"id": "1", "created_at": "2016-04-07T21:42:07.416Z","repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc","state": "executing","owner": "u-fa35bbc5f","automation_id": "6","automation_name": "Chef_test"},
{"id": "2","created_at": "2016-04-07T21:42:17.416Z","repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc","state": "failed","owner": "u-fa35bbc5f","automation_id": "6","automation_name": "Chef_test"},
{"id": "3","created_at": "2016-04-07T21:42:27.416Z","repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc","state": "completed","owner": "u-fa35bbc5f","automation_id": "6","automation_name": "Chef_test"}]`

	source := []map[string]interface{}{}
	err := json.Unmarshal([]byte(responseBody), &source)
	if err != nil {
		t.Error(err.Error())
		return
	}

	response := []map[string]interface{}{}
	err = json.Unmarshal([]byte(resulter.Output), &response)
	if err != nil {
		t.Error(err.Error())
		return
	}

	eq := reflect.DeepEqual(source, response)
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}

func runPaginationServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		page := r.URL.Query().Get("page")
		if page == "1" {
			w.Header().Set("Pagination-Page", "1")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id": "1", "created_at": "2016-04-07T21:42:07.416Z","repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc","state": "executing","owner": "u-fa35bbc5f","automation_id": "6","automation_name": "Chef_test"}]`)
		} else if page == "2" {
			w.Header().Set("Pagination-Page", "2")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id": "2","created_at": "2016-04-07T21:42:17.416Z","repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc","state": "failed","owner": "u-fa35bbc5f","automation_id": "6","automation_name": "Chef_test"}]`)
		} else if page == "3" {
			w.Header().Set("Pagination-Page", "3")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id": "3","created_at": "2016-04-07T21:42:27.416Z","repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc","state": "completed","owner": "u-fa35bbc5f","automation_id": "6","automation_name": "Chef_test"}]`)
		}
	}))
	return server
}
