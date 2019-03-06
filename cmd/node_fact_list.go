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
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

var NodeFactListCmd = &cobra.Command{
	Use:   "list",
	Short: locales.CmdShortDescription("arc-node-fact-list"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required node id
		if len(viper.GetString("arc-fact-list-node-id")) == 0 {
			return errors.New(locales.ErrorMessages("node-id-missing"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		response, err := nodeFactList(viper.GetString("arc-fact-list-node-id"))
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
	NodeFactCmd.AddCommand(NodeFactListCmd)
	initNodeFactListCmdFlags()
}

func initNodeFactListCmdFlags() {
	NodeFactListCmd.Flags().StringP(FLAG_ARC_NODE_ID, "", "", locales.AttributeDescription(FLAG_ARC_NODE_ID))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("arc-fact-list-node-id", NodeFactListCmd.Flags().Lookup(FLAG_ARC_NODE_ID)), "BindPFlag:")
}

func nodeFactList(id string) (string, error) {
	arcService := RestClient.Services["arc"]
	response, _, err := arcService.Get(path.Join("agents", id, "facts"), url.Values{}, false)
	if err != nil {
		return "", err
	}

	return response, nil
}
