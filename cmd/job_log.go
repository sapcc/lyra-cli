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
	"fmt"
	"net/url"
	"path"

	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var JobLogCmd = &cobra.Command{
	Use:   "log",
	Short: locales.CmdShortDescription("job-log"),
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
		fmt.Println(response)

		return nil
	},
}

func init() {
	JobCmd.AddCommand(JobLogCmd)
	initJobLogCmdFlags()
}

func initJobLogCmdFlags() {
	JobLogCmd.Flags().StringP(FLAG_JOB_ID, "", "", locales.AttributeDescription("job-id"))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("log-job-id", JobLogCmd.Flags().Lookup(FLAG_JOB_ID)), "BindPFlag:")
}

func jobLog(id string) (string, error) {
	arcService := RestClient.Services["arc"]
	response, _, err := arcService.Get(path.Join("jobs", id, "log"), url.Values{}, false)
	if err != nil {
		return "", err
	}

	return response, nil
}
