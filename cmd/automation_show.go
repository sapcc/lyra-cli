// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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
	"fmt"
	"net/url"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var automationId string

// showCmd represents the show command
var AutomationShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show a specific automation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

		// check required automation id
		if len(automationId) == 0 {
			log.Fatalf("Error: no automation id given.")
		}

		// setup rest client
		setupRestClient()

		show()
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationShowCmd)
	AutomationShowCmd.Flags().StringVarP(&automationId, "id", "i", "", "Id of the automation to show.")
}

func show() {
	response, err := RestClient.Get(path.Join("automations", automationId), url.Values{})
	if err != nil {
		log.Fatalf("%s\n", err.Error())
	}

	fmt.Println(response)
}
