package cmd

import (
	"fmt"
	"net/url"
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

	// add test command
	RootCmd.AddCommand(testCmd)
	// reset stuff
	ResetFlags()
	// run commando
	resulter := FullCmdTester(RootCmd, fmt.Sprintf("lyra test --lyra-service-endpoint=%s --arc-service-endpoint=%s --token=%s --debug", "http://somewhere.com", server.URL, "token123"))

	// request object
	if !strings.Contains(resulter.ErrorOutput, "User-Agent") {
		t.Error(fmt.Sprintf("Debug request object missing."))
		return
	}

	// response object
	if !strings.Contains(resulter.ErrorOutput, "Content-Length") {
		t.Error(fmt.Sprintf("Debug response object missing"))
	}

}
