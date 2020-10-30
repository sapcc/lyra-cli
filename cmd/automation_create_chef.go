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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var AutomationCreateChefCmd = &cobra.Command{
	Use:   "chef",
	Short: locales.CmdShortDescription("automation-create-chef"),
	RunE: func(cmd *cobra.Command, args []string) error {
		chef = Chef{
			Automation: Automation{
				Name:                  viper.GetString("automation-create-chef-name"),
				Repository:            viper.GetString("automation-create-chef-repository"),
				RepositoryRevision:    viper.GetString("automation-create-chef-repository-revision"),
				RepositoryCredentials: viper.GetString("automation-create-chef-repository-credentials"),
				Timeout:               viper.GetInt("automation-create-chef-timeout"),
			},
			ChefVersion: viper.GetString("automation-create-chef-version"),
			Debug:       viper.GetBool("automation-create-chef-debug"),
		}

		// setup automation create chef attributes
		err := setupAutomationChefAttr(&chef)
		if err != nil {
			return err
		}

		// create automation
		response, err := automationCreateChef(&chef)
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
	AutomationCreateCmd.AddCommand(AutomationCreateChefCmd)
	initAutomationCreateChefCmdFlags()
}

func initAutomationCreateChefCmdFlags() {
	// flags
	AutomationCreateChefCmd.Flags().StringP("name", "", "", locales.AttributeDescription("automation-name"))
	AutomationCreateChefCmd.Flags().StringP("repository", "", "", locales.AttributeDescription("automation-repository"))
	AutomationCreateChefCmd.Flags().StringP("repository-credentials", "", "", locales.AttributeDescription("automation-repository-credentials"))
	AutomationCreateChefCmd.Flags().StringP("repository-revision", "", "master", locales.AttributeDescription("automation-repository-revision"))
	AutomationCreateChefCmd.Flags().IntP("timeout", "", 3600, locales.AttributeDescription("automation-timeout"))
	AutomationCreateChefCmd.Flags().BoolP("chef-debug", "", false, locales.AttributeDescription("automation-debug"))
	AutomationCreateChefCmd.Flags().StringP("runlist", "", "", locales.AttributeDescription("automation-runlist"))
	AutomationCreateChefCmd.Flags().StringP("chef-version", "", "", locales.AttributeDescription("automation-chef-version"))
	AutomationCreateChefCmd.Flags().StringP("attributes", "", "", locales.AttributeDescription("automation-attributes"))
	AutomationCreateChefCmd.Flags().StringP("attributes-from-file", "", "", locales.AttributeDescription("automation-attributes-from-file"))

	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-name", AutomationCreateChefCmd.Flags().Lookup("name")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-repository", AutomationCreateChefCmd.Flags().Lookup("repository")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-repository-credentials", AutomationCreateChefCmd.Flags().Lookup("repository-credentials")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-repository-revision", AutomationCreateChefCmd.Flags().Lookup("repository-revision")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-timeout", AutomationCreateChefCmd.Flags().Lookup("timeout")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-debug", AutomationCreateChefCmd.Flags().Lookup("chef-debug")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-version", AutomationCreateChefCmd.Flags().Lookup("chef-version")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-runlist", AutomationCreateChefCmd.Flags().Lookup("runlist")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-attributes", AutomationCreateChefCmd.Flags().Lookup("attributes")), "BindPFlag:")
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("automation-create-chef-attributes-from-file", AutomationCreateChefCmd.Flags().Lookup("attributes-from-file")), "BindPFlag:")
}

// private

func setupAutomationChefAttr(chef *Chef) error {
	chef.Runlist = helpers.StringToArray(viper.GetString("automation-create-chef-runlist"))

	// read attributes
	if len(viper.GetString("automation-create-chef-attributes")) > 0 {
		err := helpers.JSONStringToStructure(viper.GetString("automation-create-chef-attributes"), &chef.Attributes)

		if err != nil {
			return err
		}
	} else {
		attr, err := helpers.ReadFromFile(viper.GetString("automation-create-chef-attributes-from-file"))
		if err != nil {
			return err
		}
		if len(attr) > 0 {
			err = helpers.JSONStringToStructure(attr, &chef.Attributes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func automationCreateChef(chef *Chef) (string, error) {
	// add the type
	chef.AutomationType = "Chef"
	// convert to Json
	body, err := json.Marshal(chef)
	if err != nil {
		return "", err
	}

	automationClient := RestClient.Services["automation"]
	response, _, err := automationClient.Post("automations", url.Values{}, http.Header{}, string(body))
	if err != nil {
		return "", err
	}

	return response, nil
}
