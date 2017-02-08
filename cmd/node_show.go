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
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

var NodeShowCmd = &cobra.Command{
	Use:   "show",
	Short: locales.CmdShortDescription("arc-node-show"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required node id
		if len(viper.GetString("arc-show-node-id")) == 0 {
			return errors.New(locales.ErrorMessages("node-id-missing"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		response, err := nodeShow(viper.GetString("arc-show-node-id"))
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

		// print response
		cmd.Println(bodyPrint)

		return nil
	},
}

func init() {
	NodeCmd.AddCommand(NodeShowCmd)
	initNodeShowCmdFlags()
}

func initNodeShowCmdFlags() {
	NodeShowCmd.Flags().StringP(FLAG_ARC_NODE_ID, "", "", locales.AttributeDescription(FLAG_ARC_NODE_ID))
	viper.BindPFlag("arc-show-node-id", NodeShowCmd.Flags().Lookup(FLAG_ARC_NODE_ID))
}

func nodeShow(id string) (string, error) {
	arcService := RestClient.Services["arc"]
	response, _, err := arcService.Get(path.Join("agents", id), url.Values{}, false)
	if err != nil {
		return "", err
	}

	return response, nil
}
