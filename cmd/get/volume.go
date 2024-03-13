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

package get

import (
	"fmt"

	"github.com/openebs/openebsctl/pkg/util"
	"github.com/openebs/openebsctl/pkg/volume"
	"github.com/spf13/cobra"
)

// NewCmdGetVolume displays status of OpenEBS Volume(s)
func NewCmdGetVolume() *cobra.Command {
	var openebsNs string
	var casType string
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"vol", "v", "volumes"},
		Short:   "Displays status information about Volume(s)",
		Run: func(cmd *cobra.Command, args []string) {
			openebsNS, _ := cmd.Flags().GetString("openebs-namespace")
			casType, _ := cmd.Flags().GetString("cas-type")
			util.CheckErr(volume.Get(args, openebsNS, casType), util.Fatal)
		},
	}
	cmd.PersistentFlags().StringVarP(&openebsNs, "openebs-namespace", "", "", "to read the openebs namespace from user.\nIf not provided it is determined from components.")
	cmd.PersistentFlags().StringVarP(&casType, "cas-type", "", "", fmt.Sprintf("the type of the engine %s, %s", util.LVMCasType, util.ZFSCasType))
	return cmd
}
