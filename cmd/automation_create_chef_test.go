package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/helpers"
)

func resetAutomationCreateChefFlagVars() {
	// reset flag
	ResetFlags()
}

func TestAutomationCreateChefCmdWithAuthenticationFlags(t *testing.T) {
	responseBody := `{"id": 40,"type": "Chef","name": "test","project_id": "p-9597d2775","repository": "https://github.com/userId0123456789/automation-test.git","repository_revision": "master","timeout": 3600,"tags": null,"created_at": "2016-05-19T12:48:51.629Z","updated_at": "2016-05-19T12:48:51.629Z","run_list": ["recipe[App:Nginx]"],"chef_attributes": null,"log_level": null,"chef_version": "1.2.3","path": null,"arguments": null,"environment": null}`
	testServer := TestServer(200, responseBody, map[string]string{})
	defer testServer.Close()
	// mock interface for authenticationt test to return mocked endopoints and tokens and test method can use user authentication params to run
	auth.AuthenticationV3 = newMockAuthenticationV3(testServer)
	want := `+---------------------+---------------------------------------------------------+
|         KEY         |                          VALUE                          |
+---------------------+---------------------------------------------------------+
| arguments           | <nil>                                                   |
| chef_attributes     | <nil>                                                   |
| chef_version        | 1.2.3                                                   |
| created_at          | 2016-05-19T12:48:51.629Z                                |
| environment         | <nil>                                                   |
| id                  | 40                                                      |
| log_level           | <nil>                                                   |
| name                | test                                                    |
| path                | <nil>                                                   |
| project_id          | p-9597d2775                                             |
| repository          | https://github.com/userId0123456789/automation-test.git |
| repository_revision | master                                                  |
| run_list            | [recipe[App:Nginx]]                                     |
| tags                | <nil>                                                   |
| timeout             | 3600                                                    |
| type                | Chef                                                    |
| updated_at          | 2016-05-19T12:48:51.629Z                                |
+---------------------+---------------------------------------------------------+`

	// reset stuff
	resetAutomationCreateChefFlagVars()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra automation create chef --auth-url=%s --user-id=%s --project-id=%s --password=%s --name=%s --repository=%s --runlist=%s", "some_test_url", "miau", "bup", "123456789", "chef_test", "http://some_repository", "recipe[nginx]"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
	}
}

func TestAutomationCreateChefShouldSetMinimumAttributes(t *testing.T) {
	// set test server
	responseBody := `{"id": 40,"type": "Chef","name": "test","project_id": "p-9597d2775","repository": "https://github.com/userId0123456789/automation-test.git","repository_revision": "master","timeout": 3600,"tags": null,"created_at": "2016-05-19T12:48:51.629Z","updated_at": "2016-05-19T12:48:51.629Z","run_list": ["recipe[nginx]"],"chef_attributes": null,"log_level": null,"debug":false,"chef_version": null,"path": null,"arguments": null,"environment": null}`
	want := `+---------------------+---------------------------------------------------------+
|         KEY         |                          VALUE                          |
+---------------------+---------------------------------------------------------+
| arguments           | <nil>                                                   |
| chef_attributes     | <nil>                                                   |
| chef_version        | <nil>                                                   |
| created_at          | 2016-05-19T12:48:51.629Z                                |
| debug               | false                                                   |
| environment         | <nil>                                                   |
| id                  | 40                                                      |
| log_level           | <nil>                                                   |
| name                | test                                                    |
| path                | <nil>                                                   |
| project_id          | p-9597d2775                                             |
| repository          | https://github.com/userId0123456789/automation-test.git |
| repository_revision | master                                                  |
| run_list            | [recipe[nginx]]                                         |
| tags                | <nil>                                                   |
| timeout             | 3600                                                    |
| type                | Chef                                                    |
| updated_at          | 2016-05-19T12:48:51.629Z                                |
+---------------------+---------------------------------------------------------+`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	resetAutomationCreateChefFlagVars()
	resulter := FullCmdTester(RootCmd,
		fmt.Sprintf("lyra automation create chef --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --name=%s --repository=%s --runlist=%s",
			server.URL,
			server.URL,
			"token123",
			"chef_test",
			"http://some_repository",
			"recipe[nginx]"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
		return
	}
}

func TestAutomationCreateChefShouldSetAttributes(t *testing.T) {
	// set test server
	responseBody := `{
	"id":932,
	"type":"Chef",
	"name":"test",
	"project_id":"abcdefghijklmnopqrstuwxyz1234567",
	"repository":"https://github.com/userId0123456789/automation-test.git",
	"repository_authentication_enabled": true,
	"repository_revision":"master",
	"timeout":4000,
	"tags":null,
	"created_at":"2017-06-19T11:54:30.329Z",
	"updated_at":"2017-06-19T11:54:30.329Z",
	"run_list":["recipe[nginx]","recipe[test]"],
	"chef_attributes":{"test":"test"},
	"log_level":null,
	"debug":true,
	"chef_version":"1.2.3",
	"path":null,
	"arguments":[],
	"environment":{}
}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()
	want := `+-----------------------------------+---------------------------------------------------------+
|                KEY                |                          VALUE                          |
+-----------------------------------+---------------------------------------------------------+
| arguments                         | []                                                      |
| chef_attributes                   | map[test:test]                                          |
| chef_version                      | 1.2.3                                                   |
| created_at                        | 2017-06-19T11:54:30.329Z                                |
| debug                             | true                                                    |
| environment                       | map[]                                                   |
| id                                | 932                                                     |
| log_level                         | <nil>                                                   |
| name                              | test                                                    |
| path                              | <nil>                                                   |
| project_id                        | abcdefghijklmnopqrstuwxyz1234567                        |
| repository                        | https://github.com/userId0123456789/automation-test.git |
| repository_authentication_enabled | true                                                    |
| repository_revision               | master                                                  |
| run_list                          | [recipe[nginx]                                          |
|                                   | recipe[test]]                                           |
| tags                              | <nil>                                                   |
| timeout                           | 4000                                                    |
| type                              | Chef                                                    |
| updated_at                        | 2017-06-19T11:54:30.329Z                                |
+-----------------------------------+---------------------------------------------------------+`

	resetAutomationCreateChefFlagVars()
	resulter := FullCmdTester(RootCmd,
		fmt.Sprintf("lyra automation create chef --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --name=%s --repository=%s --repository-credentials=%s --repository-revision=%s --timeout=%d --runlist=%s --attributes=%s --chef-version=%s --chef-debug=%t",
			server.URL,
			server.URL,
			"token123",
			"chef_test",
			"http://some_repository",
			"ex_github_personal_access_toke",
			"master",
			3600,
			"recipe[nginx],recipe[test]",
			`{"test":"test"}`,
			`1.2.3`,
			true))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}

	if !strings.Contains(resulter.Output, want) {
		diffString := StringDiff(resulter.Output, want)
		t.Error(fmt.Sprintf("Command response body doesn't match. \n \n %s", diffString))
		return
	}
	if !strings.Contains(chef.AutomationType, "Chef") {
		t.Error(`Command create chef expected to have Chef type'`)
	}
	if !strings.Contains(chef.Name, "chef_test") {
		t.Error(`Command create chef expected to have same name'`)
	}
	if !strings.Contains(chef.Repository, "http://some_repository") {
		t.Error(`Command create chef expected to have same repository'`)
	}
	if !strings.Contains(*chef.RepositoryCredentials, "ex_github_personal_access_toke") {
		t.Error(`Command create chef expected to have credentials set'`)
	}
	if !strings.Contains(chef.RepositoryRevision, "master") {
		t.Error(`Command create chef expected to have same repository revision'`)
	}
	if chef.Timeout != 3600 {
		t.Error(`Command create chef expected to have same timeout'`)
	}
	if len(chef.Runlist) != 2 {
		t.Error(`Command create chef expected to have runlist.'`)
	}
	testString, _ := helpers.StructureToJSON(chef.Attributes)
	if !strings.Contains(testString, `{"test":"test"}`) {
		t.Error(`Command create chef expected to have same attributes`)
	}
	if !chef.Debug {
		t.Error(`Command create chef expected to have same debug flag`)
	}
}

func TestAutomationCreateChefShouldSetAttributesFromFile(t *testing.T) {
	// set test server
	responseBody := `{"id":"1","name":"Chef_test1","repository":"https://github.com/userId0123456789/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// path to the current directory
	pwd, _ := os.Getwd()
	file := fmt.Sprint(pwd, "/../examples/example1.JSON")
	txt, _ := ioutil.ReadFile(file)

	resetAutomationCreateChefFlagVars()
	resulter := FullCmdTester(RootCmd,
		fmt.Sprintf("lyra automation create chef --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --attributes-from-file=%s",
			server.URL,
			server.URL,
			"token123",
			file))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if len(txt) == 0 {
		t.Error(`Command create chef expected to find an attributes file with content'`)
	}

	// convert interface to string and compare
	testString, _ := helpers.StructureToJSON(chef.Attributes)
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, txt); err != nil {
		t.Errorf("Failed to compact json: %s", err)
	}
	if !strings.Contains(testString, buffer.String()) {
		t.Error(`Command create chef expected to have same attributes'`)
	}
}

func TestAutomationCreateChefShouldSetAttributesFromStdInput(t *testing.T) {
	// set test server
	responseBody := `{"id":"1","name":"Chef_test1","repository":"https://github.com/userId0123456789/automation-test.git","repository_revision":"master","run_list":"[recipe[nginx]]","chef_attributes":{"test":"test"},"log_level":"info","arguments":"{}"}`
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// path to the current directory
	pwd, _ := os.Getwd()
	file := fmt.Sprint(pwd, "/../examples/example1.JSON")
	// read example
	txt, _ := ioutil.ReadFile(file)

	// keep backup of the real stdout
	oldStdout := os.Stdout

	// write passowrd
	_, err := pipeToStdin(string(txt))
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

	resetAutomationCreateChefFlagVars()
	resulter := FullCmdTester(RootCmd,
		fmt.Sprintf("lyra automation create chef --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --attributes-from-file=%s",
			server.URL,
			"http://some_nice_url",
			"token123",
			"-"))

	// flush, restore close
	os.Stdout = oldStdout
	flushStdin()
	w.Close()

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	// convert interface to string and compare
	testString, err := helpers.StructureToJSON(chef.Attributes)
	if err != nil {
		t.Errorf("Failed to convert to JSON: %s", err)
	}
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, txt); err != nil {
		t.Errorf("Failed to compact JSON: %s", err)
	}
	if !strings.Contains(testString, buffer.String()) {
		t.Error(`Command create chef expected to have same attributes'`)
	}
}
