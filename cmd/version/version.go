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

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

func getValidVersions(version string) string {
	if version != "" {
		return version
	}

	return "Not Installed"
}

// NewCmdVersion shows OpenEBSCTL version
func NewCmdVersion(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show client version",
		Run: func(cmd *cobra.Command, args []string) {
			k, _ := client.NewK8sClient("openebs")
			nsMap, err := k.GetVersionMapOfComponents()

			if err != nil {
				fmt.Println("Client Version: " + rootCmd.Version)
				fmt.Println("\nError getting Components Version")
				return
			}

			header := []string{"Components", "Version"}
			rows := [][]string{
				{"Client", rootCmd.Version},
				{"OpenEBS", getValidVersions(nsMap[util.OpenEBSProvisioner])},
				{"OpenEBS CStor", getValidVersions(nsMap[util.CstorCasType])},
				{"OpenEBS Jiva", getValidVersions(nsMap[util.JivaCasType])},
				{"OpenEBS ZFS LocalPV", getValidVersions(nsMap[util.ZFSCasType])},
			}

			util.PrintToTable(header, rows)
		},
	}
	return cmd
}
