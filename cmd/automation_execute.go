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
	"math/rand"
	"net/http"
	"net/url"
	"time"

	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/print"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		return setupAutomationRun()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var response string
		if viper.GetBool("watch") {
			// keep the auth options for reauthentication
			ExecuteAuthOps = auth.AuthOptions{
				IdentityEndpoint:            viper.GetString(ENV_VAR_AUTH_URL),
				Username:                    viper.GetString(ENV_VAR_USERNAME),
				UserId:                      viper.GetString(ENV_VAR_USER_ID),
				Password:                    viper.GetString(ENV_VAR_PASSWORD),
				ProjectName:                 viper.GetString(ENV_VAR_PROJECT_NAME),
				ProjectId:                   viper.GetString(ENV_VAR_PROJECT_ID),
				UserDomainName:              viper.GetString(ENV_VAR_USER_DOMAIN_NAME),
				UserDomainId:                viper.GetString(ENV_VAR_USER_DOMAIN_ID),
				ProjectDomainName:           viper.GetString(ENV_VAR_PROJECT_DOMAIN_NAME),
				ProjectDomainId:             viper.GetString(ENV_VAR_PROJECT_DOMAIN_ID),
				ApplicationCredentialID:     viper.GetString(ENV_VAR_APPLICATION_CREDENTIAL_ID),
				ApplicationCredentialName:   viper.GetString(ENV_VAR_APPLICATION_CREDENTIAL_NAME),
				ApplicationCredentialSecret: viper.GetString(ENV_VAR_APPLICATION_CREDENTIAL_SECRET),
			}
			ExecuteAuthV3 = auth.AuthenticationV3(ExecuteAuthOps)
			// force reauthenticate with password and keep values
			err := setupRestClient(cmd, &ExecuteAuthV3, true)
			if err != nil {
				return err
			}

			response, err = automationRunWait(cmd)
			if err != nil {
				return err
			}

		} else {
			// get std authentication
			err := setupRestClient(cmd, nil, false)
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
		var bodyPrint string
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
		fmt.Println(bodyPrint)

		return nil
	},
}

func init() {
	AutomationCmd.AddCommand(AutomationExecuteCmd)
	initAutomationExecuteCmdFlags()
	rand.New(rand.NewSource((time.Now().UnixNano())))
}

func initAutomationExecuteCmdFlags() {
	//flags
	AutomationExecuteCmd.Flags().StringP(FLAG_AUTOMATION_ID, "", "", locales.AttributeDescription("automation-id"))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag(FLAG_AUTOMATION_ID, AutomationExecuteCmd.Flags().Lookup(FLAG_AUTOMATION_ID)), "BindPFlag:")
	AutomationExecuteCmd.Flags().StringP(FLAG_SELECTOR, "", "", locales.AttributeDescription("selector"))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag(FLAG_SELECTOR, AutomationExecuteCmd.Flags().Lookup(FLAG_SELECTOR)), "BindPFlag:")
	AutomationExecuteCmd.Flags().BoolP("watch", "", false, locales.AttributeDescription("watch"))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("watch", AutomationExecuteCmd.Flags().Lookup("watch")), "BindPFlag:")
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

func automationRunWait(cmd *cobra.Command) (string, error) {
	//run automation
	var runData string
	var err error
	err = retry(5, 30*time.Second, func() error {
		runData, err = automationRun()
		return err
	})
	// retry error
	if err != nil {
		return "", err
	}

	// convert data to struct
	automationRun := AutomationRun{}
	err = helpers.JSONStringToStructure(runData, &automationRun)
	if err != nil {
		return "", err
	}

	cmd.Printf("Automation run is created with id %s\n", automationRun.Id)
	cmd.Printf("Automation run state %s\n", automationRun.State)

	// update data
	tickChan := time.NewTicker(time.Second * 5)
	defer tickChan.Stop()

	runningJobs := []string{}
	jobsState := map[string]string{}
	for ; true; <-tickChan.C {
		// check if the token still valid
		isExpired, err := tokenExpired()
		if err != nil {
			return "", err
		}
		if isExpired {
			cmd.Println("WARNING: token expired.")
			// reauthenticate
			err = retry(5, 30*time.Second, func() error {
				return setupRestClient(cmd, &ExecuteAuthV3, true)
			})
			// retry error
			if err != nil {
				return "", err
			}
		}

		// get new run update
		var runUpdate AutomationRun
		err = retry(5, 30*time.Second, func() error {
			runUpdate, err = getAutomationRun(automationRun.Id)
			return err
		})
		// error from retry
		if err != nil {
			return "", err
		}

		switch automationRun.State {
		case RunPreparing:
			// add new jobs
			if len(runUpdate.Jobs) != len(automationRun.Jobs) {
				// update jobs
				automationRun.Jobs = runUpdate.Jobs
				cmd.Printf("Scheduled %d jobs:\n", len(runUpdate.Jobs))
				for _, v := range runUpdate.Jobs {
					// save them to keep track
					runningJobs = append(automationRun.Jobs, v)
					jobsState[v] = ""
					cmd.Printf("%s\n", v)
				}
			}
		case RunExecuting:
			stillrunningJobs := []string{}
			for _, v := range runningJobs {
				// get job update
				var stateStr string
				err = retry(5, 30*time.Second, func() error {
					stateStr, err = getJobStateUpdate(v)
					return err
				})
				// error from retry
				if err != nil {
					return "", err
				}

				if stateStr != jobsState[v] {
					cmd.Printf("Job %s is %s\n", v, stateStr)
					jobsState[v] = stateStr
				}
				// if state is failed or complete then remove entry
				if stateStr != JobFailed && stateStr != JobComplete {
					stillrunningJobs = append(stillrunningJobs, v)
				}
			}
			runningJobs = stillrunningJobs
		}

		// did the run state change?
		if runUpdate.State != automationRun.State {
			// update state
			automationRun.State = runUpdate.State
			switch automationRun.State {
			case RunFailed:
				cmd.Printf("Automation run %s %s. %d of %d jobs failed\n", automationRun.Id, automationRun.State, jobsFailed(jobsState), len(automationRun.Jobs))
				// get the last state of the run
				runStr, err := runShow(automationRun.Id)
				if err != nil {
					return "", err
				}
				// force return error
				return runStr, errors.New(locales.ErrorMessages("automation-run-failed"))
			case RunCompleted:
				cmd.Printf("Automation run %s %s. %d jobs succeeded\n", automationRun.Id, automationRun.State, len(automationRun.Jobs))
				// return the last state of the run
				return runShow(automationRun.Id)
			}
			cmd.Printf("Automation run state %s\n", automationRun.State)
		}
	}
	return "", nil
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
	if err != nil {
		return "", err
	}
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
		return "", fmt.Errorf("error converting job state to string")
	}
	return stateStr, nil
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}
