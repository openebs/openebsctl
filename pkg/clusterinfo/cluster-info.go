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

package clusterinfo

import (
	"fmt"
	"strings"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// ShowClusterInfo shows the openebs components and their status and versions
func ShowClusterInfo() error {
	k := client.NewK8sClient()
	err := compute(k)
	return err
}

func compute(k *client.K8sClient) error {
	var clusterInfoRows []metav1.TableRow
	for casType, componentNames := range util.CasTypeToComponentNamesMap {
		componentDataMap, err := getComponentDataByComponents(k, componentNames, casType)
		if err == nil && len(componentDataMap) != 0 {
			status, working := getStatus(componentDataMap)
			version := getVersion(componentDataMap)
			namespace := getNamespace(componentDataMap)
			clusterInfoRows = append(
				clusterInfoRows,
				metav1.TableRow{Cells: []interface{}{casType, namespace, version, working, util.ColorStringOnStatus(status)}},
			)
		}
	}
	if len(clusterInfoRows) == 0 {
		return fmt.Errorf("none Of the OpenEBS Storage Engines are installed in this cluster")
	}
	util.TablePrinter(util.ClusterInfoColumnDefinitions, clusterInfoRows, printers.PrintOptions{})
	return nil
}

func getComponentDataByComponents(k *client.K8sClient, componentNames string, casType string) (map[string]util.ComponentData, error) {
	var podList *corev1.PodList
	componentDataMap := make(map[string]util.ComponentData)
	podList, _ = k.GetPods(fmt.Sprintf("openebs.io/component-name in (%s)", componentNames), "", "")
	if len(podList.Items) != 0 {
		for _, item := range podList.Items {
			if val, ok := componentDataMap[item.Labels["openebs.io/component-name"]]; ok {
				// Update only if the status of the component is not running.
				if val.Status != string(v1.PodRunning) {
					componentDataMap[item.Labels["openebs.io/component-name"]] = util.ComponentData{
						Namespace: item.Namespace,
						Status:    string(item.Status.Phase),
						Version:   item.Labels["openebs.io/version"],
						CasType:   casType,
					}
				}
			} else {
				componentDataMap[item.Labels["openebs.io/component-name"]] = util.ComponentData{
					Namespace: item.Namespace,
					Status:    string(item.Status.Phase),
					Version:   item.Labels["openebs.io/version"],
					CasType:   casType,
				}
			}
		}

		for _, item := range strings.Split(componentNames, ",") {
			if _, ok := componentDataMap[item]; !ok {
				componentDataMap[item] = util.ComponentData{}
			}
		}

		return componentDataMap, nil
	}
	return nil, fmt.Errorf("components for %s engine are not installed", casType)
}

func getStatus(componentDataMap map[string]util.ComponentData) (string, string) {
	totalComponents := len(componentDataMap)
	healthyComponents := 0
	for _, val := range componentDataMap {
		if val.Status == string(v1.PodRunning) {
			healthyComponents++
		}
	}
	if healthyComponents == totalComponents {
		return "Healthy", fmt.Sprintf("%d/%d", healthyComponents, totalComponents)
	} else if healthyComponents < totalComponents && healthyComponents != 0 {
		return "Degraded", fmt.Sprintf("%d/%d", healthyComponents, totalComponents)
	} else {
		return "Unhealthy", fmt.Sprintf("%d/%d", 0, totalComponents)
	}
}

func getVersion(componentDataMap map[string]util.ComponentData) string {
	for _, val := range componentDataMap {
		if val.Version != "" && val.Status == "Running" {
			return val.Version
		}
	}
	return ""
}

func getNamespace(componentDataMap map[string]util.ComponentData) string {
	for _, val := range componentDataMap {
		if val.Namespace != "" {
			return val.Namespace
		}
	}
	return ""
}
