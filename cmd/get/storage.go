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
	"github.com/openebs/openebsctl/pkg/storage"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	storageListCommandHelpText = `This command lists of all/specific known storages in the Cluster.

Usage:
  kubectl openebs get storage [flags]

Flags:
  -h, --help                           help for openebs get command
      --openebs-namespace string       filter by a fixed OpenEBS namespace
                                       If not provided it is determined from components.
      --cas-type                       to specify the cas-type of the engine, for engine based filtering.
                                       ex- cstor, jiva, localpv-lvm, localpv-zfs.
`
)

// NewCmdGetStorage displays status of OpenEBS Pool(s)
func NewCmdGetStorage() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "storage",
		Aliases: []string{"storages", "s"},
		Short:   "Displays status information about Storage(s)",
		Run: func(cmd *cobra.Command, args []string) {
			openebsNS, _ := cmd.Flags().GetString("openebs-namespace")
			casType, _ := cmd.Flags().GetString("cas-type")
			util.CheckErr(storage.Get(args, openebsNS, casType), util.Fatal)
		},
	}
	cmd.SetUsageTemplate(storageListCommandHelpText)
	return cmd
}
