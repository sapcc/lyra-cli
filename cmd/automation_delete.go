// Copyright © 2016 Arturo Reuschenbach Puncernau <a.reuschenbach.puncernau@sap.com>
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
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
)

var AutomationDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: locales.CmdShortDescription("automation-delete"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required node id
		if len(viper.GetString("automation-delete-id")) == 0 {
			return errors.New(locales.ErrorMessages("automation-id-missing"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := automationDelete(viper.GetString("automation-delete-id"))
		if err != nil {
			return err
		}
		// Print response
		cmd.Print("Automation with id ", viper.GetString("automation-delete-id"), " deleted.")

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationDeleteCmd)
	initAutomationDeleteCmdFlags()
}

func initAutomationDeleteCmdFlags() {
	AutomationDeleteCmd.Flags().StringP(FLAG_AUTOMATION_ID, "", "", locales.AttributeDescription(FLAG_AUTOMATION_ID))
	viper.BindPFlag("automation-delete-id", AutomationDeleteCmd.Flags().Lookup(FLAG_AUTOMATION_ID))
}

func automationDelete(id string) (string, error) {
	lyraService := RestClient.Services["automation"]
	response, _, err := lyraService.Delete(path.Join("automations", id), url.Values{})
	if err != nil {
		return "", err
	}
	return response, nil
}
