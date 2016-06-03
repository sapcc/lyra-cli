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
	"os"

	"github.com/spf13/cobra"
)

// bash-completionCmd represents the bash-completion command
var bashCompletionCmd = &cobra.Command{
	Use:   "bash-completion",
	Short: "Generate completions for bash",
	Long:  `Add $(lyra bash-completion) to your .bashrc to enable tab completion for lyra`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// DO NOT REMOVE. SHOULD OVERRIDE THE ROOT PersistentPreRunE
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		RootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	RootCmd.AddCommand(bashCompletionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bash-completionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// bash-completionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
