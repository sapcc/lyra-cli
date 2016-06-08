package cmd

import (
	"strings"
	"testing"
)

func resetBashCompletion() {
	// reset automation flag vars
	ResetFlags()
}

func TestBashCompletionAuthenticationNotRequired(t *testing.T) {
	resetBashCompletion()
	// check
	resulter := FullCmdTester(RootCmd, "lyra bash-completion")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, AutomationUpdateCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
