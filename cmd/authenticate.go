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
	// "bufio"
	"errors"
	"fmt"
	"os"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v3/endpoints"
	"github.com/rackspace/gophercloud/openstack/identity/v3/services"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"github.com/rackspace/gophercloud/pagination"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/spf13/cobra"
)

var ENV_VAR_USERNAME = "USERNAME"
var ENV_VAR_USERID = "USERID"
var ENV_VAR_PASSWORD = "PASSWORD"
var identityEndpoint, username, userId, password, projectName, projectId, userDomainName, userDomainId, projectDomainName, projectDomainId string

// authenticateCmd represents the authenticate command
var authenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Get an authentication token project based.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// setup
		err := setupAuthentication()
		if err != nil {
			return err
		}
		// authenticate
		response, err := authenticate()
		if err != nil {
			return err
		}
		// print response
		cmd.Print(response)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(authenticateCmd)

	username_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_USERNAME)
	userid_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_USERID)
	password_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_PASSWORD)

	authenticateCmd.Flags().StringVarP(&identityEndpoint, "identity-endpoint", "e", "", "Endpoint entities represent URL endpoints for OpenStack web services.")
	authenticateCmd.Flags().StringVar(&username, "username", "", fmt.Sprint("Name of the user that wants to log in. (default ", username_default_env_name, ")"))
	authenticateCmd.Flags().StringVar(&userId, "userId", "", fmt.Sprint("Name of the user that wants to log in. (default ", userid_default_env_name, ")"))
	authenticateCmd.Flags().StringVarP(&password, "password", "p", "", fmt.Sprint("Password of the user that wants to log in. If not given the environment variable ", password_default_env_name, " will be checkt. If no environment variable found then will promtp from terminal."))

	authenticateCmd.Flags().StringVar(&projectName, "project-name", "", "Name of the project.")
	authenticateCmd.Flags().StringVar(&projectId, "project-id", "", "Id of the project.")

	authenticateCmd.Flags().StringVar(&userDomainName, "user-domain-name", "", "Name of the domain where the user is created.")
	authenticateCmd.Flags().StringVar(&userDomainId, "user-domain-id", "", "Id of the domain where the user is created.")

	authenticateCmd.Flags().StringVar(&projectDomainName, "project-domain-name", "", "Name of the domain where the project is created. If no project domain name is given, then the token will be scoped in the user domain.")
	authenticateCmd.Flags().StringVar(&projectDomainId, "project-domain-id", "", "Id of the domain where the project is created. If no project domain id is given, then the token will be scoped in the user domain.")
}

func setupAuthentication() error {
	// setup flags with environment variablen
	if len(username) == 0 {
		username = os.Getenv(ENV_VAR_USERNAME)
	}
	if len(userId) == 0 {
		userId = os.Getenv(ENV_VAR_USERID)
	}
	// check we have user name or id
	if len(username) == 0 && len(userId) == 0 {
		return errors.New("Username or userid not given.")
	}
	// check password
	if len(password) == 0 {
		if len(os.Getenv(ENV_VAR_PASSWORD)) == 0 {
			// ask the user for the password
			fmt.Print("Enter password: ")
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				log.Fatalf(err.Error())
			}
			fmt.Print("\n")
			password = string(bytePassword)

		} else {
			password = os.Getenv(ENV_VAR_PASSWORD)
		}
	}

	return nil
}

func authenticate() (string, error) {
	// add auth options
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: identityEndpoint,
		Username:         username,
		UserID:           userId,
		Password:         password,
		DomainName:       userDomainName,
		DomainID:         userDomainId,
	}

	// get provider client struct
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return "", err
	}

	// Creates a ServiceClient that may be used to access the v3 identity service
	client := openstack.NewIdentityV3(provider)

	// get the automation entry id from the catalog
	opts := services.ListOpts{ServiceType: "automation", Page: 1, PerPage: 1}
	servicesPager := services.List(client, opts)
	automationServiceId := ""
	err = servicesPager.EachPage(func(page pagination.Page) (bool, error) {
		servicesList, err := services.ExtractServices(page)
		if err != nil {
			return false, err
		}
		if len(servicesList) != 1 {
			return false, fmt.Errorf("No service automation found in catalog.")
		}
		if len(servicesList[0].ID) == 0 {
			return false, fmt.Errorf("No service automation id found in catalog.")
		}
		// save the automation id
		automationServiceId = servicesList[0].ID
		return true, nil
	})
	if err != nil {
		return "", err
	}

	// get automation service endpoints
	endpointsOpts := endpoints.ListOpts{ServiceID: automationServiceId, Page: 1, PerPage: 1}
	endpointsPager := endpoints.List(client, endpointsOpts)
	automationPublicEndpoint := ""
	err = endpointsPager.EachPage(func(page pagination.Page) (bool, error) {
		endpointList, err := endpoints.ExtractEndpoints(page)
		if err != nil {
			return false, err
		}
		if len(endpointList) == 0 {
			return false, fmt.Errorf("No endpoints for service automation found in catalog.")
		}
		for _, e := range endpointList {
			if e.Availability == "public" {
				automationPublicEndpoint = e.URL
				break
			}
		}
		if len(automationPublicEndpoint) == 0 {
			return false, fmt.Errorf("No service automation public url found in catalog.")
		}
		return true, nil
	})
	if err != nil {
		return "", err
	}

	// set the project scope
	scope := tokens.Scope{
		ProjectName: projectName,
		ProjectID:   projectId,
		DomainName:  projectDomainName,
		DomainID:    projectDomainId,
	}

	// get the token
	token, err := tokens.Create(client, authOpts, &scope).Extract()
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("export %s=%s\nexport %s=%s\n", ENV_VAR_AUTOMATION_ENDPOINT_NAME, automationPublicEndpoint, ENV_VAR_TOKEN_NAME, token.ID), nil
}
