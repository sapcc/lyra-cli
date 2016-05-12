package cmd

import (
	"fmt"
	"testing"

	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/restclient"
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
	responseBody := `{"miau":"bup"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	resetAutomationList()
	CheckCmdWorksWithEndpointAndTokenFlag(t, RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --token=%s", server.URL, "token123"), responseBody)
}

func TestAutomationListCmdWithPagination(t *testing.T) {
	// set test server
	responseBody := `{"miau":"test"}`
	server := TestServer(200, responseBody, map[string]string{"Pagination-Page": "1", "Pagination-Per-Page": "2", "Pagination-Pages": "3"})
	defer server.Close()

	resetAutomationList()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation list --lyra-service-endpoint=%s --token=%s", server.URL, "token123"))

	pagData := restclient.PagResp{}
	helpers.JSONStringToStructure(string(resulter.Output), &pagData)

	if pagData.Pagination.Page != 1 {
		t.Error(`Automation list command pagination response doesn't match.'`)
	}
	if pagData.Pagination.PerPage != 2 {
		t.Error(`Automation list command pagination response doesn't match.'`)
	}
	if pagData.Pagination.Pages != 3 {
		t.Error(`Automation list command pagination response doesn't match.'`)
	}
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
}
