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
	"fmt"
	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	"k8s.io/cli-runtime/pkg/printers"
	"os"
	"time"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	cstorVolInfoTemplate = `
{{.Name}} Details :
-----------------
NAME            : {{.Name}}
ACCESS MODE     : {{.AccessMode}}
CSI DRIVER      : {{.CSIDriver}}
STORAGE CLASS   : {{.StorageClass}}
VOLUME PHASE    : {{.VolumePhase }}
VERSION         : {{.Version}}
CSPC            : {{.CSPC}}
SIZE            : {{.Size}}
STATUS          : {{.Status}}
REPLICA COUNT	: {{.ReplicaCount}}

`

	cstorPortalTemplate = `
Portal Details :
------------------
IQN              :  {{.IQN}}
VOLUME NAME      :  {{.VolumeName}}
TARGET NODE NAME :  {{.TargetNodeName}}
PORTAL           :  {{.Portal}}
TARGET IP        :  {{.TargetIP}}

`
)

// GetCStor returns a list of CStor volumes
func GetCStor(c *client.K8sClient, pvList *corev1.PersistentVolumeList, openebsNS string) ([]metav1.TableRow, error) {
	var (
		cvMap  map[string]v1.CStorVolume
		cvaMap map[string]v1.CStorVolumeAttachment
	)
	// no need to proceed if CVs/CVAs don't exist
	var err error
	_, cvMap, err = c.GetCVs(nil, util.Map, "", util.MapOptions{
		Key: util.Name})
	if err != nil {
		return nil, fmt.Errorf("failed to list CStorVolumes")
	}
	_, cvaMap, err = c.GetCVAs(util.Map, "", util.MapOptions{
		Key:      util.Label,
		LabelKey: "Volname"})
	if err != nil {
		return nil, fmt.Errorf("failed to list CStorVolumeAttachments")
	}
	var rows []metav1.TableRow
	// 3. Show the required ones
	for _, pv := range pvList.Items {
		var attachedNode, storageVersion, customStatus, ns string
		// TODO: Estimate the cas-type and decide to print it out
		// Should all AccessModes be shown in a csv format, or the highest be displayed ROO < RWO < RWX?
		// 2. For eligible PVs fetch the custom-resource to add more info
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == util.CStorCSIDriver {
			// For all CSI CStor PV, there exist a CV
			cv, ok := cvMap[pv.Name]
			if !ok {
				// condition not possible
				_, _ = fmt.Fprintf(os.Stderr, "couldn't find cv "+pv.Name)
			}
			ns = cv.Namespace
			if openebsNS != "" && openebsNS != ns {
				continue
			}
			customStatus = string(cv.Status.Phase)
			storageVersion = cv.VersionDetails.Status.Current
			cva, cvaOk := cvaMap[pv.Name]
			if cvaOk {
				attachedNode = cva.Labels["nodeID"]
			}
			// TODO: What should be done for multiple AccessModes
			accessMode := pv.Spec.AccessModes[0]
			rows = append(rows, metav1.TableRow{
				Cells: []interface{}{
					ns, pv.Name, customStatus, storageVersion, pv.Spec.Capacity.Storage(), pv.Spec.StorageClassName, pv.Status.Phase,
					accessMode, attachedNode}})
		}
	}
	return rows, nil
}

func DescribeCstorVolume(c *client.K8sClient, vol corev1.PersistentVolume) error  {
	// Fetch all details of a volume is called to get the volume controller's
	// info such as controller's IP, status, iqn, replica IPs etc.
	// 1. cStor volume info
	volumeInfo, err := c.GetCV(vol.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get cStorVolume for %s\n", vol.Name)
		return err
	}
	// 2. cStor Volume Config
	cvcInfo, err := c.GetCVC(vol.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get cStor Volume config for %s\n", vol.Name)
		return err
	}
	// 3. Get Node for Target Pod from the openebs-ns
	node, err := c.GetCVA(util.CVAVolnameKey + "=" + vol.Name)
	var nodeName string
	if err != nil {
		nodeName = util.NotAttached
		fmt.Fprintf(os.Stderr, "failed to get CStorVolumeAttachments for %s\n", vol.Name)
	} else {
		nodeName = node.Spec.Volume.OwnerNodeID
	}

	// 5. cStor Volume Replicas
	cvrInfo, err := c.GetCVRs(cstortypes.PersistentVolumeLabelKey + "=" + vol.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get cStor Volume Replicas for %s\n", vol.Name)
	}

	cSPCLabel := cstortypes.CStorPoolClusterLabelKey
	volume := util.VolumeInfo{
		AccessMode:              util.AccessModeToString(vol.Spec.AccessModes),
		Capacity:                volumeInfo.Status.Capacity.String(),
		CSPC:                    cvcInfo.Labels[cSPCLabel],
		CSIDriver:               vol.Spec.CSI.Driver,
		CSIVolumeAttachmentName: vol.Spec.CSI.VolumeHandle,
		Name:                    volumeInfo.Name,
		Namespace:               volumeInfo.Namespace,
		PVC:                     vol.Spec.ClaimRef.Name,
		ReplicaCount:            volumeInfo.Spec.ReplicationFactor,
		VolumePhase:             vol.Status.Phase,
		StorageClass:            vol.Spec.StorageClassName,
		Version:                 util.CheckVersion(volumeInfo.VersionDetails),
		Size:                    util.ConvertToIBytes(volumeInfo.Status.Capacity.String()),
		Status:                  string(volumeInfo.Status.Phase),
	}

	// Print the output for the portal status info
	err = util.PrintByTemplate("volume", cstorVolInfoTemplate, volume)
	if err != nil {
		return err
	}

	portalInfo := util.PortalInfo{
		IQN:            volumeInfo.Spec.Iqn,
		VolumeName:     volumeInfo.Name,
		Portal:         volumeInfo.Spec.TargetPortal,
		TargetIP:       volumeInfo.Spec.TargetIP,
		TargetNodeName: nodeName,
	}

	// Print the output for the portal status info
	err = util.PrintByTemplate("PortalInfo", cstorPortalTemplate, portalInfo)
	if err != nil {
		return err
	}

	replicaCount := volumeInfo.Spec.ReplicationFactor
	// This case will occur only if user has manually specified zero replica.
	// or if none of the replicas are healthy & running
	if replicaCount == 0 || len(volumeInfo.Status.ReplicaStatuses) == 0 {
		fmt.Printf("None of the replicas are running\n")
		//please check the volume pod's status by running [kubectl describe pvc -l=openebs/replica --all-namespaces]\Oor try again later.")
		return nil
	}

	// Print replica details
	if cvrInfo != nil && len(cvrInfo.Items) > 0 {
		fmt.Printf("\nReplica Details :\n-----------------\n")
		var rows []metav1.TableRow
		for _, cvr := range cvrInfo.Items {
			rows = append(rows, metav1.TableRow{Cells: []interface{}{
				cvr.Name,
				util.ConvertToIBytes(cvr.Status.Capacity.Total),
				util.ConvertToIBytes(cvr.Status.Capacity.Used),
				cvr.Status.Phase,
				util.Duration(time.Since(cvr.ObjectMeta.CreationTimestamp.Time))}})
		}
		util.TablePrinter(util.CstorReplicaColumnDefinations, rows, printers.PrintOptions{Wide: true})
	}

	cStorBackupList, err := c.GetCVBackups(vol.Name)
	if cStorBackupList != nil {
		fmt.Printf("\nCstor Backup Details :\n" + "---------------------\n")
		var rows []metav1.TableRow
		for _, item := range cStorBackupList.Items {
			rows = append(rows, metav1.TableRow{Cells: []interface{}{
				item.ObjectMeta.Name,
				item.Spec.BackupName,
				item.Spec.VolumeName,
				item.Spec.BackupDest,
				item.Spec.SnapName,
				item.Status,
			}})
		}
		util.TablePrinter(util.CstorBackupColumnDefinations, rows, printers.PrintOptions{Wide: true})
	}

	cstorCompletedBackupList, err := c.GetCVCompletedBackups(vol.Name)
	if cstorCompletedBackupList != nil {
		fmt.Printf("\nCstor Completed Backup Details :" + "\n-------------------------------\n")
		var rows []metav1.TableRow
		for _, item := range cstorCompletedBackupList.Items {
			rows = append(rows, metav1.TableRow{Cells: []interface{}{
				item.Name,
				item.Spec.BackupName,
				item.Spec.VolumeName,
				item.Spec.LastSnapName,
			}})
		}
		util.TablePrinter(util.CstorCompletedBackupColumnDefinations, rows, printers.PrintOptions{Wide: true})
	}

	cStorRestoreList, err := c.GetCVRestores(vol.Name)
	if cStorRestoreList != nil {
		fmt.Printf("\nCstor Restores Details :" + "\n-----------------------\n")
		var rows []metav1.TableRow
		for _, item := range cStorRestoreList.Items {
			rows = append(rows, metav1.TableRow{Cells: []interface{}{
				item.ObjectMeta.Name,
				item.Spec.RestoreName,
				item.Spec.VolumeName,
				item.Spec.RestoreSrc,
				item.Spec.StorageClass,
				item.Spec.Size.String(),
				item.Status,
			}})
		}
		util.TablePrinter(util.CstorRestoreColumnDefinations, rows, printers.PrintOptions{Wide: true})
	}
	fmt.Println()
	return nil
}