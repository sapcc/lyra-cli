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
	"net/url"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/restclient"
)

var (
	cfgFile                          string
	Token                            string
	AutomationUrl                    string
	ArcUrl                           string
	RestClient                       *restclient.Client
	ENV_VAR_TOKEN_NAME               = "OS_TOKEN"
	ENV_VAR_AUTOMATION_ENDPOINT_NAME = "LYRA_SERVICE_ENDPOINT"
	ENV_VAR_ARC_ENDPOINT_NAME        = "ARC_SERVICE_ENDPOINT"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "lyra",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		// fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lyra-cli.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "g", false, "Help message for toggle")

	// Custom flags
	token_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_TOKEN_NAME)
	automation_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_AUTOMATION_ENDPOINT_NAME)
	arc_default_env_name := fmt.Sprintf("[$%s]", ENV_VAR_ARC_ENDPOINT_NAME)
	RootCmd.PersistentFlags().StringVarP(&Token, "token", "t", "", fmt.Sprint("Authentication token. To create a token run the authenticate command. (default ", token_default_env_name, ")"))
	RootCmd.PersistentFlags().StringVarP(&AutomationUrl, "lyra-service-endpoint", "l", "", fmt.Sprint("Automation service endpoint. To get the automation endpoint run the authenticate command. (default ", automation_default_env_name, ")"))
	RootCmd.PersistentFlags().StringVarP(&ArcUrl, "arc-service-endpoint", "a", "", fmt.Sprint("Arc service endpoint. To get the arc endpoint run the authenticate command. (default ", arc_default_env_name, ")"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".lyra-cli") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// setup the rest client
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
			return errors.New("Automation endpoint not given. To get the automation endpoint run the authenticate command.")
		} else {
			AutomationUrl = os.Getenv(ENV_VAR_AUTOMATION_ENDPOINT_NAME)
		}
	}

	if len(ArcUrl) == 0 {
		if len(os.Getenv(ENV_VAR_ARC_ENDPOINT_NAME)) == 0 {
			return errors.New("Arc endpoint not given. To get the arc endpoint run the authenticate command.")
		} else {
			ArcUrl = os.Getenv(ENV_VAR_ARC_ENDPOINT_NAME)
		}
	}

	// add api version to the automation url
	autoUri, err := url.Parse(AutomationUrl)
	if err != nil {
		return err
	}
	autoUri.Path = path.Join(autoUri.Path, "/api/v1/")

	// add api version to the arc url
	arcUri, err := url.Parse(ArcUrl)
	if err != nil {
		return err
	}
	arcUri.Path = path.Join(arcUri.Path, "/api/v1/")

	services := restclient.Services{
		Automation: restclient.Endpoint{Url: autoUri.String()},
		Arc:        restclient.Endpoint{Url: arcUri.String()},
	}

	// init rest client
	RestClient = restclient.NewClient(services, Token)

	return nil
}
