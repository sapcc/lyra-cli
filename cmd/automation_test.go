package cmd

import (
	"strings"
	"testing"
)

func resetAutomation() {
	// reset automation flag vars
	ResetFlags()
}

func TestAutomationShouldDoNothing(t *testing.T) {
	resetAutomation()
	// check
	resulter := FullCmdTester(RootCmd, "lyra automation")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, AutomationCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
