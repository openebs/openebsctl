/*
Copyright 2020-2021 The OpenEBS Authors

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

package get

import (
	"github.com/openebs/openebsctl/pkg/cstor/get"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	poolListCommandHelpText = `
This command lists of all known pools in the Cluster.

Usage:
$ kubectl openebs get pool [options]
`
)

// NewCmdGetPool displays status of OpenEBS Pool(s)
func NewCmdGetPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Aliases: []string{"pools", "p"},
		Short:   "Displays status information about Pool(s)",
		Long:    poolListCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			ns, err := cmd.Flags().GetString("namespace")
			if err != nil {
				ns = "openebs"
			}
			// TODO: De-couple CLI code, logic code, API code
			util.CheckErr(get.RunPoolsList(cmd, args, ns), util.Fatal)
		},
	}
	return cmd
}
