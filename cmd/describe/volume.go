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
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/openebs/openebsctl/pkg/volume"
	"github.com/spf13/cobra"
)

var (
	volumeInfoCommandHelpText = `This command fetches information and status of the various
aspects of a cStor Volume such as ISCSI, Controller, and Replica.

Usage:
  kubectl openebs describe volume [...names] [flags]

Flags:
  -h, --help                           help for openebs
  -n, --namespace string               to read the namespace for the pvc.
      --openebs-namespace string       to read the openebs namespace from user.
                                       If not provided it is determined from components.
`
)

// NewCmdDescribeVolume displays OpenEBS Volume information.
func NewCmdDescribeVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"volumes", "vol", "v"},
		Short:   "Displays Openebs information",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Get this from flags, pflag, etc
			openebsNS, _ := cmd.Flags().GetString("openebs-namespace")
			util.CheckErr(volume.Describe(args, openebsNS), util.Fatal)
		},
	}
	cmd.SetUsageTemplate(volumeInfoCommandHelpText)
	return cmd
}
