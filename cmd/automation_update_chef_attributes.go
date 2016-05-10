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
	"errors"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/helpers"
)

// updateCmd represents the update command
var AutomationUpdateChefAttributesCmd = &cobra.Command{
	Use:   "attributes",
	Short: "Updates chef attributes",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// setup update chef attributes
		err := setupAutomationUpdateChefAttributes()
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// update automation
		response, err := automationUpdateChefAttributes()
		if err != nil {
			return err
		}
		// Print response
		cmd.Print(response)

		return nil
	},
}

func init() {
	AutomationUpdateChefCmd.AddCommand(AutomationUpdateChefAttributesCmd)
	//flags
	AutomationUpdateChefAttributesCmd.Flags().StringVarP(&attributes, "attributes", "", "", "Attributes are JSON format.")
	AutomationUpdateChefAttributesCmd.Flags().StringVarP(&attributesFromFile, "attributes-from-file", "", "", "Path to the file containing the chef attributes in JSON format. Giving a dash '-' will be read from standard input.")
}

func setupAutomationUpdateChefAttributes() error {
	// read attributes
	if len(attributes) > 0 {
		chef.Attributes = attributes
	} else {
		attr, err := helpers.ReadFromFile(attributesFromFile)
		if err != nil {
			return err
		}
		chef.Attributes = attr
	}

	return nil
}

func automationUpdateChefAttributes() (string, error) {
	response, code, err := RestClient.Get(path.Join("automations", automationId), url.Values{})
	if err != nil {
		return "", err
	}

	if int(code) >= 400 {
		return "", errors.New(response)
	}

	// get the existing data
	chef = Chef{}
	respByt := []byte(response)
	if err := json.Unmarshal(respByt, &chef); err != nil {
		return "", err
	}

	// change attributres
	chef.Attributes = attributes

	// convert to Json
	body, err := json.Marshal(chef)
	if err != nil {
		return "", err
	}

	// send data back
	newResp, _, err := RestClient.Put(path.Join("automations", automationId), url.Values{}, string(body))
	if err != nil {
		return "", err
	}

	return newResp, nil
}
