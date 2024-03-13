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

package persistentvolumeclaim

import (
	"sort"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

// Describe manages various implementations of PersistentVolumeClaim Describing
func Describe(pvcs []string, namespace string, openebsNs string) error {
	if len(pvcs) == 0 || pvcs == nil {
		return errors.New("please provide atleast one pvc name to describe")
	}
	// Clienset creation
	k := client.NewK8sClient(openebsNs)

	// 1. Get a list of required PersistentVolumeClaims
	var pvcList *corev1.PersistentVolumeClaimList
	pvcList, err := k.GetPVCs(namespace, pvcs, "")
	if err != nil {
		return errors.New("no pvcs found corresponding to the names")
	}
	// 2. Get the namespaces
	nsMap, _ := k.GetOpenEBSNamespaceMap()
	// 3. Get all Pods to find PVC mount Pods
	var nsPods []corev1.Pod
	podList, err := k.GetAllPods(namespace)
	if err != nil {
		nsPods = []corev1.Pod{}
	} else {
		nsPods = podList.Items
	}
	// 4. Range over the list of PVCs
	for _, pvc := range pvcList.Items {
		// 5. Fetch the storage class, used to get the cas-type
		sc, _ := k.GetSC(*pvc.Spec.StorageClassName)
		pv, _ := k.GetPV(pvc.Spec.VolumeName)
		// 6. Get cas type
		casType := util.GetCasType(pv, sc)
		mountPods := PodsToString(SortPods(GetMountPods(pvc.Name, nsPods)))
		// 7. Assign a namespace corresponding to the engine
		if openebsNs == "" {
			if val, ok := nsMap[casType]; ok {
				k.Ns = val
			}
		}
		// 8. Describe the volume based on its casType
		if desc, ok := CasDescribeMap()[casType]; ok {
			err = desc(k, &pvc, pv, mountPods)
			if err != nil {
				continue
			}
		} else {
			// Describe volume with some generic stuffs if casType is not understood
			err := DescribeGenericVolumeClaim(&pvc, pv, casType, mountPods)
			if err != nil {
				continue
			}
		}
	}
	return nil
}

// CasDescribeMap returns a map cas-types to functions for persistentvolumeclaim describing
func CasDescribeMap() map[string]func(*client.K8sClient, *corev1.PersistentVolumeClaim, *corev1.PersistentVolume, string) error {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, *corev1.PersistentVolumeClaim, *corev1.PersistentVolume, string) error{
		util.LVMCasType: DescribeLVMVolumeClaim,
		util.ZFSCasType: DescribeZFSVolumeClaim,
	}
}

// GetMountPods filters the array of Pods and returns an array of Pods that mount the PersistentVolumeClaim
func GetMountPods(pvcName string, nsPods []corev1.Pod) []corev1.Pod {
	var pods []corev1.Pod
	for _, pod := range nsPods {
		volumes := pod.Spec.Volumes
		for _, volume := range volumes {
			pvc := volume.VolumeSource.PersistentVolumeClaim
			if pvc != nil && pvc.ClaimName == pvcName {
				pods = append(pods, pod)
				break
			}
		}
	}
	return pods
}

// SortPods sorts the array of Pods by name
func SortPods(pods []corev1.Pod) []corev1.Pod {
	sort.Slice(pods, func(i, j int) bool {
		cmpKey := func(pod corev1.Pod) string {
			return pod.Name
		}
		return cmpKey(pods[i]) < cmpKey(pods[j])
	})
	return pods
}

// PodsToString Flattens the array of Pods and returns a string fit to display in the output
func PodsToString(pods []corev1.Pod) string {
	if len(pods) == 0 {
		return "none"
	}
	str := ""
	for _, pod := range pods {
		str += pod.Name + " "
	}
	return str
}
