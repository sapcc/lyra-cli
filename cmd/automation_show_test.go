package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func resetAutomationShow() {
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	AutomationCmd.ResetCommands()
	AutomationShowCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(AutomationCmd)
	AutomationCmd.AddCommand(AutomationShowCmd)
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
	responseBody := `{"id":"1","name":"Chef_test1","repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}`
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
| repository          | https://github.com/user123/automation-test.git |
| repository_revision | master                                                  |
| run_list            | [recipe[nginx]]                                         |
+---------------------+---------------------------------------------------------+`

	// reset stuff
	resetAutomationShow()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s -i=%s", server.URL, "https://somewhere.com", "token123", "automation_id"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationShowCmdWithResultJSON(t *testing.T) {
	// set test server
	responseBody := `{"id":"1","name":"Chef_test1","repository":"https://github.com/user123/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetAutomationShow()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s -i=%s --json", server.URL, "https://somewhere.com", "token123", "automation_id"))

	source := map[string]interface{}{}
	err := json.Unmarshal([]byte(responseBody), &source)
	if err != nil {
		t.Error(err.Error())
		return
	}
	response := map[string]interface{}{}
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
