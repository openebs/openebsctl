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
	"log"
	"os"
	"path/filepath"
	"strings"

	lvmclient "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset"
	"github.com/openebs/openebsctl/pkg/util"
	zfsclient "github.com/openebs/zfs-localpv/pkg/generated/clientset/internalclientset"
	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

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
	// LVMCS is the client for accessing OpenEBS LVM components
	LVMCS lvmclient.Interface
	// ZFCS is the client for accessing OpenEBS ZFS components
	ZFCS zfsclient.Interface
}

/*
	CLIENT CREATION METHODS AND RELATED OPERATIONS
*/

// NewK8sClient is a wrapper around newK8sClient to handle errors in
// creating clients implicitilty and simulating namespace as an optional parameter
// ns: kubernetes namespace
func NewK8sClient(ns ...string) *K8sClient {
	// If more than one-namespace is provided as a function param, throw error and exit
	if len(ns) > 1 {
		log.Fatal("Only one namespace arg is allowed")
	}

	namespace := ""
	if len(ns) == 1 {
		namespace = ns[0]
	}

	k, err := newK8sClient(namespace)
	if err != nil {
		log.Fatal("error creating kubernetes client: ", err)
	}

	return k
}

// newK8sClient creates a new K8sClient
// TODO: improve K8sClientset instantiation. for example remove the Ns from
// K8sClient struct
func newK8sClient(ns string) (*K8sClient, error) {
	// get the appropriate clientsets & set the kubeconfig accordingly
	// TODO: The kubeconfig should ideally be initialized in the CLI depending on various flags
	GetOutofClusterKubeConfig()
	config := os.Getenv("KUBECONFIG")
	k8sCS, err := getK8sClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build Kubernetes clientset")
	}
	lv, _ := getLVMclient(config)
	zf, _ := getZFSclient(config)
	return &K8sClient{
		Ns:    ns,
		K8sCS: k8sCS,
		LVMCS: lv,
		ZFCS:  zf,
	}, nil
}

// GetOutofClusterKubeConfig creates returns a clientset for the kubeconfig &
// sets the env variable for the same
func GetOutofClusterKubeConfig() {
	var kubeconfig *string
	// config file not provided, auto detect from the host OS
	if util.Kubeconfig == "" {
		if home := homeDir(); home != "" {
			cfg := filepath.Join(home, ".kube", "config")
			kubeconfig = &cfg
		} else {
			log.Fatal(`kubeconfig not provided, Please provide config file path with "--kubeconfig" flag`)
		}
	} else {
		// Get the kubeconfig file from CLI args
		kubeconfig = &util.Kubeconfig
	}
	err := os.Setenv("KUBECONFIG", *kubeconfig)
	if err != nil {
		log.Fatal(err)
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
// NOTE: This will not work correctly if CSI controller pod runs in kube-system NS
func (k K8sClient) GetOpenEBSNamespace(casType string) (string, error) {
	pods, err := k.K8sCS.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Running", LabelSelector: fmt.Sprintf("openebs.io/component-name=%s", util.CasTypeAndComponentNameMap[strings.ToLower(casType)])})
	if err != nil || len(pods.Items) == 0 {
		return "", fmt.Errorf("unable to determine openebs namespace, err: %v", err)
	}
	return pods.Items[0].Namespace, nil
}

// GetOpenEBSNamespaceMap maps the cas-type to it's namespace, e.g. n[zfs] = zfs-ns
// NOTE: This will not work correctly if CSI controller pod runs in kube-system NS
func (k K8sClient) GetOpenEBSNamespaceMap() (map[string]string, error) {
	label := "openebs.io/component-name in ("
	for _, v := range util.CasTypeAndComponentNameMap {
		label = label + v + ","
	}
	label += ")"
	pods, err := k.K8sCS.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Running", LabelSelector: label})
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

// Get Versions of different components running in K8s
func (k K8sClient) GetVersionMapOfComponents() (map[string]string, error) {
	label := "openebs.io/component-name in ("
	for _, v := range util.CasTypeAndComponentNameMap {
		label = label + v + ","
	}
	label += ")"

	pods, err := k.K8sCS.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Running", LabelSelector: label})

	if err != nil {
		return nil, err
	}

	versionMap := make(map[string]string)

	for _, pod := range pods.Items {
		podLabels := pod.ObjectMeta.Labels
		labelName := podLabels["openebs.io/component-name"]
		version := podLabels["openebs.io/version"]

		cas, okCas := util.ComponentNameToCasTypeMap[labelName]
		if okCas {
			versionMap[cas] = version
		}
	}

	return versionMap, nil
}

/**
CORE RESOURCES
*/
// GetPods returns the corev1 Pods based on the label and field selectors
func (k K8sClient) GetPods(labelSelector string, fieldSelector string, namespace string) (*corev1.PodList, error) {
	pods, err := k.K8sCS.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector, FieldSelector: fieldSelector})
	if err != nil {
		return nil, fmt.Errorf("error getting pods : %v", err)
	}
	return pods, nil
}

// GetAllPods returns all corev1 Pods
func (k K8sClient) GetAllPods(namespace string) (*corev1.PodList, error) {
	pods, err := k.K8sCS.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting pods : %v", err)
	}
	return pods, nil
}

// GetSC returns a StorageClass object using the scName passed.
func (k K8sClient) GetSC(scName string) (*v1.StorageClass, error) {
	sc, err := k.K8sCS.StorageV1().StorageClasses().Get(context.TODO(), scName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "error while getting storage class")
	}
	return sc, nil
}

// GetCSIControllerSTS returns the CSI controller sts with a specific
// openebs-component-name label key
func (k K8sClient) GetCSIControllerSTS(name string) (*appsv1.StatefulSet, error) {
	if sts, err := k.K8sCS.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("openebs.io/component-name=%s", name),
	}); err == nil && len(sts.Items) == 1 {
		return &sts.Items[0], nil
	} else if sts != nil {
		return nil, fmt.Errorf("got %d statefulsets with the label openebs.io/component-name=%s", len(sts.Items), name)
	} else {
		return nil, fmt.Errorf("got 0 statefulsets with the label openebs.io/component-name=%s", name)
	}
}

// GetEvents returns the corev1 events based on the fieldSelectors
func (k K8sClient) GetEvents(fieldSelector string) (*corev1.EventList, error) {
	events, err := k.K8sCS.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{FieldSelector: fieldSelector})
	if err != nil {
		return nil, fmt.Errorf("error getting events for the resource : %v", err)
	}
	return events, nil
}

/*
	PERSISTENT VOLUMES AND CLAIMS
*/

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
	if len(volNames) == 0 {
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

// GetPvByCasType returns a list of PersistentVolumes based on cas-type slice
// casTypes slice if is nil or empty, it returns all the PVs in the cluster.
// casTypes slice if is not nil or not empty, it return the PVs with cas-types present in the slice.
// labelselector takes the label(key+value) and makes an api call with this filter applied. Can be empty string if label filtering is not needed.
func (k K8sClient) GetPvByCasType(casTypes []string, labelselector string) (*corev1.PersistentVolumeList, error) {
	pvs, err := k.K8sCS.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return nil, err
	}

	if len(casTypes) == 0 {
		return pvs, nil
	}

	var list []corev1.PersistentVolume

	for _, vol := range pvs.Items {
		for _, casType := range casTypes {
			if CSIProvisioner, ok := util.CasTypeToCSIProvisionerMap[casType]; ok {
				if vol.Spec.CSI != nil && vol.Spec.CSI.Driver == CSIProvisioner {
					list = append(list, vol)
				}
			}
		}
	}

	// No volumes with given cas-type found
	if len(list) == 0 {
		casTypesString := strings.Join(casTypes, ",")
		return nil, fmt.Errorf("couldn't find volumes of cas-type(s) %s", casTypesString)
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
	if len(pvcNames) == 0 {
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

// GetDeploymentList returns the deployment-list with a specific
// label selector query
func (k K8sClient) GetDeploymentList(labelSelector string) (*appsv1.DeploymentList, error) {
	if pv, err := k.K8sCS.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	}); err == nil && len(pv.Items) >= 1 {
		return pv, nil
	}
	return nil, fmt.Errorf("got 0 deployments with label-Selector as %s", labelSelector)
}

// GetNodes returns a list of nodes with the name of nodes
func (k K8sClient) GetNodes(nodes []string, label, field string) (*corev1.NodeList, error) {
	// 1. Get all nodes
	n, err := k.K8sCS.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field})
	if len(nodes) == 0 {
		return n, err
	}
	// 2. Put them in a map[string]corev1.Node
	nodeMap := make(map[string]corev1.Node)
	for _, item := range n.Items {
		nodeMap[item.Name] = item
	}
	// 3. Get the nodes by the name nodes
	var items []corev1.Node
	for _, name := range nodes {
		if _, ok := nodeMap[name]; ok {
			items = append(items, nodeMap[name])
		}
	}
	return &corev1.NodeList{
		Items: items,
	}, nil
}
