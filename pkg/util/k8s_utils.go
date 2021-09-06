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

package util

import (
	"strconv"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
)

// GetUsedCapacityFromCVR as the healthy replicas would have the correct used capacity details
func GetUsedCapacityFromCVR(cvrList *cstorv1.CStorVolumeReplicaList) string {
	for _, item := range cvrList.Items {
		if item.Status.Phase == Healthy {
			return item.Status.Capacity.Used
		}
	}
	return ""
}

// GetCasType from the v1pv and v1sc, this is a fallback checker method, it checks
// both the resource only if the castype is not found.
func GetCasType(v1PV *corev1.PersistentVolume, v1SC *v1.StorageClass) string {
	if v1PV != nil {
		if val := GetCasTypeFromPV(v1PV); val != Unknown {
			return val
		}
	}
	if v1SC != nil {
		if val := GetCasTypeFromSC(v1SC); val != Unknown {
			return val
		}
	}
	return Unknown
}

// GetCasTypeFromPV from the passed PersistentVolume or the Stora
func GetCasTypeFromPV(v1PV *corev1.PersistentVolume) string {
	if v1PV != nil {
		if v1PV.ObjectMeta.Labels != nil {
			if val, ok := v1PV.ObjectMeta.Labels[OpenEBSCasTypeKey]; ok {
				return val
			}
		}
		if v1PV.ObjectMeta.Annotations != nil {
			if val, ok := v1PV.ObjectMeta.Annotations[OpenEBSCasTypeKey]; ok {
				return val
			}
		}
		if v1PV.Spec.CSI != nil && v1PV.Spec.CSI.VolumeAttributes != nil {
			if val, ok := v1PV.Spec.CSI.VolumeAttributes[OpenEBSCasTypeKey]; ok {
				return val
			}
		}
		if v1PV.Spec.CSI != nil {
			if val, ok := ProvsionerAndCasTypeMap[v1PV.Spec.CSI.Driver]; ok {
				return val
			}
		}
	}
	return Unknown
}

// GetCasTypeFromSC by passing the storage class
func GetCasTypeFromSC(v1SC *v1.StorageClass) string {
	if v1SC != nil {
		if v1SC.Parameters != nil {
			if val, ok := v1SC.Parameters[OpenEBSCasTypeKeySc]; ok {
				return val
			}
		}
		if val, ok := ProvsionerAndCasTypeMap[v1SC.Provisioner]; ok {
			return val
		}
	}
	return Unknown
}

// GetReadyContainers to show the number of ready vs total containers of pod i.e 2/3
func GetReadyContainers(containers []corev1.ContainerStatus) string {
	total := len(containers)
	ready := 0
	for _, item := range containers {
		if item.Ready {
			ready++
		}
	}
	return strconv.Itoa(ready) + "/" + strconv.Itoa(total)
}

// IsValidCasType to return true if the casType is supported
func IsValidCasType(casType string) bool {
	return casType == CstorCasType || casType == JivaCasType || casType == LVMCasType || casType == ZFSCasType
}
