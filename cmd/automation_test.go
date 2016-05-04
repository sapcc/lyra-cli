package cmd

import (
	"strings"
	"testing"
)

func resetAutomationFlagVars() {
	Token = ""
	AutomationUrl = ""
}

func resetAutomation() {
	resetAutomationFlagVars()
	AutomationCmd.ResetCommands()
}

func TestAutomationShouldDoNothing(t *testing.T) {
	resetAutomation()
	// check
	resulter := FullCmdTester(AutomationCmd, "automation")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, AutomationCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
