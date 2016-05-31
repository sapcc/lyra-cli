package cmd

import (
	"strings"
	"testing"
)

func resetJob() {
	// reset automation flag vars
	ResetFlags()
}

func TestJobShouldDoNothing(t *testing.T) {
	resetJob()
	// check
	resulter := FullCmdTester(RootCmd, "lyra job")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, JobCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
