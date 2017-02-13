package cmd

import (
	"strings"
	"testing"
)

func TestNodeFactCmdShouldDoNothing(t *testing.T) {
	ResetFlags()
	// check
	resulter := FullCmdTester(RootCmd, "lyra node fact")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, NodeFactCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
