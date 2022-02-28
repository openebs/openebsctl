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

	jiva "github.com/openebs/jiva-operator/pkg/apis/openebs/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

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
	if len(volNames) == 0 {
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
		return nil, errors.New("The controller and replica pod for the volume was not found")
	}
	return pods, nil
}
