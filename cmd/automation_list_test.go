package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestAutomationWithNoEndpointTokenCommand(t *testing.T) {
	// clean env variablen
	os.Unsetenv(ENV_VAR_TOKEN_NAME)
	os.Unsetenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)

	AutomationCmd.AddCommand(AutomationListCmd)
	resulter := FullCmdTester(AutomationCmd, "automation list")

	if resulter.Error == nil {
		t.Error(`Automation list expected to get an error with wrong flag`)
	}
}

func TestAutomationWithEnvEndpointTokenCommand(t *testing.T) {}

func TestAutomationEndpointTokenFlagListCommand(t *testing.T) {
	responseBody := "Miau"
	server := TestServer(200, responseBody)
	defer server.Close()

	AutomationCmd.AddCommand(AutomationListCmd)
	resulter := FullCmdTester(AutomationCmd, fmt.Sprintf("automation list --automation-endpoint=%s --token=%s", server.URL, "token123"))

	if !strings.Contains(resulter.Output, responseBody) {
		t.Error(`Automation list response body doesn't match.'`)
	}

	if resulter.Error != nil {
		t.Error(`Automation list expected to not get an error`)
	}
}

func TestAutomationWrongFlagsListCommand(t *testing.T) {
	responseBody := "Miau"
	server := TestServer(200, responseBody)
	defer server.Close()

	AutomationCmd.AddCommand(AutomationListCmd)
	resulter := FullCmdTester(AutomationCmd, fmt.Sprintf("automation list --automation-endpoint=%s --token=%s --wrong-flog=something", server.URL, "token123"))

	if resulter.Error == nil {
		t.Error(`Automation list expected to get an error with wrong flag`)
	}
}
