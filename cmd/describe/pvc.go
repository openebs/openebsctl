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

var (
	pvcInfoCommandHelpText = `This command fetches information and status  of  the  various  aspects 
of  the  PersistentVolumeClaims  and  its underlying related resources 
in the provided namespace. If no namespace is provided it uses default
namespace for execution.

Usage:
  kubectl openebs describe pvc [...names] [flags]

Flags:
  -h, --help                           help for openebs
  -n, --namespace string               to read the namespace for the pvc.
      --openebs-namespace string       to read the openebs namespace from user.
                                       If not provided it is determined from components.
      --debug                          to launch the debugging mode for cstor pvcs.
`
)

// NewCmdDescribePVC Displays the pvc describe details
func NewCmdDescribePVC() *cobra.Command {
	var debug bool
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
			if debug {
				util.CheckErr(persistentvolumeclaim.Debug(args, pvNs, openebsNamespace), util.Fatal)
			} else {
				util.CheckErr(persistentvolumeclaim.Describe(args, pvNs, openebsNamespace), util.Fatal)
			}

		},
	}
	cmd.SetUsageTemplate(pvcInfoCommandHelpText)
	cmd.Flags().BoolVar(&debug, "debug", false, "Debug cstor volume")
	return cmd
}
