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
	"strconv"
	"strings"

	"github.com/openebs/openebsctl/pkg/generate"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdGenerate provides options for generating
func NewCmdGenerate() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "generate",
		Short:     "Generate one or more OpenEBS resource like cspc",
		ValidArgs: []string{"cspc"},
	}
	cmd.AddCommand(NewCmdGenerateCStorStoragePoolCluster())
	return cmd
}

// NewCmdGenerateCStorStoragePoolCluster provides options for generating cspc
// NOTE: When other custom resources need to be generated, the function
// should be renamed appropriately, as of now it made no sense to generically
// state pools when other pools aren't supported.
func NewCmdGenerateCStorStoragePoolCluster() *cobra.Command {
	var nodes, raidType, cap string
	var devices int
	cmd := &cobra.Command{
		Use:   "cspc",
		Short: "Generates cspc resources YAML/configuration which can be used to provision cStor storage pool clusters",
		Run: func(cmd *cobra.Command, args []string) {
			node, _ := cmd.Flags().GetString("nodes")
			raid, _ := cmd.Flags().GetString("raidtype")
			capacity, _ := cmd.Flags().GetString("capacity")
			devs := numDevices(cmd)
			nodeList := strings.Split(node, ",")
			util.CheckErr(generate.CSPC(nodeList, devs, raid, capacity), util.Fatal)
		},
	}
	cmd.PersistentFlags().StringVarP(&nodes, "nodes", "", "",
		"comma separated set of nodes for pool creation --nodes=node1,node2,node3,node4")
	_ = cmd.MarkPersistentFlagRequired("nodes")
	cmd.PersistentFlags().StringVarP(&raidType, "raidtype", "", "stripe",
		"allowed RAID configuration such as, stripe, mirror, raid, raidz2")
	cmd.PersistentFlags().StringVarP(&cap, "capacity", "", "10Gi",
		"allowed RAID configuration such as, stripe, mirror, raid, raidz2")
	cmd.PersistentFlags().IntVar(&devices, "number-of-devices", 1, "number of devices per node, selects default based on raid-type")
	return cmd
}

// numDevices figures out the number of devices based on the raid type
func numDevices(cmd *cobra.Command) int {
	// if number-of-devices is not set, set it to appropriate value
	if !cmd.Flag("number-of-devices").Changed {
		var devCount = map[string]int{
			"stripe": 1,
			"mirror": 2,
			"raidz":  3,
			"raidz2": 4}
		switch cmd.Flag("raidtype").Value.String() {
		case "stripe", "mirror", "raidz", "raidz2":
			c := devCount[cmd.Flag("raidtype").Value.String()]
			err := cmd.Flags().Set("number-of-devices", strconv.Itoa(c))
			if err != nil {
				return 1
			}
			return c
		}
	} else {
		d, _ := cmd.Flags().GetInt("number-of-devices")
		return d
	}
	// setting default value to 1
	return 1
}
