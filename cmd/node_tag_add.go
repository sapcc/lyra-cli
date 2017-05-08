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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
)

var NodeTagAddCmd = &cobra.Command{
	Use:   "add",
	Short: locales.CmdShortDescription("arc-node-tag-add"),
	Long:  locales.CmdLongDescription("arc-node-tag-add"),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required node id
		if len(viper.GetString("arc-tag-add-node-id")) == 0 {
			return errors.New(locales.ErrorMessages("node-id-missing"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// parse arguments
		body, err := parseArgs(args)
		if err != nil {
			return err
		}

		// post tags
		_, err = nodeTagAdd(viper.GetString("arc-tag-add-node-id"), body)
		if err != nil {
			return err
		}

		// Print response
		fmt.Fprintln(os.Stderr, "Tags added successfully to the node with id ", viper.GetString("arc-tag-add-node-id"))

		return nil
	},
}

func init() {
	NodeTagCmd.AddCommand(NodeTagAddCmd)
	initNodeTagAddCmdFlags()
}

func initNodeTagAddCmdFlags() {
	NodeTagAddCmd.Flags().StringP(FLAG_ARC_NODE_ID, "", "", locales.AttributeDescription(FLAG_ARC_NODE_ID))
	viper.BindPFlag("arc-tag-add-node-id", NodeTagAddCmd.Flags().Lookup(FLAG_ARC_NODE_ID))
}

func parseArgs(args []string) (string, error) {
	keyValuePairs := map[string]string{}
	for _, element := range args {
		data := regexp.MustCompile("=|:").Split(element, 2)
		if len(data) != 2 {
			continue
		}
		keyValuePairs[data[0]] = data[1]
	}

	jsonString, err := json.Marshal(keyValuePairs)
	if err != nil {
		return "", err
	}

	return string(jsonString), nil
}

func nodeTagAdd(id, body string) (string, error) {
	arcService := RestClient.Services["arc"]
	response, _, err := arcService.Post(path.Join("agents", id, "tags"), url.Values{}, http.Header{}, body)
	if err != nil {
		return "", err
	}

	return response, nil
}
