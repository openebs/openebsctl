package client

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/klog"

	storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"

	cstorv1 "github.com/openebs/api/pkg/apis/cstor/v1"
	cstorv1CS "github.com/openebs/api/pkg/client/clientset/versioned/typed/cstor/v1"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	util "github.com/vaniisgh/mayactl/kubectl-mayactl/cli/util"
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

	// cs refers to the Clientset capable of communicating
	// with the K8s cluster
	cs *kubernetes.Clientset

	// sc is the client to interact with the CSI driver/nodes/volumes etc.
	sc *storagev1.StorageV1Client

	// cstor Clientset
	cStorCS *cstorv1CS.CstorV1Client

	kubeconfig string
}

// NewK8sClient creates a new K8sClient
func NewK8sClient(ns string) (*K8sClient, error) {
	// get the appropriate clientset
	cs := GetOutofClusterCS()

	config := os.Getenv("KUBECONFIG")

	sc := getStorageClient(config)

	cStorCS := getCStorClient(config)

	return &K8sClient{
		ns:         ns,
		cs:         cs,
		sc:         sc,
		cStorCS:    cStorCS,
		kubeconfig: config,
	}, nil

}

// GetOutofClusterCS creates returns a clientset for the kubeconfig &
// sets the env variable for the same
func GetOutofClusterCS() (client *kubernetes.Clientset) {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.
			Join(home, ".kube", "config"), "absolute path to kubeconfig")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)

	os.Setenv("KUBECONFIG", *kubeconfig)

	return clientset
}

// getStorageClass is a function that returns the storage client for the
// config
func getStorageClient(kubeconfig string) *storagev1.StorageV1Client {

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)

	sc, err := storagev1.NewForConfig(config)

	if err != nil {
		klog.Error(err)
	}

	return sc
}

func getCStorClient(kubeconfig string) *cstorv1CS.CstorV1Client {

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	client := cstorv1CS.NewForConfigOrDie(config)

	return client
}

// GetService fetches the K8s Service with the provided name
func (k *K8sClient) GetService(name string) {
	client := k.cs
	svcLabelSelector := "name" + "=" + name
	sops, err := client.CoreV1().Services(k.ns).
		List(metav1.ListOptions{LabelSelector: svcLabelSelector})

	if err != nil {
		klog.Error("Error while accessing " + name + " in namespaces: " + k.ns)
		klog.Error(err)

	}

	if len(sops.Items) == 0 {
		klog.Error("No services" + name + "in namespaces :" + k.ns)
		return
	}

	fmt.Println(sops.Items[0])
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("KUBECONFIG")
}

// GetStorageClass using the K8sClient's storage class client
func (k K8sClient) GetStorageClass(driver string) *v1.StorageClass {

	scs, err := k.sc.StorageClasses().Get(driver, metav1.GetOptions{})

	if err != nil {
		klog.Errorf("Error while while getting storage class: %s\n", err)
		os.Exit(1)
	}

	return scs
}

// GetCSIVolume using the K8sClient's storage class client
func (k K8sClient) GetCSIVolume(volname string) *v1.VolumeAttachment {
	vol, err := k.sc.VolumeAttachments().Get(volname, metav1.GetOptions{})

	if err != nil {
		klog.Errorf("Error while while getting volumes: %s\n", err)
		os.Exit(1)
	}

	return vol
}

// GetcStorVolumes using the K8sClient's storage class client
func (k K8sClient) GetcStorVolumes() cstorv1.CStorVolumeList {

	cStorVols, err := k.cStorCS.CStorVolumes("").List(metav1.ListOptions{})

	if err != nil {
		klog.Errorf("Error while while getting volumes: %s\n", err)
		os.Exit(1)
	}

	return *cStorVols

}

// GetcStorVolume fetches the volume object of the given name in the given namespace
func (k K8sClient) GetcStorVolume(volName string) *cstorv1.CStorVolume {
	vols := k.cStorCS.CStorVolumes(k.ns)
	volInfo, err := vols.Get(volName, metav1.GetOptions{})

	if err != nil {
		klog.Errorf("Error while while getting volume %s in %s namespace %s\n",
			volName, k.ns, err)
		os.Exit(1)
	}

	return volInfo
}

// GetcStorPVCs used to get the infor for the underlying
// PVC
func (k K8sClient) GetcStorPVCs(node string) map[string]*util.Volume {

	volumes := make(map[string]*util.Volume)

	PVCs, err := k.sc.VolumeAttachments().List(metav1.ListOptions{})
	if err != nil {
		klog.Errorf("Error while while getting storage volume attachments on %s", node, err)
		os.Exit(1)
	}

	for _, i := range PVCs.Items {
		vol := &util.Volume{
			StorageClass:            i.Spec.Attacher,
			Node:                    i.Spec.NodeName,
			PVC:                     *i.Spec.Source.PersistentVolumeName,
			CSIVolumeAttachmentName: i.Name,
			AttachementStatus:       util.CheckVolAttachmentError(i.Status),
			// first fetch access modes & then convert to string
			AccessMode: util.AccessModeToString(k.GetPV(*i.Spec.Source.PersistentVolumeName).Spec.AccessModes),
		}
		volumes[vol.PVC] = vol
	}
	return volumes
}

// GetPV returns a PV object after querying Kubernetes API
func (k K8sClient) GetPV(name string) *corev1.PersistentVolume {

	vol, err := k.cs.CoreV1().PersistentVolumes().Get(name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("Error while getting volume %s", name)
	}

	return vol
}

// GetCVC used to get cStor Volume Config information for cStor a given volume using a cStorClient
func (k K8sClient) GetCVC(name string) *cstorv1.CStorVolumeConfig {

	cStorVolumeConfig, err := k.cStorCS.CStorVolumeConfigs(k.ns).Get(name, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("Error while getting cStor Volume Config for  %s in %s", name, k.ns)
	}

	return cStorVolumeConfig
}

// GetCVR used to get cStor Volume Replicas for a given cStor volumes using cStor Client
func (k K8sClient) GetCVR(name string) []cstorv1.CStorVolumeReplica {

	label := "cstorvolume.openebs.io/name" + "=" + name

	CStorVolumeReplicas, err := k.cStorCS.CStorVolumeReplicas("").List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		klog.Errorf("Error while getting cStor Volume Replica for  %s in %s", name, k.ns)
	}

	if len(CStorVolumeReplicas.Items) == 0 {
		klog.Errorf("Error while getting cStor Volume Replica for  %s , couldnot fild any replicas", name)
	}

	return CStorVolumeReplicas.Items
}

// NodeForVolume used to get NodeName for the volume from the Kubernetes API
func (k K8sClient) NodeForVolume(volName string) string {

	label := "openebs.io/persistent-volume" + "=" + volName

	podInfo, err := k.cs.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: label})
	if err != nil {
		klog.Errorf("Error while getting target Pod for volume %s", volName)
	}

	if len(podInfo.Items) != 1 {
		klog.Errorf("Error invalid number of Pods for volume %s", volName)
	}

	return podInfo.Items[0].Spec.NodeName
}
