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

// GetLVMLocalPV returns a list of LVM-LocalPV volumes
func GetLVMLocalPV(c *client.K8sClient, pvList *corev1.PersistentVolumeList, openebsNS string) ([]metav1.TableRow, error) {
	var rows []metav1.TableRow
	var version string
	if CSIctrl, err := c.GetCSIControllerSTS("openebs-lvm-controller"); err == nil {
		version = CSIctrl.Labels["openebs.io/version"]
	}
	if version == "" {
		version = "N/A"
	}
	for _, pv := range pvList.Items {
		var attachedNode, customStatus, ns string
		_, lvmVolMap, err := c.GetLVMvol(nil, util.Map, "", util.MapOptions{Key: util.Name})
		if err != nil {
			return nil, fmt.Errorf("failed to list LVMVolumes")
		}
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == util.LocalPVLVMCSIDriver {
			lvmVol, ok := lvmVolMap[pv.Name]
			if !ok {
				// condition not possible
				_, _ = fmt.Fprintf(os.Stderr, "couldn't find LVM volume "+pv.Name)
			}
			ns = lvmVol.Namespace
			if openebsNS != "" && openebsNS != ns {
				continue
			}
			accessMode := pv.Spec.AccessModes[0]
			customStatus = lvmVol.Status.State
			attachedNode = lvmVol.Spec.OwnerNodeID
			rows = append(rows, metav1.TableRow{
				Cells: []interface{}{
					ns, pv.Name, customStatus, version, pv.Spec.Capacity.Storage().String(), pv.Spec.StorageClassName, pv.Status.Phase,
					accessMode, attachedNode}})
		}
	}
	return rows, nil
}
