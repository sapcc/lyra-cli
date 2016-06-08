package cmd

import (
	"strings"
	"testing"
)

func resetVersion() {
	// reset automation flag vars
	ResetFlags()
}

func TestVersionAuthenticationNotRequired(t *testing.T) {
	resetVersion()
	// check
	resulter := FullCmdTester(RootCmd, "lyra version")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, AutomationUpdateCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
