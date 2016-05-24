package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rackspace/gophercloud"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func resetAuthenticate() {
	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	authenticateCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(authenticateCmd)
	// mock interface
	AuthenticationV3 = newMockAuthenticationV3
}

func TestAuthenticationUserIdOrNameRequired(t *testing.T) {
	// reset params
	resetAuthenticate()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --identity-endpoint=%s --user-id=%s --project-id=%s --password=%s", "http://some_test_url", "", "bup", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
	// reset params
	resetAuthenticate()
	// run cmd
	resulter = FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --identity-endpoint=%s --username=%s --project-id=%s --password=%s", "http://some_test_url", "", "bup", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestAuthenticationPasswordfromStdInputWhenEmpty(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=test_public_endpoint
export ARC_SERVICE_ENDPOINT=test_public_endpoint
export OS_TOKEN=test_token_id`

	// keep backup of the real stdout
	oldStdout := os.Stdout

	// write passowrd
	_, err := pipeToStdin("password!!!\n")
	if err != nil {
		t.Error(err.Error())
		return
	}

	// pipe std out
	_, w, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
	}
	os.Stdout = w

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --identity-endpoint=%s --user-id=%s --project-id=%s", "http://some_test_url", "miau", "bup", "123456789"))

	// flush, restore close
	os.Stdout = oldStdout
	flushStdin()
	w.Close()

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationWithAllFlags(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=test_public_endpoint
export ARC_SERVICE_ENDPOINT=test_public_endpoint
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --identity-endpoint=%s --user-id=%s --username=%s --password=%s --project-id=%s --project-name=%s --user-domain-name=%s --user-domain-id=%s --project-domain-name=%s --project-domain-id=%s", "http://some_test_url", "userid", "username", "passwrod", "projectid", "projectname", "userdomainname", "userdomainid", "projectdomainid", "projectdomainname"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationResultString(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=test_public_endpoint
export ARC_SERVICE_ENDPOINT=test_public_endpoint
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --identity-endpoint=%s --user-id=%s --project-id=%s --password=%s", "http://some_test_url", "miau", "bup", "123456789"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationResultJSON(t *testing.T) {
	want := `{"ARC_SERVICE_ENDPOINT": "test_public_endpoint","LYRA_SERVICE_ENDPOINT": "test_public_endpoint","OS_TOKEN": "test_token_id"}`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --identity-endpoint=%s --user-id=%s --project-id=%s --password=%s --json", "http://some_test_url", "miau", "bup", "123456789"))

	wantSource := map[string]string{}
	err := json.Unmarshal([]byte(want), &wantSource)
	if err != nil {
		t.Error(err.Error())
		return
	}

	response := map[string]string{}
	err = json.Unmarshal([]byte(resulter.Output), &response)
	if err != nil {
		t.Error(err.Error())
		return
	}

	eq := reflect.DeepEqual(wantSource, response)
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}

}

//
// Helpers
//

func pipeToStdin(s string) (int, error) {
	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		fmt.Println("Error getting os pipes:", err)
		os.Exit(1)
	}
	os.Stdin = pipeReader
	w, err := pipeWriter.WriteString(s)
	pipeWriter.Close()
	return w, err
}

// flushStdin reads from stdin for .5 seconds to ensure no bytes are left on
// the buffer.  Returns the number of bytes read.
func flushStdin() int {
	ch := make(chan byte)
	go func(ch chan byte) {
		reader := bufio.NewReader(os.Stdin)
		for {
			b, err := reader.ReadByte()
			if err != nil { // Maybe log non io.EOF errors, if you want
				close(ch)
				return
			}
			ch <- b
		}
		close(ch)
	}(ch)

	numBytes := 0
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return numBytes
			}
			numBytes++
		case <-time.After(500 * time.Millisecond):
			return numBytes
		}
	}
	return numBytes
}

//
// Mock authentication interface
//

type MockV3 struct {
	AuthOpts LyraAuthOps
	client   *gophercloud.ServiceClient
}

func newMockAuthenticationV3(authOpts LyraAuthOps) Authentication {
	return &MockV3{AuthOpts: authOpts}
}

func (a *MockV3) GetToken() (string, error) {
	return "test_token_id", nil
}

func (a *MockV3) GetServicePublicEndpoint(serviceId string) (string, error) {
	return "test_public_endpoint", nil
}

func (a *MockV3) GetServiceId(serviceType string) (string, error) {
	return "Test_service_id", nil
}
