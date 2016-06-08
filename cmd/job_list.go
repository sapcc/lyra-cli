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
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

var JobListCmd = &cobra.Command{
	Use:   "list",
	Short: locales.CmdShortDescription("job-list"),
	RunE: func(cmd *cobra.Command, args []string) error {
		// show automation
		response, err := jobList()
		if err != nil {
			return err
		}

		printer := print.Print{Data: response}
		tablePrint := ""
		if viper.GetBool("json") {
			tablePrint, err = printer.JSON()
			if err != nil {
				return err
			}
		} else {
			tablePrint, err = printer.TableList([]string{"request_id", "status", "action", "agent", "user_id", "created_at"})
			if err != nil {
				return err
			}
		}

		// print response
		cmd.Println(tablePrint)

		return nil
	},
}

func init() {
	JobCmd.AddCommand(JobListCmd)
	initJobListCmdFlags()
}

func initJobListCmdFlags() {
}

func jobList() (interface{}, error) {
	response, _, err := RestClient.Services.Arc.GetList("jobs", url.Values{})
	if err != nil {
		return "", err
	}
	return response, nil
}
