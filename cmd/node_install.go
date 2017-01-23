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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sapcc/lyra-cli/locales"
)

const (
	ARC_INSTALL_NODE_IDENTITY_FLAG = "node-identity"
	ARC_INSTALL_FORMAT_FLAG        = "install-format"
)

var NodeInstallCmd = &cobra.Command{
	Use:   "install",
	Short: locales.CmdShortDescription("arc-install"),
	RunE: func(cmd *cobra.Command, args []string) error {

		err := checkArcInstallParams()
		if err != nil {
			return err
		}

		script, err := generateScript()
		if err != nil {
			return err
		}

		// print response
		cmd.Println(script)

		return nil
	},
}

func init() {
	NodeCmd.AddCommand(NodeInstallCmd)
	initNodeInstallCmdFlags()
}

func checkArcInstallParams() error {

	switch viper.GetString(ARC_INSTALL_FORMAT_FLAG) {
	case "linux":
	case "windows":
	case "cloud-config":
	case "json":
	default:
		return fmt.Errorf("Invalid %#v given. Valid: windows,linux,cloud-config,json", ARC_INSTALL_FORMAT_FLAG)
	}

	return nil
}

func initNodeInstallCmdFlags() {
	NodeInstallCmd.Flags().StringP(ARC_INSTALL_NODE_IDENTITY_FLAG, "", "", locales.AttributeDescription(ARC_INSTALL_NODE_IDENTITY_FLAG))
	viper.BindPFlag(ARC_INSTALL_NODE_IDENTITY_FLAG, NodeInstallCmd.Flags().Lookup(ARC_INSTALL_NODE_IDENTITY_FLAG))
	NodeInstallCmd.Flags().StringP(ARC_INSTALL_FORMAT_FLAG, "", "json", locales.AttributeDescription(ARC_INSTALL_FORMAT_FLAG))
	viper.BindPFlag(ARC_INSTALL_FORMAT_FLAG, NodeInstallCmd.Flags().Lookup(ARC_INSTALL_FORMAT_FLAG))
}

func generateScript() (string, error) {

	requestBody, err := json.Marshal(&map[string]string{"CN": viper.GetString(ARC_INSTALL_NODE_IDENTITY_FLAG)})
	if err != nil {
		return "", errors.New("Failed to marshel request body")
	}
	arcService := RestClient.Services["arc"]

	acceptHeader := "application/json"
	switch viper.GetString(ARC_INSTALL_FORMAT_FLAG) {
	case "linux", "shell":
		acceptHeader = "text/x-shellscript"
	case "windows", "powershell":
		acceptHeader = "text/x-powershellscript"
	case "cloud-config":
		acceptHeader = "text/cloud-config"
	}
	response, status, err := arcService.Post("pki/token", url.Values{}, http.Header{"Accept": []string{acceptHeader}}, string(requestBody))
	if err != nil {
		return "", err
	}
	if status >= 400 {
		return "", fmt.Errorf("Received %d reponse", status)
	}

	return response, nil
}

type PkiResult struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}
