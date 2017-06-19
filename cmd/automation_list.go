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
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

// automation/listCmd represents the automation/list command
var AutomationListCmd = &cobra.Command{
	Use:   "list",
	Short: locales.CmdShortDescription("automation-list"),
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		response, err := automationList()
		if err != nil {
			return err
		}

		printer := print.Print{Data: response}
		tablePrint := ""
		if viper.GetBool("json") {
			tablePrint, err = printer.JSON()
			if err != nil {
				return err
			}
		} else {
			tablePrint, err = printer.TableList([]string{"id", "name", "type", "repository", "repository_revision", "timeout", "run_list", "chef_version", "debug"})
			if err != nil {
				return err
			}
		}

		// print response
		fmt.Println(tablePrint)

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationListCmd)
	initAutomationListCmdFlags()
}

func initAutomationListCmdFlags() {
}

func automationList() (interface{}, error) {
	// collect all automations do the pagination
	automationService := RestClient.Services["automation"]
	response, _, err := automationService.GetList("automations", url.Values{})
	if err != nil {
		return "", err
	}

	return response, nil
}
