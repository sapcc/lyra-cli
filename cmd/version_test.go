package cmd

func resetVersion() {
	// reset automation flag vars
	ResetFlags()
}

// func TestVersionAuthenticationNotRequired(t *testing.T) {
// 	resetVersion()
// 	// check
// 	resulter := FullCmdTester(RootCmd, "lyra version")
// 	if resulter.Error != nil {
// 		t.Error(fmt.Sprint(`Command expected to not get an error. `, resulter.Error))
// 	}
//
// 	if !strings.Contains(resulter.Output, AutomationUpdateCmd.Long) {
// 		diffString := StringDiff(resulter.Output, AutomationUpdateCmd.Long)
// 		t.Error(fmt.Sprintf("Command response doesn't match. \n \n %s", diffString))
// 	}
// }
