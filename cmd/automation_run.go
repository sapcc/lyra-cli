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

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/locales"
)

// updateCmd represents the update command
var AutomationRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs an exsiting automation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// setup automation run attributes
		err := setupAutomationRun()
		if err != nil {
			return err
		}

		// setup rest client
		err = setupRestClient()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// run automation
		response, err := automationRun()
		if err != nil {
			return err
		}
		// Print response
		cmd.Print(response)

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationRunCmd)
	//flags
	run = Run{}
	AutomationRunCmd.Flags().StringVarP(&run.AutomationId, "id", "", "", locales.AttributeDescription("automation-id"))
	AutomationRunCmd.Flags().StringVarP(&run.Selector, "selector", "", "", locales.AttributeDescription("automation-selector"))
}

func setupAutomationRun() error {
	// check required automation id
	if len(run.AutomationId) == 0 {
		return errors.New(locales.ErrorMessages("automation-id-missing"))
	}
	// check selector
	if len(run.Selector) == 0 {
		return errors.New(locales.ErrorMessages("automation-selector-missing"))
	}

	return nil
}

func automationRun() (string, error) {
	// convert to Json
	body, err := json.Marshal(run)
	if err != nil {
		return "", err
	}
	// send request
	response, _, err := RestClient.Post("runs", url.Values{}, string(body))
	if err != nil {
		return "", err
	}

	return response, nil
}
