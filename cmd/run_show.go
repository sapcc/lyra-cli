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
	"errors"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/locales"
)

// showCmd represents the show command
var RunShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a specific automation run",
	Long:  `A longer description for automation run show.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required automation id
		if len(runId) == 0 {
			return errors.New(locales.ErrorMessages("run-id-missing"))
		}
		// setup rest client
		err := setupRestClient()
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// show automation
		response, err := runShow()
		if err != nil {
			return err
		}
		// print response
		cmd.Println(response)

		return nil
	},
}

func init() {
	RunCmd.AddCommand(RunShowCmd)
	RunShowCmd.Flags().StringVar(&runId, "run-id", "", locales.AttributeDescription("run-id"))
}

func runShow() (string, error) {
	response, _, err := RestClient.Get(path.Join("runs", runId), url.Values{})
	if err != nil {
		return "", err
	}
	return response, nil
}
