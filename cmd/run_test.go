package cmd

import (
	"strings"
	"testing"
)

func resetRun() {
	run = Run{}
	runId = ""
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	RunCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(RunCmd)
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
