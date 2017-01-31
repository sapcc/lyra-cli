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
	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/print"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
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
		options := auth.AuthOptions{
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

		// authentication object
		authV3 := auth.AuthenticationV3(options)

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

func authenticate(authV3 auth.Authentication) (map[string]string, error) {
	// do the check params inside do that authenticate is being called from other places
	err := checkAuthenticateAuthParams(authV3.GetOptions())
	if err != nil {
		return map[string]string{}, err
	}
	// get the token result
	token, err := authV3.GetToken()
	if err != nil {
		return map[string]string{}, err
	}

	// arc endpoint
	arcEndpoint, err := authV3.GetServiceEndpoint("arc", viper.GetString(ENV_VAR_REGION), "public")
	if err != nil {
		return map[string]string{}, err
	}

	// automation endpoint
	automationEndpoint, err := authV3.GetServiceEndpoint("automation", viper.GetString(ENV_VAR_REGION), "public")
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

func checkAuthenticateAuthParams(opts *auth.AuthOptions) error {
	// check some params
	if len(opts.UserId) == 0 && len(opts.Username) == 0 {
		return fmt.Errorf("Flag %s or '%s not given.", FLAG_USER_ID, FLAG_USERNAME)
	}

	if len(opts.ProjectId) == 0 && len(opts.ProjectName) == 0 {
		return fmt.Errorf("Flag %s or %s not given.", FLAG_PROJECT_ID, FLAG_PROJECT_NAME)
	}

	if len(opts.IdentityEndpoint) == 0 {
		return fmt.Errorf("Flag %s not given.", FLAG_AUTH_URL)
	}

	// check password and prompt
	if len(opts.Password) == 0 {
		// ask the user for the password
		fmt.Print("Enter password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			return err
		}
		opts.Password = string(pass)
	}

	return nil
}
