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
	"os"

	"github.com/howeyc/gopass"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v3/endpoints"
	"github.com/rackspace/gophercloud/openstack/identity/v3/services"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"github.com/rackspace/gophercloud/pagination"

	"github.com/spf13/cobra"
)

type Auth struct {
	IdentityEndpoint  string
	Username          string
	UserId            string
	Password          string
	ProjectName       string
	ProjectId         string
	UserDomainName    string
	UserDomainId      string
	ProjectDomainName string
	ProjectDomainId   string
}

var (
	ENV_VAR_USERNAME = "USERNAME"
	ENV_VAR_USERID   = "USERID"
	ENV_VAR_PASSWORD = "PASSWORD"
	lyraAuthOps      = Auth{}
)

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
		cmd.Println(response)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(authenticateCmd)

	username_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_USERNAME)
	userid_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_USERID)
	password_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_PASSWORD)

	authenticateCmd.Flags().StringVar(&lyraAuthOps.IdentityEndpoint, "identity-endpoint", "", "Endpoint entities represent URL endpoints for OpenStack web services.")
	authenticateCmd.Flags().StringVar(&lyraAuthOps.Username, "username", "", fmt.Sprint("Name of the user that wants to log in. (default ", username_default_env_name, ")"))
	authenticateCmd.Flags().StringVar(&lyraAuthOps.UserId, "userId", "", fmt.Sprint("Id of the user that wants to log in. (default ", userid_default_env_name, ")"))
	authenticateCmd.Flags().StringVar(&lyraAuthOps.Password, "password", "", fmt.Sprint("Password of the user that wants to log in. If not given the environment variable ", password_default_env_name, " will be checkt. If no environment variable found then will promtp from terminal."))

	authenticateCmd.Flags().StringVar(&lyraAuthOps.ProjectName, "project-name", "", "Name of the project.")
	authenticateCmd.Flags().StringVar(&lyraAuthOps.ProjectId, "project-id", "", "Id of the project.")

	authenticateCmd.Flags().StringVar(&lyraAuthOps.UserDomainName, "user-domain-name", "", "Name of the domain where the user is created.")
	authenticateCmd.Flags().StringVar(&lyraAuthOps.UserDomainId, "user-domain-id", "", "Id of the domain where the user is created.")

	authenticateCmd.Flags().StringVar(&lyraAuthOps.ProjectDomainName, "project-domain-name", "", "Name of the domain where the project is created. If no project domain name is given, then the token will be scoped in the user domain.")
	authenticateCmd.Flags().StringVar(&lyraAuthOps.ProjectDomainId, "project-domain-id", "", "Id of the domain where the project is created. If no project domain id is given, then the token will be scoped in the user domain.")
}

func setupAuthentication() error {
	// setup flags with environment variablen
	if len(lyraAuthOps.Username) == 0 {
		lyraAuthOps.Username = os.Getenv(ENV_VAR_USERNAME)
	}
	if len(lyraAuthOps.UserId) == 0 {
		lyraAuthOps.UserId = os.Getenv(ENV_VAR_USERID)
	}
	// check we have user name or id
	if len(lyraAuthOps.Username) == 0 && len(lyraAuthOps.UserId) == 0 {
		return errors.New("Username or userid not given.")
	}
	// check password
	if len(lyraAuthOps.Password) == 0 {
		if len(os.Getenv(ENV_VAR_PASSWORD)) == 0 {
			// ask the user for the password
			fmt.Print("Enter password: ")
			pass, err := gopass.GetPasswd()
			if err != nil {
				return err
			}
			lyraAuthOps.Password = string(pass)

		} else {
			lyraAuthOps.Password = os.Getenv(ENV_VAR_PASSWORD)
		}
	}

	return nil
}

func authenticate() (string, error) {
	// add auth options
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: lyraAuthOps.IdentityEndpoint,
		Username:         lyraAuthOps.Username,
		UserID:           lyraAuthOps.UserId,
		Password:         lyraAuthOps.Password,
		DomainName:       lyraAuthOps.UserDomainName,
		DomainID:         lyraAuthOps.UserDomainId,
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
		ProjectName: lyraAuthOps.ProjectName,
		ProjectID:   lyraAuthOps.ProjectId,
		DomainName:  lyraAuthOps.ProjectDomainName,
		DomainID:    lyraAuthOps.ProjectDomainId,
	}

	// get the token
	token, err := tokens.Create(client, authOpts, &scope).Extract()
	if err != nil {
		return "", nil
	}

	return fmt.Sprintf("export %s=%s\nexport %s=%s", ENV_VAR_AUTOMATION_ENDPOINT_NAME, automationPublicEndpoint, ENV_VAR_TOKEN_NAME, token.ID), nil
}
