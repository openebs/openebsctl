/*
Copyright 2020 The OpenEBS Authors

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

package pool

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	poolCommandHelpText = `# List Pools:
	$ kubectl openebs cStor pool list
 `
)

// NewCmdPool provides options for managing OpenEBS Pool
func NewCmdPool(rootCmd *cobra.Command) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "pool",
		Short: "Provides operations related to a Pool",
		Long:  poolCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(poolCommandHelpText)
		},
	}

	cmd.AddCommand(
		NewCmdPoolsList(),
	)

	return cmd
}
