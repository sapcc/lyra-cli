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
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v3/tokens"
	"github.com/sapcc/lyra-cli/print"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
)

type LyraAuthOps struct {
	IdentityEndpoint  string
	Region            string
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
	AuthenticationV3 = newAuthenticationV3
)

// authenticateCmd represents the authenticate command
var AuthenticateCmd = &cobra.Command{
	Use:   "authenticate",
	Short: locales.CmdShortDescription("authenticate"),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// DO NOT REMOVE. SHOULD OVERRIDE THE ROOT PersistentPreRunE
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// set authentication params
		lyraAuthOps := LyraAuthOps{
			IdentityEndpoint:  viper.GetString(ENV_VAR_AUTH_URL),
			Region:            viper.GetString(ENV_VAR_REGION),
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

		// authentication object
		authV3 := AuthenticationV3(lyraAuthOps)

		// authenticate
		response, err := authenticate(authV3)
		if err != nil {
			return err
		}

		// remove token expires at
		delete(response, TOKEN_EXPIRES_AT)

		// print the data out
		printer := print.Print{Data: response}
		bodyPrint := ""
		if viper.GetBool("json") {
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
	RootCmd.AddCommand(AuthenticateCmd)
	initAuthenticationCmdFlags()
}

func initAuthenticationCmdFlags() {
	// authenticate flags are global
}

func authenticate(authV3 Authentication) (map[string]string, error) {
	// check for user and password
	err := authV3.CheckAuthenticationParams()
	if err != nil {
		return map[string]string{}, err
	}

	// get the token result
	token, err := authV3.GetToken()
	if err != nil {
		return map[string]string{}, err
	}

	// arc endpoint
	arcEndpoint, err := authV3.GetServicePublicEndpoint("arc")
	if err != nil {
		return map[string]string{}, err
	}

	// automation endpoint
	automationEndpoint, err := authV3.GetServicePublicEndpoint("automation")
	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		ENV_VAR_AUTOMATION_ENDPOINT_NAME: automationEndpoint,
		ENV_VAR_ARC_ENDPOINT_NAME:        arcEndpoint,
		ENV_VAR_TOKEN_NAME:               token.ID,
		TOKEN_EXPIRES_AT:                 token.ExpiresAt.String(),
	}, nil
}

//
// Interface Authentication V3
//
type Authentication interface {
	CheckAuthenticationParams() error
	GetToken() (*tokens.Token, error)
	GetServicePublicEndpoint(serviceType string) (string, error)
}

type V3 struct {
	AuthOpts     LyraAuthOps
	client       *gophercloud.ServiceClient
	commonResult *tokens.CreateResult
}

func newAuthenticationV3(authOpts LyraAuthOps) Authentication {
	return &V3{AuthOpts: authOpts}
}

func checkAuthenticationParams(authOpts *LyraAuthOps) error {
	// check some params
	if len(authOpts.UserId) == 0 && len(authOpts.Username) == 0 {
		return fmt.Errorf("Flag %s or '%s not given.", FLAG_USER_ID, FLAG_USERNAME)
	}

	if len(authOpts.ProjectId) == 0 && len(authOpts.ProjectName) == 0 {
		return fmt.Errorf("Flag %s or %s not given.", FLAG_PROJECT_ID, FLAG_PROJECT_NAME)
	}

	if len(authOpts.IdentityEndpoint) == 0 {
		return fmt.Errorf("Flag %s not given.", FLAG_AUTH_URL)
	}

	// check password and prompt
	if len(authOpts.Password) == 0 {
		// ask the user for the password
		fmt.Print("Enter password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return err
		}
		authOpts.Password = string(pass)
	}

	return nil
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

func (a *V3) CheckAuthenticationParams() error {
	return checkAuthenticationParams(&a.AuthOpts)
}

func (a *V3) getClient() (*gophercloud.ServiceClient, error) {
	// get provider client struct
	provider, err := openstack.AuthenticatedClient(a.getAuthOptions())
	if err != nil {
		return nil, err
	}
	return openstack.NewIdentityV3(provider), nil
}

func (a *V3) GetToken() (*tokens.Token, error) {
	var err error
	if a.commonResult == nil {
		err = a.createTokenCommonResult()
		if err != nil {
			return nil, err
		}
	}

	token, err := a.commonResult.Extract()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *V3) createTokenCommonResult() error {
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
			return err
		}
	}

	// get common result
	result := tokens.Create(a.client, a.getAuthOptions(), &scope)
	// save common result
	a.commonResult = &result

	return nil
}

func (a *V3) GetServicePublicEndpoint(serviceType string) (string, error) {
	// get common result
	var err error
	if a.commonResult == nil {
		err = a.createTokenCommonResult()
		if err != nil {
			return "", err
		}
	}

	// get catalog
	var catalog *tokens.ServiceCatalog
	if a.commonResult != nil {
		catalog, err = a.commonResult.ExtractServiceCatalog()
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf("Authenticate: GetServicePublicEndpoint: Could not get token common result.")
	}

	// get entry from catalog
	serviceEntry, err := getServiceEntry(serviceType, catalog)
	if err != nil {
		return "", err
	}

	// get endpoint
	endpoint, err := getServicePublicEndpoint(a.AuthOpts.Region, serviceEntry)
	if err != nil {
		return "", err
	}

	return endpoint, nil
}

func getServicePublicEndpoint(region string, entry *tokens.CatalogEntry) (string, error) {
	if entry != nil && len(entry.Endpoints) > 0 {
		var endpoint string
		for _, ep := range entry.Endpoints {
			if region != "" {
				if ep.Interface == "public" && ep.Region == region {
					endpoint = ep.URL
					break
				}
			} else {
				if ep.Interface == "public" {
					endpoint = ep.URL
					break
				}
			}
		}
		return endpoint, nil
	} else {
		return "", fmt.Errorf("Authenticate: getServicePublicEndpoint: entry nil or no endpoints found for %+v.", entry)
	}
	return "", nil
}

func getServiceEntry(serviceType string, catalog *tokens.ServiceCatalog) (*tokens.CatalogEntry, error) {
	if catalog != nil && len(catalog.Entries) > 0 {
		serviceEntry := tokens.CatalogEntry{}
		for _, service := range catalog.Entries {
			if service.Type == serviceType {
				serviceEntry = service
				break
			}
		}
		return &serviceEntry, nil
	} else {
		return nil, fmt.Errorf("Authenticate: GetServicePublicEndpoint: catalog nil or emtpy.")
	}

	return nil, nil
}
