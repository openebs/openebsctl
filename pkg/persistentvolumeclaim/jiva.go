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
	"time"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/openebs/openebsctl/pkg/volume"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	jivaPvcInfoTemplate = `
{{.Name}} Details  :
-------------------
NAME               : {{.Name}}
NAMESPACE          : {{.Namespace}}
CAS TYPE           : {{.CasType}}
BOUND VOLUME       : {{.BoundVolume}}
ATTACHED TO NODE   : {{.AttachedToNode}}
JIVA VOLUME POLICY : {{.JVP}}
STORAGE CLASS      : {{.StorageClassName}}
SIZE               : {{.Size}}
JV STATUS          : {{.JVStatus}}
PV STATUS          : {{.PVStatus}}
`
)

// DescribeJivaVolumeClaim describes a jiva storage engine PersistentVolumeClaim
func DescribeJivaVolumeClaim(c *client.K8sClient, pvc *corev1.PersistentVolumeClaim, vol *corev1.PersistentVolume) error {
	// 1. Get the JivaVolume Corresponding to the pvc name
	jv, err := c.GetJV(pvc.Spec.VolumeName)
	if err != nil {
		fmt.Printf("failed to get JivaVolume for %s\n", pvc.Spec.VolumeName)
	}
	// 2. Fill in Jiva Volume Claim related details
	jivaPvcInfo := util.JivaPVCInfo{
		Name:             pvc.Name,
		Namespace:        pvc.Namespace,
		CasType:          util.JivaCasType,
		BoundVolume:      pvc.Spec.VolumeName,
		StorageClassName: *pvc.Spec.StorageClassName,
		Size:             util.ConvertToIBytes(jv.Spec.Capacity),
	}
	if jv != nil {
		jivaPvcInfo.AttachedToNode = jv.Labels["nodeID"]
		jivaPvcInfo.JVP = jv.Annotations["openebs.io/volume-policy"]
		jivaPvcInfo.JVStatus = jv.Status.Status
	}
	if vol != nil {
		jivaPvcInfo.PVStatus = vol.Status.Phase
	}
	// 3. Print the Jiva Volume Claim information
	_ = util.PrintByTemplate("jivaPvcInfo", jivaPvcInfoTemplate, jivaPvcInfo)
	// 4. Print the Portal Information
	if jv != nil {
		util.TemplatePrinter(volume.JivaPortalTemplate, jv)
	}
	// 5. Fetch the Jiva controller and replica pod details
	podList, err := c.GetJVTargetPod(jv.Name)
	if err == nil {
		fmt.Printf("Controller and Replica Pod Details :" + "\n-----------------------------------\n")
		var rows []metav1.TableRow
		for _, pod := range podList.Items {
			rows = append(rows, metav1.TableRow{Cells: []interface{}{
				pod.Namespace, pod.Name,
				util.GetReadyContainers(pod.Status.ContainerStatuses),
				pod.Status.Phase, util.Duration(time.Since(pod.ObjectMeta.CreationTimestamp.Time)),
				pod.Status.PodIP, pod.Spec.NodeName}})
		}
		util.TablePrinter(util.PodDetailsColumnDefinations, rows, printers.PrintOptions{Wide: true})
	} else {
		fmt.Printf("Controller and Replica Pod Details :\n-----------------------------------\nNo Controller and Replica pod exists for the JivaVolume\n")
	}
	// 6. Fetch the replica PVCs and create rows for cli-runtime
	var rows []metav1.TableRow
	pvcList, err := c.GetPVCs(jv.Namespace, nil, "openebs.io/component=jiva-replica,openebs.io/persistent-volume="+jv.Name)
	if err != nil || len(pvcList.Items) == 0 {
		fmt.Printf("No replicas found for the JivaVolume %s", jv.Name)
		return nil
	}
	for _, pvc := range pvcList.Items {
		rows = append(rows, metav1.TableRow{Cells: []interface{}{
			pvc.Name,
			pvc.Status.Phase,
			pvc.Spec.VolumeName,
			util.ConvertToIBytes(pvc.Spec.Resources.Requests.Storage().String()),
			*pvc.Spec.StorageClassName,
			util.Duration(time.Since(pvc.ObjectMeta.CreationTimestamp.Time)),
			pvc.Spec.VolumeMode}})
	}
	// 6. Print the replica details if present
	fmt.Printf("\n\nReplica Details :\n-----------------\n")
	util.TablePrinter(util.JivaReplicaPVCColumnDefinations, rows, printers.PrintOptions{Wide: true})
	return nil
}
