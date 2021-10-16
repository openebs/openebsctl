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

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// GetPods returns the corev1 Pods based on the label and field selectors
func (k K8sClient) GetPods(labelSelector string, fieldSelector string, namespace string) (*corev1.PodList, error) {
	pods, err := k.K8sCS.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector, FieldSelector: fieldSelector})
	if err != nil {
		return nil, fmt.Errorf("error getting pods : %v", err)
	}
	return pods, nil
}

// GetDeploymentList returns the deployment-list with a specific
// label selector query
func (k K8sClient) GetDeploymentList(labelSelector string) (*appsv1.DeploymentList, error) {
	if pv, err := k.K8sCS.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	}); err == nil && len(pv.Items) >= 1 {
		return pv, nil
	} else {
		return nil, fmt.Errorf("got 0 deployments with label-Selector as %s", labelSelector)
	}
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
