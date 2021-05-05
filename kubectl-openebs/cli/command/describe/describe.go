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
	volumeCommandHelpText = `# List Volumes:
$ kubectl openebs describe [volumes|pools] [name]

# Info of a Volume:
$ kubectl openebs describe volume pvc-abcd -n [namespace]`
)

// NewCmdDescribe provides options for managing OpenEBS Volume
func NewCmdDescribe(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Provides operations related to a Volume",
		Long:  volumeCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(volumeCommandHelpText)
		},
	}

	cmd.AddCommand(
		NewCmdVolumeInfo(),
		// TODO: Add Pool Info command
	)

	return cmd
}
