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
	volumeInfoCommandHelpText = `
This command fetches information and status of the various
aspects of a cStor Volume such as ISCSI, Controller, and Replica.

#
$ kubectl openebs describe [pool|volume] [name]

`
)

// NewCmdDescribeVolume displays OpenEBS Volume information.
func NewCmdDescribeVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"volumes", "vol", "v"},
		Short:   "Displays Openebs information",
		Long:    volumeInfoCommandHelpText,
		Example: `kubectl openebs describe volume [vol]`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Get this from flags, pflag, etc
			var ns string
			if ns, _ = cmd.Flags().GetString("namespace"); ns == "" {
				// NOTE: The error comes as nil even when the ns flag is not specified
				ns = "openebs"
			}
			util.CheckErr(describe.RunVolumeInfo(args, ns), util.Fatal)
		},
	}
	return cmd
}
