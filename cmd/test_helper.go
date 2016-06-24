package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/foize/go.sgr"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/restclient"
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
	os.Unsetenv(ENV_VAR_ARC_ENDPOINT_NAME)

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
	os.Unsetenv(ENV_VAR_ARC_ENDPOINT_NAME)

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
	os.Setenv(ENV_VAR_ARC_ENDPOINT_NAME, server.URL)

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

func StringDiff(text1, text2 string) string {
	// find diffs
	dmp := diffmatchpatch.New()
	test := dmp.DiffMain(text1, text2, false)
	diffs := dmp.DiffCleanupSemantic(test)

	// output with colors
	var buffer bytes.Buffer
	for _, v := range diffs {
		// scape text
		v.Text = strings.Replace(v.Text, "[", "[[", -1)
		v.Text = strings.Replace(v.Text, "]", "]]", -1)

		if v.Type == 0 {
			buffer.WriteString(v.Text)
		} else if v.Type == -1 {
			buffer.WriteString("[bg-red bold]")
			buffer.WriteString(v.Text)
			buffer.WriteString("[reset]")
		} else if v.Type == 1 {
			buffer.WriteString("[bg-blue bold]")
			buffer.WriteString(v.Text)
			buffer.WriteString("[reset]")
		}
	}
	// parse to set colors
	colorDiff := sgr.MustParseln(buffer.String())
	return colorDiff
}

func ResetFlags() {
	// reset other stuff
	RestClient = restclient.NewClient([]restclient.Endpoint{}, "")

	// Remove env variablen
	os.Unsetenv(ENV_VAR_TOKEN_NAME)
	os.Unsetenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)
	os.Unsetenv(ENV_VAR_ARC_ENDPOINT_NAME)
	// Reset viper
	viper.Reset()
	// Reset command flags
	RootCmd.ResetFlags()
	AuthenticateCmd.ResetFlags()
	AutomationCreateChefCmd.ResetFlags()
	AutomationCreateScriptCmd.ResetFlags()
	AutomationCreateCmd.ResetFlags()
	AutomationExecuteCmd.ResetFlags()
	AutomationListCmd.ResetFlags()
	AutomationShowCmd.ResetFlags()
	AutomationUpdateChefAttributesCmd.ResetFlags()
	AutomationUpdateChefCmd.ResetFlags()
	AutomationUpdateCmd.ResetFlags()
	AutomationCmd.ResetFlags()
	JobListCmd.ResetFlags()
	JobLogCmd.ResetFlags()
	JobShowCmd.ResetFlags()
	JobCmd.ResetFlags()
	NodeCmd.ResetFlags()
	NodeInstallCmd.ResetFlags()
	RunListCmd.ResetFlags()
	RunShowCmd.ResetFlags()
	RunCmd.ResetFlags()
	// set flags again
	initRootCmdFlags()
	initAuthenticationCmdFlags()
	initAutomationCreateChefCmdFlags()
	initAutomationCreateScriptCmdFlags()
	initAutomationCreateCmdFlags()
	initAutomationExecuteCmdFlags()
	initAutomationListCmdFlags()
	initAutomationShowCmdFlags()
	initAutomationUpdateChefAttributesCmdFlags()
	initAutomationUpdateChefCmdFlags()
	initAutomationUpdateCmdFlags()
	initAutomationCmdFlags()
	initJobListCmdFlags()
	initJobLogCmdFlags()
	initJobShowCmdFlags()
	initJobCmdFlags()
	initNodeCmdFlags()
	initNodeInstallCmdFlags()
	initRunListCmdFlags()
	initRunShowCmdFlags()
	initRunCmdFlags()
}

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
