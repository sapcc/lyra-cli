package cmd

import (
	"fmt"
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
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra-cli automation list")
	resetAutomationList()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra-cli automation list")
	resetAutomationList()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra-cli automation list")
}

func TestAutomationListCmdWithEndpointTokenFlag(t *testing.T) {
	// set test server
	responseBody := "Miau"
	server := TestServer(200, responseBody)
	defer server.Close()

	resetAutomationList()
	CheckCmdWorksWithEndpointAndTokenFlag(t, RootCmd, fmt.Sprintf("lyra-cli automation list --lyra-service-endpoint=%s --token=%s", server.URL, "token123"), responseBody)
}
