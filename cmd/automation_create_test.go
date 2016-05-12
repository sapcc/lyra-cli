package cmd

import (
	"strings"
	"testing"
)

func resetAutomationCreateFlagVars() {
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	AutomationCmd.ResetCommands()
	AutomationCreateCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(AutomationCmd)
	AutomationCmd.AddCommand(AutomationCreateCmd)
}

func TestAutomationCreateShouldDoNothing(t *testing.T) {
	resetAutomationCreateFlagVars()
	// check
	resulter := FullCmdTester(RootCmd, "lyra automation create")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, AutomationCreateCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
