package client

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/klog"

	storagev1 "k8s.io/client-go/kubernetes/typed/storage/v1"

	cstorv1 "github.com/openebs/api/pkg/apis/cstor/v1"
	cstorv1CS "github.com/openebs/api/pkg/client/clientset/versioned/typed/cstor/v1"
	openebsv1 "github.com/openebs/api/pkg/client/clientset/versioned/typed/openebs.io/v1alpha1"

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

	openebsCS *openebsv1.OpenebsV1alpha1Client
	//Service v1.ServiceInterface

	kubeconfig string
}

// NewK8sClient creates a new K8sClient
func NewK8sClient(ns string) (*K8sClient, error) {
	// get the appropriate clientset
	cs := GetOutofClusterCS()

	config := os.Getenv("KUBECONFIG")

	sc := getStorageClient(config)

	cStorCS := getCStorClient(config)

	openebsCS := getOpenEBSClient(config)
	//discoveryCS := getDiscoveryCS(*configFlags, config)

	return &K8sClient{
		ns:         ns,
		cs:         cs,
		sc:         sc,
		cStorCS:    cStorCS,
		openebsCS:  openebsCS,
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

func getOpenEBSClient(kubeconfig string) *openebsv1.OpenebsV1alpha1Client {

	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	client := openebsv1.NewForConfigOrDie(config)

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
func (k K8sClient) GetcStorVolume(volName string, namespace string) *cstorv1.CStorVolume {
	vols := k.cStorCS.CStorVolumes(namespace)
	volInfo, err := vols.Get(volName, metav1.GetOptions{})

	if err != nil {
		klog.Errorf("Error while while getting volume %s in %s namespace %s\n",
			volName, namespace, err)
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
			AccessMode:              util.CheckIfAccessable(i),
		}
		volumes[vol.PVC] = vol
	}
	return volumes
}
