package cmd

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
)

func resetAutomationCreateScriptFlagVars() {
	// reset flag
	ResetFlags()
}

func TestAutomationCreateScriptCmdWithAuthenticationFlags(t *testing.T) {
	responseBody := `{"arguments": null,
	"chef_attributes": null,
	"chef_version": null,
	"created_at": "2016-06-01T08:34:11.761Z",
	"environment": null,
	"id": 45,
	"log_level": null,
	"name": "test_script_cli",
	"path": "path_to_the_file",
	"project_id": "p-9597d2775",
	"repository": "https://github.com/user123/automation-test.git",
	"repository_revision": "master",
	"run_list": null,
	"tags": {
		"name": "arturo"
	},
	"timeout": 3600,
	"type": "Script",
	"updated_at": "2016-06-01T08:34:11.761Z"
}`
	testServer := TestServer(200, responseBody, map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)

	want := `+---------------------+---------------------------------------------------------+
|         KEY         |                          VALUE                          |
+---------------------+---------------------------------------------------------+
| arguments           | <nil>                                                   |
| chef_attributes     | <nil>                                                   |
| chef_version        | <nil>                                                   |
| created_at          | 2016-06-01T08:34:11.761Z                                |
| environment         | <nil>                                                   |
| id                  | 45                                                      |
| log_level           | <nil>                                                   |
| name                | test_script_cli                                         |
| path                | path_to_the_file                                        |
| project_id          | p-9597d2775                                             |
| repository          | https://github.com/user123/automation-test.git |
| repository_revision | master                                                  |
| run_list            | <nil>                                                   |
| tags                | map[name:arturo]                                        |
| timeout             | 3600                                                    |
| type                | Script                                                  |
| updated_at          | 2016-06-01T08:34:11.761Z                                |
+---------------------+---------------------------------------------------------+`

	// reset stuff
	resetAutomationCreateScriptFlagVars()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation create script --name=test_script_cli --repository=http://some_repository --path=path_to_the_file --auth-url=http://some_auth_url --user-id=u-519166a05  --project-id=p-9597d2775 --password=123456789"))

	if resulter.Error != nil {
		t.Errorf(`Command expected to not get an error: %s`, resulter.Error)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationCreateScriptShouldSetMinimumAttributes(t *testing.T) {
	// set test server
	responseBody := `{
  "arguments": null,
  "chef_attributes": null,
  "chef_version": null,
  "created_at": "2016-06-01T08:34:11.761Z",
  "environment": null,
  "id": 45,
  "log_level": null,
  "name": "test",
  "path": "path_to_the_file",
  "project_id": "p-9597d2775",
  "repository": "https://github.com/user123/automation-test.git",
  "repository_revision": "master",
  "run_list": null,
  "tags": null,
  "timeout": 3600,
  "type": "Script",
  "updated_at": "2016-06-01T08:34:11.761Z"
}`
	want := `+---------------------+---------------------------------------------------------+
|         KEY         |                          VALUE                          |
+---------------------+---------------------------------------------------------+
| arguments           | <nil>                                                   |
| chef_attributes     | <nil>                                                   |
| chef_version        | <nil>                                                   |
| created_at          | 2016-06-01T08:34:11.761Z                                |
| environment         | <nil>                                                   |
| id                  | 45                                                      |
| log_level           | <nil>                                                   |
| name                | test                                                    |
| path                | path_to_the_file                                        |
| project_id          | p-9597d2775                                             |
| repository          | https://github.com/user123/automation-test.git |
| repository_revision | master                                                  |
| run_list            | <nil>                                                   |
| tags                | <nil>                                                   |
| timeout             | 3600                                                    |
| type                | Script                                                  |
| updated_at          | 2016-06-01T08:34:11.761Z                                |
+---------------------+---------------------------------------------------------+`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	resetAutomationCreateScriptFlagVars()
	resulter := FullCmdTester(RootCmd,
		fmt.Sprintf("lyra automation create script --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --name=%s --repository=%s --path=%s",
			server.URL,
			server.URL,
			"token123",
			"test",
			"http://some_repository",
			"some_nice_path"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
		return
	}
}

func TestAutomationCreateScriptShouldSetAttributes(t *testing.T) {
	// set test server
	responseBody := `{
  "arguments": null,
  "chef_attributes": null,
  "chef_version": null,
  "created_at": "2016-06-01T08:34:11.761Z",
  "environment": null,
  "id": 45,
  "log_level": null,
  "name": "script_test",
  "path": "path_to_the_file",
  "project_id": "p-9597d2775",
  "repository": "https://github.com/user123/automation-test.git",
  "repository_revision": "master",
  "run_list": null,
  "tags": null,
  "timeout": 3600,
  "type": "Script",
  "updated_at": "2016-06-01T08:34:11.761Z"
}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+---------------------+---------------------------------------------------------+
|         KEY         |                          VALUE                          |
+---------------------+---------------------------------------------------------+
| arguments           | <nil>                                                   |
| chef_attributes     | <nil>                                                   |
| chef_version        | <nil>                                                   |
| created_at          | 2016-06-01T08:34:11.761Z                                |
| environment         | <nil>                                                   |
| id                  | 45                                                      |
| log_level           | <nil>                                                   |
| name                | script_test                                             |
| path                | path_to_the_file                                        |
| project_id          | p-9597d2775                                             |
| repository          | https://github.com/user123/automation-test.git |
| repository_revision | master                                                  |
| run_list            | <nil>                                                   |
| tags                | <nil>                                                   |
| timeout             | 3600                                                    |
| type                | Script                                                  |
| updated_at          | 2016-06-01T08:34:11.761Z                                |
+---------------------+---------------------------------------------------------+`

	resetAutomationCreateChefFlagVars()
	resulter := FullCmdTester(RootCmd,
		fmt.Sprintf("lyra automation create script --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --name=%s --repository=%s --repository-revision=%s --timeout=%d --path=%s --arg=%s --arg=%s --env=%s --env=%s",
			server.URL,
			server.URL,
			"token123",
			"script_test",
			"http://some_repository",
			"master",
			3600,
			"some_nice_path",
			`arg1`,
			`arg2,with,commas`,
			`PROXY:test1`,
			`NO_PROXY:test2,test4`))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error: %s`, resulter.Error)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
		return
	}
	if !strings.Contains(script.AutomationType, "Script") {
		t.Error(`Command create script expected to have Script type'`)
	}
	if !strings.Contains(script.Name, "script_test") {
		t.Error(`Command create script expected to have same name'`)
	}
	if !strings.Contains(script.Repository, "http://some_repository") {
		t.Error(`Command create script expected to have same repository'`)
	}
	if !strings.Contains(script.RepositoryRevision, "master") {
		t.Error(`Command create script expected to have same repository revision'`)
	}
	if script.Timeout != 3600 {
		t.Error(`Command create script expected to have same timeout'`)
	}
	if !strings.Contains(script.Path, "some_nice_path") {
		t.Error(`Command create script expected to have same path'`)
	}
	expectedArgs := []string{"arg1", "arg2,with,commas"}
	if !reflect.DeepEqual(expectedArgs, script.Arguments) {
		t.Errorf(`Expected arguments: %#v, Got %#v`, expectedArgs, script.Arguments)
	}
	if !strings.Contains(script.Environment["PROXY"], "test1") {
		t.Error(`Command create script expected to have same environment'`)
	}
	if !strings.Contains(script.Environment["NO_PROXY"], "test2,test4") {
		t.Error(`Command create script expected to have same environment'`)
	}
}
