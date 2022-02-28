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
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/openebs/openebsctl/pkg/volume"
	corev1 "k8s.io/api/core/v1"
)

const (
	zfsPvcInfoTemplate = `
{{.Name}} Details  :
-------------------
NAME               : {{.Name}}
NAMESPACE          : {{.Namespace}}
CAS TYPE           : {{.CasType}}
BOUND VOLUME       : {{.BoundVolume}}
STORAGE CLASS      : {{.StorageClassName}}
SIZE               : {{.Size}}
PVC STATUS         : {{.PVCStatus}}
MOUNTED BY         : {{.MountPods}}
`
)

// DescribeZFSVolumeClaim describes a ZFS storage engine PersistentVolumeClaim
func DescribeZFSVolumeClaim(c *client.K8sClient, pvc *corev1.PersistentVolumeClaim, pv *corev1.PersistentVolume, mountPods string) error {
	zfsPVCinfo := util.ZFSPVCInfo{
		Name:             pvc.Name,
		Namespace:        pvc.Namespace,
		CasType:          util.ZFSCasType,
		BoundVolume:      pvc.Spec.VolumeName,
		StorageClassName: *pvc.Spec.StorageClassName,
		Size:             pvc.Spec.Resources.Requests.Storage().String(),
		PVCStatus:        pvc.Status.Phase,
		MountPods:        mountPods,
	}

	if pv != nil {
		_ = util.PrintByTemplate("zfsPvc", zfsPvcInfoTemplate, zfsPVCinfo)
		_ = volume.DescribeZFSLocalPVs(c, pv)
	} else {
		_ = util.PrintByTemplate("ZFSPvc", zfsPvcInfoTemplate, zfsPVCinfo)
	}

	return nil
}
