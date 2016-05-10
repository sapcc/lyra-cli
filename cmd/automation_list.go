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
	"github.com/spf13/cobra"
	"net/url"
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
		// print response
		cmd.Print(response)

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationListCmd)
}

func automationList() (string, error) {
	response, _, err := RestClient.Get("automations", url.Values{})
	if err != nil {
		return "", err
	}
	return response, nil
}
