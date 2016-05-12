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
	"net/url"

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/helpers"
)

// createCmd represents the create command
var AutomationCreateChefCmd = &cobra.Command{
	Use:   "chef",
	Short: "Create a new chef automation.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// setup automation create chef attributes
		err := setupAutomationCreateChef()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// create automation
		response, err := automationCreateChef()
		if err != nil {
			return err
		}
		// Print response
		cmd.Println(response)

		return nil
	},
}

func init() {
	AutomationCreateCmd.AddCommand(AutomationCreateChefCmd)
	// flags
	chef = Chef{}
	AutomationCreateChefCmd.Flags().StringVarP(&chef.Name, "name", "", "", "Describes the template. Should be short and alphanumeric without white spaces.")
	AutomationCreateChefCmd.Flags().StringVarP(&chef.Repository, "repository", "", "", "Describes the place where the automation is being described. Git ist the only suported repository type. Ex: https://github.com/user123/automation-test.git.")
	AutomationCreateChefCmd.Flags().StringVarP(&chef.RepositoryRevision, "repository-revision", "", "master", "Describes the repository branch.")
	AutomationCreateChefCmd.Flags().IntVarP(&chef.Timeout, "timeout", "", 3600, "Describes the time elapsed before a timeout is being triggered.")
	AutomationCreateChefCmd.Flags().StringVarP(&tags, "tags", "", "", "Are key value pairs. Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.")
	AutomationCreateChefCmd.Flags().StringVarP(&runlist, "runlist", "", "", "Describes the sequence of recipes should be executed. Runlist is an array of strings. Array of strings are separated by ','.")
	AutomationCreateChefCmd.Flags().StringVarP(&attributes, "attributes", "", "", "Attributes are JSON based.")
	AutomationCreateChefCmd.Flags().StringVarP(&attributesFromFile, "attributes-from-file", "", "", "Path to the file containing the chef attributes in JSON format. Giving a dash '-' will be read from standard input.")
	AutomationCreateChefCmd.Flags().StringVarP(&chef.LogLevel, "log-level", "", "", "Describe the level should be used when logging.")
}

// private

func setupAutomationCreateChef() error {
	chef.Tags = helpers.StringTokeyValueMap(tags)
	chef.Runlist = helpers.StringToArray(runlist)

	// read attributes
	if len(attributes) > 0 {
		err := helpers.JSONStringToStructure(attributes, &chef.Attributes)
		if err != nil {
			return err
		}
	} else {
		attr, err := helpers.ReadFromFile(attributesFromFile)
		if err != nil {
			return err
		}
		err = helpers.JSONStringToStructure(attr, &chef.Attributes)
		if err != nil {
			return err
		}
	}

	return nil
}

func automationCreateChef() (string, error) {
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
