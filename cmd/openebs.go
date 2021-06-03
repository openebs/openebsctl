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

	"github.com/openebs/openebsctl/cmd/completion"
	"github.com/openebs/openebsctl/cmd/describe"
	"github.com/openebs/openebsctl/cmd/get"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// Version is the version of the openebsctl binary, info filled by goreleaser
var Version = "dev"

// NewOpenebsCommand creates the `openebs` command and its nested children.
func NewOpenebsCommand() *cobra.Command {
	var openebsNs string
	cmd := &cobra.Command{
		Use:       "openebs",
		ValidArgs: []string{"get", "describe", "completion"},
		Short:     "openebs is a a kubectl plugin for interacting with OpenEBS storage components",
		Long: `openebs is a a kubectl plugin for interacting with OpenEBS storage components
Find out more about OpenEBS on https://docs.openebs.io/`,
		Version: Version,
	}
	// TODO: Check if this brings in the flags from kubectl binary to this one via exec for all platforms
	kubernetesConfigFlags := genericclioptions.NewConfigFlags(true)
	kubernetesConfigFlags.AddFlags(cmd.PersistentFlags())
	cmd.AddCommand(
		// Add a helper command to show what version of X is installed
		completion.NewCmdCompletion(cmd),
		get.NewCmdGet(cmd),
		describe.NewCmdDescribe(cmd),
		v.NewCmdVersion(cmd),
	)
	cmd.PersistentFlags().StringVarP(&openebsNs, "openebs-namespace", "", "", "to read the openebs namespace from user.\nIf not provided it is determined from components.")
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	_ = flag.CommandLine.Parse([]string{})
	_ = viper.BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace"))
	return cmd
}
