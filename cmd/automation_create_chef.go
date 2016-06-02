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
	"net/url"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/print"
)

// createCmd represents the create command
var AutomationCreateChefCmd = &cobra.Command{
	Use:   "chef",
	Short: "Create a new chef automation.",
	RunE: func(cmd *cobra.Command, args []string) error {
		chef = Chef{
			Name:               viper.GetString("automation-create-chef-name"),
			Repository:         viper.GetString("automation-create-chef-repository"),
			RepositoryRevision: viper.GetString("automation-create-chef-repository-revision"),
			Timeout:            viper.GetInt("automation-create-chef-timeout"),
			LogLevel:           viper.GetString("automation-create-chef-log-level"),
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
		bodyPrint := ""
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
		cmd.Println(bodyPrint)

		return nil
	},
}

func init() {
	AutomationCreateCmd.AddCommand(AutomationCreateChefCmd)
	initAutomationCreateChefCmdFlags()
}

func initAutomationCreateChefCmdFlags() {
	// flags
	AutomationCreateChefCmd.Flags().StringP("name", "", "", "Describes the template. Should be short and alphanumeric without white spaces.")
	AutomationCreateChefCmd.Flags().StringP("repository", "", "", "Describes the place where the automation is being described. Git ist the only suported repository type. Ex: https://github.com/user123/automation-test.git.")
	AutomationCreateChefCmd.Flags().StringP("repository-revision", "", "master", "Describes the repository branch.")
	AutomationCreateChefCmd.Flags().IntP("timeout", "", 3600, "Describes the time elapsed before a timeout is being triggered.")
	AutomationCreateChefCmd.Flags().StringP("log-level", "", "", "Describe the level should be used when logging.")
	AutomationCreateChefCmd.Flags().StringP("tags", "", "", "Are key value pairs. Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.")
	AutomationCreateChefCmd.Flags().StringP("runlist", "", "", "Describes the sequence of recipes should be executed. Runlist is an array of strings. Array of strings are separated by ','.")
	AutomationCreateChefCmd.Flags().StringP("attributes", "", "", "Attributes are JSON based.")
	AutomationCreateChefCmd.Flags().StringP("attributes-from-file", "", "", "Path to the file containing the chef attributes in JSON format. Giving a dash '-' will be read from standard input.")
	viper.BindPFlag("automation-create-chef-name", AutomationCreateChefCmd.Flags().Lookup("name"))
	viper.BindPFlag("automation-create-chef-repository", AutomationCreateChefCmd.Flags().Lookup("repository"))
	viper.BindPFlag("automation-create-chef-repository-revision", AutomationCreateChefCmd.Flags().Lookup("repository-revision"))
	viper.BindPFlag("automation-create-chef-timeout", AutomationCreateChefCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("automation-create-chef-log-level", AutomationCreateChefCmd.Flags().Lookup("log-level"))
	viper.BindPFlag("automation-create-chef-tags", AutomationCreateChefCmd.Flags().Lookup("tags"))
	viper.BindPFlag("automation-create-chef-runlist", AutomationCreateChefCmd.Flags().Lookup("runlist"))
	viper.BindPFlag("automation-create-chef-attributes", AutomationCreateChefCmd.Flags().Lookup("attributes"))
	viper.BindPFlag("automation-create-chef-attributes-from-file", AutomationCreateChefCmd.Flags().Lookup("attributes-from-file"))
}

// private

func createViperKey(flag string) string {
	_, filename, _, _ := runtime.Caller(1)

	fmt.Println("°°°")
	fmt.Println(filename)
	fmt.Println("°°°")

	return fmt.Sprint(filename, "_", flag)
}

func setupAutomationChefAttr(chef *Chef) error {
	chef.Tags = helpers.StringTokeyValueMap(viper.GetString("automation-create-chef-tags"))
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

	response, _, err := RestClient.Services.Automation.Post("automations", url.Values{}, string(body))
	if err != nil {
		return "", err
	}

	return response, nil
}
