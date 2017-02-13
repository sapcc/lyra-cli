package cmd

import (
	"fmt"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetRunShow() {
	// reset automation flag vars
	ResetFlags()
}

func TestRunShowCmdWithWrongEnvEndpointsAndTokenSet(t *testing.T) {
	resetRunShow()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra run show")
	resetRunShow()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra run show")
	resetRunShow()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra run show")
}

func TestRunShowCmdWithAuthenticationFlags(t *testing.T) {
	responseBody := `{
	  "id": "30",
	  "log": "Selecting nodes with selector @identity='0128e993-c709-4ce1-bccf-e06eb10900a0'\nSelected nodes:\n0128e993-c709-4ce1-bccf-e06eb10900a0 mo-85b92ea6f",
	  "created_at": "2016-04-07T21:42:07.416Z",
	  "updated_at": "2016-04-07T21:42:14.294Z",
	  "repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc",
	  "state": "executing",
	  "jobs": [
	    "b843bbe9-fa95-4a0b-9329-aed05d1de8b8"
	  ],
	  "owner": "u-fa35bbc5f",
	  "automation_id": "6",
	  "automation_name": "Chef_test",
	  "selector": "@identity='0128e993-c709-4ce1-bccf-e06eb10900a0'"
	}`
	testServer := TestServer(200, responseBody, map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)

	want := `+---------------------+--------------------------------------------------+
|         KEY         |                      VALUE                       |
+---------------------+--------------------------------------------------+
| automation_id       | 6                                                |
| automation_name     | Chef_test                                        |
| created_at          | 2016-04-07T21:42:07.416Z                         |
| id                  | 30                                               |
| jobs                | [b843bbe9-fa95-4a0b-9329-aed05d1de8b8]           |
| log                 | Selecting nodes with selector                    |
|                     | @identity='0128e993-c709-4ce1-bccf-e06eb10900a0' |
|                     | Selected nodes:                                  |
|                     | 0128e993-c709-4ce1-bccf-e06eb10900a0             |
|                     | mo-85b92ea6f                                     |
| owner               | u-fa35bbc5f                                      |
| repository_revision | 0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc         |
| selector            | @identity='0128e993-c709-4ce1-bccf-e06eb10900a0' |
| state               | executing                                        |
| updated_at          | 2016-04-07T21:42:14.294Z                         |
+---------------------+--------------------------------------------------+`

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra run show --auth-url=%s --user-id=%s --project-id=%s --password=%s --run-id=123456789", "some_test_url", "miau", "bup", "123456789"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestRunShowCmdWithResultTable(t *testing.T) {
	// set test server
	responseBody := `{
  "id": "30",
  "log": "Selecting nodes with selector @identity='0128e993-c709-4ce1-bccf-e06eb10900a0'\nSelected nodes:\n0128e993-c709-4ce1-bccf-e06eb10900a0 mo-85b92ea6f",
  "created_at": "2016-04-07T21:42:07.416Z",
  "updated_at": "2016-04-07T21:42:14.294Z",
  "repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc",
  "state": "executing",
  "jobs": [
    "b843bbe9-fa95-4a0b-9329-aed05d1de8b8"
  ],
  "owner": "u-fa35bbc5f",
  "automation_id": "6",
  "automation_name": "Chef_test",
  "selector": "@identity='0128e993-c709-4ce1-bccf-e06eb10900a0'"
}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+---------------------+--------------------------------------------------+
|         KEY         |                      VALUE                       |
+---------------------+--------------------------------------------------+
| automation_id       | 6                                                |
| automation_name     | Chef_test                                        |
| created_at          | 2016-04-07T21:42:07.416Z                         |
| id                  | 30                                               |
| jobs                | [b843bbe9-fa95-4a0b-9329-aed05d1de8b8]           |
| log                 | Selecting nodes with selector                    |
|                     | @identity='0128e993-c709-4ce1-bccf-e06eb10900a0' |
|                     | Selected nodes:                                  |
|                     | 0128e993-c709-4ce1-bccf-e06eb10900a0             |
|                     | mo-85b92ea6f                                     |
| owner               | u-fa35bbc5f                                      |
| repository_revision | 0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc         |
| selector            | @identity='0128e993-c709-4ce1-bccf-e06eb10900a0' |
| state               | executing                                        |
| updated_at          | 2016-04-07T21:42:14.294Z                         |
+---------------------+--------------------------------------------------+`

	// reset stuff
	resetRunShow()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra run show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --run-id=%s", server.URL, "https://somewhere1.com", "token123", "run_id"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestRunShowCmdWithResultJSON(t *testing.T) {
	// set test server
	responseBody := `{
  "id": "30",
  "log": "Selecting nodes with selector @identity='0128e993-c709-4ce1-bccf-e06eb10900a0'\nSelected nodes:\n0128e993-c709-4ce1-bccf-e06eb10900a0 mo-85b92ea6f",
  "created_at": "2016-04-07T21:42:07.416Z",
  "updated_at": "2016-04-07T21:42:14.294Z",
  "repository_revision": "0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc",
  "state": "executing",
  "jobs": [
    "b843bbe9-fa95-4a0b-9329-aed05d1de8b8"
  ],
  "owner": "u-fa35bbc5f",
  "automation_id": "6",
  "automation_name": "Chef_test",
  "selector": "@identity='0128e993-c709-4ce1-bccf-e06eb10900a0'"
}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetRunShow()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra run show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --run-id=%s --json", server.URL, "https://somewhere2.com", "token123", "run_id"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
		return
	}

	eq, err := JsonDiff(responseBody, resulter.Output)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}
