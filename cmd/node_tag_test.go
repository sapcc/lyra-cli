package cmd

import (
	"strings"
	"testing"
)

func TestNodeTagCmdShouldDoNothing(t *testing.T) {
	ResetFlags()
	// check
	resulter := FullCmdTester(RootCmd, "lyra node tag")
	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, NodeTagCmd.Long) {
		t.Error(`Command response body doesn't match.'`)
	}
}
