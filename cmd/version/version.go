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
		Short: "Show client version",
		Run: func(cmd *cobra.Command, args []string) {
			k, _ := client.NewK8sClient("openebs")
			nsMap, err := k.GetVersionMapOfComponents()

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
					Cells: []interface{}{"OpenEBS", getValidVersion(nsMap[util.OpenEBSProvisioner])},
				},
				{
					Cells: []interface{}{"OpenEBS CStor", getValidVersion(nsMap[util.CstorCasType])},
				},
				{
					Cells: []interface{}{"OpenEBS Jiva", getValidVersion(nsMap[util.JivaCasType])},
				},
				{
					Cells: []interface{}{"OpenEBS ZFS LocalPV", getValidVersion(nsMap[util.ZFSCasType])},
				},
			}

			util.TablePrinter(util.VersionColumnDefinition, rows, printers.PrintOptions{Wide: true})
		},
	}
	return cmd
}
