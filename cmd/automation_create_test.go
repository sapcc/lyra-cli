package cmd

import (
	"strings"
	"testing"
)

func resetAutomationCreateFlagVars() {
	resetAutomationFlagVars()
	AutomationCmd.ResetCommands()
	AutomationCmd.AddCommand(AutomationCreateCmd)
}

func TestAutomationCreateShouldDoNothing(t *testing.T) {
	resetAutomationCreateFlagVars()
	// check
	resulter := FullCmdTester(AutomationCmd, "automation create")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, AutomationCreateCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
