package cmd

import (
	"testing"
)

func resetJobShow() {
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	JobCmd.ResetCommands()
	JobShowCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(JobCmd)
	JobCmd.AddCommand(JobShowCmd)
}

func TestJobShowCmdWithWrongEnvEndpointAndTokenSet(t *testing.T) {
	resetJobShow()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra job show")
	resetJobShow()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra job show")
	resetJobShow()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra job show")
}
