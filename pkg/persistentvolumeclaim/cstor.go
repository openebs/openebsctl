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
	"fmt"
	"time"

	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	cstorPvcInfoTemplate = `
{{.Name}} Details :
------------------
NAME             : {{.Name}}
NAMESPACE        : {{.Namespace}}
CAS TYPE         : {{.CasType}}
BOUND VOLUME     : {{.BoundVolume}}
ATTACHED TO NODE : {{.AttachedToNode}}
POOL             : {{.Pool}}
STORAGE CLASS    : {{.StorageClassName}}
SIZE             : {{.Size}}
USED             : {{.Used}}
CV STATUS	 : {{.CVStatus}}
PV STATUS        : {{.PVStatus}}
MOUNTED BY       : {{.MountPods}}
`

	detailsFromCVC = `
Additional Details from CVC :
-----------------------------
NAME          : {{ .metadata.name }}
REPLICA COUNT : {{ .spec.provision.replicaCount }}
POOL INFO     : {{ .status.poolInfo}}
VERSION       : {{ .versionDetails.status.current}}
UPGRADING     : {{if eq .versionDetails.status.current .versionDetails.desired}}false{{else}}true{{end}}
`
)

// DescribeCstorVolumeClaim describes a cstor storage engine PersistentVolumeClaim
func DescribeCstorVolumeClaim(c *client.K8sClient, pvc *corev1.PersistentVolumeClaim, pv *corev1.PersistentVolume, mountPods string) error {
	// Create Empty template objects and fill gradually when underlying sub CRs are identified.
	pvcInfo := util.CstorPVCInfo{}

	pvcInfo.Name = pvc.Name
	pvcInfo.Namespace = pvc.Namespace
	pvcInfo.BoundVolume = pvc.Spec.VolumeName
	pvcInfo.CasType = util.CstorCasType
	pvcInfo.StorageClassName = *pvc.Spec.StorageClassName
	pvcInfo.MountPods = mountPods

	if pv != nil {
		pvcInfo.PVStatus = pv.Status.Phase
	}

	// fetching the underlying CStorVolume for the PV, to get the phase and size and notify the user
	// if the CStorVolume is not found.
	cv, err := c.GetCV(pvc.Spec.VolumeName)
	if err != nil {
		fmt.Println("Underlying CstorVolume is not found for: ", pvc.Name)
	} else {
		pvcInfo.Size = util.ConvertToIBytes(cv.Spec.Capacity.String())
		pvcInfo.CVStatus = cv.Status.Phase
	}

	// fetching the underlying CStorVolumeConfig for the PV, to get the cvc info and Pool Name and notify the user
	// if the CStorVolumeConfig is not found.
	cvc, err := c.GetCVC(pvc.Spec.VolumeName)
	if err != nil {
		fmt.Println("Underlying CstorVolumeConfig is not found for: ", pvc.Name)
	} else {
		pvcInfo.Pool = cvc.Labels[cstortypes.CStorPoolClusterLabelKey]
	}

	// fetching the underlying CStorVolumeAttachment for the PV, to get the attached to node and notify the user
	// if the CStorVolumeAttachment is not found.
	cva, err := c.GetCVA(util.CVAVolnameKey + "=" + pvc.Spec.VolumeName)
	if err != nil {
		pvcInfo.AttachedToNode = util.NotAvailable
		fmt.Println("Underlying CstorVolumeAttachment is not found for: ", pvc.Name)
	} else {
		pvcInfo.AttachedToNode = cva.Spec.Volume.OwnerNodeID
	}

	// fetching the underlying CStorVolumeReplicas for the PV, to list their details and notify the user
	// none of the replicas are running if the CStorVolumeReplicas are not found.
	cvrs, err := c.GetCVRs(cstortypes.PersistentVolumeLabelKey + "=" + pvc.Spec.VolumeName)
	if err == nil && len(cvrs.Items) > 0 {
		pvcInfo.Used = util.ConvertToIBytes(util.GetUsedCapacityFromCVR(cvrs))
	}

	// Printing the Filled Details of the Cstor PVC
	_ = util.PrintByTemplate("pvc", cstorPvcInfoTemplate, pvcInfo)

	// fetching the underlying TargetPod for the PV, to display its relevant details and notify the user
	// if the TargetPod is not found.
	tgtPod, err := c.GetCVTargetPod(pvc.Name, pvc.Spec.VolumeName)
	if err == nil {
		fmt.Printf("\nTarget Details :\n----------------\n")
		var rows []metav1.TableRow
		rows = append(rows, metav1.TableRow{Cells: []interface{}{
			tgtPod.Namespace, tgtPod.Name,
			util.GetReadyContainers(tgtPod.Status.ContainerStatuses),
			tgtPod.Status.Phase, util.Duration(time.Since(tgtPod.ObjectMeta.CreationTimestamp.Time)),
			tgtPod.Status.PodIP, tgtPod.Spec.NodeName}})
		util.TablePrinter(util.PodDetailsColumnDefinations, rows, printers.PrintOptions{Wide: true})
	} else {
		fmt.Printf("\nTarget Details :\n----------------\nNo target pod exists for the CstorVolume\n")
	}

	// If CVRs are found list them and show relevant details else notify the user none of the replicas are
	// running if not found
	if cvrs != nil && len(cvrs.Items) > 0 {
		fmt.Printf("\nReplica Details :\n-----------------\n")
		var rows []metav1.TableRow
		for _, cvr := range cvrs.Items {
			rows = append(rows, metav1.TableRow{Cells: []interface{}{cvr.Name,
				util.ConvertToIBytes(cvr.Status.Capacity.Total),
				util.ConvertToIBytes(cvr.Status.Capacity.Used),
				cvr.Status.Phase,
				util.Duration(time.Since(cvr.ObjectMeta.CreationTimestamp.Time))}})
		}
		util.TablePrinter(util.CstorReplicaColumnDefinations, rows, printers.PrintOptions{Wide: true})
	} else {
		fmt.Printf("\nReplica Details :\n-----------------\nNo running replicas found\n")
	}

	if cvc != nil {
		util.TemplatePrinter(detailsFromCVC, cvc)
	}

	return nil
}
