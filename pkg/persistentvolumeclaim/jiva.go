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
	"strings"
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
		fmt.Println(fmt.Sprintf("failed to get JivaVolume for %s", pvc.Spec.VolumeName))
	}
	// 2. Fill in Jiva Volume Claim related details
	jivaPvcInfo := util.JivaPVCInfo{
		Name:             pvc.Name,
		Namespace:        pvc.Namespace,
		CasType:          util.JivaCasType,
		BoundVolume:      pvc.Spec.VolumeName,
		StorageClassName: *pvc.Spec.StorageClassName,
		Size:             pvc.Spec.Resources.Requests.Storage().String(),
	}
	if jv != nil {
		jivaPvcInfo.AttachedToNode = jv.Labels["nodeID"]
		jivaPvcInfo.JVP = jv.Annotations["openebs.io/volume-policy"]
		jivaPvcInfo.JVStatus = jv.Status.Status
	}
	if vol != nil {
		jivaPvcInfo.PVStatus = vol.Status.Phase
		// 3. Print the Jiva Volume Claim information
		_ = util.PrintByTemplate("jivaPvcInfo", jivaPvcInfoTemplate, jivaPvcInfo)
	} else {
		_ = util.PrintByTemplate("jivaPvcInfo", jivaPvcInfoTemplate, jivaPvcInfo)
		fmt.Println(fmt.Sprintf("PersistentVolume %s, doesnot exist", pvc.Spec.VolumeName))
		return nil
	}
	// 4. Print the Portal Information
	replicaPodIPAndModeMap := make(map[string]string)
	if jv != nil {
		util.TemplatePrinter(volume.JivaPortalTemplate, jv)
		// Create Replica IP to Mode Map
		if jv.Status.ReplicaStatuses != nil && len(jv.Status.ReplicaStatuses) != 0 {
			for _, replicaStatus := range jv.Status.ReplicaStatuses {
				replicaPodIPAndModeMap[strings.Split(replicaStatus.Address, ":")[1][2:]] = replicaStatus.Mode
			}
		}
	}
	// 5. Fetch the Jiva controller and replica pod details
	podList, err := c.GetJVTargetPod(vol.Name)
	if err == nil {
		fmt.Println("Controller and Replica Pod Details :")
		fmt.Println("-----------------------------------")
		var rows []metav1.TableRow
		for _, pod := range podList.Items {
			if strings.Contains(pod.Name, "-ctrl-") {
				rows = append(rows, metav1.TableRow{Cells: []interface{}{
					pod.Namespace, pod.Name, jv.Status.Status,
					pod.Spec.NodeName, pod.Status.Phase, pod.Status.PodIP,
					util.GetReadyContainers(pod.Status.ContainerStatuses),
					util.Duration(time.Since(pod.ObjectMeta.CreationTimestamp.Time))}})
			} else {
				if val, ok := replicaPodIPAndModeMap[pod.Status.PodIP]; ok {
					rows = append(rows, metav1.TableRow{Cells: []interface{}{
						pod.Namespace, pod.Name, val,
						pod.Spec.NodeName, pod.Status.Phase, pod.Status.PodIP,
						util.GetReadyContainers(pod.Status.ContainerStatuses),
						util.Duration(time.Since(pod.ObjectMeta.CreationTimestamp.Time))}})
				}
			}
		}
		util.TablePrinter(util.JivaPodDetailsColumnDefinations, rows, printers.PrintOptions{Wide: true})
	} else {
		fmt.Printf("Controller and Replica Pod Details :")
		fmt.Println("-----------------------------------")
		fmt.Println("No Controller and Replica pod exists for the JivaVolume")
	}
	// 6. Fetch the replica PVCs and create rows for cli-runtime
	var rows []metav1.TableRow
	pvcList, err := c.GetPVCs(c.Ns, nil, "openebs.io/component=jiva-replica,openebs.io/persistent-volume="+jv.Name)
	if err != nil || len(pvcList.Items) == 0 {
		fmt.Printf("No replicas found for the JivaVolume %s", vol.Name)
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
	fmt.Println()
	fmt.Println("Replica Data Volume Details :")
	fmt.Println("-----------------------------")
	util.TablePrinter(util.JivaReplicaPVCColumnDefinations, rows, printers.PrintOptions{Wide: true})
	return nil
}
