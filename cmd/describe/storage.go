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
	"strings"

	"github.com/openebs/openebsctl/pkg/storage"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

var (
	storageInfoCommandHelpText = `
This command fetches information and status of the various aspects 
of the openebs storage and its underlying related resources in the openebs namespace.

$ kubectl openebs describe storage [storage-name-1...n] [options]

Options:
--------
Filter volumes by cas-type
--cas-type=[cstor]

Advanced:
Filter by a fixed OpenEBS namespace
--openebs-namespace=[...]
`
)

// NewCmdDescribeStorage displays OpenEBS storage related information.
func NewCmdDescribeStorage() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "storage",
		Aliases: []string{"storages", "s"},
		Short:   "Displays storage related information",
		Long:    storageInfoCommandHelpText,
		Example: `kubectl openebs describe storage storage-1`,
		Run: func(cmd *cobra.Command, args []string) {
			openebsNs, _ := cmd.Flags().GetString("openebs-namespace")
			casType, _ := cmd.Flags().GetString("cas-type")
			casType = strings.ToLower(casType)
			util.CheckErr(storage.Describe(args, openebsNs, casType), util.Fatal)
		},
	}
	return cmd
}
