package cmd

import (
	"fmt"
	"testing"
)

func resetAutomationShow() {
	resetAutomationFlagVars()
	AutomationCmd.ResetCommands()
	AutomationCmd.AddCommand(AutomationShowCmd)
}

func TestAutomationShowCmdWithWrongEnvEndpointAndTokenSet(t *testing.T) {
	resetAutomationShow()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, AutomationCmd, "automation show")
	resetAutomationShow()
	CheckhErrorWhenNoEnvEndpointSet(t, AutomationCmd, "automation show")
	resetAutomationShow()
	CheckhErrorWhenNoEnvTokenSet(t, AutomationCmd, "automation show")
}

func TestAutomationShowCmdWithEndpointTokenFlag(t *testing.T) {
	// set test server
	responseBody := "Miau"
	server := TestServer(200, responseBody)
	defer server.Close()

	resetAutomationShow()
	CheckCmdWorksWithEndpointAndTokenFlag(t, AutomationCmd, fmt.Sprintf("automation show --automation-endpoint=%s --token=%s -i=%s", server.URL, "token123", "automation_id"), responseBody)
}
