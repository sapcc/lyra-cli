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

var JobLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Shows job log",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required job id
		if len(jobId) == 0 {
			return errors.New(locales.ErrorMessages("job-id-missing"))
		}
		// setup rest client
		err := setupRestClient()
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		response, err := jobLog()
		if err != nil {
			return err
		}

		// print response
		cmd.Println(response)

		return nil
	},
}

func init() {
	JobCmd.AddCommand(JobLogCmd)
	JobLogCmd.Flags().StringVar(&jobId, "job-id", "", locales.AttributeDescription("job-id"))
}

func jobLog() (string, error) {
	response, _, err := RestClient.Services.Arc.Get(path.Join("jobs", jobId, "log"), url.Values{}, false)
	if err != nil {
		return "", err
	}

	return response, nil
}
