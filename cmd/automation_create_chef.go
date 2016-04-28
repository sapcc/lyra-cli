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

var (
	name              string
	repository        string
	repositoryVersion string
	timeout           int
	tags              string // JSON (1 level key value)
	runlist           string // JSON (1 level array)
	attributes        string // JSON
	logLevel          string
)

// createCmd represents the create command
var AutomationCreateChefCmd = &cobra.Command{
	Use:   "chef",
	Short: "Create a new chef automation.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	AutomationCreateCmd.AddCommand(AutomationCreateChefCmd)
	AutomationCreateChefCmd.Flags().StringVarP(&name, "name", "", "", "Describes the template. Should be short and alphanumeric without white spaces.")
	AutomationCreateChefCmd.Flags().StringVarP(&repository, "repository", "", "", "Describes the place where the automation is being described. Git ist the only suported repository type. Ex: https://github.com/user123/automation-test.git.")
	AutomationCreateChefCmd.Flags().StringVarP(&repositoryVersion, "repository-version", "", "master", "Describes the repository branch.")
	AutomationCreateChefCmd.Flags().IntVarP(&timeout, "timeout", "", 3600, "Describes the time elapsed before a timeout is being triggered.")
	AutomationCreateChefCmd.Flags().StringVarP(&tags, "tags", "", "", "Are key value pairs. Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.")
	AutomationCreateChefCmd.Flags().StringVarP(&runlist, "runlist", "", "", "Describes the sequence of recipes should be executed. Runlist is an array of strings. Array of strings are separated by ','.")
	AutomationCreateChefCmd.Flags().StringVarP(&attributes, "attributes", "", "", "Attributes are JSON based. JSON Schema defines seven primitive types for JSON values: array, boolean, integer, number, null, object and string. Example: {'title':'root'}")
}
