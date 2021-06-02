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
	"fmt"

	"github.com/spf13/cobra"
)

const (
	volumeCommandHelpText = `# Show detail of a specific OpenEBS resource:
$ kubectl openebs describe [volumes|pools|pvc] [name]

# Describe a Volume:
$ kubectl openebs describe volume pvc-abcd -n [namespace]

# Describe PVCs present in the same namespace:
$ kubectl openebs describe pvc [name1] [name2] ... [nameN] -n [namespace]

# Describe a cStor Pool Instance:
$ kubectl openebs describe pool [cspi-name] -n [namespace]
`
)

// NewCmdDescribe provides options for managing OpenEBS Volume
func NewCmdDescribe(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Provide detailed information about an OpenEBS resource",
		Long:  volumeCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(volumeCommandHelpText)
		},
	}
	cmd.AddCommand(
		NewCmdDescribeVolume(),
		NewCmdDescribePVC(),
		NewCmdDescribePool(),
	)

	return cmd
}
