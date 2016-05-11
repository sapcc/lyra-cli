package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type resulter struct {
	Error   error
	Output  string
	Command *cobra.Command
}

var cmdTestRootNoRun = &cobra.Command{
	Use:   "lyra-test",
	Short: "The root can run its own function",
	Long:  "The root description for help",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
}

func FullCmdTester(testCommand *cobra.Command, input string) resulter {
	buf := new(bytes.Buffer)
	c := cmdTestRootNoRun
	c.SetOutput(buf)
	c.AddCommand(testCommand)
	c.SetArgs(strings.Split(input, " "))
	err := c.Execute()
	output := buf.String()
	return resulter{err, output, c}
}

func TestServer(code int, body string, headers map[string]string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(code) // keep the code after setting headers. If not they will disapear...
		fmt.Fprintln(w, body)
	}))
	return server
}

func CheckhErrorWhenNoEnvEndpointAndTokenSet(t *testing.T, cmd *cobra.Command, input string) {
	// clean env variablen
	os.Unsetenv(ENV_VAR_TOKEN_NAME)
	os.Unsetenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)

	// check
	resulter := FullCmdTester(cmd, input)
	if resulter.Error == nil {
		t.Error(`Command expected to get an error because of missing token and endpoint`)
	}
}

func CheckhErrorWhenNoEnvEndpointSet(t *testing.T, cmd *cobra.Command, input string) {
	// just token
	os.Setenv(ENV_VAR_TOKEN_NAME, "token_123")
	os.Unsetenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)

	// check
	resulter := FullCmdTester(cmd, input)
	if resulter.Error == nil {
		t.Error(`Command expected to get an error because of missing endpoint`)
	}
}

func CheckhErrorWhenNoEnvTokenSet(t *testing.T, cmd *cobra.Command, input string) {
	// set test server
	responseBody := "Miau"
	server := TestServer(200, responseBody, map[string]string{})
	defer server.Close()

	// just endpoitn
	os.Unsetenv(ENV_VAR_TOKEN_NAME)
	os.Setenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME, server.URL)

	// check
	resulter := FullCmdTester(cmd, input)
	if resulter.Error == nil {
		t.Error(`Command expected to get an error because of missing token`)
	}
}

func CheckCmdWorksWithEndpointAndTokenFlag(t *testing.T, cmd *cobra.Command, input string, responseBody string) {
	// clean env variablen
	os.Unsetenv(ENV_VAR_TOKEN_NAME)
	os.Unsetenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)

	resulter := FullCmdTester(cmd, input)

	buffer := new(bytes.Buffer)
	json.Compact(buffer, []byte(resulter.Output))
	if !strings.Contains(buffer.String(), responseBody) {
		t.Error(`Command response body doesn't match.'`)
	}

	if resulter.Error != nil {
		t.Error(`Command expected to not get an error`)
	}
}

func resetRootFlagVars() {
	Token = ""
	AutomationUrl = ""
}
