package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetAutomationList() {
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	AutomationCmd.ResetCommands()
	AutomationListCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(AutomationCmd)
	AutomationCmd.AddCommand(AutomationListCmd)
}

func TestAutomationListCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	resetAutomationList()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra automation list")
	resetAutomationList()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra automation list")
	resetAutomationList()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra automation list")
}

func TestAutomationListCmdWithEndpointTokenFlag(t *testing.T) {
	// set test server
	responseBody := `[{"name":"bup"}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "token123"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
}

func TestAutomationListCmdResultTable(t *testing.T) {
	// set test server
	responseBody := `[{"id":"6","name":"Chef_test","repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+----+-----------+---------------------------------------------------------+---------------------+-----------------+-----------------+-----------+-----------+
| ID |   NAME    |                       REPOSITORY                        | REPOSITORY REVISION |    RUN LIST     | CHEF ATTRIBUTES | LOG LEVEL | ARGUMENTS |
+----+-----------+---------------------------------------------------------+---------------------+-----------------+-----------------+-----------+-----------+
| 6  | Chef_test | https://github.com/user123/automation-test.git | master              | [recipe[nginx]] | map[test:test]  | info      | {}        |
+----+-----------+---------------------------------------------------------+---------------------+-----------------+-----------------+-----------+-----------+`

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "http://somewhere.com", "token123"))

	if !strings.Contains(resulter.Output, want) {
		t.Error(`Command response body doesn't match.'`)
	}
}

func TestAutomationListCmdResultTableExtraCustomColumns(t *testing.T) {
	// pagData := restclient.PagResp{}
	// helpers.JSONStringToStructure(string(resulter.Output), &pagData)
	//
	// if pagData.Pagination.Page != 1 {
	//   t.Error(`Automation list command pagination response doesn't match.'`)
	// }
	// if pagData.Pagination.PerPage != 2 {
	//   t.Error(`Automation list command pagination response doesn't match.'`)
	// }
	// if pagData.Pagination.Pages != 3 {
	//   t.Error(`Automation list command pagination response doesn't match.'`)
	// }
	// if resulter.Error != nil {
	//   t.Error(`Command expected to not get an error`)
	// }
}

func TestAutomationListCmdWithPaginationResultTable(t *testing.T) {
	// set test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		page := r.URL.Query().Get("page")
		if page == "1" {
			w.Header().Set("Pagination-Page", "1")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id":"1","name":"Chef_test1","repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`)
		} else if page == "2" {
			w.Header().Set("Pagination-Page", "2")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id":"2","name":"Chef_test2","repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`)
		} else if page == "3" {
			w.Header().Set("Pagination-Page", "3")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id":"3","name":"Chef_test3","repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`)
		}
	}))
	defer server.Close()

	want := `+----+------------+---------------------------------------------------------+---------------------+-----------------+-----------------+-----------+-----------+
| ID |    NAME    |                       REPOSITORY                        | REPOSITORY REVISION |    RUN LIST     | CHEF ATTRIBUTES | LOG LEVEL | ARGUMENTS |
+----+------------+---------------------------------------------------------+---------------------+-----------------+-----------------+-----------+-----------+
| 1  | Chef_test1 | https://github.com/user123/automation-test.git | master              | [recipe[nginx]] | map[test:test]  | info      | {}        |
| 2  | Chef_test2 | https://github.com/user123/automation-test.git | master              | [recipe[nginx]] | map[test:test]  | info      | {}        |
| 3  | Chef_test3 | https://github.com/user123/automation-test.git | master              | [recipe[nginx]] | map[test:test]  | info      | {}        |
+----+------------+---------------------------------------------------------+---------------------+-----------------+-----------------+-----------+-----------+`

	resetAutomationList()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "http://somewhere.com", "token123"))

	if !strings.Contains(resulter.Output, want) {
		t.Error(`Command response body doesn't match.'`)
	}
}

func TestAutomationListCmdWithPaginationResultJSON(t *testing.T) {}
