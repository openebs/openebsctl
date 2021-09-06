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

package cluster_info

import (
	cluster_info "github.com/openebs/openebsctl/pkg/cluster-info"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
)

const (
	clusterInfoCmdHelp = `Usage:
  kubectl openebs cluster-info
Flags:
  -h, --help                           help for openebs get command
`
)

// NewCmdClusterInfo shows OpenEBSCTL cluster-info
func NewCmdClusterInfo(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster-info",
		Short: "Show component version, status and running components for each installed engine",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(cluster_info.ShowClusterInfo(), util.Fatal)
		},
	}
	cmd.SetUsageTemplate(clusterInfoCmdHelp)
	return cmd
}
