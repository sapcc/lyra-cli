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

// updateCmd represents the update command
var AutomationUpdateChefAttributesCmd = &cobra.Command{
	Use:   "attributes",
	Short: locales.CmdShortDescription("automation-update-chef-attributes"),
	Long:  locales.CmdLongDescription("automation-update-chef-attributes"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required automation id
		if len(viper.GetString("automation-update-chef-attributes-automation-id")) == 0 {
			return errors.New(locales.ErrorMessages("automation-id-missing"))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		chef = Chef{}

		// setup update chef attributes
		err := setupAutomationUpdateChefAttributes(&chef)
		if err != nil {
			return err
		}

		// update automation
		response, err := automationUpdateChefAttributes(&chef)
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

		// print response
		fmt.Println(bodyPrint)

		return nil
	},
}

func init() {
	AutomationUpdateChefCmd.AddCommand(AutomationUpdateChefAttributesCmd)
	initAutomationUpdateChefAttributesCmdFlags()
}

func initAutomationUpdateChefAttributesCmdFlags() {
	AutomationUpdateChefAttributesCmd.Flags().StringP("attributes", "", "", locales.AttributeDescription("automation-attributes"))
	AutomationUpdateChefAttributesCmd.Flags().StringP("attributes-from-file", "", "", locales.AttributeDescription("automation-attributes-from-file"))
	AutomationUpdateChefAttributesCmd.Flags().StringP(FLAG_AUTOMATION_ID, "", "", locales.AttributeDescription("automation-id"))

	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-update-chef-attributes", AutomationUpdateChefAttributesCmd.Flags().Lookup("attributes")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-update-chef-attributes-from-file", AutomationUpdateChefAttributesCmd.Flags().Lookup("attributes-from-file")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-update-chef-attributes-automation-id", AutomationUpdateChefAttributesCmd.Flags().Lookup("automation-id")), "BindPFlag:")
}

func setupAutomationUpdateChefAttributes(chefObj *Chef) error {
	// read attributes
	if len(viper.GetString("automation-update-chef-attributes")) > 0 {
		err := helpers.JSONStringToStructure(viper.GetString("automation-update-chef-attributes"), &chefObj.Attributes)
		if err != nil {
			return err
		}
	} else {
		attr, err := helpers.ReadFromFile(viper.GetString("automation-update-chef-attributes-from-file"))
		if err != nil {
			return err
		}
		if len(attr) > 0 {
			err = helpers.JSONStringToStructure(attr, &chefObj.Attributes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func automationUpdateChefAttributes(chefObj *Chef) (string, error) {
	automationService := RestClient.Services["automation"]
	response, code, err := automationService.Get(path.Join("automations", viper.GetString("automation-update-chef-attributes-automation-id")), url.Values{}, false)
	if err != nil {
		return "", err
	}

	if int(code) >= 400 {
		return "", errors.New(response)
	}

	// map response to the automation object
	oldChef := Chef{}
	if err := oldChef.Unmarshal(response); err != nil {
		return "", err
	}

	// change attributres
	oldChef.Attributes = chefObj.Attributes

	body, err := oldChef.Marshal()
	if err != nil {
		return "", err
	}

	// send data back
	newResp, _, err := automationService.Put(path.Join("automations", viper.GetString("automation-update-chef-attributes-automation-id")), url.Values{}, string(body))
	if err != nil {
		return "", err
	}

	return newResp, nil
}
