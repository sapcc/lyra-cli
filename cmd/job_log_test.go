package cmd

import (
	"fmt"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetJobLog() {
	// reset automation flag vars
	ResetFlags()
}

func TestJobLogCmdWithAuthenticationFlags(t *testing.T) {
	testServer := TestServer(200, `This is a job log`, map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	want := `This is a job log`

	// reset stuff
	resetJobLog()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job log --auth-url=%s --user-id=%s --project-id=%s --password=%s --job-id=123456789", "some_test_url", "miau", "bup", "123456789"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Errorf("Command response body doesn't match. \n \n %s", diffString)
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
		t.Errorf("Command response body doesn't match. \n \n %s", diffString)
	}
}
