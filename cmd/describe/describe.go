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
	"github.com/spf13/cobra"
)

const (
	volumeCommandHelpText = `Show detailed description of a specific OpenEBS resource:

Usage:
  kubectl openebs describe [volume|storage|pvc] [...names] [flags]

Describe a Volume:
  kubectl openebs describe volume [...names] [flags]

Describe PVCs present in the same namespace:
  kubectl openebs describe pvc [...names] [flags]

Describe a Storage :
  kubectl openebs describe storage [...names] [flags]

Flags:
  -h, --help                           help for openebs
  -n, --namespace string               to read the namespace for the pvc.
      --openebs-namespace string       to read the openebs namespace from user.
                                       If not provided it is determined from components.
      --cas-type                       to specify the cas-type of the engine, for engine based filtering.
                                       ex- cstor, jiva, localpv-lvm, localpv-zfs.
      --debug                          to launch the debugging mode for cstor pvcs.
`
)

// NewCmdDescribe provides options for managing OpenEBS Volume
func NewCmdDescribe(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "describe",
		ValidArgs: []string{"pool", "volume", "pvc"},
		Short:     "Provide detailed information about an OpenEBS resource",
	}
	cmd.AddCommand(
		NewCmdDescribeVolume(),
		NewCmdDescribePVC(),
		NewCmdDescribeStorage(),
	)
	cmd.SetUsageTemplate(volumeCommandHelpText)
	return cmd
}
