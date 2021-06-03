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
	getCmdHelp = `Display one or many OpenEBS resources like volumes, pools

$ kubectl openebs get [volumes|pools] [-n example-namespace]

# Get volumes
$ kubectl openebs get volume

# Get pools
$ kubectl openebs get pool
`
)

// NewCmdGet provides options for managing OpenEBS Volume
func NewCmdGet(rootCmd *cobra.Command) *cobra.Command {
	var casType string
	cmd := &cobra.Command{
		Use:       "get",
		Short:     "Provides operations related to a Volume",
		ValidArgs: []string{"pool", "volume"},
		Long:      getCmdHelp,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getCmdHelp)
		},
	}
	cmd.PersistentFlags().StringVarP(&casType, "cas-type", "", "", "the cas-type filter option for fetching resources")
	cmd.AddCommand(
		NewCmdGetVolume(),
		NewCmdGetPool(),
	)
	return cmd
}
