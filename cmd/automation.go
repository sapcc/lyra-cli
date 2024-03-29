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

	"github.com/sapcc/lyra-cli/locales"
	"github.com/spf13/cobra"
)

// removed tags since no use case yet
// Tags               map[string]string `json:"tags,omitempty"`      // JSON
type Automation struct {
	Id                              int     `json:"id"`
	Name                            string  `json:"name"`                // required
	Repository                      string  `json:"repository"`          // required
	RepositoryRevision              string  `json:"repository_revision"` // required
	RepositoryCredentials           *string `json:"repository_credentials,omitempty"`
	RepositoryAuthenticationEnabled *bool   `json:"repository_authentication_enabled,omitempty"`
	Timeout                         int     `json:"timeout"` // required
}

type Chef struct {
	Automation
	AutomationType string      `json:"type"`
	Runlist        []string    `json:"run_list,omitempty"`        // required, JSON
	Attributes     interface{} `json:"chef_attributes,omitempty"` // JSON
	LogLevel       string      `json:"log_level,omitempty"`
	Debug          bool        `json:"debug,omitempty"`
	ChefVersion    string      `json:"chef_version,omitempty"`
}

type Script struct {
	Automation
	AutomationType string            `json:"type"`
	Path           string            `json:"path"`
	Arguments      []string          `json:"arguments"`   // array of strings
	Environment    map[string]string `json:"environment"` // JSON
}

var (
	chef   Chef
	script Script
)

// automationCmd represents the automation command
var AutomationCmd = &cobra.Command{
	Use:   "automation",
	Short: locales.CmdShortDescription("automation"),
}

func init() {
	RootCmd.AddCommand(AutomationCmd)
	initAutomationCmdFlags()
}

func initAutomationCmdFlags() {
}

// Unmarshal map to chef struct
func (c *Chef) Unmarshal(response string) error {
	respByt := []byte(response)
	if err := json.Unmarshal(respByt, &c); err != nil {
		return err
	}
	return nil
}

// Unmarshal map to chef struct
func (s *Script) Unmarshal(response string) error {
	respByt := []byte(response)
	if err := json.Unmarshal(respByt, &s); err != nil {
		return err
	}
	return nil
}

// Marshal map chef json
func (c *Chef) Marshal() (string, error) {
	// convert to json
	body, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Marshal map chef json
func (s *Script) Marshal() (string, error) {
	// convert to json
	body, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
