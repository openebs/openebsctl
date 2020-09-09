/*
Copyright 2020 The OpenEBS Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	"github.com/spf13/cobra"

	"github.com/openebs/openebsctl/kubectl-openebs/cli/command/cstor"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
)

// NewOpenebsCommand creates the `openebs` command and its nested children.
func NewOpenebsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl-openebs",
		Short: "OpenEBSctl is a a tool for interacting with OpenEBS storage components",
		Long:  `OpenEBSctl is a kubectl plugin to interact with OpenEBS container Attached Storage components. `,
	}

	cmd.AddCommand(
		util.NewCmdCompletion(cmd),
		cstor.NewCmdcStor(),
	)

	return cmd
}
