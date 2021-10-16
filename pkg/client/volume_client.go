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

package client

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// GetPV returns a PersistentVolume object using the pv name passed.
func (k K8sClient) GetPV(name string) (*corev1.PersistentVolume, error) {
	pv, err := k.K8sCS.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while getting persistent volume")
	}
	return pv, nil
}

// GetPVs returns a list of PersistentVolumes based on the values of volNames slice.
// volNames slice if is nil or empty, it returns all the PVs in the cluster.
// volNames slice if is not nil or not empty, it return the PVs whose names are present in the slice.
// labelselector takes the label(key+value) and makes an api call with this filter applied. Can be empty string if label filtering is not needed.
func (k K8sClient) GetPVs(volNames []string, labelselector string) (*corev1.PersistentVolumeList, error) {
	pvs, err := k.K8sCS.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return nil, err
	}
	volMap := make(map[string]corev1.PersistentVolume)
	for _, vol := range pvs.Items {
		volMap[vol.Name] = vol
	}
	var list []corev1.PersistentVolume
	if len(volNames) == 0 {
		return pvs, nil
	}
	for _, name := range volNames {
		if pool, ok := volMap[name]; ok {
			list = append(list, pool)
		} else {
			fmt.Printf("Error from server (NotFound): PV %s not found\n", name)
		}
	}
	return &corev1.PersistentVolumeList{
		Items: list,
	}, nil
}

// GetPVC returns a PersistentVolumeClaim object using the pvc name passed.
func (k K8sClient) GetPVC(name string, namespace string) (*corev1.PersistentVolumeClaim, error) {
	pvc, err := k.K8sCS.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while getting persistent volume claim")
	}
	return pvc, nil
}

// GetPVCs returns a list of PersistentVolumeClaims based on the values of pvcNames slice.
// namespace takes the namespace in which PVCs are present.
// pvcNames slice if is nil or empty, it returns all the PVCs in the cluster, in the namespace.
// pvcNames slice if is not nil or not empty, it return the PVCs whose names are present in the slice, in the namespace.
// labelselector takes the label(key+value) and makes an api call with this filter applied. Can be empty string if label filtering is not needed.
func (k K8sClient) GetPVCs(namespace string, pvcNames []string, labelselector string) (*corev1.PersistentVolumeClaimList, error) {
	pvcs, err := k.K8sCS.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return nil, err
	}
	if len(pvcNames) == 0 {
		return pvcs, nil
	}
	pvcNamePVCmap := make(map[string]corev1.PersistentVolumeClaim)
	for _, item := range pvcs.Items {
		pvcNamePVCmap[item.Name] = item
	}
	var items = make([]corev1.PersistentVolumeClaim, 0)
	for _, name := range pvcNames {
		if _, ok := pvcNamePVCmap[name]; ok {
			items = append(items, pvcNamePVCmap[name])
		}
	}
	return &corev1.PersistentVolumeClaimList{
		Items: items,
	}, nil
}
