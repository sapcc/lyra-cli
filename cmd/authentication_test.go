package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/locales"
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

	errorMsg := locales.ErrorMessages("flag-missing")
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
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

	errorMsg := locales.ErrorMessages("flag-missing")
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationPasswordfromStdInputWhenEmpty(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra-app-staging
export ARC_SERVICE_ENDPOINT=https://arc-app-staging/public
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
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s", "http://some_test_url", "miau", "123456789"))

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
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra-app-staging
export ARC_SERVICE_ENDPOINT=https://arc-app-staging/public
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --username=%s --password=%s --project-id=%s --project-name=%s --user-domain-name=%s --user-domain-id=%s --project-domain-name=%s --project-domain-id=%s --region=%s", "http://some_test_url", "userid", "username", "passwrod", "projectid", "projectname", "userdomainname", "userdomainid", "projectdomainid", "projectdomainname", "staging"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationApplicationCredential(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra-app-staging
export ARC_SERVICE_ENDPOINT=https://arc-app-staging/public
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --application-credential-name=%s --application-credential-secret=%s --user-id=%s", "http://some_test_url", "miau", "123456789", "1234567890"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationApplicationCredentialAllRequiredDomainName(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra-app-staging
export ARC_SERVICE_ENDPOINT=https://arc-app-staging/public
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --application-credential-name=%s --application-credential-secret=%s --username=%s --user-domain-name=%s", "http://some_test_url", "miau", "123456789", "1234567890", "default"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationApplicationCredentialAllRequiredDomainId(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra-app-staging
export ARC_SERVICE_ENDPOINT=https://arc-app-staging/public
export OS_TOKEN=test_token_id`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --application-credential-name=%s --application-credential-secret=%s --username=%s --user-domain-id=%s", "http://some_test_url", "miau", "123456789", "1234567890", "0987654321"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationApplicationCredentialDomainIdOrNameRequired(t *testing.T) {
	// reset params
	resetAuthenticate()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --application-credential-name=%s --application-credential-secret=%s --username=%s", "http://some_test_url", "miau", "123456789", "user"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}

	errorMsg := locales.ErrorMessages("flag-missing")
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationApplicationCredentialUserIdOrNameRequired(t *testing.T) {
	// reset params
	resetAuthenticate()
	// run cmd
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --application-credential-name=%s --application-credential-secret=%s", "http://some_test_url", "miau", "123456789"))
	if resulter.Error == nil {
		t.Error(`Command expected to get an error`)
	}

	errorMsg := locales.ErrorMessages("flag-missing")
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
	}
}

func TestAuthenticationResultString(t *testing.T) {
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra-app-staging
export ARC_SERVICE_ENDPOINT=https://arc-app-staging/public
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
	want := `{"ARC_SERVICE_ENDPOINT": "https://arc-app-staging/public","LYRA_SERVICE_ENDPOINT": "https://lyra-app-staging","OS_TOKEN": "test_token_id"}`

	// reset params
	resetAuthenticate()
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra authenticate --auth-url=%s --user-id=%s --project-id=%s --password=%s --json", "http://some_test_url", "miau", "bup", "123456789"))

	eq, err := JsonDiff(want, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !eq {
		t.Error("Json response body and print out Json do not match.")
	}
}

//
// Test region
//

func TestAuthenticationNotGivenARegion(t *testing.T) {
	// should return first entry from endpoints
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra-app-staging
export ARC_SERVICE_ENDPOINT=https://arc-app-staging/public
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
	want := `export LYRA_SERVICE_ENDPOINT=https://lyra-app-staging
export ARC_SERVICE_ENDPOINT=https://arc-app-staging/public
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
