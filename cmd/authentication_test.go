package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetAuthenticate() {
	// reset automation flag vars
	ResetFlags()
	// mock interface
	auth.AuthenticationV3 = auth.NewMockAuthenticationV3
}

func TestAuthenticationUserIdOrNameRequired(t *testing.T) {
	// reset params
	resetAuthenticate()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --project-id=%s --password=%s", "http://some_test_url", "bup", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestAuthenticationProjectIdOrNameRequired(t *testing.T) {
	// reset params
	resetAuthenticate()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --password=%s", "http://some_test_url", "bup", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}
}

func TestAuthenticationPasswordfromStdInputWhenEmpty(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***/public
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
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s", "http://some_test_url", "miau", "bup", "123456789"))

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
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***/public
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --username=%s --password=%s --project-id=%s --project-name=%s --user-domain-name=%s --user-domain-id=%s --project-domain-name=%s --project-domain-id=%s --region=%s", "http://some_test_url", "userid", "username", "passwrod", "projectid", "projectname", "userdomainname", "userdomainid", "projectdomainid", "projectdomainname", "staging"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationResultString(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***/public
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s", "http://some_test_url", "miau", "bup", "123456789"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationResultJSON(t *testing.T) {
	want := `{"ARC_SERVICE_ENDPOINT": "https://arc.staging.***REMOVED***/public","LYRA_SERVICE_ENDPOINT": "https://lyra.staging.***REMOVED***","OS_TOKEN": "test_token_id"}`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s --json", "http://some_test_url", "miau", "bup", "123456789"))

	eq, err := JsonDiff(want, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}

//
// Test region
//

func TestAuthenticationNotGivenARegion(t *testing.T) {
	// should return first entry from endpoints
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***/public
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s", "http://some_test_url", "miau", "bup", "123456789"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationGivenARegion(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra.staging.***REMOVED***
export ARC_SERVICE_ENDPOINT=https://arc.staging.***REMOVED***/public
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s --region=%s", "http://some_test_url", "miau", "bup", "123456789", "staging"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationGivenAWrongRegion(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=
export ARC_SERVICE_ENDPOINT=
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s --region=%s", "http://some_test_url", "miau", "bup", "123456789", "wrong"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}
