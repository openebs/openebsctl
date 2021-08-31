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
	"github.com/spf13/cobra"
)

const (
	getCmdHelp = `Display one or many OpenEBS resources like volumes, storages, blockdevices.

Usage:
  kubectl openebs get [volume|storage|bd] [flags]

Get volumes:
  kubectl openebs get volume [flags]

Get storages:
  kubectl openebs get storage [flags]

Get blockdevices:
  kubectl openebs get bd

Flags:
  -h, --help                           help for openebs get command
      --openebs-namespace string       filter by a fixed OpenEBS namespace
                                       If not provided it is determined from components.
      --cas-type                       to specify the cas-type of the engine, for engine based filtering.
                                       ex- cstor, jiva, localpv-lvm, localpv-zfs.
`
)

// NewCmdGet provides options for managing OpenEBS Volume
func NewCmdGet(rootCmd *cobra.Command) *cobra.Command {
	var casType string
	cmd := &cobra.Command{
		Use:       "get",
		Short:     "Provides fetching operations related to a Volume/Storage",
		ValidArgs: []string{"storage", "volume", "bd"},
	}
	cmd.SetUsageTemplate(getCmdHelp)
	cmd.PersistentFlags().StringVarP(&casType, "cas-type", "", "", "the cas-type filter option for fetching resources")
	cmd.AddCommand(
		NewCmdGetVolume(),
		NewCmdGetStorage(),
		NewCmdGetBD(),
	)
	return cmd
}
