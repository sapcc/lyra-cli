package cmd

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use: "test",
	RunE: func(cmd *cobra.Command, args []string) error {
		arcService := RestClient.Services["arc"]
		_, _, err := arcService.Get("", url.Values{}, false)
		if err != nil {
			return err
		}
		return nil
	},
}

func TestRootDebugFlag(t *testing.T) {
	server := TestServer(200, "Arc API", map[string]string{})
	defer server.Close()

	// keep backup of the real stdout
	oldStdout := os.Stdout

	// pipe std out
	r, w, err := os.Pipe()
	if err != nil {
		t.Error(err)
	}
	os.Stdout = w

	// add test command
	RootCmd.AddCommand(testCmd)
	// reset stuff
	ResetFlags()
	// run commando
	FullCmdTester(RootCmd, fmt.Sprintf("lyra test --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --debug", "http://somewhere.com", server.URL, "token123"))

	// flush, restore close
	os.Stdout = oldStdout
	w.Close()

	// read log
	loggedStuff, _ := ioutil.ReadAll(r)

	// request object
	if !strings.Contains(string(loggedStuff), "User-Agent") {
		t.Error(fmt.Sprintf("Debug request object missing."))
	}

	// response object
	if !strings.Contains(string(loggedStuff), "Content-Length") {
		t.Error(fmt.Sprintf("Debug response object missing"))
	}

}
