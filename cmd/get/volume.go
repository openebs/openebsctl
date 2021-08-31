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
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/openebs/openebsctl/pkg/volume"
	"github.com/spf13/cobra"
)

var (
	volumesListCommandHelpText = `Usage: 
  kubectl openebs get volume [flags]

Flags:
  -h, --help                           help for openebs get command
      --openebs-namespace string       filter by a fixed OpenEBS namespace
                                       If not provided it is determined from components.
      --cas-type                       to specify the cas-type of the engine, for engine based filtering.
                                       ex- cstor, jiva, localpv-lvm, localpv-zfs.
`
)

// NewCmdGetVolume displays status of OpenEBS Volume(s)
func NewCmdGetVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"vol", "v", "volumes"},
		Short:   "Displays status information about Volume(s)",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Should this method create the k8sClient object
			openebsNS, _ := cmd.Flags().GetString("openebs-namespace")
			casType, _ := cmd.Flags().GetString("cas-type")
			util.CheckErr(volume.Get(args, openebsNS, casType), util.Fatal)
		},
	}
	cmd.SetUsageTemplate(volumesListCommandHelpText)
	return cmd
}
