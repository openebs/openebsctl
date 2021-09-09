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

package generate

import (
	"strings"

	"github.com/openebs/openebsctl/pkg/generate"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdGenerate provides options for generating
func NewCmdGenerate() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "generate",
		ValidArgs: []string{"cspc"},
		Short:     "Generate one or more OpenEBS resources based on flags",
	}
	cmd.AddCommand(NewCmdGenerateStorage())
	return cmd
}

func NewCmdGenerateStorage() *cobra.Command {
	var nodes, raidType string
	var devices int
	cmd := &cobra.Command{
		Use:   "cspc",
		Short: "Displays status information about Volume(s)",
		Run: func(cmd *cobra.Command, args []string) {
			node, _ := cmd.Flags().GetString("nodes")
			devs, _ := cmd.Flags().GetInt("number-of-devices")
			raid, _ := cmd.Flags().GetString("raidtype")
			Nodes := strings.Split(node, ",")
			util.CheckErr(generate.Pool(Nodes, devs, raid), util.Fatal)
		},
	}
	cmd.PersistentFlags().StringVarP(&nodes, "nodes", "", "",
		"comma separated set of nodes for pool creation --nodes=node1,node2,node3,node4")
	cmd.MarkPersistentFlagRequired("nodes")
	cmd.PersistentFlags().IntVar(&devices, "number-of-devices", 1, "--number-of-devices=2")
	cmd.PersistentFlags().StringVarP(&raidType, "raidtype", "", "stripe", "--raidtype=mirrored")
	return cmd
}
