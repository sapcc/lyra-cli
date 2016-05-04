package cmd

import (
	"fmt"
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
	// set test server
	responseBody := "Miau"
	server := TestServer(200, responseBody)
	defer server.Close()

	resetAutomationList()
	CheckCmdWorksWithEndpointAndTokenFlag(t, AutomationCmd, fmt.Sprintf("automation list --automation-endpoint=%s --token=%s", server.URL, "token123"), responseBody)
}
