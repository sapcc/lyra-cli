// Copyright © 2016 Arturo Reuschenbach <a.reuschenbach.puncernau@sap.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/print"
)

// automation/listCmd represents the automation/list command
var AutomationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available automations",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// setup rest client
		err := setupRestClient()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		response, err := automationList()
		if err != nil {
			return err
		}

		printer := print.Print{Data: response, Writer: os.Stdout}
		tablePrint, err := printer.TableList([]string{"id", "name", "repository", "repository_revision", "run_list", "chef_attributes", "log_level", "arguments"})
		if err != nil {
			return err
		}

		// print response
		cmd.Println(tablePrint)

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationListCmd)
}

func automationList() (interface{}, error) {
	// collect all automations do the pagination
	response, _, err := RestClient.Services.Automation.GetList("automations", url.Values{})
	if err != nil {
		return "", err
	}

	return response, nil
}
