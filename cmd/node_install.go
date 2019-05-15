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

	"github.com/sapcc/lyra-cli/helpers"
	"github.com/sapcc/lyra-cli/locales"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var NodeInstallCmd = &cobra.Command{
	Use:   "install",
	Short: locales.CmdShortDescription("arc-node-install"),
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
		fmt.Println(script)

		return nil
	},
}

func init() {
	NodeCmd.AddCommand(NodeInstallCmd)
	initNodeInstallCmdFlags()
}

func initNodeInstallCmdFlags() {
	NodeInstallCmd.Flags().StringP(FLAG_ARC_NODE_ID, "", "", locales.AttributeDescription(FLAG_ARC_NODE_ID))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("arc-node-id", NodeInstallCmd.Flags().Lookup(FLAG_ARC_NODE_ID)), "BindPFlag:")
	NodeInstallCmd.Flags().StringP(FLAG_ARC_INSTALL_FORMAT, "", "json", locales.AttributeDescription(FLAG_ARC_INSTALL_FORMAT))
	helpers.CheckErrAndPrintToStdErr(viper.BindPFlag("arc-install-format", NodeInstallCmd.Flags().Lookup(FLAG_ARC_INSTALL_FORMAT)), "BindPFlag:")
}

func checkArcInstallParams() error {
	switch viper.GetString("arc-install-format") {
	case "linux":
	case "windows":
	case "cloud-config":
	case "json":
	default:
		return fmt.Errorf("invalid %#v given. Valid: windows,linux,cloud-config,json", "arc-install-format")
	}

	return nil
}

func generateScript() (string, error) {
	requestBody, err := json.Marshal(&map[string]string{"CN": viper.GetString("arc-node-id")})
	if err != nil {
		return "", errors.New("failed to marshel request body")
	}
	arcService := RestClient.Services["arc"]

	acceptHeader := "application/json"
	switch viper.GetString("arc-install-format") {
	case "linux", "shell":
		acceptHeader = "text/x-shellscript"
	case "windows", "powershell":
		acceptHeader = "text/x-powershellscript"
	case "cloud-config":
		acceptHeader = "text/cloud-config"
	}
	response, status, err := arcService.Post("agents/init", url.Values{}, http.Header{"Accept": []string{acceptHeader}}, string(requestBody))
	if err != nil {
		return "", err
	}
	if status >= 400 {
		return "", fmt.Errorf("received %d reponse", status)
	}

	return response, nil
}

type PkiResult struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}
