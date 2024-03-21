/*
Copyright 2020-2022 The OpenEBS Authors

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
	"github.com/openebs/openebsctl/pkg/persistentvolumeclaim"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdDescribePVC Displays the pvc describe details
func NewCmdDescribePVC() *cobra.Command {
	var openebsNs string
	var pvNs string
	cmd := &cobra.Command{
		Use:     "pvc",
		Aliases: []string{"pvcs", "persistentvolumeclaims", "persistentvolumeclaim"},
		Short:   "Displays PersistentVolumeClaim information",
		Run: func(cmd *cobra.Command, args []string) {
			var pvNs, openebsNamespace string
			if pvNs, _ = cmd.Flags().GetString("namespace"); pvNs == "" {
				pvNs = "default"
			}
			openebsNamespace, _ = cmd.Flags().GetString("openebs-namespace")
			util.CheckErr(persistentvolumeclaim.Describe(args, pvNs, openebsNamespace), util.Fatal)
		},
	}
	cmd.PersistentFlags().StringVarP(&openebsNs, "openebs-namespace", "", "", "to read the openebs namespace from user.\nIf not provided it is determined from components.")
	cmd.PersistentFlags().StringVarP(&pvNs, "namespace", "n", "", "to read the namespace of the pvc from the user. If not provided defaults to default namespace.")
	return cmd
}
