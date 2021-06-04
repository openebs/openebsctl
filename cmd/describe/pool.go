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

package describe

import (
	"github.com/openebs/openebsctl/pkg/cstor/describe"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	poolInfoCommandHelpText = `
This command fetches information and status of the various aspects 
of the cStor Pool Instance and its underlying related resources in the provided namespace.
If no namespace is provided it uses default namespace for execution.
$ kubectl openebs describe pool [cspi-name] -n [namespace]
`
)

// NewCmdDescribePool displays OpenEBS cStor pool instance information.
func NewCmdDescribePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Aliases: []string{"pools", "p"},
		Short:   "Displays cStorPoolInstance information",
		Long:    poolInfoCommandHelpText,
		Example: `kubectl openebs describe pool cspi-one -n openebs`,
		Run: func(cmd *cobra.Command, args []string) {
			openebsNs, _ := cmd.Flags().GetString("openebs-namespace")
			util.CheckErr(describe.RunPoolInfo(cmd, args, openebsNs), util.Fatal)
		},
	}
	return cmd
}
