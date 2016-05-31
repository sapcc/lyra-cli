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
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
)

var JobLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Shows job log",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required job id
		if len(viper.GetString("log-job-id")) == 0 {
			return errors.New(locales.ErrorMessages("job-id-missing"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		response, err := jobLog(viper.GetString("log-job-id"))
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
	initJobLogCmdFlags()
}

func initJobLogCmdFlags() {
	JobLogCmd.Flags().StringP(FLAG_JOB_ID, "", "", locales.AttributeDescription(FLAG_JOB_ID))
	viper.BindPFlag("log-job-id", JobLogCmd.Flags().Lookup(FLAG_JOB_ID))
}

func jobLog(id string) (string, error) {
	response, _, err := RestClient.Services.Arc.Get(path.Join("jobs", id, "log"), url.Values{}, false)
	if err != nil {
		return "", err
	}

	return response, nil
}
