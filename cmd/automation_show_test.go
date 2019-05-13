package cmd

import (
	"fmt"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetAutomationShow() {
	// reset flags
	ResetFlags()
}

func TestAutomationShowCmdWithUserAuthenticationFlags(t *testing.T) {
	responseBody := `{"id":"1","name":"Chef_test1","repository":"https://github.com/userId0123456789/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}`
	testServer := TestServer(200, responseBody, map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	want := `+---------------------+---------------------------------------------------------+
|         KEY         |                          VALUE                          |
+---------------------+---------------------------------------------------------+
| arguments           | {}                                                      |
| chef_attributes     | map[test:test]                                          |
| id                  | 1                                                       |
| log_level           | info                                                    |
| name                | Chef_test1                                              |
| repository          | https://github.com/userId0123456789/automation-test.git |
| repository_revision | master                                                  |
| run_list            | [recipe[nginx]]                                         |
+---------------------+---------------------------------------------------------+`

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation show --auth-url=%s --user-id=%s --project-id=%s --password=%s --automation-id=%s", "some_test_url", "miau", "bup", "123456789", "automation_id"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationShowCmdWithWrongEnvEndpointAndTokenSet(t *testing.T) {
	resetAutomationShow()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra automation show")
	resetAutomationShow()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra automation show")
	resetAutomationShow()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra automation show")
}

func TestAutomationShowCmdWithResultTable(t *testing.T) {
	// set test server
	responseBody := `{"id":"1","name":"Chef_test1","repository":"https://github.com/userId0123456789/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+---------------------+---------------------------------------------------------+
|         KEY         |                          VALUE                          |
+---------------------+---------------------------------------------------------+
| arguments           | {}                                                      |
| chef_attributes     | map[test:test]                                          |
| id                  | 1                                                       |
| log_level           | info                                                    |
| name                | Chef_test1                                              |
| repository          | https://github.com/userId0123456789/automation-test.git |
| repository_revision | master                                                  |
| run_list            | [recipe[nginx]]                                         |
+---------------------+---------------------------------------------------------+`

	// reset stuff
	resetAutomationShow()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --automation-id=%s", server.URL, "https://somewhere.com", "token123", "automation_id"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationShowCmdWithResultJSON(t *testing.T) {
	// set test server
	responseBody := `{"id":"1","name":"Chef_test1","repository":"https://github.com/userId0123456789/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetAutomationShow()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --automation-id=%s --json", server.URL, "https://somewhere.com", "token123", "automation_id"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
		return
	}

	eq, err := JsonDiff(responseBody, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !eq {
		t.Error("Json response body and print out Json do not match.")
	}
}
