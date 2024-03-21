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

package cmd

import (
	"flag"

	"github.com/openebs/openebsctl/cmd/clusterinfo"
	"github.com/openebs/openebsctl/cmd/completion"
	"github.com/openebs/openebsctl/cmd/describe"
	"github.com/openebs/openebsctl/cmd/get"
	v "github.com/openebs/openebsctl/cmd/version"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version is the version of the openebsctl binary, info filled by go-releaser
var Version = "dev"

// NewOpenebsCommand creates the `openebs` command and its nested children.
func NewOpenebsCommand() *cobra.Command {
	//var openebsNs string
	cmd := &cobra.Command{
		Use:       "openebs",
		ValidArgs: []string{"get", "describe", "completion"},
		Short:     "kubectl openebs is a a kubectl plugin for interacting with OpenEBS storage components",
		Long: `openebs is a a kubectl plugin for interacting with OpenEBS storage components such as storage(zfspools, volumegroups), volumes, pvcs.
Find out more about OpenEBS on https://openebs.io/`,
		Version:          Version,
		TraverseChildren: true,
	}
	cmd.AddCommand(
		completion.NewCmdCompletion(cmd),
		get.NewCmdGet(cmd),
		describe.NewCmdDescribe(cmd),
		v.NewCmdVersion(cmd),
		clusterinfo.NewCmdClusterInfo(cmd),
	)
	cmd.PersistentFlags().StringVarP(&util.Kubeconfig, "kubeconfig", "c", "", "path to config file")
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	_ = flag.CommandLine.Parse([]string{})
	_ = viper.BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace"))
	_ = viper.BindPFlag("openebs-namespace", cmd.PersistentFlags().Lookup("openebs-namespace"))
	return cmd
}
