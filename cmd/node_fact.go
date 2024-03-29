// Copyright © 2016 Arturo Reuschenbach Puncernau <a.reuschenbach.puncernau@sap.com>
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
	"github.com/sapcc/lyra-cli/locales"
	"github.com/spf13/cobra"
)

var NodeFactCmd = &cobra.Command{
	Use:   "fact",
	Short: locales.CmdShortDescription("arc-node-fact"),
}

func init() {
	NodeCmd.AddCommand(NodeFactCmd)
}
