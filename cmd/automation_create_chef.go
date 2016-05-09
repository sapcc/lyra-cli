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
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/sapcc/lyra-cli/helpers"
)

type Chef struct {
	Name               string            `json:"name"`                // required
	Repository         string            `json:"repository"`          // required
	RepositoryRevision string            `json:"repository_revision"` // required
	Timeout            int               `json:"timeout"`             // required
	Tags               map[string]string `json:"tags,omitempty"`      // JSON
	AutomationType     string            `json:"type"`
	Runlist            []string          `json:"run_list,omitempty"`        // required, JSON
	Attributes         string            `json:"chef_attributes,omitempty"` // JSON
	LogLevel           string            `json:"log_level,omitempty"`
}

var (
	chef               = Chef{}
	tags               string // JSON (1 level key value)
	runlist            string // JSON (1 level array)
	attributes         string // JSON
	attributesFromFile string // paht to a file
)

// createCmd represents the create command
var AutomationCreateChefCmd = &cobra.Command{
	Use:   "chef",
	Short: "Create a new chef automation.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// setup
		err := setupRestClient()
		if err != nil {
			return err
		}
		// setup attributes
		err = setupCreateChef()
		if err != nil {
			return nil
		}
		// create automation
		response, err := AutomationCreateChef()
		if err != nil {
			return err
		}
		// Print response
		cmd.Print(response)

		return nil
	},
}

func init() {
	AutomationCreateCmd.AddCommand(AutomationCreateChefCmd)
	AutomationCreateChefCmd.Flags().StringVarP(&chef.Name, "name", "", "", "Describes the template. Should be short and alphanumeric without white spaces.")
	AutomationCreateChefCmd.Flags().StringVarP(&chef.Repository, "repository", "", "", "Describes the place where the automation is being described. Git ist the only suported repository type. Ex: https://github.com/user123/automation-test.git.")
	AutomationCreateChefCmd.Flags().StringVarP(&chef.RepositoryRevision, "repository-revision", "", "master", "Describes the repository branch.")
	AutomationCreateChefCmd.Flags().IntVarP(&chef.Timeout, "timeout", "", 3600, "Describes the time elapsed before a timeout is being triggered.")
	AutomationCreateChefCmd.Flags().StringVarP(&tags, "tags", "", "", "Are key value pairs. Key-value pairs are separated by ':' or '='. Following this pattern: 'key1:value1,key2=value2...'.")
	AutomationCreateChefCmd.Flags().StringVarP(&runlist, "runlist", "", "", "Describes the sequence of recipes should be executed. Runlist is an array of strings. Array of strings are separated by ','.")
	AutomationCreateChefCmd.Flags().StringVarP(&attributes, "attributes", "", "", "Attributes are JSON based.")
	AutomationCreateChefCmd.Flags().StringVarP(&attributesFromFile, "attributes-from-file", "", "", "Path to the file containing the chef attributes in JSON format. Giving a dash '-' will be read from standard input.")
	AutomationCreateChefCmd.Flags().StringVarP(&chef.LogLevel, "log-level", "", "", "Describe the level should be used when logging.")
}

func setupCreateChef() error {
	chef.Tags = helpers.StringTokeyValueMap(tags)
	chef.Runlist = helpers.StringToArray(runlist)

	// read attributes
	if len(attributes) > 0 {
		chef.Attributes = attributes
	} else {
		// check for a dash
		if len(attributesFromFile) == 1 && attributesFromFile == "-" {
			// read from input
			var buffer bytes.Buffer
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				buffer.WriteString(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return err
			}
			chef.Attributes = buffer.String()
		} else if len(attributesFromFile) > 1 {
			// read file
			dat, err := ioutil.ReadFile(attributesFromFile)
			if err != nil {
				return err
			}
			chef.Attributes = string(dat)
		}
	}

	return nil
}

func AutomationCreateChef() (string, error) {
	// add the type
	chef.AutomationType = "Chef"
	// convert to Json
	body, err := json.Marshal(chef)
	if err != nil {
		return "", err
	}

	response, err := RestClient.Post("automations", url.Values{}, string(body))
	if err != nil {
		return "", err
	}

	return response, nil
}
