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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/klog"

	"github.com/pkg/errors"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	openebsclientset "github.com/openebs/api/v2/pkg/client/clientset/versioned"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// K8sAPIVersion represents valid kubernetes api version of a native or custom
// resource
type K8sAPIVersion string

// K8sClient provides the necessary utility to operate over
// various K8s Kind objects
type K8sClient struct {
	// ns refers to K8s namespace where the operation
	// will be performed
	ns string
	// K8sCS refers to the Clientset capable of communicating
	// with the K8s cluster
	K8sCS kubernetes.Interface
	// OpenebsClientset capable of accessing the OpenEBS
	// components
	OpenebsCS openebsclientset.Interface
}

// NewK8sClient creates a new K8sClient
// TODO: improve K8sClientset instantiation. for example remove the ns from
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
		ns:        ns,
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
	os.Setenv("KUBECONFIG", *kubeconfig)
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

// GetStorageClass using the K8sClient's storage class client
func (k K8sClient) GetStorageClass(driver string) (*v1.StorageClass, error) {
	scs, err := k.K8sCS.StorageV1().StorageClasses().Get(context.TODO(), driver, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while while getting storage class")
	}
	return scs, nil
}

// GetCSIVolume using the K8sClient's storage class client
func (k K8sClient) GetCSIVolume(volname string) (*cstorv1.CStorVolumeAttachment, error) {
	vol, err := k.OpenebsCS.CstorV1().CStorVolumeAttachments("").Get(context.TODO(), volname, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while while getting storage csi volume")
	}
	return vol, nil
}

// GetcStorVolumes using the K8sClient's storage class client
func (k K8sClient) GetcStorVolumes() (*cstorv1.CStorVolumeList, error) {
	cStorVols, err := k.OpenebsCS.CstorV1().CStorVolumes("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while while getting volumes")
	}
	return cStorVols, nil
}

// GetcStorVolume fetches the volume object of the given name in the given namespace
func (k K8sClient) GetcStorVolume(volName string) (*cstorv1.CStorVolume, error) {
	volInfo, err := k.OpenebsCS.CstorV1().CStorVolumes(k.ns).Get(context.TODO(), volName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error while while getting volume %s", volName)
	}
	return volInfo, nil
}

// GetCStorVolumeInfoMap used to get the info for for the underlying
// PVC
func (k K8sClient) GetCStorVolumeInfoMap(node string) (map[string]*util.Volume, error) {
	volumes := make(map[string]*util.Volume)
	cstorVA, err := k.OpenebsCS.CstorV1().CStorVolumeAttachments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return volumes, errors.Wrap(err, "error while while getting storage volume attachments")
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
			Node:                    i.ObjectMeta.OwnerReferences[0].Name,
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

// GetPV returns a PV object after querying Kubernetes API
func (k K8sClient) GetPV(name string) (*corev1.PersistentVolume, error) {
	vol, err := k.K8sCS.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while while getting persistant volume")
	}
	return vol, nil
}

// GetCVC used to get cStor Volume Config information for cStor a given volume using a cStorClient
func (k K8sClient) GetCVC(name string) (*cstorv1.CStorVolumeConfig, error) {
	cStorVolumeConfig, err := k.OpenebsCS.CstorV1().CStorVolumeConfigs(k.ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting cStor Volume Config for  %s in %s", name, k.ns)
	}
	return cStorVolumeConfig, nil
}

// GetCVR used to get cStor Volume Replicas for a given cStor volumes using cStor Client
func (k K8sClient) GetCVR(name string) (*cstorv1.CStorVolumeReplicaList, error) {
	label := "cstorvolume.openebs.io/name" + "=" + name
	CStorVolumeReplicas, err := k.OpenebsCS.CstorV1().CStorVolumeReplicas("").List(context.TODO(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting cStor Volume Replica for volume %s", name)
	}
	if len(CStorVolumeReplicas.Items) == 0 {
		klog.Errorf("Error while getting cStor Volume Replica for  %s , couldnot fild any replicas", name)
	}
	return CStorVolumeReplicas, nil
}

// NodeForVolume used to get NodeName for the volume from the Kubernetes API
func (k K8sClient) NodeForVolume(volName string) (string, error) {
	label := cstortypes.PersistentVolumeLabelKey + "=" + volName
	podInfo, err := k.K8sCS.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return "", errors.Wrapf(err, "error while getting target Pod for volume %s", volName)
	}
	if len(podInfo.Items) != 1 {
		klog.Errorf("Error invalid number of Pods %d for volume %s", len(podInfo.Items), volName)
	}
	return podInfo.Items[0].Spec.NodeName, nil
}

// GetcStorPools using the OpenEBS's Client
func (k K8sClient) GetcStorPools() (*cstorv1.CStorPoolInstanceList, error) {
	cStorPools, err := k.OpenebsCS.CstorV1().CStorPoolInstances("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while while getting cspc")
	}
	return cStorPools, nil
}

// GetPVCs list from the passed list of PVC names and the namespace
func (k K8sClient) GetPVCs(namespace string, pvcNames []string) (*corev1.PersistentVolumeClaimList, error) {
	pvcs, err := k.K8sCS.CoreV1().PersistentVolumeClaims("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	pvcNamePVCmap := make(map[string]corev1.PersistentVolumeClaim)
	for _, item := range pvcs.Items {
		pvcNamePVCmap[item.Name] = item
	}
	var items = make([]corev1.PersistentVolumeClaim, 0)
	for _, name := range pvcNames {
		items = append(items, pvcNamePVCmap[name])
	}
	return &corev1.PersistentVolumeClaimList{
		TypeMeta: metav1.TypeMeta{},
		ListMeta: metav1.ListMeta{},
		Items:    items,
	}, nil
}

func GetCasType(v1PV *corev1.PersistentVolume, v1SC *v1.StorageClass) string {
	if val := GetCasTypeFromPV(v1PV); val != util.UNKNOWN {
		return val
	}
	if val := GetCasTypeFromSC(v1SC); val != util.UNKNOWN {
		return val
	}
	return util.UNKNOWN
}

// GetCasTypeFromPV from the passed PersistentVolume or the Stora
func GetCasTypeFromPV(v1PV *corev1.PersistentVolume) string {
	if v1PV.ObjectMeta.Labels != nil {
		if _, ok := v1PV.ObjectMeta.Labels[util.OPENEBS_CAS_TYPE_KEY]; ok {
			return v1PV.ObjectMeta.Labels[util.OPENEBS_CAS_TYPE_KEY]
		}
	} else if v1PV.ObjectMeta.Annotations != nil {
		if _, ok := v1PV.ObjectMeta.Annotations[util.OPENEBS_CAS_TYPE_KEY]; ok {
			return v1PV.ObjectMeta.Annotations[util.OPENEBS_CAS_TYPE_KEY]
		}
	} else if v1PV.Spec.CSI != nil && v1PV.Spec.CSI.VolumeAttributes != nil {
		if _, ok := v1PV.Spec.CSI.VolumeAttributes[util.OPENEBS_CAS_TYPE_KEY]; ok {
			return v1PV.Spec.CSI.VolumeAttributes[util.OPENEBS_CAS_TYPE_KEY]
		}
	}
	return util.UNKNOWN
}

// GetCasTypeFromSC by passing the storage class
func GetCasTypeFromSC(v1SC *v1.StorageClass) string {
	if v1SC.Parameters != nil {
		if _, ok := v1SC.Parameters[util.OPENBEBS_CAS_TYPE_KEY_SC]; ok {
			return v1SC.Parameters[util.OPENBEBS_CAS_TYPE_KEY_SC]
		}
	}
	return util.UNKNOWN
}

// GetCVA from the passed cstorvolume name
func (k K8sClient) GetCVA(volumeName string) (*cstorv1.CStorVolumeAttachment, error) {
	cvaList, err := k.OpenebsCS.CstorV1().CStorVolumeAttachments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range cvaList.Items {
		if item.Spec.Volume.Name == volumeName {
			return &item, nil
		}
	}
	return nil, errors.New("Couldn't find the CVA for the passed volume")
}

// GetUsedCapacityFromCVR as the healthy replicas would have the correct used capacity details
func GetUsedCapacityFromCVR(cvrList *cstorv1.CStorVolumeReplicaList) string {
	for _, item := range cvrList.Items {
		if item.Status.Phase == util.Healthy {
			return item.Status.Capacity.Used
		}
	}
	return ""
}

// GetCstorVolumeTargetPod for the passed volume to show details
func (k K8sClient) GetCstorVolumeTargetPod(volumeName string) (*corev1.Pod, error) {
	pods, err := k.K8sCS.CoreV1().Pods(k.ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, item := range pods.Items {
		if strings.Contains(item.Name, volumeName+"-target-") {
			return &item, nil
		}
	}
	return nil, errors.New("The target pod for the given cstor volume was not found")
}

// GetReadyContainers to show the number of ready bs total containers of pod
func GetReadyContainers(containers []corev1.ContainerStatus) string {
	total := len(containers)
	ready := 0
	if total > 0 {
		for _, item := range containers {
			if item.Ready == true {
				ready++
			}
		}
	}
	return strconv.Itoa(ready) + "/" + strconv.Itoa(total)
}
