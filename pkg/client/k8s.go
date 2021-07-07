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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	openebsclientset "github.com/openebs/api/v2/pkg/client/clientset/versioned"
	jiva "github.com/openebs/jiva-operator/pkg/apis/openebs/v1alpha1"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// K8sAPIVersion represents valid kubernetes api version of a native or custom
// resource
type K8sAPIVersion string

// K8sClient provides the necessary utility to operate over
// various K8s Kind objects
type K8sClient struct {
	// Ns refers to K8s namespace where the operation
	// will be performed
	Ns string
	// K8sCS refers to the Clientset capable of communicating
	// with the K8s cluster
	K8sCS kubernetes.Interface
	// OpenebsClientset capable of accessing the OpenEBS
	// components
	OpenebsCS openebsclientset.Interface
}

/*
	CLIENT CREATION METHODS AND RELATED OPERATIONS
*/

// NewK8sClient creates a new K8sClient
// TODO: improve K8sClientset instantiation. for example remove the Ns from
// K8sClient struct
func NewK8sClient(ns string) (*K8sClient, error) {
	// get the appropriate clientsets & set the kubeconfig accordingly
	// TODO: The kubeconfig should ideally be initialized in the CLI depending on various flags
	GetOutofClusterKubeConfig()
	config := os.Getenv("KUBECONFIG")
	k8sCS, err := getK8sClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build Kubernetes clientset")
	}
	openebsCS, err := getOpenEBSClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build OpenEBS clientset")
	}
	return &K8sClient{
		Ns:        ns,
		K8sCS:     k8sCS,
		OpenebsCS: openebsCS,
	}, nil
}

// GetOutofClusterKubeConfig creates returns a clientset for the kubeconfig &
// sets the env variable for the same
func GetOutofClusterKubeConfig() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.
			Join(home, ".kube", "config"), "absolute path to kubeconfig")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig")
	}
	flag.Parse()
	err := os.Setenv("KUBECONFIG", *kubeconfig)
	if err != nil {
		return
	}
}

// getK8sClient returns K8s clientset by taking kubeconfig as an argument
func getK8sClient(kubeconfig string) (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "Could not build config from flags")
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get new config")
	}
	return clientset, nil
}

// getOpenEBSClient returns OpenEBS clientset by taking kubeconfig as an
// argument
func getOpenEBSClient(kubeconfig string) (*openebsclientset.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "Could not build config from flags")
	}
	client, err := openebsclientset.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get new config")
	}
	return client, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("KUBECONFIG")
}

/*
	NAMESPACE DETERMINATION METHODS
*/

// GetOpenEBSNamespace from the specific engine component based on cas-type
func (k K8sClient) GetOpenEBSNamespace(casType string) (string, error) {
	pods, err := k.K8sCS.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("openebs.io/component-name=%s", util.CasTypeAndComponentNameMap[strings.ToLower(casType)])})
	if err != nil || len(pods.Items) == 0 {
		return "", errors.New("unable to determine openebs namespace")
	}
	return pods.Items[0].Namespace, nil
}

// GetOpenEBSNamespaceMap maps the cas-type to it's namespace, e.g. n[cstor] = cstor-ns
func (k K8sClient) GetOpenEBSNamespaceMap() (map[string]string, error) {
	label := "openebs.io/component-name in ("
	for _, v := range util.CasTypeAndComponentNameMap {
		label = label + v + ","
	}
	label += ")"
	pods, err := k.K8sCS.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: label})
	if err != nil || pods == nil || len(pods.Items) == 0 {
		return nil, errors.New("unable to determine openebs namespace")
	}
	NSmap := make(map[string]string)
	for _, pod := range pods.Items {
		ns := pod.Namespace
		cas, ok := util.ComponentNameToCasTypeMap[pod.Labels["openebs.io/component-name"]]
		if ok {
			NSmap[cas] = ns
		}
	}
	return NSmap, nil
}

/*
	NATIVE RESOURCE FETCHING METHODS
*/

// GetSC returns a StorageClass object using the scName passed.
func (k K8sClient) GetSC(scName string) (*v1.StorageClass, error) {
	sc, err := k.K8sCS.StorageV1().StorageClasses().Get(context.TODO(), scName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while getting storage class")
	}
	return sc, nil
}

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
	if volNames == nil || len(volNames) == 0 {
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
	if pvcNames == nil || len(pvcNames) == 0 {
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

/*
	OPENEBS RESOURCE FETCHING METHODS
*/

// GetBD returns the BlockDevice passed as name with OpenEBS's Client
func (k K8sClient) GetBD(bd string) (*v1alpha1.BlockDevice, error) {
	blockDevice, err := k.OpenebsCS.OpenebsV1alpha1().BlockDevices(k.Ns).Get(context.TODO(), bd, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting block device")
	}
	return blockDevice, nil
}

// GetBDs returns a list of BlockDevices based on the values of bdNames slice.
// bdNames slice if is nil or empty, it returns all the BDs in the cluster.
// bdNames slice if is not nil or not empty, it return the BDs whose names are present in the slice.
// labelselector takes the label(key+value) and makes an api call with this filter applied. Can be empty string if label filtering is not needed.
func (k K8sClient) GetBDs(bdNames []string, labelselector string) (*v1alpha1.BlockDeviceList, error) {
	bds, err := k.OpenebsCS.OpenebsV1alpha1().BlockDevices(k.Ns).List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting block device")
	}
	if bdNames == nil || len(bdNames) == 0 {
		return bds, nil
	}
	bdNameBDmap := make(map[string]v1alpha1.BlockDevice)
	for _, item := range bds.Items {
		bdNameBDmap[item.Name] = item
	}
	var items = make([]v1alpha1.BlockDevice, 0)
	for _, name := range bdNames {
		if _, ok := bdNameBDmap[name]; ok {
			items = append(items, bdNameBDmap[name])
		}
	}
	return &v1alpha1.BlockDeviceList{
		Items: items,
	}, nil
}

/*
	CSTOR STORAGE ENGINE SPECIFIC METHODS
*/

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
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Error while getting volumes")
	}
	var list []cstorv1.CStorVolume
	if volNames == nil || len(volNames) == 0 {
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
	if err != nil {
		return nil, nil, err
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
	if cspiNames == nil || len(cspiNames) == 0 {
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
		} else {
			fmt.Printf("Error from server (NotFound): pool %s not found\n", name)
		}
	}
	return &cstorv1.CStorPoolInstanceList{
		Items: list,
	}, nil
}

/*
	JIVA STORAGE ENGINE SPECIFIC METHODS
*/

// GetJV returns the JivaVolume passed as name with REST Client
func (k K8sClient) GetJV(jv string) (*jiva.JivaVolume, error) {
	var j jiva.JivaVolume
	err := k.K8sCS.Discovery().RESTClient().Get().Namespace(k.Ns).Name(jv).AbsPath("/apis/openebs.io/v1alpha1").
		Resource("jivavolumes").Do(context.TODO()).Into(&j)
	if err != nil {
		return nil, err
	}
	return &j, nil
}

// GetJVs returns a list or map of JivaVolumes based on the values of volNames slice, and options.
// volNames slice if is nil or empty, it returns all the JVs in the cluster.
// volNames slice if is not nil or not empty, it return the JVs whose names are present in the slice.
// rType takes the return type of the method, can either List or Map.
// labelselector takes the label(key+value) and makes an api call with this filter applied, can be empty string if label filtering is not needed.
// options takes a MapOptions object which defines how to create a map, refer to types for more info. Can be empty in case of rType is List.
// Only one type can be returned at a time, please define the other type as '_' while calling.
func (k K8sClient) GetJVs(volNames []string, rType util.ReturnType, labelSelector string, options util.MapOptions) (*jiva.JivaVolumeList, map[string]jiva.JivaVolume, error) {
	jvs := jiva.JivaVolumeList{}
	// NOTE: The resource name must be plural and the API-group should be present for getting CRs
	err := k.K8sCS.Discovery().RESTClient().Get().AbsPath("/apis/openebs.io/v1alpha1").
		Resource("jivavolumes").Do(context.TODO()).Into(&jvs)
	if err != nil {
		return nil, nil, err
	}
	var list []jiva.JivaVolume
	if volNames == nil || len(volNames) == 0 {
		list = jvs.Items
	} else {
		jvsMap := make(map[string]jiva.JivaVolume)
		for _, jv := range jvs.Items {
			jvsMap[jv.Name] = jv
		}
		for _, name := range volNames {
			if jv, ok := jvsMap[name]; ok {
				list = append(list, jv)
			} else {
				fmt.Printf("Error from server (NotFound): jivavolume %s not found\n", name)
			}
		}
	}
	if rType == util.List {
		return &jiva.JivaVolumeList{
			Items: list,
		}, nil, nil
	}
	if rType == util.Map {
		jvMap := make(map[string]jiva.JivaVolume)
		switch options.Key {
		case util.Label:
			for _, jv := range list {
				if vol, ok := jv.Labels[options.LabelKey]; ok {
					jvMap[vol] = jv
				}
			}
			return nil, jvMap, nil
		case util.Name:
			for _, jv := range list {
				jvMap[jv.Name] = jv
			}
			return nil, jvMap, nil
		default:
			return nil, nil, errors.New("invalid map options")
		}
	}
	return nil, nil, errors.New("invalid return type")
}

// GetJVTargetPod returns the Jiva Volume Controller and Replica Pods, corresponding to the volumeName.
func (k K8sClient) GetJVTargetPod(volumeName string) (*corev1.PodList, error) {
	pods, err := k.K8sCS.CoreV1().Pods(k.Ns).List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("openebs.io/cas-type=jiva,openebs.io/persistent-volume=%s", volumeName)})
	if err != nil || len(pods.Items) == 0 {
		return nil, errors.New("The target pod for the volume was not found")
	}
	return pods, nil
}
