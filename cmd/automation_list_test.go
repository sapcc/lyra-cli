package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetAutomationList() {
	// reset automation flag vars
	ResetFlags()
}

func newMockAuthenticationV3AutomationList(authOpts auth.AuthOptions) auth.Authentication {
	// set test server
	responseBody := `[{"id":"6","name":"Chef_test","type":"Chef", "timeout":"3600", "repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`
	server := TestServer(200, responseBody, map[string]string{})

	return &auth.MockV3{Options: authOpts, TestServer: server}
}

func TestAutomationListCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	resetAutomationList()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra automation list")
	resetAutomationList()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra automation list")
	resetAutomationList()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra automation list")
}

func TestAutomationListCmdWithEndpointsTokenFlag(t *testing.T) {
	// set test server
	responseBody := `[{"name":"bup"}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, server.URL, "token123"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
}

func TestAutomationListCmdWithAuthenticationFlags(t *testing.T) {
	// mock interface for authenticationt test
	auth.AuthenticationV3 = newMockAuthenticationV3AutomationList
	want := `+----+-----------+------+---------------------------------------------------------+---------------------+-------+---------+
| ID |   NAME    | TYPE |                       REPOSITORY                        | REPOSITORY REVISION | TAGS  | TIMEOUT |
+----+-----------+------+---------------------------------------------------------+---------------------+-------+---------+
| 6  | Chef_test | Chef | https://github.com/user123/automation-test.git | master              | <nil> | 3600    |
+----+-----------+------+---------------------------------------------------------+---------------------+-------+---------+`

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --auth-url=%s --user-id=%s --project-id=%s --password=%s", "some_test_url", "miau", "bup", "123456789"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationListCmdResultTable(t *testing.T) {
	// set test server
	responseBody := `[{"id":"6","name":"Chef_test", "type":"Chef", "timeout":"3600", "repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+----+-----------+------+---------------------------------------------------------+---------------------+-------+---------+
| ID |   NAME    | TYPE |                       REPOSITORY                        | REPOSITORY REVISION | TAGS  | TIMEOUT |
+----+-----------+------+---------------------------------------------------------+---------------------+-------+---------+
| 6  | Chef_test | Chef | https://github.com/user123/automation-test.git | master              | <nil> | 3600    |
+----+-----------+------+---------------------------------------------------------+---------------------+-------+---------+`

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "http://somewhere.com", "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationListCmdWithResultJSON(t *testing.T) {
	// set test server
	responseBody := `[{"id":"6","name":"Chef_test","repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --json", server.URL, "http://somewhere.com", "token123"))

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

func TestAutomationListCmdResultTableExtraCustomColumns(t *testing.T) {}

func TestAutomationListCmdWithPaginationResultTable(t *testing.T) {
	// set test server
	server := automationPaginationServer()
	defer server.Close()

	want := `+----+------------+------+---------------------------------------------------------+---------------------+-------+---------+
| ID |    NAME    | TYPE |                       REPOSITORY                        | REPOSITORY REVISION | TAGS  | TIMEOUT |
+----+------------+------+---------------------------------------------------------+---------------------+-------+---------+
| 1  | Chef_test1 | Chef | https://github.com/user123/automation-test.git | master              | <nil> | 3600    |
| 2  | Chef_test2 | Chef | https://github.com/user123/automation-test.git | master              | <nil> | 3600    |
| 3  | Chef_test3 | Chef | https://github.com/user123/automation-test.git | master              | <nil> | 3600    |
+----+------------+------+---------------------------------------------------------+---------------------+-------+---------+`

	resetAutomationList()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s", server.URL, "http://somewhere.com", "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationListCmdWithPaginationResultJSON(t *testing.T) {
	// set test server
	server := automationPaginationServer()
	defer server.Close()

	resetAutomationList()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --json", server.URL, "http://somewhere.com", "token123"))

	responseBody := `[{"id":"1","name":"Chef_test1", "type":"Chef", "timeout":"3600", "repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"},
{"id":"2","name":"Chef_test2", "type":"Chef", "timeout":"3600", "repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"},
{"id":"3","name":"Chef_test3", "type":"Chef", "timeout":"3600", "repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`

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

func automationPaginationServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		page := r.URL.Query().Get("page")
		if page == "1" {
			w.Header().Set("Pagination-Page", "1")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id":"1","name":"Chef_test1", "type":"Chef", "timeout":"3600", "repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`)
		} else if page == "2" {
			w.Header().Set("Pagination-Page", "2")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id":"2","name":"Chef_test2", "type":"Chef", "timeout":"3600", "repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`)
		} else if page == "3" {
			w.Header().Set("Pagination-Page", "3")
			w.Header().Set("Pagination-Per-Page", "1")
			w.Header().Set("Pagination-Pages", "3")
			w.WriteHeader(200) // keep the code after setting headers. If not they will disapear...
			fmt.Fprintln(w, `[{"id":"3","name":"Chef_test3", "type":"Chef", "timeout":"3600", "repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}]`)
		}
	}))
	return server
}
