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

package pool

import (
	"flag"
	"fmt"
	"time"

	"github.com/openebs/openebsctl/client"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	poolListCommandHelpText = `
	This command lists of all known pools in the Cluster.
	
	Usage: kubectl openebs cStor pool list [options]
	`
	namespace string
)

// NewCmdPoolsList displays status of OpenEBS Pool(s)
func NewCmdPoolsList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Displays status information about Pool(s)",
		Long:  poolListCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(RunPoolsList(cmd), util.Fatal)
		},
	}

	cmd.Flags().StringVarP(&namespace, "namespace", "n", "openebs",
		"namespace name, required if pool is not in the `default` namespace")

	flag.CommandLine.Parse([]string{})

	return cmd
}

//RunPoolsList fetchs & lists the pools
func RunPoolsList(cmd *cobra.Command) error {

	client, err := client.NewK8sClient(namespace)
	util.CheckErr(err, util.Fatal)

	cpools, err := client.GetcStorPools()
	if err != nil {
		return errors.Wrap(err, "error listing pools")
	}

	out := make([]string, len(cpools.Items)+2)
	out[0] = "Name|Namespace|HealthyInstances|ProvisionedInstances|DesiredInstances|Age"
	out[1] = "----|---------|----------------|--------------------|----------------|---"
	for i, item := range cpools.Items {
		out[i+2] = fmt.Sprintf("%s|%s|%d|%d|%d|%s",
			item.ObjectMeta.Name,
			item.ObjectMeta.Namespace,
			item.Status.HealthyInstances,
			item.Status.ProvisionedInstances,
			item.Status.DesiredInstances,
			util.Duration(time.Since(item.ObjectMeta.CreationTimestamp.Time)),
		)
	}
	if len(out) == 2 {
		fmt.Println("No Pools are found")
		return nil
	}

	fmt.Println(util.FormatList(out))

	return nil
}
