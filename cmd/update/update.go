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

package update

import (
	"fmt"

	"github.com/openebs/openebsctl/pkg/storage"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

const (
	updateCmdHelp = `Modify one or many OpenEBS resources like storage

$ kubectl openebs update [resource-type] [resource-names]  [arguments...] [options...]
`
	cspcPatchHelp = `Update CSPC to move pools from an old node to the newer
node after the blockdevices have moved to the newer node.`
)

// NewCmdUpdate provides options for managing OpenEBS resources
func NewCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "update",
		Short:     "Provides update operations related to an OpenEBS resource",
		ValidArgs: []string{"storage"},
		Long:      updateCmdHelp,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(updateCmdHelp)
		},
	}
	cmd.AddCommand(
		NewCmdCSPCNodePatch(),
	)
	return cmd
}

// NewCmdCSPCNodePatch provides option to patch CStor CSPC
func NewCmdCSPCNodePatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "storage",
		Aliases: []string{"s", "storages", "pool"},
		Short:   "Updates information about a storage",
		Long:    updateCmdHelp,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 3 {
				ns, _ := cmd.Flags().GetString("openebs-namespace")
				util.CheckErr(storage.Update(ns, args[0], args[1], args[2]), util.Fatal)
			} else {
				fmt.Println(cspcPatchHelp)
			}
		},
	}
	return cmd
}
