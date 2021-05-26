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

package describe

import (
	"github.com/openebs/openebsctl/client"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	poolInfoCommandHelpText = `
This command fetches information and status of the various aspects 
of the cStor Pool Instance and its underlying related resources in the provided namespace.
If no namespace is provided it uses default namespace for execution.
$ kubectl openebs describe pool [cspi-name] -n [namespace]
`
)

const (
	cStorPoolInstanceInfoTemplate = `
{{.Name}} Details :
----------------
Name             : {{.Name}}
Hostname         : {{.HostName}}
Size             : {{.Size}}
Free Capacity    : {{.FreeCapacity}}
Read Only Status : {{.ReadOnlyStatus}}
Status	         : {{.Status}}
RAID Type        : {{.RaidType}}

`

	blockDevicesInfoFromCSPI = `
Block Device Details :
----------------
Name     : {{.Name}}
Capacity : {{.Capacity}}
State   : {{.State}}

`

	provisionedReplicasInfoFromCSPI = `
Provisioned Replicas Details :
----------------
Name     : {{.Name}}
PVC Name : {{.PvcName}}
Size     : {{.Size}}
Status	 : {{.Status}}

`
)

// NewCmdDescribePool displays OpenEBS cStor pool instance information.
func NewCmdDescribePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pool",
		Aliases: []string{"pools", "p"},
		Short:   "Displays cStorPoolInstance information",
		Long:    poolInfoCommandHelpText,
		Example: `kubectl openebs describe pool cspi-one -n openebs`,
		Run: func(cmd *cobra.Command, args []string) {
			var namespace string // This namespace belongs to the CSPI entered
			if namespace, _ = cmd.Flags().GetString("namespace"); namespace == "" {
				// NOTE: The error comes as nil even when the ns flag is not specified
				namespace = "openebs"
			}
			util.CheckErr(RunPoolInfo(cmd, args, namespace), util.Fatal)
		},
	}
	return cmd
}

// RunPoolInfo method runs info command and make call to DisplayPoolInfo to display the results
func RunPoolInfo(cmd *cobra.Command, pools []string, ns string) error {
	if len(pools) != 1 {
		return errors.New("Please give one cspi name to describe")
	}

	clientset, err := client.NewK8sClient(ns)
	if err != nil {
		return errors.Wrap(err, "Failed to execute pool info command")
	}

	// Fetch the CSPI object by passing the name of CSPI taken through CLI in ns namespace
	poolName := pools[0]
	poolInfo, err := clientset.GetcStorPool(poolName)
	if err != nil {
		return errors.Wrap(err, "Error getting pool info")
	}

	poolDetails := util.PoolInfo{
		Name:           poolInfo.Name,
		HostName:       poolInfo.Spec.HostName,
		Size:           poolInfo.Status.Capacity.Total.String(),
		FreeCapacity:   poolInfo.Status.Capacity.Free.String(),
		ReadOnlyStatus: poolInfo.Status.ReadOnly,
		Status:         poolInfo.Status.Phase,
		RaidType:       poolInfo.Spec.PoolConfig.DataRaidGroupType,
	}
	// Fetch all the raid groups in the CSPI
	RaidGroupsInPool := poolInfo.GetAllRaidGroups()

	// Fetch all the block devices in the raid groups associated to the CSPI
	var BlockDevicesInPool []string
	for _, item := range RaidGroupsInPool {
		BlockDevicesInPool = append(BlockDevicesInPool, item.GetBlockDevices()...)
	}

	// Fetch info for every block device
	var BdInfo []util.BlockDevicesInfoInPool
	for _, item := range BlockDevicesInPool {
		bd, err := clientset.GetBlockDevice(item)
		if err != nil {
			return errors.Wrap(err, "Error getting block device info")
		}

		BdInfo = append(BdInfo, util.BlockDevicesInfoInPool{
			Name:     bd.Name,
			Capacity: bd.Spec.Capacity.Storage,
			State:    bd.Status.State,
		})
	}

	// Fetch info for provisional replica
	CVRsInPool, err := clientset.GetCVRByPoolName(poolName)
	if err != nil {
		return errors.Wrap(err, "Error getting block device info")
	}

	var CVRInfoInPool []util.CVRInfo
	for _, cvr := range CVRsInPool.Items {
		CVRInfoInPool = append(CVRInfoInPool, util.CVRInfo{
			Name:    cvr.Name,
			PvcName: clientset.GetPVCNameByCVR(cvr.Labels["openebs.io/persistent-volume"]),
			Size:    cvr.Status.Capacity.Total,
			Status:  cvr.Status.Phase,
		})
	}

	// Printing the filled details of the Pool
	err = util.PrintByTemplate("pool", cStorPoolInstanceInfoTemplate, poolDetails)
	if err != nil {
		return err
	}

	for _, bd := range BdInfo {
		err = util.PrintByTemplate("bd", blockDevicesInfoFromCSPI, bd)
		if err != nil {
			return err
		}
	}

	for _, cvr := range CVRInfoInPool {
		err = util.PrintByTemplate("cvr", provisionedReplicasInfoFromCSPI, cvr)
		if err != nil {
			return err
		}
	}

	return nil
}
