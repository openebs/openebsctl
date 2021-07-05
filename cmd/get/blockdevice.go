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
	"github.com/openebs/openebsctl/pkg/blockdevice"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	bdListCommandHelpText = `
This command displays status of available OpenEBS BlockDevice(s).

Usage: kubectl openebs get bd [options]

Advanced:
Filter by a fixed OpenEBS namespace
--openebs-namespace=[...]
`
)

// NewCmdGetBD displays status of OpenEBS BlockDevice(s)
func NewCmdGetBD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bd",
		Aliases: []string{"bds", "blockdevice", "blockdevices"},
		Short:   "Displays status information about BlockDevice(s)",
		Long:    bdListCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Should this method create the k8sClient object
			openebsNS, _ := cmd.Flags().GetString("openebs-namespace")
			util.CheckErr(blockdevice.Get(args, openebsNS), util.Fatal)
		},
	}
	return cmd
}
