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

	"github.com/sapcc/lyra-cli/locales"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var NodeTagDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: locales.CmdShortDescription("arc-node-tag-delete"),
	Long:  locales.CmdLongDescription("arc-node-tag-delete"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required node id
		if len(viper.GetString("arc-tag-delete-node-id")) == 0 {
			return errors.New(locales.ErrorMessages("node-id-missing"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// loop over the tag keys to delete
		for _, element := range args {
			_, err := nodeTagdelete(viper.GetString("arc-tag-delete-node-id"), element)
			if err != nil {
				return err
			}

			// Print response
			cmd.Println("Tag from agent with id ", viper.GetString("arc-tag-delete-node-id"), " and value ", element, " is deleted.")
		}

		return nil
	},
}

func init() {
	NodeTagCmd.AddCommand(NodeTagDeleteCmd)
	initNodeTagDeleteCmdFlags()
}

func initNodeTagDeleteCmdFlags() {
	NodeTagDeleteCmd.Flags().StringP(FLAG_ARC_NODE_ID, "", "", locales.AttributeDescription(FLAG_ARC_NODE_ID))
	viper.BindPFlag("arc-tag-delete-node-id", NodeTagDeleteCmd.Flags().Lookup(FLAG_ARC_NODE_ID))
}

func nodeTagdelete(id, tagKey string) (string, error) {
	arcService := RestClient.Services["arc"]
	response, _, err := arcService.Delete(path.Join("agents", id, "tags", tagKey), url.Values{})
	if err != nil {
		return "", err
	}
	return response, nil
}
