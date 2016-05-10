package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/sapcc/lyra-cli/helpers"
)

func resetAutomationCreateChefFlagVars() {
	chef = Chef{}
	tags = ""
	runlist = ""
	attributes = ""
	attributesFromFile = ""

	// reset automation flag vars
	resetRootFlagVars()
	// reset commands
	RootCmd.ResetCommands()
	AutomationCmd.ResetCommands()
	AutomationCreateCmd.ResetCommands()
	// build commands
	RootCmd.AddCommand(AutomationCmd)
	AutomationCmd.AddCommand(AutomationCreateCmd)
	AutomationCreateCmd.AddCommand(AutomationCreateChefCmd)
}

func TestAutomationCreateChefShouldSetAttributes(t *testing.T) {
	// set test server
	responseBody := "Miau"
	server := TestServer(200, responseBody)
	defer server.Close()

	resetAutomationCreateChefFlagVars()
	resulter := FullCmdTester(RootCmd,
		fmt.Sprintf("lyra-cli automation create chef --lyra-service-endpoint=%s --token=%s --name=%s --repository=%s --repository-revision=%s --timeout=%d --tags=%s --runlist=%s --attributes=%s --log-level=%s",
			server.URL,
			"token123",
			"chef_test",
			"http://some_repository",
			"master",
			3600,
			"name:test,tag1=test",
			"recipe[nginx],recipe[test]",
			`{"test":"test"}`,
			"info"))

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
	if !strings.Contains(resulter.Output, responseBody) {
		t.Error(`Command response body doesn't match.'`)
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
	if !strings.Contains(chef.RepositoryRevision, "master") {
		t.Error(`Command create chef expected to have same repository revision'`)
	}
	if chef.Timeout != 3600 {
		t.Error(`Command create chef expected to have same time out'`)
	}
	if len(chef.Tags) != 2 {
		t.Error(`Command create chef expected to have tags.'`)
	}
	if len(chef.Runlist) != 2 {
		t.Error(`Command create chef expected to have runlist.'`)
	}
	testString, _ := helpers.StructureToJSON(chef.Attributes)
	if !strings.Contains(testString, `{"test":"test"}`) {
		t.Error(`Command create chef expected to have same attributes'`)
	}
	if !strings.Contains(chef.LogLevel, "info") {
		t.Error(`Command create chef expected to have same log level'`)
	}
}

func TestAutomationCreateChefShouldSetAttributesFromFile(t *testing.T) {
	// set test server
	responseBody := "Miau"
	server := TestServer(200, responseBody)
	defer server.Close()

	// path to the current directory
	pwd, _ := os.Getwd()
	file := fmt.Sprint(pwd, "/../examples/example1.JSON")
	txt, _ := ioutil.ReadFile(file)

	resetAutomationCreateChefFlagVars()
	resulter := FullCmdTester(RootCmd,
		fmt.Sprintf("lyra-cli automation create chef --lyra-service-endpoint=%s --token=%s --attributes-from-file=%s",
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
	json.Compact(buffer, txt)
	if !strings.Contains(testString, buffer.String()) {
		t.Error(`Command create chef expected to have same attributes'`)
	}
}

// func TestAutomationCreateChefShouldSetAttributesFromStdInput(t *testing.T) {
//   // set test server
//   responseBody := "Miau"
//   server := TestServer(200, responseBody)
//   defer server.Close()
//
//   // path to the current directory
//   pwd, _ := os.Getwd()
//   file := fmt.Sprint(pwd, "/../examples/example1.JSON")
//   // read example
//   txt, _ := ioutil.ReadFile(file)
//
//   // NOT WORKING
//   // pipe std in and out
//   oldStdout := os.Stdout
//   oldStdin := os.Stdin
//   r, w, _ := os.Pipe() //type *os.File
//   os.Stdout = w
//   os.Stdin = r
//   // writing in stdout
//   fmt.Fprintf(w, string(txt))
//
//   resetAutomationCreateChefFlagVars()
//   resulter := FullCmdTester(RootCmd,
//     fmt.Sprintf("lyra-cli automation create chef --lyra-service-endpoint=%s --token=%s --attributes-from-file=%s",
//       server.URL,
//       "token123",
//       "-"))
//
//   fmt.Println("#####")
//   fmt.Println(chef.Attributes)
//   fmt.Println("#####")
//
//   // back to normal state
//   w.Close()
//   //r.Close()
//   // restoring the real stdout
//   os.Stdout = oldStdout
//   os.Stdin = oldStdin
//
//   if resulter.Error != nil {
//     t.Error(`Command expected to not get an error`)
//   }
//   if len(txt) == 0 {
//     t.Error(`Command create chef expected to find an attributes file with content'`)
//   }
//   if !strings.Contains(chef.Attributes, string(txt)) {
//     t.Error(`Command create chef expected to have same attributes'`)
//   }
// }
