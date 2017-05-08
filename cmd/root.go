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
	"net/url"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/restclient"
)

var (
	cfgFile string

	RestClient *restclient.Client
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "lyra",
	Short:        locales.CmdShortDescription("root"),
	Long:         locales.CmdLongDescription("root"),
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// setup rest client
		err := setupRestClient(cmd, nil, false)
		if err != nil {
			return err
		}
		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	initRootCmdFlags()
}

func initRootCmdFlags() {
	// set command standard output to the stdout
	RootCmd.SetOutput(os.Stderr)

	// Cobra flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lyra-cli.yaml)")
	RootCmd.Flags().BoolP("toggle", "g", false, "Help message for toggle")

	// Custom flags
	// Results as JSON format
	RootCmd.PersistentFlags().BoolP("json", "j", false, fmt.Sprint("Print out the data in JSON format."))
	viper.BindPFlag("json", RootCmd.PersistentFlags().Lookup("json"))

	// Authentication with token und services flags
	RootCmd.PersistentFlags().StringP(FLAG_TOKEN, "t", "", fmt.Sprint("Authentication token. To create a token run the authenticate command. (default ", fmt.Sprintf("[$%s]", ENV_VAR_TOKEN_NAME), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_LYRA_SERVICE_ENDPOINT, "l", "", fmt.Sprint("Automation service endpoint. To get the automation endpoint run the authenticate command. (default ", fmt.Sprintf("[$%s]", ENV_VAR_AUTOMATION_ENDPOINT_NAME), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_ARC_SERVICE_ENDPOINT, "a", "", fmt.Sprint("Arc service endpoint. To get the arc endpoint run the authenticate command. (default ", fmt.Sprintf("[$%s]", ENV_VAR_ARC_ENDPOINT_NAME), ")"))
	// bind to env variables
	viper.BindPFlag(ENV_VAR_TOKEN_NAME, RootCmd.PersistentFlags().Lookup("token"))
	viper.BindEnv(ENV_VAR_TOKEN_NAME)
	viper.BindPFlag(ENV_VAR_AUTOMATION_ENDPOINT_NAME, RootCmd.PersistentFlags().Lookup("lyra-service-endpoint"))
	viper.BindEnv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)
	viper.BindPFlag(ENV_VAR_ARC_ENDPOINT_NAME, RootCmd.PersistentFlags().Lookup("arc-service-endpoint"))
	viper.BindEnv(ENV_VAR_ARC_ENDPOINT_NAME)

	// Authentication user flags
	RootCmd.PersistentFlags().StringP(FLAG_AUTH_URL, "", "", fmt.Sprint("Endpoint entities represent URL endpoints for OpenStack web services. (default ", fmt.Sprintf("[$%s]", ENV_VAR_AUTH_URL), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_REGION, "", "", fmt.Sprint("A region is a general division of an OpenStack deployment. (default ", fmt.Sprintf("[$%s]", ENV_VAR_REGION), " or the first entry found in catalog)"))
	RootCmd.PersistentFlags().StringP(FLAG_USER_ID, "", "", fmt.Sprint("Id of the user that wants to log in. (default ", fmt.Sprintf("[$%s]", ENV_VAR_USER_ID), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_USERNAME, "", "", fmt.Sprint("Name of the user that wants to log in. (default ", fmt.Sprintf("[$%s]", ENV_VAR_USERNAME), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_PASSWORD, "", "", fmt.Sprint("Password of the user that wants to log in. If not given the environment variable ", fmt.Sprintf("[$%s]", ENV_VAR_PASSWORD), " will be checked. If no environment variable found then will prompt from terminal."))
	RootCmd.PersistentFlags().StringP(FLAG_PROJECT_ID, "", "", fmt.Sprint("Id of the project. (default ", fmt.Sprintf("[$%s]", ENV_VAR_PROJECT_ID), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_PROJECT_NAME, "", "", fmt.Sprint("Name of the project. (default ", fmt.Sprintf("[$%s]", ENV_VAR_PROJECT_NAME), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_USER_DOMAIN_ID, "", "", fmt.Sprint("Id of the domain where the user is created. (default ", fmt.Sprintf("[$%s]", ENV_VAR_USER_DOMAIN_ID), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_USER_DOMAIN_NAME, "", "", fmt.Sprint("Name of the domain where the user is created. (default ", fmt.Sprintf("[$%s]", ENV_VAR_USER_DOMAIN_NAME), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_PROJECT_DOMAIN_ID, "", "", fmt.Sprint("Id of the domain where the project is created. If no project domain id is given, then the token will be scoped in the user domain. (default ", fmt.Sprintf("[$%s]", ENV_VAR_PROJECT_DOMAIN_ID), ")"))
	RootCmd.PersistentFlags().StringP(FLAG_PROEJECT_DOMAIN_NAME, "", "", fmt.Sprint("Name of the domain where the project is created. If no project domain name is given, then the token will be scoped in the user domain. (default ", fmt.Sprintf("[$%s]", ENV_VAR_PROJECT_DOMAIN_NAME), ")"))
	// bind to env variablen
	viper.BindPFlag(ENV_VAR_AUTH_URL, RootCmd.PersistentFlags().Lookup(FLAG_AUTH_URL))
	viper.BindEnv(ENV_VAR_AUTH_URL)
	viper.BindPFlag(ENV_VAR_REGION, RootCmd.PersistentFlags().Lookup(FLAG_REGION))
	viper.BindEnv(ENV_VAR_REGION)
	viper.BindPFlag(ENV_VAR_USER_ID, RootCmd.PersistentFlags().Lookup(FLAG_USER_ID))
	viper.BindEnv(ENV_VAR_USER_ID)
	viper.BindPFlag(ENV_VAR_USERNAME, RootCmd.PersistentFlags().Lookup(FLAG_USERNAME))
	viper.BindEnv(ENV_VAR_USERNAME)
	viper.BindPFlag(ENV_VAR_PASSWORD, RootCmd.PersistentFlags().Lookup(FLAG_PASSWORD))
	viper.BindEnv(ENV_VAR_PASSWORD)
	viper.BindPFlag(ENV_VAR_PROJECT_ID, RootCmd.PersistentFlags().Lookup(FLAG_PROJECT_ID))
	viper.BindEnv(ENV_VAR_PROJECT_ID)
	viper.BindPFlag(ENV_VAR_PROJECT_NAME, RootCmd.PersistentFlags().Lookup(FLAG_PROJECT_NAME))
	viper.BindEnv(ENV_VAR_PROJECT_NAME)
	viper.BindPFlag(ENV_VAR_USER_DOMAIN_ID, RootCmd.PersistentFlags().Lookup(FLAG_USER_DOMAIN_ID))
	viper.BindEnv(ENV_VAR_USER_DOMAIN_ID)
	viper.BindPFlag(ENV_VAR_USER_DOMAIN_NAME, RootCmd.PersistentFlags().Lookup(FLAG_USER_DOMAIN_NAME))
	viper.BindEnv(ENV_VAR_USER_DOMAIN_NAME)
	viper.BindPFlag(ENV_VAR_PROJECT_DOMAIN_ID, RootCmd.PersistentFlags().Lookup(FLAG_PROJECT_DOMAIN_ID))
	viper.BindEnv(ENV_VAR_PROJECT_DOMAIN_ID)
	viper.BindPFlag(ENV_VAR_PROJECT_DOMAIN_NAME, RootCmd.PersistentFlags().Lookup(FLAG_PROEJECT_DOMAIN_NAME))
	viper.BindEnv(ENV_VAR_PROJECT_DOMAIN_NAME)
	// debug flag
	RootCmd.PersistentFlags().BoolP(FLAG_DEBUG, "", false, "Print out request and response objects.")
	viper.BindPFlag(FLAG_DEBUG, RootCmd.PersistentFlags().Lookup(FLAG_DEBUG))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".lyra-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	// viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file: ", viper.ConfigFileUsed())
	}
}

// setup the rest client
func setupRestClient(cmd *cobra.Command, authV3 *auth.Authentication, forceReauthenticate bool) error {
	if len(viper.GetString(ENV_VAR_TOKEN_NAME)) == 0 || len(viper.GetString(ENV_VAR_AUTOMATION_ENDPOINT_NAME)) == 0 || len(viper.GetString(ENV_VAR_ARC_ENDPOINT_NAME)) == 0 || forceReauthenticate {
		fmt.Fprintln(os.Stderr, "Using password authentication.")

		// authentication object
		if authV3 == nil {
			lyraAuthOps := auth.AuthOptions{
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

			newAuthV3 := auth.AuthenticationV3(lyraAuthOps)
			authV3 = &newAuthV3
		}

		// authenticate
		authParams, err := authenticate(cmd, *authV3)
		if err != nil {
			return err
		}

		// reset the vars to viper
		viper.Set(ENV_VAR_AUTOMATION_ENDPOINT_NAME, authParams[ENV_VAR_AUTOMATION_ENDPOINT_NAME])
		viper.Set(ENV_VAR_ARC_ENDPOINT_NAME, authParams[ENV_VAR_ARC_ENDPOINT_NAME])
		viper.Set(ENV_VAR_TOKEN_NAME, authParams[ENV_VAR_TOKEN_NAME])
		viper.Set(TOKEN_EXPIRES_AT, authParams[TOKEN_EXPIRES_AT])
	} else {
		fmt.Fprintln(os.Stderr, "Using token authentication.")
	}

	// add api version to the automation url
	autoUri, err := url.Parse(viper.GetString(ENV_VAR_AUTOMATION_ENDPOINT_NAME))
	if err != nil {
		return err
	}
	autoUri.Path = path.Join(autoUri.Path, "/api/v1/")

	// add api version to the arc url
	arcUri, err := url.Parse(viper.GetString(ENV_VAR_ARC_ENDPOINT_NAME))
	if err != nil {
		return err
	}
	arcUri.Path = path.Join(arcUri.Path, "/api/v1/")

	endpoints := []restclient.Endpoint{
		restclient.Endpoint{ID: "automation", Url: autoUri.String()},
		restclient.Endpoint{ID: "arc", Url: arcUri.String()},
	}

	// init rest client
	RestClient = restclient.NewClient(endpoints, viper.GetString(ENV_VAR_TOKEN_NAME), viper.GetBool(FLAG_DEBUG))

	return nil
}
