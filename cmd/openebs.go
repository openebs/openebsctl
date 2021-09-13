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

package cmd

import (
	"flag"
	cluster_info "github.com/openebs/openebsctl/cmd/cluster-info"
	"github.com/openebs/openebsctl/cmd/upgrade"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/openebs/openebsctl/cmd/completion"
	"github.com/openebs/openebsctl/cmd/describe"
	"github.com/openebs/openebsctl/cmd/get"
	v "github.com/openebs/openebsctl/cmd/version"
)

const (
	usageTemplate = `Usage:
  kubectl openebs [command] [resource] [...names] [flags]

Available Commands:
  completion    Outputs shell completion code for the specified shell (bash or zsh)
  describe      Provide detailed information about an OpenEBS resource
  get           Provides fetching operations related to a Volume/Pool
  help          Help about any command
  version       Shows openebs kubectl plugin's version
  cluster-info  Show component version, status and running components for each installed engine
  upgrade       Upgrade CSI Interfaces and Volumes

Flags:
  -h, --help                           help for openebs
  -n, --namespace string               If present, the namespace scope for this CLI request
      --openebs-namespace string       to read the openebs namespace from user.
                                       If not provided it is determined from components.
      --cas-type                       to specify the cas-type of the engine, for engine based filtering.
                                       ex- cstor, jiva, localpv-lvm, localpv-zfs.
      --debug                          to launch the debugging mode for cstor pvcs.

Use "kubectl openebs command --help" for more information about a command.
`
)

// Version is the version of the openebsctl binary, info filled by go-releaser
var Version = "dev"

// NewOpenebsCommand creates the `openebs` command and its nested children.
func NewOpenebsCommand() *cobra.Command {
	var openebsNs string
	cmd := &cobra.Command{
		Use:       "openebs",
		ValidArgs: []string{"get", "describe", "completion", "upgrade"},
		Short:     "openebs is a a kubectl plugin for interacting with OpenEBS storage components",
		Long: `openebs is a a kubectl plugin for interacting with OpenEBS storage components such as storage(pools, volumegroups), volumes, blockdevices, pvcs.
Find out more about OpenEBS on https://openebs.io/`,
		Version:          Version,
		TraverseChildren: true,
	}
	cmd.SetUsageTemplate(usageTemplate)
	// TODO: Check if this brings in the flags from kubectl binary to this one via exec for all platforms
	kubernetesConfigFlags := genericclioptions.NewConfigFlags(true)
	kubernetesConfigFlags.AddFlags(cmd.PersistentFlags())
	cmd.AddCommand(
		// Add a helper command to show what version of X is installed
		completion.NewCmdCompletion(cmd),
		get.NewCmdGet(cmd),
		describe.NewCmdDescribe(cmd),
		v.NewCmdVersion(cmd),
		cluster_info.NewCmdClusterInfo(cmd),
		upgrade.NewCmdVolumeUpgrade(cmd),
	)
	cmd.PersistentFlags().StringVarP(&openebsNs, "openebs-namespace", "", "", "to read the openebs namespace from user.\nIf not provided it is determined from components.")
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	_ = flag.CommandLine.Parse([]string{})
	_ = viper.BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace"))
	_ = viper.BindPFlag("openebs-namespace", cmd.PersistentFlags().Lookup("openebs-namespace"))
	return cmd
}
