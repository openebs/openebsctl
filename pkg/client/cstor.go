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

package client

import (
	"context"
	"fmt"

	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// GetCV returns the CStorVolume passed as name with OpenEBS's Client
func (k K8sClient) GetCV(volName string) (*cstorv1.CStorVolume, error) {
	volInfo, err := k.OpenebsCS.CstorV1().CStorVolumes(k.Ns).Get(context.TODO(), volName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting volume %s", volName)
	}
	return volInfo, nil
}

// GetCVs returns a list or map of CStorVolumes based on the values of volNames slice, and options.
// volNames slice if is nil or empty, it returns all the CVs in the cluster.
// volNames slice if is not nil or not empty, it return the CVs whose names are present in the slice.
// rType takes the return type of the method, can be either List or Map.
// labelselector takes the label(key+value) and makes an api call with this filter applied, can be empty string if label filtering is not needed.
// options takes a MapOptions object which defines how to create a map, refer to types for more info. Can be empty in case of rType is List.
// Only one type can be returned at a time, please define the other type as '_' while calling.
func (k K8sClient) GetCVs(volNames []string, rType util.ReturnType, labelSelector string, options util.MapOptions) (*cstorv1.CStorVolumeList, map[string]cstorv1.CStorVolume, error) {
	cVols, err := k.OpenebsCS.CstorV1().CStorVolumes("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if len(cVols.Items) == 0 {
		return nil, nil, errors.Errorf("Error while getting volumes%v", err)
	}
	var list []cstorv1.CStorVolume
	if len(volNames) == 0 {
		list = cVols.Items
	} else {
		csMap := make(map[string]cstorv1.CStorVolume)
		for _, cv := range cVols.Items {
			csMap[cv.Name] = cv
		}
		for _, name := range volNames {
			if cv, ok := csMap[name]; ok {
				list = append(list, cv)
			} else {
				fmt.Printf("Error from server (NotFound): cStorVolume %s not found\n", name)
			}
		}
	}
	if rType == util.List {
		return &cstorv1.CStorVolumeList{
			Items: list,
		}, nil, nil
	}
	if rType == util.Map {
		cvMap := make(map[string]cstorv1.CStorVolume)
		switch options.Key {
		case util.Label:
			for _, cv := range list {
				if vol, ok := cv.Labels[options.LabelKey]; ok {
					cvMap[vol] = cv
				}
			}
			return nil, cvMap, nil
		case util.Name:
			for _, cv := range list {
				cvMap[cv.Name] = cv
			}
			return nil, cvMap, nil
		default:
			return nil, nil, errors.New("invalid map options")
		}
	}
	return nil, nil, errors.New("invalid return type")
}

// GetCVA returns the CStorVolumeAttachment, corresponding to the label passed.
// Ex:- labelSelector: {cstortypes.PersistentVolumeLabelKey + "=" + pvName}
func (k K8sClient) GetCVA(labelSelector string) (*cstorv1.CStorVolumeAttachment, error) {
	vol, err := k.OpenebsCS.CstorV1().CStorVolumeAttachments("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, errors.Wrap(err, "error from server (NotFound): CVA not found")
	} else if vol == nil || len(vol.Items) == 0 {
		return nil, fmt.Errorf("error from server (NotFound): CVA not found for %s", labelSelector)
	}
	return &vol.Items[0], nil
}

// GetCVAs returns a list or map of CStorVolumeAttachments based on the values of options.
// rType takes the return type of the method, can either be List or Map.
// labelselector takes the label(key+value) and makes a api call with this filter applied, can be empty string if label filtering is not needed.
// options takes a MapOptions object which defines how to create a map, refer to types for more info. Can be empty in case of rType is List.
// Only one type can be returned at a time, please define the other type as '_' while calling.
func (k K8sClient) GetCVAs(rType util.ReturnType, labelSelector string, options util.MapOptions) (*cstorv1.CStorVolumeAttachmentList, map[string]cstorv1.CStorVolumeAttachment, error) {
	cvaList, err := k.OpenebsCS.CstorV1().CStorVolumeAttachments("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if len(cvaList.Items) == 0 {
		return nil, nil, errors.Errorf("No CVA found for %s, %v", labelSelector, err)
	}
	if rType == util.List {
		return cvaList, nil, nil
	}
	if rType == util.Map {
		cvaMap := make(map[string]cstorv1.CStorVolumeAttachment)
		switch options.Key {
		case util.Label:
			for _, cva := range cvaList.Items {
				if vol, ok := cva.Labels[options.LabelKey]; ok {
					cvaMap[vol] = cva
				}
			}
			return nil, cvaMap, nil
		case util.Name:
			for _, cva := range cvaList.Items {
				cvaMap[cva.Name] = cva
			}
			return nil, cvaMap, nil
		default:
			return nil, nil, errors.New("invalid map options")
		}
	}
	return nil, nil, errors.New("invalid return type")
}

// GetCVTargetPod returns the Cstor Volume Target Pod, corresponding to the volumeClaim and volumeName.
func (k K8sClient) GetCVTargetPod(volumeClaim string, volumeName string) (*corev1.Pod, error) {
	pods, err := k.K8sCS.CoreV1().Pods(k.Ns).List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("openebs.io/persistent-volume-claim=%s,openebs.io/persistent-volume=%s,openebs.io/target=cstor-target", volumeClaim, volumeName)})
	if err != nil || len(pods.Items) == 0 {
		return nil, errors.New("The target pod for the volume was not found")
	}
	return &pods.Items[0], nil
}

// GetCVInfoMap returns a Volume object, filled using corresponding CVA and PV.
func (k K8sClient) GetCVInfoMap() (map[string]*util.Volume, error) {
	volumes := make(map[string]*util.Volume)
	cstorVA, _, err := k.GetCVAs(util.List, "", util.MapOptions{})
	if err != nil {
		return volumes, errors.Wrap(err, "error while getting storage volume attachments")
	}
	for _, i := range cstorVA.Items {
		if i.Spec.Volume.Name == "" {
			continue
		}
		pv, err := k.GetPV(i.Spec.Volume.Name)
		if err != nil {
			klog.Errorf("Failed to get PV %s", i.ObjectMeta.Name)
			continue
		}
		vol := &util.Volume{
			StorageClass:            pv.Spec.StorageClassName,
			Node:                    i.Labels["nodeID"],
			PVC:                     pv.Spec.ClaimRef.Name,
			CSIVolumeAttachmentName: i.Name,
			AttachementStatus:       string(pv.Status.Phase),
			// first fetch access modes & then convert to string
			AccessMode: util.AccessModeToString(pv.Spec.AccessModes),
		}
		// map the pv name to the vol obj
		volumes[i.Spec.Volume.Name] = vol
	}
	return volumes, nil
}

// GetCVBackups returns the CStorVolumeBackup, corresponding to the label passed.
// Ex:- labelSelector: {cstortypes.PersistentVolumeLabelKey + "=" + pvName}
func (k K8sClient) GetCVBackups(labelselector string) (*cstorv1.CStorBackupList, error) {
	cstorBackupList, err := k.OpenebsCS.CstorV1().CStorBackups("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil || len(cstorBackupList.Items) == 0 {
		return nil, errors.New("no cstorbackups were found for the volume")
	}
	return cstorBackupList, nil
}

// GetCVCompletedBackups returns the CStorCompletedBackups, corresponding to the label passed.
// Ex:- labelSelector: {cstortypes.PersistentVolumeLabelKey + "=" + pvName}
func (k K8sClient) GetCVCompletedBackups(labelselector string) (*cstorv1.CStorCompletedBackupList, error) {
	cstorCompletedBackupList, err := k.OpenebsCS.CstorV1().CStorCompletedBackups("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil || len(cstorCompletedBackupList.Items) == 0 {
		return nil, errors.New("no cstorcompletedbackups were found for the volume")
	}
	return cstorCompletedBackupList, nil
}

// GetCVRestores returns the CStorRestores, corresponding to the label passed.
// Ex:- labelSelector: {cstortypes.PersistentVolumeLabelKey + "=" + pvName}
func (k K8sClient) GetCVRestores(labelselector string) (*cstorv1.CStorRestoreList, error) {
	cStorRestoreList, err := k.OpenebsCS.CstorV1().CStorRestores("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil || len(cStorRestoreList.Items) == 0 {
		return nil, errors.New("no cstorrestores were found for the volume")
	}
	return cStorRestoreList, nil
}

// GetCVC returns the CStorVolumeConfig for cStor volume using the PV/CV/CVC name.
func (k K8sClient) GetCVC(name string) (*cstorv1.CStorVolumeConfig, error) {
	cStorVolumeConfig, err := k.OpenebsCS.CstorV1().CStorVolumeConfigs(k.Ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting cStor Volume Config for  %s in %s", name, k.Ns)
	}
	return cStorVolumeConfig, nil
}

// GetCVRs returns the list CStorVolumeReplica, corresponding to the label passed.
// For ex:- labelselector : {"cstorvolume.openebs.io/name" + "=" + name} , {"cstorpoolinstance.openebs.io/name" + "=" + poolName}
func (k K8sClient) GetCVRs(labelselector string) (*cstorv1.CStorVolumeReplicaList, error) {
	cvrs, err := k.OpenebsCS.CstorV1().CStorVolumeReplicas("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting cStor Volume Replica for %s", labelselector)
	}
	if cvrs == nil || len(cvrs.Items) == 0 {
		fmt.Printf("Error while getting cStor Volume Replica for %s, no replicas found for \n", labelselector)
	}
	return cvrs, nil
}

// GetCSPC returns the CStorPoolCluster for cStor volume using the poolName passed.
func (k K8sClient) GetCSPC(poolName string) (*cstorv1.CStorPoolCluster, error) {
	cStorPool, err := k.OpenebsCS.CstorV1().CStorPoolClusters(k.Ns).Get(context.TODO(), poolName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting cspc")
	}
	return cStorPool, nil
}

func (k K8sClient) ListCSPC() (*cstorv1.CStorPoolClusterList, error) {
	cStorPool, err := k.OpenebsCS.CstorV1().CStorPoolClusters(k.Ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting cspc")
	}
	return cStorPool, nil
}

// GetCSPI returns the CStorPoolInstance for cStor volume using the poolName passed.
func (k K8sClient) GetCSPI(poolName string) (*cstorv1.CStorPoolInstance, error) {
	cStorPool, err := k.OpenebsCS.CstorV1().CStorPoolInstances(k.Ns).Get(context.TODO(), poolName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting cspi")
	}
	return cStorPool, nil
}

// GetCSPIs returns a list of CStorPoolInstances based on the values of cspiNames slice
// cspiNames slice if is nil or empty, it returns all the CSPIs in the cluster
// cspiNames slice if is not nil or not empty, it return the CSPIs whose names are present in the slice
// labelselector takes the label(key+value) and makes an api call with this filter applied. Can be empty string if label filtering is not needed.
func (k K8sClient) GetCSPIs(cspiNames []string, labelselector string) (*cstorv1.CStorPoolInstanceList, error) {
	cspi, err := k.OpenebsCS.CstorV1().CStorPoolInstances("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting cspi")
	}
	if len(cspiNames) == 0 {
		return cspi, nil
	}
	poolMap := make(map[string]cstorv1.CStorPoolInstance)
	for _, p := range cspi.Items {
		poolMap[p.Name] = p
	}
	var list []cstorv1.CStorPoolInstance
	for _, name := range cspiNames {
		if pool, ok := poolMap[name]; ok {
			list = append(list, pool)
		}
		// else {
		// This logging might be omitted
		// fmt.Fprintf(os.Stderr, "Error from server (NotFound): pool %s not found\n", name)
		//}
	}
	return &cstorv1.CStorPoolInstanceList{
		Items: list,
	}, nil
}
