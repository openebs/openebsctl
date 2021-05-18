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

package get

import (
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

Usage:
$ kubectl openebs get pool [options]
`
)

// NewCmdGetPool displays status of OpenEBS Pool(s)
func NewCmdGetPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Aliases: []string{"pools", "p"},
		Short:   "Displays status information about Pool(s)",
		Long:    poolListCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			ns, err := cmd.Flags().GetString("namespace")
			if err != nil {
				ns = "openebs"
			}
			// TODO: De-couple CLI code, logic code, API code
			util.CheckErr(RunPoolsList(cmd, args, ns), util.Fatal)
		},
	}
	return cmd
}

//RunPoolsList fetchs & lists the pools
func RunPoolsList(cmd *cobra.Command, pools []string, ns string) error {
	client, err := client.NewK8sClient(ns)
	util.CheckErr(err, util.Fatal)
	cpools, err := client.GetcStorPools()
	if err != nil {
		return errors.Wrap(err, "error listing pools")
	}

	out := make([]string, len(cpools.Items)+2)
	out[0] = "Name|Namespace|HostName|Free|Capacity|ReadOnly|ProvisionedReplicas|HealthyReplicas|Status|Age"
	out[1] = "----|---------|--------|----|--------|--------|-------------------|---------------|------|---"
	for i, item := range cpools.Items {
		out[i+2] = fmt.Sprintf("%s|%s|%s|%s|%s|%v|%d|%d|%s|%s",
			item.ObjectMeta.Name,
			item.ObjectMeta.Namespace,
			item.ObjectMeta.Labels["kubernetes.io/hostname"],
			item.Status.Capacity.Free.String(),
			item.Status.Capacity.Total.String(),
			item.Status.ReadOnly,
			item.Status.ProvisionedReplicas,
			item.Status.HealthyReplicas,
			item.Status.Phase,
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
