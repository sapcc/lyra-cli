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
	"bytes"
	"fmt"
	"net/url"
	"path"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	auth "github.com/sapcc/go-openstack-auth"
	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/sapcc/lyra-cli/restclient"
)

var ()

const (
	ARC_INSTALL_INSTANCE_IDENTIFIER_FLAG = "instance-identifier"
	ARC_INSTALL_INSTANCE_OS_FLAG         = "instance-os"
	ARC_INSTALL_PKI_URL_FLAG             = "pki-service-url"
	ARC_INSTALL_UPDATE_URL_FLAG          = "update-service-url"
	ARC_INSTALL_ARC_BROKER_URL_FLAG      = "arc-broker-url"
)

var ArcInstallCmd = &cobra.Command{
	Use:   "install",
	Short: locales.CmdShortDescription("arc-install"),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// DO NOT REMOVE. SHOULD OVERRIDE THE ROOT PersistentPreRunE
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// set authentication params
		options := auth.AuthV3Options{
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

		// check auth params
		err := checkArcInstallAuthParams(&options)
		if err != nil {
			return err
		}

		// check install params
		err = checkArcInstallParams()
		if err != nil {
			return err
		}

		authV3 := auth.AuthenticationV3(options)
		token, err := authV3.GetToken()
		if err != nil {
			return err
		}

		project, err := authV3.GetProject()
		if err != nil {
			return err
		}

		script, err := generateScript(token.ID, project.ID, project.DomainID)
		if err != nil {
			return err
		}

		fmt.Println("°°°")
		fmt.Printf("%+v\n", script)
		fmt.Println("°°°")

		return nil
	},
}

func init() {
	ArcCmd.AddCommand(ArcInstallCmd)
	initArcInstallCmdFlags()
}

func checkArcInstallAuthParams(opts *auth.AuthV3Options) error {
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

func checkArcInstallParams() error {
	if len(viper.GetString(ARC_INSTALL_INSTANCE_IDENTIFIER_FLAG)) == 0 {
		return fmt.Errorf("Flag %s not given.", ARC_INSTALL_INSTANCE_IDENTIFIER_FLAG)
	}
	if len(viper.GetString(ARC_INSTALL_INSTANCE_OS_FLAG)) == 0 {
		return fmt.Errorf("Flag %s not given.", ARC_INSTALL_INSTANCE_OS_FLAG)
	}

	if viper.GetString(ARC_INSTALL_INSTANCE_OS_FLAG) != "linux" && viper.GetString(ARC_INSTALL_INSTANCE_OS_FLAG) != "windows" {
		return fmt.Errorf("Flag %s value not known. Available OS are linux or windows.", ARC_INSTALL_INSTANCE_OS_FLAG)
	}

	return nil
}

func initArcInstallCmdFlags() {
	ArcInstallCmd.Flags().StringP(ARC_INSTALL_INSTANCE_IDENTIFIER_FLAG, "", "", locales.AttributeDescription("arc-install-identifier"))
	viper.BindPFlag(ARC_INSTALL_INSTANCE_IDENTIFIER_FLAG, ArcInstallCmd.Flags().Lookup(ARC_INSTALL_INSTANCE_IDENTIFIER_FLAG))
	ArcInstallCmd.Flags().StringP(ARC_INSTALL_INSTANCE_OS_FLAG, "", "", locales.AttributeDescription("arc-install-os"))
	viper.BindPFlag(ARC_INSTALL_INSTANCE_OS_FLAG, ArcInstallCmd.Flags().Lookup(ARC_INSTALL_INSTANCE_OS_FLAG))
	ArcInstallCmd.Flags().StringP(ARC_INSTALL_UPDATE_URL_FLAG, "", "https://arc-updates.***REMOVED***", locales.AttributeDescription("update-service-url"))
	viper.BindPFlag(ARC_INSTALL_UPDATE_URL_FLAG, ArcInstallCmd.Flags().Lookup(ARC_INSTALL_UPDATE_URL_FLAG))
	ArcInstallCmd.Flags().StringP(ARC_INSTALL_PKI_URL_FLAG, "", "https://arc-pki.***REMOVED***", locales.AttributeDescription("pki-service-url"))
	viper.BindPFlag(ARC_INSTALL_PKI_URL_FLAG, ArcInstallCmd.Flags().Lookup(ARC_INSTALL_PKI_URL_FLAG))
	ArcInstallCmd.Flags().StringP(ARC_INSTALL_ARC_BROKER_URL_FLAG, "", "tls://arc-broker.***REMOVED***:8883", locales.AttributeDescription("arc-broker-url"))
	viper.BindPFlag(ARC_INSTALL_ARC_BROKER_URL_FLAG, ArcInstallCmd.Flags().Lookup(ARC_INSTALL_ARC_BROKER_URL_FLAG))
}

func generateScript(token, projectId, domainId string) (string, error) {
	registrationUrl, err := registrationUrl(token, projectId, domainId)
	if err != nil {
		return "", err
	}
	return processRequest(registrationUrl)
}

// LINUX example
// curl --create-dirs -o /opt/arc/arc https://arc-updates.***REMOVED***/builds/latest/arc/linux/amd64
// chmod +x /opt/arc/arc
// /opt/arc/arc init --endpoint tls://arc-broker.***REMOVED***:8883 --update-uri https://arc-updates.***REMOVED***/updates --registration-url https://arc-pki.***REMOVED***/api/v1/sign/9f164fdb-a791-42a6-af6b-bb246aec9e00
// WINDOWS example
// mkdir C:\monsoon\arc
// powershell (new-object System.Net.WebClient).DownloadFile('https://arc-updates.***REMOVED***/builds/latest/arc/windows/amd64','C:\monsoon\arc\arc.exe')
// C:\monsoon\arc\arc.exe init --endpoint tls://arc-broker.***REMOVED***:8883 --update-uri https://arc-updates.***REMOVED***/updates --registration-url https://arc-pki.***REMOVED***/api/v1/sign/32151d55-7b58-41a5-9687-77a6ef3b05a2
func processRequest(registrationUrl string) (string, error) {
	// latest build url
	latestBuildUrl, err := url.Parse(viper.GetString(ARC_INSTALL_UPDATE_URL_FLAG))
	if err != nil {
		return "", err
	}
	latestBuildUrl.Path = path.Join(latestBuildUrl.Path, "builds/latest/arc", viper.GetString(ARC_INSTALL_INSTANCE_OS_FLAG), "amd64")

	// update url
	updateUrl, err := url.Parse(viper.GetString(ARC_INSTALL_UPDATE_URL_FLAG))
	if err != nil {
		return "", err
	}
	updateUrl.Path = path.Join(updateUrl.Path, "updates")

	// build script
	var buffer bytes.Buffer

	if viper.GetString(ARC_INSTALL_INSTANCE_OS_FLAG) == "linux" {
		buffer.WriteString(fmt.Sprint(`curl --create-dirs -o /opt/arc/arc `, latestBuildUrl.String(), "\n"))
		buffer.WriteString(fmt.Sprint(`chmod +x /opt/arc/arc`, "\n"))
		buffer.WriteString(fmt.Sprint(`/opt/arc/arc init --endpoint `, viper.GetString(ARC_INSTALL_ARC_BROKER_URL_FLAG), " --update-uri ", updateUrl, " --registration-url ", registrationUrl))
	} else if viper.GetString(ARC_INSTALL_INSTANCE_OS_FLAG) == "windows" {
		buffer.WriteString(fmt.Sprint(`mkdir C:\monsoon\arc`, "\n"))
		buffer.WriteString(fmt.Sprint(`powershell (new-object System.Net.WebClient).DownloadFile('`, latestBuildUrl, `','C:\monsoon\arc\arc.exe')`, "\n"))
		buffer.WriteString(fmt.Sprint(`C:\monsoon\arc\arc.exe init --endpoint `, viper.GetString(ARC_INSTALL_ARC_BROKER_URL_FLAG), " --update-uri ", updateUrl, " --registration-url ", registrationUrl))
	}

	return buffer.String(), nil
}

type PkiResult struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

func registrationUrl(token, projectId, domainId string) (string, error) {
	// init rest client
	body := fmt.Sprintf(`{"CN": "%s", "names": [{"OU": "%s", "O": "%s"}] }`, viper.GetString(ARC_INSTALL_INSTANCE_IDENTIFIER_FLAG), projectId, domainId)
	pkiClient := restclient.NewClient([]restclient.Endpoint{restclient.Endpoint{ID: "pki", Url: viper.GetString(ARC_INSTALL_PKI_URL_FLAG)}}, token)
	pkiService := pkiClient.Services["pki"]
	response, _, err := pkiService.Post("api/v1/token", url.Values{}, body)
	if err != nil {
		return "", fmt.Errorf("Error using pki service. %s", err.Error())
	}

	res := PkiResult{}
	err = helpers.JSONStringToStructure(response, &res)
	if err != nil {
		return "", fmt.Errorf("Error unmarshaling pki response. %s", err.Error())
	}

	return res.Url, nil
}
