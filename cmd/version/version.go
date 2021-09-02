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
	"os"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	versionCmdHelp = `Usage:
  kubectl openebs version
Flags:
  -h, --help                           help for openebs get command
`
)

// Get versions of components, return "Not Installed" on empty version
func getValidVersion(version string) string {
	if version != "" {
		return version
	}

	return "Not Installed"
}

// NewCmdVersion shows OpenEBSCTL version
func NewCmdVersion(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Shows openebs kubectl plugin's version",
		Run: func(cmd *cobra.Command, args []string) {
			k, _ := client.NewK8sClient("")
			componentVersionMap, err := k.GetVersionMapOfComponents()

			if err != nil {
				fmt.Println("Client Version: " + getValidVersion(rootCmd.Version))
				fmt.Fprintf(os.Stderr, "\nError getting Components Version...")
				return
			}

			var rows []metav1.TableRow = []metav1.TableRow{
				{
					Cells: []interface{}{"Client", getValidVersion(rootCmd.Version)},
				},
				{
					Cells: []interface{}{"OpenEBS CStor", getValidVersion(componentVersionMap[util.CstorCasType])},
				},
				{
					Cells: []interface{}{"OpenEBS Jiva", getValidVersion(componentVersionMap[util.JivaCasType])},
				},
				{
					Cells: []interface{}{"OpenEBS LVM LocalPV", getValidVersion(componentVersionMap[util.LVMCasType])},
				},
				{
					Cells: []interface{}{"OpenEBS ZFS LocalPV", getValidVersion(componentVersionMap[util.ZFSCasType])},
				},
			}

			util.TablePrinter(util.VersionColumnDefinition, rows, printers.PrintOptions{Wide: true})
		},
	}
	cmd.SetUsageTemplate(versionCmdHelp)
	return cmd
}
