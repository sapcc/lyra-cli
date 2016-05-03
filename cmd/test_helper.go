package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestServer(code int, body string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	return server
}
