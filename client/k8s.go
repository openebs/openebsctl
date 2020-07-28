package client

import (
	"flag"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/klog"

	"github.com/pkg/errors"

	cstorv1 "github.com/openebs/api/pkg/apis/cstor/v1"
	cstortypes "github.com/openebs/api/pkg/apis/types"
	openebsclientset "github.com/openebs/api/pkg/client/clientset/versioned"
	util "github.com/openebs/openebsctl/kubectl-openebs/cli/util"

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

	scs, err := k.K8sCS.StorageV1().StorageClasses().Get(driver, metav1.GetOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "error while while getting storage class")
	}

	return scs, nil
}

// GetCSIVolume using the K8sClient's storage class client
func (k K8sClient) GetCSIVolume(volname string) (*v1.VolumeAttachment, error) {
	vol, err := k.K8sCS.StorageV1().VolumeAttachments().Get(volname, metav1.GetOptions{})

	if err != nil {
		return nil, errors.Wrap(err, "error while while getting storage csi volume")
	}

	return vol, nil
}

// GetcStorVolumes using the K8sClient's storage class client
func (k K8sClient) GetcStorVolumes() (*cstorv1.CStorVolumeList, error) {

	cStorVols, err := k.OpenebsCS.CstorV1().CStorVolumes("").List(metav1.ListOptions{})

	if err != nil {
		return nil, errors.Wrapf(err, "Error while while getting volumes")
	}

	return cStorVols, nil

}

// GetcStorVolume fetches the volume object of the given name in the given namespace
func (k K8sClient) GetcStorVolume(volName string) (*cstorv1.CStorVolume, error) {
	volInfo, err := k.OpenebsCS.CstorV1().CStorVolumes(k.ns).Get(volName, metav1.GetOptions{})

	if err != nil {
		return nil, errors.Wrapf(err, "error while while getting volume %s", volName)
	}

	return volInfo, nil
}

// GetCStorVolumeInfoMap used to get the info for for the underlying
// PVC
func (k K8sClient) GetCStorVolumeInfoMap(node string) (map[string]*util.Volume, error) {

	volumes := make(map[string]*util.Volume)

	PVCs, err := k.K8sCS.StorageV1().VolumeAttachments().List(metav1.ListOptions{})
	if err != nil {
		return volumes, errors.Wrap(err, "error while while getting storage volume attachments")
	}

	for _, i := range PVCs.Items {

		if i.Spec.Source.PersistentVolumeName == nil {
			continue
		}

		pv, err := k.GetPV(*i.Spec.Source.PersistentVolumeName)

		if err != nil {
			klog.Errorf("Failed to get PV", i.ObjectMeta.Name)
			continue
		}
		vol := &util.Volume{
			StorageClass:            i.Spec.Attacher,
			Node:                    i.Spec.NodeName,
			PVC:                     *i.Spec.Source.PersistentVolumeName,
			CSIVolumeAttachmentName: i.Name,
			AttachementStatus:       util.CheckVolAttachmentError(i.Status),
			// first fetch access modes & then convert to string
			AccessMode: util.AccessModeToString(pv.Spec.AccessModes),
		}
		volumes[vol.PVC] = vol
	}
	return volumes, nil
}

// GetPV returns a PV object after querying Kubernetes API
func (k K8sClient) GetPV(name string) (*corev1.PersistentVolume, error) {

	vol, err := k.K8sCS.CoreV1().PersistentVolumes().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while while getting persistant volume")
	}

	return vol, nil
}

// GetCVC used to get cStor Volume Config information for cStor a given volume using a cStorClient
func (k K8sClient) GetCVC(name string) (*cstorv1.CStorVolumeConfig, error) {

	cStorVolumeConfig, err := k.OpenebsCS.CstorV1().CStorVolumeConfigs(k.ns).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "error while getting cStor Volume Config for  %s in %s", name, k.ns)
	}

	return cStorVolumeConfig, nil
}

// GetCVR used to get cStor Volume Replicas for a given cStor volumes using cStor Client
func (k K8sClient) GetCVR(name string) (*cstorv1.CStorVolumeReplicaList, error) {

	label := "cstorvolume.openebs.io/name" + "=" + name

	CStorVolumeReplicas, err := k.OpenebsCS.CstorV1().CStorVolumeReplicas("").List(metav1.ListOptions{LabelSelector: label})
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

	podInfo, err := k.K8sCS.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		return "", errors.Wrapf(err, "error while getting target Pod for volume %s", volName)
	}

	if len(podInfo.Items) != 1 {
		klog.Errorf("Error invalid number of Pods %d for volume %s", len(podInfo.Items), volName)
	}

	return podInfo.Items[0].Spec.NodeName, nil
}
