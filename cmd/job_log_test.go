package cmd

import (
	"fmt"
	"strings"
	"testing"
)

func resetJobLog() {
	jobId = ""
	// reset automation flag vars
	ResetFlags()
}

func TestJobLogCmdWithWrongEnvEndpointAndTokenSet(t *testing.T) {
	resetJobLog()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra job log")
	resetJobLog()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra job log")
	resetJobLog()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra job log")
}

func TestJobLogCmdResult(t *testing.T) {
	// set test server
	responseBody := `This is a job log`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetJobLog()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job log --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --job-id=123456789", "http://somewhere.com", server.URL, "token123"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
		return
	}

	if !strings.Contains(resulter.Output, responseBody) {
		diffString := StringDiff(resulter.Output, responseBody)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}
