package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func resetAutomationList() {
	resetAutomationFlagVars()
	AutomationCmd.ResetCommands()
	AutomationCmd.AddCommand(AutomationListCmd)
}

func TestAutomationListCmdWithNoEnvEndpointAndTokenSet(t *testing.T) {
	resetAutomationList()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, AutomationCmd, "automation list")
	resetAutomationList()
	CheckhErrorWhenNoEnvEndpointSet(t, AutomationCmd, "automation list")
	resetAutomationList()
	CheckhErrorWhenNoEnvTokenSet(t, AutomationCmd, "automation list")
}

func TestAutomationListCmdWithEndpointTokenFlag(t *testing.T) {
	responseBody := "Miau"
	server := TestServer(200, responseBody)
	defer server.Close()

	resetAutomationList()
	resulter := FullCmdTester(AutomationCmd, fmt.Sprintf("automation list --automation-endpoint=%s --token=%s", server.URL, "token123"))

	if !strings.Contains(resulter.Output, responseBody) {
		t.Error(`Automation list response body doesn't match.'`)
	}

	if resulter.Error != nil {
		t.Error(`Automation list expected to not get an error`)
	}
}
