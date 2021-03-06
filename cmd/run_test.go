package cmd

import (
	"strings"
	"testing"
)

func resetRun() {
	// reset automation flag vars
	ResetFlags()
}

func TestRunShouldDoNothing(t *testing.T) {
	resetRun()
	// check
	resulter := FullCmdTester(RootCmd, "lyra run")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, RunCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
