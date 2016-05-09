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
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var AutomationUpdateChefAttributesCmd = &cobra.Command{
	Use:   "attributes",
	Short: "Updates chef attributes",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	AutomationUpdateChefCmd.AddCommand(AutomationUpdateChefAttributesCmd)
	//flags
	AutomationUpdateChefAttributesCmd.Flags().StringVarP(&attributes, "attributes", "", "", "Attributes are JSON format.")
	AutomationUpdateChefAttributesCmd.Flags().StringVarP(&attributesFromFile, "attributes-from-file", "", "", "Path to the file containing the chef attributes in JSON format. Giving a dash '-' will be read from standard input.")
}
