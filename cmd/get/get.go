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
	"fmt"

	"github.com/spf13/cobra"
)

const (
	getCmdHelp = `Display one or many OpenEBS resources like volumes, storages, blockdevices

$ kubectl openebs get [volume|storage|bd] [flags]

# Get volumes
$ kubectl openebs get volume [options]

# Get storages
$ kubectl openebs get storage [options]

# Get blockdevices
$ kubectl openebs get bd

Options:
--------
Filter volumes by cas-type
--cas-type=[cstor]

Advanced:
Filter by a fixed OpenEBS namespace
--openebs-namespace=[...]
`
)

// NewCmdGet provides options for managing OpenEBS Volume
func NewCmdGet(rootCmd *cobra.Command) *cobra.Command {
	var casType string
	cmd := &cobra.Command{
		Use:       "get",
		Short:     "Provides operations related to a Volume",
		ValidArgs: []string{"storage", "volume", "bd"},
		Long:      getCmdHelp,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getCmdHelp)
		},
	}
	cmd.PersistentFlags().StringVarP(&casType, "cas-type", "", "", "the cas-type filter option for fetching resources")
	cmd.AddCommand(
		NewCmdGetVolume(),
		NewCmdGetStorage(),
		NewCmdGetBD(),
	)
	return cmd
}
