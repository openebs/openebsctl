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

package storage

import (
	"fmt"
	"time"

	"github.com/docker/go-units"

	"github.com/openebs/api/v2/pkg/apis/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"
)

const cStorPoolInstanceInfoTemplate = `
{{.Name}} Details :
----------------
NAME             : {{.Name}}
HOSTNAME         : {{.HostName}}
SIZE             : {{.Size}}
FREE CAPACITY    : {{.FreeCapacity}}
READ ONLY STATUS : {{.ReadOnlyStatus}}
STATUS	         : {{.Status}}
RAID TYPE        : {{.RaidType}}
`

// GetCstorPools lists the pools
func GetCstorPools(c *client.K8sClient, pools []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error) {
	cpools, err := c.GetCSPIs(pools, "")
	if err != nil {
		return nil, nil, errors.Wrap(err, "error listing pools")
	}
	var rows []metav1.TableRow
	for _, item := range cpools.Items {
		rows = append(rows, metav1.TableRow{Cells: []interface{}{
			item.ObjectMeta.Name,
			item.ObjectMeta.Labels["kubernetes.io/hostname"],
			util.ConvertToIBytes(item.Status.Capacity.Free.String()),
			util.ConvertToIBytes(item.Status.Capacity.Total.String()),
			item.Status.ReadOnly,
			item.Status.ProvisionedReplicas,
			item.Status.HealthyReplicas,
			string(item.Status.Phase),
			util.Duration(time.Since(item.ObjectMeta.CreationTimestamp.Time))}})
	}
	if len(cpools.Items) == 0 {
		return nil, nil, fmt.Errorf("no cstor pools are found")
	}
	return util.CstorPoolListColumnDefinations, rows, nil
}

// DescribeCstorPool method runs info command and make call to DisplayPoolInfo to display the results
func DescribeCstorPool(c *client.K8sClient, poolName string) error {
	pools, err := c.GetCSPIs([]string{poolName}, "")
	if err != nil {
		return errors.Wrap(err, "error getting pool info")
	}
	if len(pools.Items) == 0 {
		return fmt.Errorf("cstor-pool %s not found", poolName)
	}
	poolInfo := pools.Items[0]
	poolDetails := util.PoolInfo{
		Name:           poolInfo.Name,
		HostName:       poolInfo.Spec.HostName,
		Size:           util.ConvertToIBytes(poolInfo.Status.Capacity.Total.String()),
		FreeCapacity:   util.ConvertToIBytes(poolInfo.Status.Capacity.Free.String()),
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

	// Printing the filled details of the Pool
	err = util.PrintByTemplate("pool", cStorPoolInstanceInfoTemplate, poolDetails)
	if err != nil {
		return err
	}

	// Fetch info for every block device
	var bdRows []metav1.TableRow
	for _, item := range BlockDevicesInPool {
		bd, err := c.GetBD(item)
		if err != nil {
			fmt.Printf("Could not find the blockdevice : %s\n", item)
		} else {
			bdRows = append(bdRows, metav1.TableRow{Cells: []interface{}{bd.Name, units.BytesSize(float64(bd.Spec.Capacity.Storage)), bd.Status.State}})
		}
	}
	if len(bdRows) != 0 {
		fmt.Printf("\nBlockdevice details :\n" + "---------------------\n")
		util.TablePrinter(util.BDListColumnDefinations, bdRows, printers.PrintOptions{Wide: true})
	} else {
		fmt.Printf("Could not find any blockdevice that belongs to the pool\n")
	}

	// Fetch info for provisional replica
	var cvrRows []metav1.TableRow
	CVRsInPool, err := c.GetCVRs(types.CStorPoolInstanceNameLabelKey + "=" + poolName)
	if err != nil {
		fmt.Printf("None of the replicas are running")
	} else {
		for _, cvr := range CVRsInPool.Items {
			pvcName := ""
			pv, err := c.GetPV(cvr.Labels["openebs.io/persistent-volume"])
			if err == nil {
				pvcName = pv.Spec.ClaimRef.Name
			}
			cvrRows = append(cvrRows, metav1.TableRow{Cells: []interface{}{
				cvr.Name,
				pvcName,
				util.ConvertToIBytes(cvr.Status.Capacity.Total),
				cvr.Status.Phase}})
		}
	}
	if len(cvrRows) != 0 {
		fmt.Printf("\nReplica Details :\n-----------------\n")
		util.TablePrinter(util.PoolReplicaColumnDefinations, cvrRows, printers.PrintOptions{Wide: true})
	}
	return nil
}
