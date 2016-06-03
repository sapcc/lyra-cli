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
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

var JobShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Shows an especific job",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check required job id
		if len(viper.GetString("show-job-id")) == 0 {
			return errors.New(locales.ErrorMessages("job-id-missing"))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// list automation
		response, err := jobShow(viper.GetString("show-job-id"))
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
	JobCmd.AddCommand(JobShowCmd)
	initJobShowCmdFlags()
}

func initJobShowCmdFlags() {
	JobShowCmd.Flags().StringP(FLAG_JOB_ID, "", "", locales.AttributeDescription("job-id"))
	viper.BindPFlag("show-job-id", JobShowCmd.Flags().Lookup(FLAG_JOB_ID))
}

func jobShow(id string) (string, error) {
	response, _, err := RestClient.Services.Arc.Get(path.Join("jobs", id), url.Values{}, false)
	if err != nil {
		return "", err
	}

	return response, nil
}
