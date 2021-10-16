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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	lvmclient "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset"
	"github.com/openebs/openebsctl/pkg/util"
	zfsclient "github.com/openebs/zfs-localpv/pkg/generated/clientset/internalclientset"
	"github.com/pkg/errors"

	openebsclientset "github.com/openebs/api/v2/pkg/client/clientset/versioned"

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
	// OpenebsClientset capable of accessing the OpenEBS
	// components
	OpenebsCS openebsclientset.Interface
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
	openebsCS, err := getOpenEBSClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build OpenEBS clientset")
	}
	lv, _ := getLVMclient(config)
	zf, _ := getZFSclient(config)
	return &K8sClient{
		Ns:        ns,
		K8sCS:     k8sCS,
		OpenebsCS: openebsCS,
		LVMCS:     lv,
		ZFCS:      zf,
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
// NOTE: This will not work correctly if CSI controller pod runs in kube-system NS
func (k K8sClient) GetOpenEBSNamespace(casType string) (string, error) {
	pods, err := k.K8sCS.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{FieldSelector: "status.phase=Running", LabelSelector: fmt.Sprintf("openebs.io/component-name=%s", util.CasTypeAndComponentNameMap[strings.ToLower(casType)])})
	if err != nil || len(pods.Items) == 0 {
		return "", errors.New("unable to determine openebs namespace")
	}
	return pods.Items[0].Namespace, nil
}

// GetOpenEBSNamespaceMap maps the cas-type to it's namespace, e.g. n[cstor] = cstor-ns
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
