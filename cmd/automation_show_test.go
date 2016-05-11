package cmd

import (
	"fmt"
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
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra-cli automation show")
	resetAutomationShow()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra-cli automation show")
	resetAutomationShow()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra-cli automation show")
}

func TestAutomationShowCmdWithEndpointTokenFlag(t *testing.T) {
	// set test server
	responseBody := `{"miau":"bup"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	resetAutomationShow()
	CheckCmdWorksWithEndpointAndTokenFlag(t, RootCmd, fmt.Sprintf("lyra-cli automation show --lyra-service-endpoint=%s --token=%s -i=%s", server.URL, "token123", "automation_id"), responseBody)
}
