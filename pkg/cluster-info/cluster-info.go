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

package cluster_info

import (
	"fmt"
	"strings"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// ShowClusterInfo shows the openebs components and their status and versions
func ShowClusterInfo() error {
	k, _ := client.NewK8sClient("")
	var clusterInfoRows []metav1.TableRow
	for casType, componentNames := range util.CasTypeToComponentNamesMap {
		componentDataMap, _ := getComponentDataByComponents(k, componentNames, casType)
		if len(componentDataMap) != 0 {
			status, working := "", ""
			if casType == util.LocalDeviceCasType {
				var err error
				status, working, err = getLocalPVDeviceStatus(componentDataMap)
				if err != nil {
					continue
				}
			} else {
				status, working = getStatus(componentDataMap)
			}
			version := getVersion(componentDataMap)
			namespace := getNamespace(componentDataMap)
			clusterInfoRows = append(
				clusterInfoRows,
				metav1.TableRow{Cells: []interface{}{casType, namespace, version, working, util.ColorStringOnStatus(status)}},
			)
		}
	}
	if len(clusterInfoRows) == 0 {
		fmt.Println("None Of the OpenEBS Storage Engines are installed in this cluster")
	} else {
		util.TablePrinter(util.ClusterInfoColumnDefinitions, clusterInfoRows, printers.PrintOptions{})
	}
	return nil
}

func getComponentDataByComponents(k *client.K8sClient, componentNames string, casType string) (map[string]util.ComponentData, error) {
	var podList *corev1.PodList
	// Fetch Cstor Components
	componentDataMap := make(map[string]util.ComponentData)
	podList, _ = k.GetPods(fmt.Sprintf("openebs.io/component-name in (%s)", componentNames), "", "")
	if len(podList.Items) != 0 {
		for _, item := range podList.Items {
			if val, ok := componentDataMap[item.Labels["openebs.io/component-name"]]; ok {
				// Update only if the status of the component is not running.
				if val.Status != "Running" {
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
		return componentDataMap, nil
	} else {
		return nil, fmt.Errorf("components for %s engine are not installed", casType)
	}
}

func getStatus(componentDataMap map[string]util.ComponentData) (string, string) {
	totalComponents := len(componentDataMap)
	healthyComponents := 0
	for _, val := range componentDataMap {
		if val.Status == "Running" {
			healthyComponents += 1
		}
	}
	if healthyComponents == totalComponents {
		return "Healthy", fmt.Sprintf("%d/%d", healthyComponents, totalComponents)
	} else if healthyComponents < totalComponents {
		return "Degraded", fmt.Sprintf("%d/%d", healthyComponents, totalComponents)
	} else {
		return "Unhealthy", fmt.Sprintf("%d/%d", 0, totalComponents)
	}
}

func getLocalPVDeviceStatus(componentDataMap map[string]util.ComponentData) (string, string, error) {
	if ndmData, ok := componentDataMap["ndm"]; ok {
		if localPVData, ok := componentDataMap["openebs-localpv-provisioner"]; ok {
			if ndmData.Namespace == localPVData.Namespace {
				status, working := getStatus(componentDataMap)
				return status, working, nil
			}
		}
	}
	return "", "", fmt.Errorf("installed NDM is not for Device LocalPV")
}

func getVersion(componentDataMap map[string]util.ComponentData) string {
	for key, val := range componentDataMap {
		if !strings.Contains(util.NDMComponentNames, key) && val.Version != "" && val.Status == "Running"{
			return val.Version
		}
	}
	return ""
}

func getNamespace(componentDataMap map[string]util.ComponentData) string {
	for key, val := range componentDataMap {
		if !strings.Contains(util.NDMComponentNames, key) && val.Namespace != "" {
			return val.Namespace
		}
	}
	return ""
}
