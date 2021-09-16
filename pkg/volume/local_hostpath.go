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
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// LocalHostpathVolInfoTemplate to store the local-hostpath volume and pvc describe related details
	LocalHostpathVolInfoTemplate = `
{{.Name}} Details :
-----------------
NAME            : {{.Name}}
ACCESS MODE     : {{.AccessMode}}
CAS TYPE        : {{.CasType}}
STORAGE CLASS   : {{.StorageClass}}
VOLUME PHASE    : {{.VolumePhase }}
SIZE            : {{.Size}}
CAPACITY        : {{.Capacity}}
PV CLAIM        : {{.PVC}}
RECLAIM POLICY  : {{.ReclaimPolicy}}
`
)

// GetLocalHostpath returns a list of localpv-hostpath columes
func GetLocalHostpath(c *client.K8sClient, pvList *corev1.PersistentVolumeList, openebsNS string) ([]metav1.TableRow, error) {
	var rows []metav1.TableRow
	for _, pv := range pvList.Items {
		// Ignore all the other volumes that is not of cas-type local-hostpath
		// dynamic-local-provisioner has this label for PVs openebs.io/cas-type=local-hostpath
		// this might be fixed later
		if util.GetCasTypeFromPV(&pv) != util.LocalHostpathCasLabel {
			continue
		}
		name := pv.Name
		capacity := pv.Spec.Capacity.Storage()
		sc := pv.Spec.StorageClassName
		attached := pv.Status.Phase
		attachedNode, customStatus := pv.Labels["nodeID"], ""
		var storageVersion, ns string
		deploy, err := c.GetDeploymentList("openebs.io/component-name=openebs-localpv-provisioner")
		if err == nil && len(deploy.Items) == 1 {
			storageVersion = deploy.Items[0].Labels["openebs.io/version"]
			ns = deploy.Items[0].Namespace
		}

		accessMode := pv.Spec.AccessModes[0]
		rows = append(rows, metav1.TableRow{
			Cells: []interface{}{
				ns, name, customStatus, storageVersion, capacity, sc, attached,
				accessMode, attachedNode},
		})
	}
	return rows, nil
}

// DescribeLocalHostpathVolume describes a localpv-hostpath PersistentVolume
func DescribeLocalHostpathVolume(_ *client.K8sClient, vol *corev1.PersistentVolume) error {
	// Get Local-volume Information
	localHostpathVolInfo := util.LocalHostPathVolInfo{
		VolumeInfo: util.VolumeInfo{
			AccessMode:   util.AccessModeToString(vol.Spec.AccessModes),
			Capacity:     util.ConvertToIBytes(vol.Spec.Capacity.Storage().String()),
			Name:         vol.Name,
			PVC:          vol.Spec.ClaimRef.Name,
			VolumePhase:  vol.Status.Phase,
			StorageClass: vol.Spec.StorageClassName,
			Size:         util.ConvertToIBytes(vol.Spec.Capacity.Storage().String()),
		},
		Path:          vol.Spec.PersistentVolumeSource.Local.Path,
		ReclaimPolicy: string(vol.Spec.PersistentVolumeReclaimPolicy),
		CasType:       util.LocalPvHostpathCasType,
	}

	// Print the Volume information
	_ = util.PrintByTemplate("localHostpathVolumeInfo", LocalHostpathVolInfoTemplate, localHostpathVolInfo)
	return nil
}
