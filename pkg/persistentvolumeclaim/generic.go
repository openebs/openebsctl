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
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
)

const (
	genericPvcInfoTemplate = `
{{.Name}} Details :
------------------
NAME             : {{.Name}}
NAMESPACE        : {{.Namespace}}
CAS TYPE         : {{.CasType}}
BOUND VOLUME     : {{.BoundVolume}}
STORAGE CLASS    : {{.StorageClassName}}
SIZE             : {{.Size}}
PV STATUS    	 : {{.PVStatus}}
`
)

// DescribeGenericVolumeClaim describes a any PersistentVolumeClaim, if the cas type is not discovered.
func DescribeGenericVolumeClaim(pvc *corev1.PersistentVolumeClaim, pv *corev1.PersistentVolume, casType string) error {
	// Incase a not known casType pvc is entered show minimal details pertaining to the PVC
	// 1. Fill in the PVC details.
	pvcInfo := util.PVCInfo{}
	pvcInfo.Name = pvc.Name
	pvcInfo.Namespace = pvc.Namespace
	pvcInfo.StorageClassName = *pvc.Spec.StorageClassName
	quantity := pvc.Status.Capacity[util.StorageKey]
	pvcInfo.Size = util.ConvertToIBytes(quantity.String())
	if pv != nil {
		pvcInfo.BoundVolume = pvc.Spec.VolumeName
		pvcInfo.PVStatus = pv.Status.Phase
	}
	pvcInfo.CasType = casType
	// 2. Print the details
	err := util.PrintByTemplate("pvc", genericPvcInfoTemplate, pvcInfo)
	if err != nil {
		return err
	}
	return nil
}
