package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func resetJobShow() {
	// reset flags
	ResetFlags()
}

func newMockAuthenticationV3JobShow(authOpts LyraAuthOps) Authentication {
	// set test server
	responseBody := `{
  "version": 1,
  "sender": "api-7461f075665433b2bb80d4c9031bbff8-4c4ab",
  "request_id": "e24e86fa-4bbd-47f3-a4d2-1566618ef765",
  "to": "0128e993-c709-4ce1-bccf-e06eb10900a0",
  "timeout": 3600,
  "agent": "chef",
  "action": "zero",
  "payload": "{\"run_list\":[\"recipe[nginx]\"],\"recipe_url\":\"https://objectstore.***REMOVED***:443/v1/AUTH_abcdefghyjklmnopqrstuwxyz1234567/monsoon-automation/0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc-chef.tgz?temp_url_sig=bd8ad675e854210689613d735bbbd43b7bc334ee\\u0026temp_url_expires=1462899028\",\"attributes\":null,\"debug\":false}",
  "status": "failed",
  "created_at": "2016-05-10T15:50:28.286532Z",
  "updated_at": "2016-05-10T15:50:33.402484Z",
  "project": "p-9597d2775",
  "user_id": "u-519166a05"
}`
	server := TestServer(200, responseBody, map[string]string{})

	return &MockV3{AuthOpts: authOpts, TestServer: server}
}

func TestJobShowCmdWithAuthenticationFlags(t *testing.T) {
	// mock interface for authenticationt test
	AuthenticationV3 = newMockAuthenticationV3JobShow
	want := `+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+
|    KEY     |                                                                          VALUE                                                                          |
+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+
| action     | zero                                                                                                                                                    |
| agent      | chef                                                                                                                                                    |
| created_at | 2016-05-10T15:50:28.286532Z                                                                                                                             |
| payload    | {"run_list":["recipe[nginx]"],"recipe_url":"https://objectstore.***REMOVED***:443/v1/AUTH_abcdefghyjklmnopqrstuwxyz1234567/monsoon-automation       |
|            | /0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc-chef.tgz?temp_url_sig=bd8ad675e854210689613d735bbbd43b7bc334ee\u0026temp_url_expires=1462899028","attributes" |
|            | :null,"debug":false}                                                                                                                                    |
| project    | p-9597d2775                                                                                                                                             |
| request_id | e24e86fa-4bbd-47f3-a4d2-1566618ef765                                                                                                                    |
| sender     | api-7461f075665433b2bb80d4c9031bbff8-4c4ab                                                                                                              |
| status     | failed                                                                                                                                                  |
| timeout    | 3600                                                                                                                                                    |
| to         | 0128e993-c709-4ce1-bccf-e06eb10900a0                                                                                                                    |
| updated_at | 2016-05-10T15:50:33.402484Z                                                                                                                             |
| user_id    | u-519166a05                                                                                                                                             |
| version    | 1                                                                                                                                                       |
+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+`

	// reset stuff
	resetAutomationList()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job show --auth-url=%s --user-id=%s --project-id=%s --password=%s --job-id=123456789", "some_test_url", "miau", "bup", "123456789"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestJobShowCmdWithWrongEnvEndpointAndTokenSet(t *testing.T) {
	resetJobShow()
	CheckhErrorWhenNoEnvEndpointAndTokenSet(t, RootCmd, "lyra job show")
	resetJobShow()
	CheckhErrorWhenNoEnvEndpointSet(t, RootCmd, "lyra job show")
	resetJobShow()
	CheckhErrorWhenNoEnvTokenSet(t, RootCmd, "lyra job show")
}

func TestJobShowCmdResultTable(t *testing.T) {
	// set test server
	responseBody := `{
  "version": 1,
  "sender": "api-7461f075665433b2bb80d4c9031bbff8-4c4ab",
  "request_id": "e24e86fa-4bbd-47f3-a4d2-1566618ef765",
  "to": "0128e993-c709-4ce1-bccf-e06eb10900a0",
  "timeout": 3600,
  "agent": "chef",
  "action": "zero",
  "payload": "{\"run_list\":[\"recipe[nginx]\"],\"recipe_url\":\"https://objectstore.***REMOVED***:443/v1/AUTH_abcdefghyjklmnopqrstuwxyz1234567/monsoon-automation/0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc-chef.tgz?temp_url_sig=bd8ad675e854210689613d735bbbd43b7bc334ee\\u0026temp_url_expires=1462899028\",\"attributes\":null,\"debug\":false}",
  "status": "failed",
  "created_at": "2016-05-10T15:50:28.286532Z",
  "updated_at": "2016-05-10T15:50:33.402484Z",
  "project": "p-9597d2775",
  "user_id": "u-519166a05"
}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+
|    KEY     |                                                                          VALUE                                                                          |
+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+
| action     | zero                                                                                                                                                    |
| agent      | chef                                                                                                                                                    |
| created_at | 2016-05-10T15:50:28.286532Z                                                                                                                             |
| payload    | {"run_list":["recipe[nginx]"],"recipe_url":"https://objectstore.***REMOVED***:443/v1/AUTH_abcdefghyjklmnopqrstuwxyz1234567/monsoon-automation       |
|            | /0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc-chef.tgz?temp_url_sig=bd8ad675e854210689613d735bbbd43b7bc334ee\u0026temp_url_expires=1462899028","attributes" |
|            | :null,"debug":false}                                                                                                                                    |
| project    | p-9597d2775                                                                                                                                             |
| request_id | e24e86fa-4bbd-47f3-a4d2-1566618ef765                                                                                                                    |
| sender     | api-7461f075665433b2bb80d4c9031bbff8-4c4ab                                                                                                              |
| status     | failed                                                                                                                                                  |
| timeout    | 3600                                                                                                                                                    |
| to         | 0128e993-c709-4ce1-bccf-e06eb10900a0                                                                                                                    |
| updated_at | 2016-05-10T15:50:33.402484Z                                                                                                                             |
| user_id    | u-519166a05                                                                                                                                             |
| version    | 1                                                                                                                                                       |
+------------+---------------------------------------------------------------------------------------------------------------------------------------------------------+`

	// reset stuff
	resetJobShow()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --job-id=123456789", "http://somewhere.com", server.URL, "token123"))

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestJobShowCmdResultJSON(t *testing.T) {
	// set test server
	responseBody := `{
  "version": 1,
  "sender": "api-7461f075665433b2bb80d4c9031bbff8-4c4ab",
  "request_id": "e24e86fa-4bbd-47f3-a4d2-1566618ef765",
  "to": "0128e993-c709-4ce1-bccf-e06eb10900a0",
  "timeout": 3600,
  "agent": "chef",
  "action": "zero",
  "payload": "{\"run_list\":[\"recipe[nginx]\"],\"recipe_url\":\"https://objectstore.***REMOVED***:443/v1/AUTH_abcdefghyjklmnopqrstuwxyz1234567/monsoon-automation/0c2ae56428273ed2f542104b2d67ab4b4d9ed6bc-chef.tgz?temp_url_sig=bd8ad675e854210689613d735bbbd43b7bc334ee\\u0026temp_url_expires=1462899028\",\"attributes\":null,\"debug\":false}",
  "status": "failed",
  "created_at": "2016-05-10T15:50:28.286532Z",
  "updated_at": "2016-05-10T15:50:33.402484Z",
  "project": "p-9597d2775",
  "user_id": "u-519166a05"
}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// reset stuff
	resetJobShow()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra job show --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --job-id=123456789 --json", "http://somewhere.com", server.URL, "token123"))

	source := map[string]interface{}{}
	err := json.Unmarshal([]byte(responseBody), &source)
	if err != nil {
		t.Error(err.Error())
		return
	}
	response := map[string]interface{}{}
	err = json.Unmarshal([]byte(resulter.Output), &response)
	if err != nil {
		t.Error(err.Error())
		return
	}

	eq := reflect.DeepEqual(source, response)
	if eq == false {
		t.Error("Json response body and print out Json do not match.")
	}
}
