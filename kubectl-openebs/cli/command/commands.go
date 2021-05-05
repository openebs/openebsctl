/*
Copyright 2020 The OpenEBS Authors

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

package command

import (
	"flag"

	"github.com/openebs/openebsctl/kubectl-openebs/cli/command/describe"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/command/get"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewOpenebsCommand creates the `openebs` command and its nested children.
func NewOpenebsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubectl openebs [get|describe] [volume|pool|...]",
		Short: "openebs is a a kubectl plugin for interacting with OpenEBS storage components",
		Long: `Find out more about OpenEBS on https://docs.openebs.io/
Usage:
# get a list of pools
kubectl openebs get pools
# get a list of volumes
kubectl openebs get volumes
# describe a volume
kubectl openebs describe volume pvc-name
`,
	}
	// TODO: Check if this brings in the flags from kubectl binary to this one via exec for all platforms
	kubernetesConfigFlags := genericclioptions.NewConfigFlags(true)
	kubernetesConfigFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		util.NewCmdCompletion(cmd),
		get.NewCmdGet(cmd),
		describe.NewCmdDescribe(cmd),
	)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	_ = flag.CommandLine.Parse([]string{})
	_ = viper.BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace"))
	return cmd
}
