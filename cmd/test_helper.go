// nolint
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	sgr "github.com/foize/go.sgr"
	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/restclient"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type resulter struct {
	Error       error
	Output      string
	ErrorOutput string
	Command     *cobra.Command
}

var cmdTestRootNoRun = &cobra.Command{
	Use:   "lyra-test",
	Short: "The root can run its own function",
	Long:  "The root description for help",
}

func FullCmdTester(testCommand *cobra.Command, input string) resulter {
	// pipe std out
	oldStdout := os.Stdout
	rOut, wOut, err := os.Pipe()
	if err != nil {
		os.Exit(1)
	}
	os.Stdout = wOut

	// pipe std err
	oldStderr := os.Stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		os.Exit(1)
	}
	os.Stderr = wErr

	// set buf for the command
	outputBuf := new(bytes.Buffer)
	testCommand.SetOutput(outputBuf)

	// set buf for the testCommand
	c := cmdTestRootNoRun
	rootOutputBuf := new(bytes.Buffer)
	c.SetOutput(rootOutputBuf)

	// add comand and run
	c.AddCommand(testCommand)
	c.SetArgs(argsSplit(input))
	err = c.Execute()

	// read std output
	outOutC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		if _, errorCopy := io.Copy(&buf, rOut); errorCopy != nil {
			return
		}
		outOutC <- buf.String()
	}()
	os.Stdout = oldStdout
	if errClose := wOut.Close(); errClose != nil {
		return resulter{err, "", "", c}
	}

	// read std outerr
	outErrC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		if _, errCopy := io.Copy(&buf, rErr); errCopy != nil {
			return
		}
		outErrC <- buf.String()
	}()
	os.Stderr = oldStderr
	if errClose := wErr.Close(); errClose != nil {
		return resulter{err, "", "", c}
	}

	//save outputs and return
	output := <-outOutC
	outErr := fmt.Sprint(outputBuf.String(), rootOutputBuf.String(), <-outErrC)
	return resulter{err, output, outErr, c}
}

// added from
// https://gist.github.com/jmervine/d88c75329f98e09f5c87
func argsSplit(s string) []string {
	split := strings.Split(s, " ")
	var splitedArgs []string
	var argsInQuote string
	var block string
	for _, v := range split {
		if argsInQuote == "" {
			if strings.HasPrefix(v, "'") || strings.HasPrefix(v, "\"") {
				argsInQuote = string(v[0])
				block = strings.TrimPrefix(v, argsInQuote) + " "
			} else {
				splitedArgs = append(splitedArgs, v)
			}
		} else {
			if !strings.HasSuffix(v, argsInQuote) {
				block += v + " "
			} else {
				block += strings.TrimSuffix(v, argsInQuote)
				argsInQuote = ""
				splitedArgs = append(splitedArgs, block)
				block = ""
			}
		}
	}
	return splitedArgs
}

// TestServer should be closed afterwards
func TestServer(code int, body string, headers map[string]string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(code) // keep the code after setting headers. If not they will disapear...
		if _, err := fmt.Fprintln(w, body); err != nil {
			fmt.Errorf("%v", err)
		}
	}))
	return server
}

func newMockAuthenticationV3(testServer *httptest.Server) func(authOpts auth.AuthOptions) auth.Authentication {
	return func(authOpts auth.AuthOptions) auth.Authentication {
		return &auth.MockV3{Options: authOpts, TestServer: testServer}
	}
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

	errorMsg := fmt.Sprint(locales.ErrorMessages("flag-missing"), FLAG_USER_ID, ", ", FLAG_USERNAME)
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
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

	errorMsg := fmt.Sprint(locales.ErrorMessages("flag-missing"), FLAG_USER_ID, ", ", FLAG_USERNAME)
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
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

	errorMsg := fmt.Sprint(locales.ErrorMessages("flag-missing"), FLAG_USER_ID, ", ", FLAG_USERNAME)
	if !strings.Contains(resulter.ErrorOutput, errorMsg) {
		diffString := StringDiff(resulter.ErrorOutput, errorMsg)
		t.Error(fmt.Sprintf("Command error doesn't match. \n \n %s", diffString))
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

func JsonDiff(responseBody, resulterOutput string) (bool, error) {
	source := map[string]interface{}{}
	err := json.Unmarshal([]byte(responseBody), &source)
	if err != nil {
		return false, err
	}

	response := map[string]interface{}{}
	err = json.Unmarshal([]byte(resulterOutput), &response)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(source, response), nil
}

func JsonListDiff(responseBody, resulterOutput string) (bool, error) {
	source := []map[string]interface{}{}
	err := json.Unmarshal([]byte(responseBody), &source)
	if err != nil {
		return false, err
	}

	response := []map[string]interface{}{}
	err = json.Unmarshal([]byte(resulterOutput), &response)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(source, response), nil
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
	RestClient = restclient.NewClient([]restclient.Endpoint{}, "", false)

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
	AutomationDeleteCmd.ResetFlags()
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
	NodeDeleteCmd.ResetFlags()
	NodeInstallCmd.ResetFlags()
	NodeListCmd.ResetFlags()
	NodeFactListCmd.ResetFlags()
	NodeTagCmd.ResetFlags()
	NodeTagAddCmd.ResetFlags()
	NodeTagDeleteCmd.ResetFlags()
	NodeTagListCmd.ResetFlags()
	NodeShowCmd.ResetFlags()
	RunListCmd.ResetFlags()
	RunShowCmd.ResetFlags()
	RunCmd.ResetFlags()
	// set flags again
	initRootCmdFlags()
	initAuthenticationCmdFlags()
	initAutomationCreateChefCmdFlags()
	initAutomationCreateScriptCmdFlags()
	initAutomationCreateCmdFlags()
	initAutomationDeleteCmdFlags()
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
	initNodeDeleteCmdFlags()
	initNodeInstallCmdFlags()
	initNodeListCmdFlags()
	initNodeShowCmdFlags()
	initNodeFactListCmdFlags()
	initNodeTagAddCmdFlags()
	initNodeTagDeleteCmdFlags()
	initNodeTagListCmdFlags()
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
