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
	"errors"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

// showCmd represents the show command
var AutomationShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a specific automation",
	Long:  `A longer description for automation show.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required automation id
		if len(automationId) == 0 {
			return errors.New("No automation id given.")
		}
		// setup rest client
		err := setupRestClient()
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// show automation
		response, err := automationShow()
		if err != nil {
			return err
		}

		// convert data to struct
		var dataStruct map[string]interface{}
		err = helpers.JSONStringToStructure(response, &dataStruct)
		if err != nil {
			return err
		}

		// print the data out
		printer := print.Print{Data: dataStruct}
		bodyPrint := ""
		if JsonOutput {
			bodyPrint, err = printer.JSON()
			if err != nil {
				return err
			}
		} else {
			bodyPrint, err = printer.Table()
			if err != nil {
				return err
			}
		}

		// Print response
		cmd.Println(bodyPrint)

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationShowCmd)
	AutomationShowCmd.Flags().StringVar(&automationId, "automation-id", "", locales.AttributeDescription("automation-id"))
}

func automationShow() (string, error) {
	response, _, err := RestClient.Services.Automation.Get(path.Join("automations", automationId), url.Values{}, false)
	if err != nil {
		return "", err
	}
	return response, nil
}
