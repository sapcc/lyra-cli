package cmd

import (
	"fmt"
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
		t.Error(fmt.Sprint(`Command expected to not get an error. `, resulter.Error))
	}

	if !strings.Contains(resulter.Output, VersionCmd.Long) {
		diffString := StringDiff(resulter.Output, AutomationUpdateCmd.Long)
		t.Errorf("Command response doesn't match. \n \n %s", diffString)
	}
}
