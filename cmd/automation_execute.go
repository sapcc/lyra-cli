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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
)

// updateCmd represents the update command
var AutomationExecuteCmd = &cobra.Command{
	Use:   "execute",
	Short: "Runs an exsiting automation",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// setup automation run attributes
		err := setupAutomationRun()
		if err != nil {
			return err
		}

		// setup rest client
		err = setupRestClient()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		//run automation
		response, err := automationRun()
		if err != nil {
			return err
		}

		//     test := `{
		//   "id": "85",
		//   "log": null,
		//   "created_at": "2016-05-25T11:32:51.723Z",
		//   "updated_at": "2016-05-25T11:32:51.723Z",
		//   "repository_revision": null,
		//   "state": "preparing",
		//   "jobs": null,
		//   "owner": "u-519166a05",
		//   "automation_id": "6",
		//   "automation_name": "Chef_test",
		//   "selector": "@identity=\"0128e993-c709-4ce1-bccf-e06eb10900a0\""
		// }`

		// err := automationRunWait(test)
		// if err != nil {
		//   return err
		// }

		// Print response
		cmd.Println(response)

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationExecuteCmd)
	//flags
	run = Run{}
	AutomationExecuteCmd.Flags().StringVarP(&run.AutomationId, "automation-id", "", "", locales.AttributeDescription("automation-id"))
	AutomationExecuteCmd.Flags().StringVarP(&run.Selector, "selector", "", "", locales.AttributeDescription("automation-selector"))
}

func setupAutomationRun() error {
	// check required automation id
	if len(run.AutomationId) == 0 {
		return errors.New(locales.ErrorMessages("automation-id-missing"))
	}
	// check selector
	if len(run.Selector) == 0 {
		return errors.New(locales.ErrorMessages("automation-selector-missing"))
	}

	return nil
}

func automationRun() (string, error) {
	// convert to Json
	body, err := json.Marshal(run)
	if err != nil {
		return "", err
	}
	// send request
	response, _, err := RestClient.Services.Automation.Post("runs", url.Values{}, string(body))
	if err != nil {
		return "", err
	}

	return response, nil
}

type ExecRun struct {
	Id    string   `json:"id"`
	State string   `json:"state"`
	Jobs  []string `json:"jobs"`
}

func automationRunWait(runData string) error {
	// convert data to struct
	execRun := ExecRun{}
	err := helpers.JSONStringToStructure(runData, &execRun)
	if err != nil {
		return err
	}

	fmt.Printf("Run automation is created with id %s\n", execRun.Id)
	fmt.Printf("Automation state %s\n", execRun.State)

	// update data
	tickChan := time.NewTicker(time.Second * 5)
	runningJobs := []string{}
	// jobState := map[string]string{}
	for {
		select {
		case <-tickChan.C:

			var err error
			var runUpdate *ExecRun

			switch execRun.State {
			case RunPreparing:
				// get new run update
				runUpdate, err = getAutomationRun(execRun.Id)
			case RunExecuting:
				// get new run update
				runUpdate, err = getAutomationRun(execRun.Id)
			}

			// exit if error occurs
			if err != nil {
				fmt.Println(err)
				return nil
			}

			// run state changed
			if len(runUpdate.State) > 0 && runUpdate.State != execRun.State {
				// update state
				execRun.State = runUpdate.State
				// print out for run state
				fmt.Printf("Automation state %s\n", execRun.State)
				// exit if run failed
				if execRun.State == RunFailed {
					tickChan.Stop()
					return nil
				}
			}

			// add new jobs
			if len(runUpdate.Jobs) > 0 && len(runUpdate.Jobs) != len(execRun.Jobs) {
				execRun.Jobs = runUpdate.Jobs
				fmt.Println("Schedule jobs:")
				for _, v := range runUpdate.Jobs {
					runningJobs = append(execRun.Jobs, v)
					fmt.Printf("Job id %s\n", v)
				}
			}

			// update running jobs
			if len(runningJobs) > 0 {
				for _, v := range runningJobs {
					// do update to the map
					fmt.Println(v)
				}
			} else {
				// check all entries in the job are cmoplete or failed and exit
			}
		}
	}
	return nil
}

const (
	RunPreparing = "preparing"
	RunExecuting = "executing"
	RunFailed    = "failed"
	RunCompleted = "completed"
)

func getAutomationRun(id string) (*ExecRun, error) {
	// get run
	data, err := getRunUpdate(id)
	if err != nil {
		return nil, err
	}
	// convert data
	updateRun := ExecRun{}
	err = helpers.JSONStringToStructure(data, &updateRun)
	if err != nil {
		return nil, err
	}

	return &updateRun, nil
}

func getRunUpdate(id string) (string, error) {
	// get new data
	data, err := runShow(id)
	if err != nil {
		return "", err
	}

	return data, nil
}

type JobState byte

const (
	_ JobState = iota
	JobQueued
	JobExecuting
	JobFailed
	JobComplete
)

var jobStateStringMap = map[JobState]string{JobQueued: "queued", JobExecuting: "executing", JobFailed: "failed", JobComplete: "complete"}

func jobRunState(id string) {

}
