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
	"fmt"
	"net/url"
	"path"

	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showCmd represents the show command
var AutomationShowCmd = &cobra.Command{
	Use:   "show",
	Short: locales.CmdShortDescription("automation-show"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required automation id
		if len(viper.GetString("show-automation-id")) == 0 {
			return errors.New("no automation id given")
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
		var bodyPrint string
		if viper.GetBool("json") {
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
		fmt.Println(bodyPrint)

		return nil
	},
}

func init() {
	initAutomationShowCmdFlags()
}

func initAutomationShowCmdFlags() {
	AutomationCmd.AddCommand(AutomationShowCmd)
	AutomationShowCmd.Flags().StringP("automation-id", "", "", locales.AttributeDescription("automation-id"))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("show-automation-id", AutomationShowCmd.Flags().Lookup("automation-id")), "BindPFlag:")
}

func automationShow() (string, error) {
	automationService := RestClient.Services["automation"]
	response, _, err := automationService.Get(path.Join("automations", viper.GetString("show-automation-id")), url.Values{}, false)
	if err != nil {
		return "", err
	}

	return response, nil
}
