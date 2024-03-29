package cmd

import (
	"strings"
	"testing"
)

func resetAutomationUpdate() {
	// reset flags
	ResetFlags()
}

func TestAutomationUpdateShouldDoNothing(t *testing.T) {
	resetAutomationUpdate()
	// check
	resulter := FullCmdTester(RootCmd, "lyra automation update")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, AutomationUpdateCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
