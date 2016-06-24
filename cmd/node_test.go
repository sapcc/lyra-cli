package cmd

import (
	"strings"
	"testing"
)

func resetNode() {
	// reset automation flag vars
	ResetFlags()
}

func TestNodeCmdShouldDoNothing(t *testing.T) {
	resetNode()
	// check
	resulter := FullCmdTester(RootCmd, "lyra node")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, NodeCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
