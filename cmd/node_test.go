package cmd

import (
	"strings"
	"testing"
)

func TestNodeCmdShouldDoNothing(t *testing.T) {
	ResetFlags()
	// check
	resulter := FullCmdTester(RootCmd, "lyra node")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, NodeCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
