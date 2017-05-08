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
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
)

var (
	ExecuteAuthOps = auth.AuthOptions{}
	ExecuteAuthV3  = auth.AuthenticationV3(ExecuteAuthOps)
)

// updateCmd represents the update command
var AutomationExecuteCmd = &cobra.Command{
	Use:   "execute",
	Short: locales.CmdShortDescription("automation-execute"),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// DO NOT REMOVE. SHOULD OVERRIDE THE ROOT PersistentPreRunE
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// setup automation run attributes
		err := setupAutomationRun()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		response := ""
		if viper.GetBool("watch") {
			// keep the auth options for reauthentication
			ExecuteAuthOps = auth.AuthOptions{
				IdentityEndpoint:  viper.GetString(ENV_VAR_AUTH_URL),
				Username:          viper.GetString(ENV_VAR_USERNAME),
				UserId:            viper.GetString(ENV_VAR_USER_ID),
				Password:          viper.GetString(ENV_VAR_PASSWORD),
				ProjectName:       viper.GetString(ENV_VAR_PROJECT_NAME),
				ProjectId:         viper.GetString(ENV_VAR_PROJECT_ID),
				UserDomainName:    viper.GetString(ENV_VAR_USER_DOMAIN_NAME),
				UserDomainId:      viper.GetString(ENV_VAR_USER_DOMAIN_ID),
				ProjectDomainName: viper.GetString(ENV_VAR_PROJECT_DOMAIN_NAME),
				ProjectDomainId:   viper.GetString(ENV_VAR_PROJECT_DOMAIN_ID),
			}
			ExecuteAuthV3 = auth.AuthenticationV3(ExecuteAuthOps)
			// force reauthenticate with password and keep values
			err := setupRestClient(&ExecuteAuthV3, true)
			if err != nil {
				return err
			}

			response, err = automationRunWait()
			if err != nil {
				return err
			}

		} else {
			// get std authentication
			err := setupRestClient(nil, false)
			if err != nil {
				return err
			}

			//run automation
			response, err = automationRun()
			if err != nil {
				return err
			}
		}

		// convert data to struct
		var dataStruct map[string]interface{}
		err := helpers.JSONStringToStructure(response, &dataStruct)
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

		// Print response
		cmd.Println(bodyPrint)

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationExecuteCmd)
	initAutomationExecuteCmdFlags()
}

func initAutomationExecuteCmdFlags() {
	//flags
	AutomationExecuteCmd.Flags().StringP(FLAG_AUTOMATION_ID, "", "", locales.AttributeDescription("automation-id"))
	viper.BindPFlag(FLAG_AUTOMATION_ID, AutomationExecuteCmd.Flags().Lookup(FLAG_AUTOMATION_ID))
	AutomationExecuteCmd.Flags().StringP(FLAG_SELECTOR, "", "", locales.AttributeDescription("selector"))
	viper.BindPFlag(FLAG_SELECTOR, AutomationExecuteCmd.Flags().Lookup(FLAG_SELECTOR))
	AutomationExecuteCmd.Flags().BoolP("watch", "", false, locales.AttributeDescription("watch"))
	viper.BindPFlag("watch", AutomationExecuteCmd.Flags().Lookup("watch"))
}

func setupAutomationRun() error {
	// check required automation id
	if len(viper.GetString(FLAG_AUTOMATION_ID)) == 0 {
		return errors.New(locales.ErrorMessages("automation-id-missing"))
	}
	// check selector
	if len(viper.GetString(FLAG_AUTOMATION_ID)) == 0 {
		return errors.New(locales.ErrorMessages("automation-selector-missing"))
	}

	return nil
}

func automationRun() (string, error) {
	run := Run{AutomationId: viper.GetString(FLAG_AUTOMATION_ID), Selector: viper.GetString(FLAG_SELECTOR)}
	// convert to Json
	body, err := json.Marshal(run)
	if err != nil {
		return "", err
	}
	// send request
	automationService := RestClient.Services["automation"]
	response, _, err := automationService.Post("runs", url.Values{}, http.Header{}, string(body))
	if err != nil {
		return "", err
	}

	return response, nil
}

type AutomationRun struct {
	Id    string   `json:"id"`
	State string   `json:"state"`
	Jobs  []string `json:"jobs"`
}

func automationRunWait() (string, error) {
	//run automation
	runData, err := automationRun()
	if err != nil {
		return "", err
	}

	// convert data to struct
	automationRun := AutomationRun{}
	err = helpers.JSONStringToStructure(runData, &automationRun)
	if err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "Run automation is created with id %s\n", automationRun.Id)
	fmt.Fprintf(os.Stderr, "Automation state %s\n", automationRun.State)

	// update data
	tickChan := time.NewTicker(time.Second * 5)
	defer tickChan.Stop()

	runningJobs := []string{}
	jobsState := map[string]string{}
	for {
		select {
		case <-tickChan.C:

			// check if the token still valid
			isExpired, err := tokenExpired()
			if err != nil {
				return "", err
			}
			if isExpired {
				fmt.Fprintln(os.Stderr, "WARNING: token expired.")
				// reauthenticate
				err = setupRestClient(&ExecuteAuthV3, true)
				if err != nil {
					return "", err
				}
			}

			var runUpdate AutomationRun
			switch automationRun.State {
			case RunPreparing:
				// get new run update
				runUpdate, err = getAutomationRun(automationRun.Id)
			case RunExecuting:
				// get new run update
				runUpdate, err = getAutomationRun(automationRun.Id)
			}

			// exit if error occurs
			if err != nil {
				return "", err
			}

			// run state changed
			if len(runUpdate.State) > 0 && runUpdate.State != automationRun.State {
				// update state
				automationRun.State = runUpdate.State

				if automationRun.State != RunFailed && automationRun.State != RunCompleted {
					// print out for run state
					fmt.Fprintf(os.Stderr, "Automation %s\n", automationRun.State)
				}
			}

			// add new jobs
			if len(runUpdate.Jobs) > 0 && len(runUpdate.Jobs) != len(automationRun.Jobs) {
				// update jobs
				automationRun.Jobs = runUpdate.Jobs
				fmt.Fprintf(os.Stderr, "Scheduled %d jobs:\n", len(runUpdate.Jobs))
				for _, v := range runUpdate.Jobs {
					// save them to keep track
					runningJobs = append(automationRun.Jobs, v)
					jobsState[v] = ""
					fmt.Fprintf(os.Stderr, "%s\n", v)
				}
			}

			// update running jobs
			if len(runningJobs) > 0 {
				stillrunningJobs := []string{}
				for _, v := range runningJobs {
					// get job update
					stateStr, err := getJobStateUpdate(v)
					if err != nil {
						return "", err
					}
					if stateStr != jobsState[v] {
						fmt.Fprintf(os.Stderr, "Job %s is %s\n", v, stateStr)
						jobsState[v] = stateStr
					}
					// if state is failed or complete then remove entry
					if stateStr != JobFailed && stateStr != JobComplete {
						stillrunningJobs = append(stillrunningJobs, v)
					}
				}
				runningJobs = stillrunningJobs
			} else {
				if automationRun.State == RunCompleted {
					fmt.Fprintf(os.Stderr, "Automation run %s %s\n.", automationRun.Id, automationRun.State)
				}
				if automationRun.State == RunFailed {
					fmt.Fprintf(os.Stderr, "Automation run %s %s. %d of %d jobs failed\n", automationRun.Id, automationRun.State, jobsFailed(jobsState), len(automationRun.Jobs))
				}

				resultRun, err := runShow(automationRun.Id)
				if err != nil {
					return "", err
				}
				return resultRun, nil
			}
		}
	}
}

func jobsFailed(jobsState map[string]string) int {
	jobs := 0

	for _, state := range jobsState {
		if state == JobFailed {
			jobs++
		}
	}

	return jobs
}

func tokenExpired() (bool, error) {
	layout := "2006-01-02 15:04:05.999 -0700 MST"
	expiresAt, err := time.Parse(layout, viper.GetString(TOKEN_EXPIRES_AT))
	if err != nil {
		return false, err
	}

	now := time.Now().In(expiresAt.Location())
	delta := now.Sub(expiresAt)

	if delta.Seconds() <= -60 {
		return false, nil
	} else {
		return true, nil
	}
}

const (
	RunPreparing = "preparing"
	RunExecuting = "executing"
	RunFailed    = "failed"
	RunCompleted = "completed"
)

func getAutomationRun(id string) (AutomationRun, error) {
	// get run
	data, err := getRunUpdate(id)
	if err != nil {
		return AutomationRun{}, err
	}
	// convert data
	updateRun := AutomationRun{}
	err = helpers.JSONStringToStructure(data, &updateRun)
	if err != nil {
		return AutomationRun{}, err
	}

	return updateRun, nil
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
	JobQueued    = "queued"
	JobExecuting = "executing"
	JobFailed    = "failed"
	JobComplete  = "complete"
)

func getJobStateUpdate(id string) (string, error) {
	// get job update
	job, err := jobShow(id)
	// convert data to struct
	var jobStruct map[string]interface{}
	err = helpers.JSONStringToStructure(job, &jobStruct)
	if err != nil {
		return "", err
	}

	// update state
	state := jobStruct["status"]
	stateStr, ok := state.(string)
	if !ok {
		return "", fmt.Errorf("Error converting job state to string")
	}
	return stateStr, nil
}
