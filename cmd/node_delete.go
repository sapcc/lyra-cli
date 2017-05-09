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

var NodeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: locales.CmdShortDescription("arc-node-delete"),
	Long:  locales.CmdLongDescription("arc-node-delete"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required node id
		if len(viper.GetString("arc-delete-node-id")) == 0 {
			return errors.New(locales.ErrorMessages("node-id-missing"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		_, err := nodeDelete(viper.GetString("arc-delete-node-id"))
		if err != nil {
			return err
		}

		// Print response to std error. No response got
		cmd.Println("Node with id ", viper.GetString("arc-delete-node-id"), " deleted.")

		return nil
	},
}

func init() {
	NodeCmd.AddCommand(NodeDeleteCmd)
	initNodeDeleteCmdFlags()
}

func initNodeDeleteCmdFlags() {
	NodeDeleteCmd.Flags().StringP(FLAG_ARC_NODE_ID, "", "", locales.AttributeDescription("arc-node-id"))
	viper.BindPFlag("arc-delete-node-id", NodeDeleteCmd.Flags().Lookup(FLAG_ARC_NODE_ID))
}

func nodeDelete(id string) (string, error) {
	arcService := RestClient.Services["arc"]
	response, _, err := arcService.Delete(path.Join("agents", id), url.Values{})
	if err != nil {
		return "", err
	}

	return response, nil
}
