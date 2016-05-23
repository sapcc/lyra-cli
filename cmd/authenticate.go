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
	"github.com/sapcc/lyra-cli/print"

	"github.com/spf13/cobra"
)

type LyraAuthOps struct {
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
	lyraAuthOps      = LyraAuthOps{}
	AuthenticationV3 = newAuthenticationV3
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

		// authentication object
		authV3 := AuthenticationV3(lyraAuthOps)

		// authenticate
		response, err := authenticate(authV3)
		if err != nil {
			return err
		}

		// print the data out
		printer := print.Print{Data: response}
		bodyPrint := ""
		if JsonOutput {
			bodyPrint, err = printer.JSON()
			if err != nil {
				return err
			}
		} else {
			bodyPrint = fmt.Sprintf("export %s=%s\nexport %s=%s\nexport %s=%s", ENV_VAR_AUTOMATION_ENDPOINT_NAME, response[ENV_VAR_AUTOMATION_ENDPOINT_NAME], ENV_VAR_ARC_ENDPOINT_NAME, response[ENV_VAR_ARC_ENDPOINT_NAME], ENV_VAR_TOKEN_NAME, response[ENV_VAR_TOKEN_NAME])
		}

		// Print response
		cmd.Println(bodyPrint)

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
	authenticateCmd.Flags().StringVar(&lyraAuthOps.UserId, "user-id", "", fmt.Sprint("Id of the user that wants to log in. (default ", userid_default_env_name, ")"))
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

func authenticate(authV3 Authentication) (map[string]string, error) {
	// get automation service id from the catalog
	automationServiceId, err := authV3.GetServiceId("automation")
	if err != nil {
		return map[string]string{}, err
	}
	// get automation service endpoints from catalog
	automationPublicEndpoint, err := authV3.GetServicePublicEndpoint(automationServiceId)
	if err != nil {
		return map[string]string{}, err
	}

	// get automation service id from the catalog
	arcServiceId, err := authV3.GetServiceId("arc")
	if err != nil {
		return map[string]string{}, err
	}
	// get automation service endpoints from catalog
	arcPublicEndpoint, err := authV3.GetServicePublicEndpoint(arcServiceId)
	if err != nil {
		return map[string]string{}, err
	}

	// get the token
	token, err := authV3.GetToken()
	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		ENV_VAR_AUTOMATION_ENDPOINT_NAME: automationPublicEndpoint,
		ENV_VAR_ARC_ENDPOINT_NAME:        arcPublicEndpoint,
		ENV_VAR_TOKEN_NAME:               token,
	}, nil
}

//
// Interface Authentication V3
//
type Authentication interface {
	GetToken() (string, error)
	GetServicePublicEndpoint(serviceId string) (string, error)
	GetServiceId(serviceType string) (string, error)
}

type V3 struct {
	AuthOpts LyraAuthOps
	client   *gophercloud.ServiceClient
}

func newAuthenticationV3(authOpts LyraAuthOps) Authentication {
	return &V3{AuthOpts: authOpts}
}

func (a *V3) getAuthOptions() gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: a.AuthOpts.IdentityEndpoint,
		Username:         a.AuthOpts.Username,
		UserID:           a.AuthOpts.UserId,
		Password:         a.AuthOpts.Password,
		DomainName:       a.AuthOpts.UserDomainName,
		DomainID:         a.AuthOpts.UserDomainId,
	}
}

func (a *V3) getClient() (*gophercloud.ServiceClient, error) {
	// get provider client struct
	provider, err := openstack.AuthenticatedClient(a.getAuthOptions())
	if err != nil {
		return nil, err
	}

	// Creates a ServiceClient that may be used to access the v3 identity service
	return openstack.NewIdentityV3(provider), nil
}

func (a *V3) GetToken() (string, error) {
	scope := tokens.Scope{
		ProjectName: a.AuthOpts.ProjectName,
		ProjectID:   a.AuthOpts.ProjectId,
		DomainName:  a.AuthOpts.ProjectDomainName,
		DomainID:    a.AuthOpts.ProjectDomainId,
	}

	// init the v3 client
	var err error
	if a.client == nil {
		a.client, err = a.getClient()
		if err != nil {
			return "", err
		}
	}

	// get the token
	token, err := tokens.Create(a.client, a.getAuthOptions(), &scope).Extract()
	if err != nil {
		return "", err
	}

	return token.ID, nil
}

func (a *V3) GetServicePublicEndpoint(serviceId string) (string, error) {
	// init the v3 client
	var err error
	if a.client == nil {
		a.client, err = a.getClient()
		if err != nil {
			return "", err
		}
	}

	// get the endpoints
	publicEndpoint := ""
	endpointsOpts := endpoints.ListOpts{ServiceID: serviceId, Page: 1, PerPage: 1}
	endpointsPager := endpoints.List(a.client, endpointsOpts)

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
				publicEndpoint = e.URL
				break
			}
		}
		if len(publicEndpoint) == 0 {
			return false, fmt.Errorf("No service automation public url found in catalog.")
		}
		return true, nil
	})
	if err != nil {
		return "", err
	}

	return publicEndpoint, nil
}

func (a *V3) GetServiceId(serviceType string) (string, error) {
	// init the v3 client
	var err error
	if a.client == nil {
		a.client, err = a.getClient()
		if err != nil {
			return "", err
		}
	}

	// get the service
	serviceId := ""
	opts := services.ListOpts{ServiceType: serviceType, Page: 1, PerPage: 1}
	servicesPager := services.List(a.client, opts)

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
		serviceId = servicesList[0].ID
		return true, nil
	})
	if err != nil {
		return "", err
	}

	return serviceId, nil
}
