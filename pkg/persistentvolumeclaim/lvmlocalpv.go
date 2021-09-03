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

package persistentvolumeclaim

import (
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/openebs/openebsctl/pkg/volume"
	corev1 "k8s.io/api/core/v1"
)

const (
	lvmPvcInfoTemplate = `
{{.Name}} Details  :
-------------------
NAME               : {{.Name}}
NAMESPACE          : {{.Namespace}}
CAS TYPE           : {{.CasType}}
BOUND VOLUME       : {{.BoundVolume}}
STORAGE CLASS      : {{.StorageClassName}}
SIZE               : {{.Size}}
PVC STATUS         : {{.PVCStatus}}
`
)

// DescribeLVMVolumeClaim describes a LVM storage engine PersistentVolumeClaim
func DescribeLVMVolumeClaim(c *client.K8sClient, pvc *corev1.PersistentVolumeClaim, pv *corev1.PersistentVolume) error {
	// 1. Fill in the PVC information
	lvmPVCinfo := util.LVMPVCInfo{
		Name:             pvc.Name,
		Namespace:        pvc.Namespace,
		CasType:          util.LVMCasType,
		BoundVolume:      pvc.Spec.VolumeName,
		StorageClassName: *pvc.Spec.StorageClassName,
		Size:             pvc.Spec.Resources.Requests.Storage().String(),
		PVCStatus:        pvc.Status.Phase,
	}

	// 2. If PV is present Describe the LVM Volume
	if pv != nil {
		_ = util.PrintByTemplate("lvmPvc", lvmPvcInfoTemplate, lvmPVCinfo)
		volume.DescribeLVMLocalPVs(c, pv)
	} else {
		// Show only PVC details if volume is not found.
		_ = util.PrintByTemplate("lvmPvc", lvmPvcInfoTemplate, lvmPVCinfo)
	}

	return nil
}
