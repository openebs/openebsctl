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

package upgrade

import (
	"github.com/openebs/openebsctl/pkg/upgrade/status"
	"github.com/spf13/cobra"	
)

func NewCmdUpgradeStatus() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"Status"},
		Short:   "Display Upgrade-status for a running upgrade-job",
		Run: func(cmd *cobra.Command, args []string) {
			status.GetJobStatus()
		},
	}
	cmd.PersistentFlags().BoolVar(&status.WaitFlag, "wait", false, "Wait for the logs stream")

	return cmd
}
