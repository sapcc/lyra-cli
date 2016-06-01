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
	"github.com/spf13/cobra"
)

type Chef struct {
	Id                 int               `json:"id"`
	Name               string            `json:"name"`                // required
	Repository         string            `json:"repository"`          // required
	RepositoryRevision string            `json:"repository_revision"` // required
	Timeout            int               `json:"timeout"`             // required
	Tags               map[string]string `json:"tags,omitempty"`      // JSON
	AutomationType     string            `json:"type"`
	Runlist            []string          `json:"run_list,omitempty"`        // required, JSON
	Attributes         interface{}       `json:"chef_attributes,omitempty"` // JSON
	LogLevel           string            `json:"log_level,omitempty"`
}

type Script struct {
	Id                 int               `json:"id"`
	Name               string            `json:"name"`                // required
	Repository         string            `json:"repository"`          // required
	RepositoryRevision string            `json:"repository_revision"` // required
	Timeout            int               `json:"timeout"`             // required
	Tags               map[string]string `json:"tags,omitempty"`      // JSON
	AutomationType     string            `json:"type"`
	Path               string            `json:"path"`
	Arguments          []string          `json:"arguments"`   // array of strings
	Environment        map[string]string `json:"environment"` // JSON
}

var (
	chef   Chef
	script Script
)

// automationCmd represents the automation command
var AutomationCmd = &cobra.Command{
	Use:   "automation",
	Short: "Automation service.",
	Long:  `A longer description for automation.`,
}

func init() {
	RootCmd.AddCommand(AutomationCmd)
	initAutomationCmdFlags()
}

func initAutomationCmdFlags() {
}
