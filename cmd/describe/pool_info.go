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
	"fmt"

	"github.com/openebs/api/v2/pkg/apis/types"

	"github.com/dustin/go-humanize"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
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
NAME             : {{.Name}}
HOSTNAME         : {{.HostName}}
SIZE             : {{.Size}}
FREE CAPACITY    : {{.FreeCapacity}}
READ ONLY STATUS : {{.ReadOnlyStatus}}
STATUS	         : {{.Status}}
RAID TYPE        : {{.RaidType}}

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
			openebsNs, _ := cmd.Flags().GetString("openebs-namespace")
			util.CheckErr(RunPoolInfo(cmd, args, openebsNs), util.Fatal)
		},
	}
	return cmd
}

// RunPoolInfo method runs info command and make call to DisplayPoolInfo to display the results
func RunPoolInfo(cmd *cobra.Command, pools []string, openebsNs string) error {
	if len(pools) != 1 {
		return errors.New("Please give one cspi name to describe")
	}
	clientset, err := client.NewK8sClient(openebsNs)
	if err != nil {
		return errors.Wrap(err, "Failed to execute pool info command")
	}

	if openebsNs == "" {
		nsFromCli, err := clientset.GetOpenEBSNamespace(util.CstorCasType)
		if err != nil {
			//return errors.Wrap(err, "Error determining the openebs namespace, please specify using \"--openebs-namespace\" flag")
			return errors.New("no cstor pools found in the cluster")
		}
		clientset.Ns = nsFromCli
	}

	// Fetch the CSPI object by passing the name of CSPI taken through CLI in ns namespace
	poolName := pools[0]
	poolInfo, err := clientset.GetCSPI(poolName)
	if err != nil {
		return errors.Wrap(err, "Error getting pool info")
	}

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
		bd, err := clientset.GetBD(item)
		if err != nil {
			fmt.Printf("Could not find the blockdevice : %s\n", item)
		} else {
			bdRows = append(bdRows, metav1.TableRow{Cells: []interface{}{bd.Name, humanize.IBytes(bd.Spec.Capacity.Storage), bd.Status.State}})
		}
	}
	if len(bdRows) != 0 {
		fmt.Printf("Blockdevice details :\n" + "---------------------\n")
		util.TablePrinter(util.BDListColumnDefinations, bdRows, printers.PrintOptions{Wide: true})
	} else {
		fmt.Printf("Could not find any blockdevice that belongs to the pool\n")
	}

	// Fetch info for provisional replica
	var cvrRows []metav1.TableRow
	CVRsInPool, err := clientset.GetCVRs(types.CStorPoolInstanceNameLabelKey + "=" + poolName)
	if err != nil {
		fmt.Printf("None of the replicas are running")
	} else {
		for _, cvr := range CVRsInPool.Items {
			pvcName := ""
			pv, err := clientset.GetPV(cvr.Labels["openebs.io/persistent-volume"])
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
