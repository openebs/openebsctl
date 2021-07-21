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

package volume

import (
	"fmt"
	"os"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetZFSLocalPVs returns a list of ZFSVolumes
func GetZFSLocalPVs(c *client.K8sClient, pvList *corev1.PersistentVolumeList, openebsNS string) ([]metav1.TableRow, error) {
	// 1. Fetch all relevant volume CRs without worrying about openebsNS
	_, zvolMap, err := c.GetZFSVols(nil, util.Map, "", util.MapOptions{Key: util.Name})
	if err != nil {
		return nil, fmt.Errorf("failed to list ZFSVolumes")
	}
	var version string
	if CSIctrl, err := c.GetCSIControllerSTS("openebs-zfs-controller"); err == nil {
		version = CSIctrl.Labels["openebs.io/version"]
	}
	if version == "" {
		version = "N/A"
	}

	var rows []metav1.TableRow
	// 3. Show the required ones
	for _, pv := range pvList.Items {
		name := pv.Name
		capacity := util.ConvertToIBytes(pv.Spec.Capacity.Storage().String())
		sc := pv.Spec.StorageClassName
		attached := pv.Status.Phase
		var attachedNode, customStatus, ns string
		// TODO: Estimate the cas-type and decide to print it out
		// Should all AccessModes be shown in a csv format, or the highest be displayed ROO < RWO < RWX?
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == util.ZFSCSIDriver {
			zvol, ok := zvolMap[pv.Name]
			if !ok {
				_, _ = fmt.Fprintln(os.Stderr, "couldn't find zfs localpv volume "+pv.Name)
			}
			ns = zvol.Namespace
			if openebsNS != "" && openebsNS != ns {
				continue
			}
			customStatus = zvol.Status.State // RW, RO, etc
			attachedNode = zvol.Labels["kubernetes.io/nodename"]
		} else {
			// Skip non-ZFS options
			continue
		}
		accessMode := pv.Spec.AccessModes[0]
		rows = append(rows, metav1.TableRow{
			Cells: []interface{}{
				ns, name, customStatus, version, capacity, sc, attached,
				accessMode, attachedNode},
		})
	}
	return rows, nil
}
