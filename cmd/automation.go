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
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/restclient"
)

var Token, AutomationUrl string
var RestClient *restclient.Client

// automationCmd represents the automation command
var AutomationCmd = &cobra.Command{
	Use:   "automation",
	Short: "Automation service.",
	Long:  `A longer description for automation.`,
}

func init() {
	RootCmd.AddCommand(AutomationCmd)

	token_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_TOKEN_NAME)
	automation_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_AUTOMATION_ENDPOINT_NAME)
	AutomationCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", fmt.Sprint("Authentication token. To create a token run the authenticate command. (default ", token_default_env_name, ")"))
	AutomationCmd.PersistentFlags().StringVarP(&AutomationUrl, "automation-endpoint", "a", "", fmt.Sprint("Automation endpoint. To get the automation endpoint run the authenticate command. (default ", automation_default_env_name, ")"))
}

func setupRestClient() error {
	// setup flags with environment variablen
	if len(Token) == 0 {
		if len(os.Getenv(ENV_VAR_TOKEN_NAME)) == 0 {
			return errors.New("Token not given. To create a token you can use the authenticate command.")
		} else {
			Token = os.Getenv(ENV_VAR_TOKEN_NAME)
		}
	}
	if len(AutomationUrl) == 0 {
		if len(os.Getenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)) == 0 {
			return errors.New("Endpoint not given. To get the automation endpoint run the authenticate command.")
		} else {
			AutomationUrl = os.Getenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)
		}
	}

	// add to the endpoint the api version
	u, err := url.Parse(AutomationUrl)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, "/api/v1/")

	// init rest client
	RestClient = restclient.NewClient(u.String(), Token)

	return nil
}
