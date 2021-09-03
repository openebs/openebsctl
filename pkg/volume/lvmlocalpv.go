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

const lvmVolInfo = `
{{.Name}} Details :
------------------
NAME              : {{.Name}}
NAMESPACE         : {{.Namespace}}
ACCESS MODE       : {{.AccessMode}}
CSI DRIVER        : {{.CSIDriver}}
CAPACITY          : {{.Capacity}}
PVC NAME          : {{.PVC}}
VOLUME PHASE      : {{.VolumePhase}}
STORAGE CLASS     : {{.StorageClass}}
VERSION           : {{.Version}}
LVM VOLUME STATUS : {{.Status}}
VOLUME GROUP      : {{.VolumeGroup}}
SHARED            : {{.Shared}}
THIN PROVISIONED  : {{.ThinProvisioned}}
NODE ID           : {{.NodeID}}   
`

// GetLVMLocalPV returns a list of LVM-LocalPV volumes
func GetLVMLocalPV(c *client.K8sClient, pvList *corev1.PersistentVolumeList, openebsNS string) ([]metav1.TableRow, error) {
	var rows []metav1.TableRow
	var version string
	if CSIctrl, err := c.GetCSIControllerSTS(util.LVMLocalPVcsiControllerLabelValue); err == nil {
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
					ns, pv.Name, customStatus, version, util.ConvertToIBytes(pv.Spec.Capacity.Storage().String()), pv.Spec.StorageClassName, pv.Status.Phase,
					accessMode, attachedNode}})
		}
	}
	return rows, nil
}

// DescribeLVMLocalPVs describes a single lvm-localpv volume
func DescribeLVMLocalPVs(c *client.K8sClient, vol *corev1.PersistentVolume) error {
	if vol == nil {
		return fmt.Errorf("LVM volume nil")
	}
	// 1. Fetch the version from the CSI Controller STS labels
	var version string
	if CSIctrl, err := c.GetCSIControllerSTS(util.LVMLocalPVcsiControllerLabelValue); err == nil {
		version = CSIctrl.Labels["openebs.io/version"]
	}
	// Assign N/A if not found
	if version == "" {
		version = "N/A"
	}
	// 2. Fill the details using the Persistent Volume
	v := util.LVMVolDesc{
		AccessMode: util.AccessModeToString(vol.Spec.AccessModes),
		Capacity:   vol.Spec.Capacity.Storage().String(),
		CSIDriver:  vol.Spec.CSI.Driver,
		Name:       vol.Name,
		// assuming that LVMPVs aren't static-ally provisioned
		PVC:          vol.Spec.ClaimRef.Name,
		VolumePhase:  vol.Status.Phase,
		StorageClass: vol.Spec.StorageClassName,
		Version:      version,
	}

	// 3. Fetch the corresponding LVM Volume CR and fill in the other details
	printErr := false
	lVols, _, err := c.GetLVMvol([]string{vol.Name}, util.List, "", util.MapOptions{})
	if err != nil || len(lVols.Items) == 0 {
		// For printing the error at last, otherwise this would come in between two sections in PVC describe
		printErr = true
	} else {
		lVol := lVols.Items[0]
		v.Namespace = lVol.Namespace
		v.Status = lVol.Status.State
		v.VolumeGroup = lVol.Spec.VolGroup
		v.Shared = lVol.Spec.Shared
		v.ThinProvisioned = lVol.Spec.ThinProvision
		v.NodeID = lVol.Spec.OwnerNodeID
	}
	// 4. Print the data
	_ = util.PrintByTemplate("volume", lvmVolInfo, v)
	if printErr {
		// 5. Print the error is any
		fmt.Println()
		fmt.Fprintf(os.Stderr, "The LVMVol for %s doesnot exist", vol.Name)
		fmt.Println()
	}
	return nil
}
