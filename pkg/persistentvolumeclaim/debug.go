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
	"fmt"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

// Debug manages various implementations of PersistentVolumeClaim Describing
func Debug(pvcs []string, namespace string, openebsNs string) error {
	if len(pvcs) == 0 || pvcs == nil {
		return errors.New("please provide atleast one pvc name to describe")
	}
	// Clienset creation
	k := client.NewK8sClient(openebsNs)

	// 1. Get a list of required PersistentVolumeClaims
	var pvcList *corev1.PersistentVolumeClaimList
	pvcList, err := k.GetPVCs(namespace, pvcs, "")
	if len(pvcList.Items) == 0 || err != nil {
		return errors.New("no pvcs found corresponding to the names")
	}
	// 2. Get the namespaces
	nsMap, _ := k.GetOpenEBSNamespaceMap()
	// 3. Range over the list of PVCs
	for _, pvc := range pvcList.Items {
		// 4. Fetch the storage class, used to get the cas-type
		sc, _ := k.GetSC(*pvc.Spec.StorageClassName)
		pv, _ := k.GetPV(pvc.Spec.VolumeName)
		// 5. Get cas type
		casType := util.GetCasType(pv, sc)
		// 6. Assign a namespace corresponding to the engine
		if openebsNs == "" {
			if val, ok := nsMap[casType]; ok {
				k.Ns = val
			}
		}
		// 7. Debug the volume based on its casType
		if desc, ok := CasDebugMap()[casType]; ok {
			err = desc(k, &pvc, pv)
			if err != nil {
				continue
			}
		} else {
			fmt.Printf("Debugging is currently not supported for %s Cas Type PVCs\n", casType)
		}
	}
	return nil
}

// CasDebugMap returns a map cas-types to functions for persistentvolumeclaim debugging
func CasDebugMap() map[string]func(*client.K8sClient, *corev1.PersistentVolumeClaim, *corev1.PersistentVolume) error {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, *corev1.PersistentVolumeClaim, *corev1.PersistentVolume) error{
		util.CstorCasType: DebugCstorVolumeClaim,
	}
}
