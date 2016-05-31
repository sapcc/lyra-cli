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

func newMockAuthenticationV3JobLog(authOpts LyraAuthOps) Authentication {
	// set test server
	responseBody := `This is a job log`
	server := TestServer(200, responseBody, map[string]string{})

	return &MockV3{AuthOpts: authOpts, TestServer: server}
}

func TestJobLogCmdWithAuthenticationFlags(t *testing.T) {
	// mock interface for authenticationt test
	AuthenticationV3 = newMockAuthenticationV3JobLog
	want := `This is a job log`

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job log --auth-url=%s --user-id=%s --project-id=%s --password=%s --job-id=123456789", "some_test_url", "miau", "bup", "123456789"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
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
