package cmd

import (
	"strings"
	"testing"
)

func resetAutomation() {
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	AutomationCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(AutomationCmd)
}

func TestAutomationShouldDoNothing(t *testing.T) {
	resetAutomation()
	// check
	resulter := FullCmdTester(RootCmd, "lyra-cli automation")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, AutomationCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
