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
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

var NodeListCmd = &cobra.Command{
	Use:   "list",
	Short: locales.CmdShortDescription("arc-node-list"),
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		response, err := nodeList()
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
			tablePrint, err = printer.TableList([]string{"agent_id", "organization", "project", "created_at", "updated_at", "updated_by", "updated_with"})
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
	NodeCmd.AddCommand(NodeListCmd)
	initNodeListCmdFlags()
}

func initNodeListCmdFlags() {
	//flags
	NodeListCmd.Flags().StringP(FLAG_SELECTOR, "", "", locales.AttributeDescription("node-selector"))
	viper.BindPFlag("node-selector", NodeListCmd.Flags().Lookup(FLAG_SELECTOR))
}

func nodeList() (interface{}, error) {
	// collect all automations do the pagination
	arcService := RestClient.Services["arc"]
	urlValues := url.Values{}
	if viper.GetString("node-selector") != "" {
		urlValues.Set("q", viper.GetString("node-selector"))
	}
	response, _, err := arcService.GetList("agents", urlValues)
	if err != nil {
		return "", err
	}

	return response, nil
}
