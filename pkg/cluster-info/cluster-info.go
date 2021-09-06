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
	err := compute(k)
	return err
}

func compute(k *client.K8sClient) error {
	var clusterInfoRows []metav1.TableRow
	for casType, componentNames := range util.CasTypeToComponentNamesMap {
		componentDataMap, err := getComponentDataByComponents(k, componentNames, casType)
		if err == nil && len(componentDataMap) != 0 {
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
		return fmt.Errorf("none Of the OpenEBS Storage Engines are installed in this cluster")
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

		// Below is to handle corner cases in case of cstor and jiva, as they use components like ndm and localpv provisioner
		// which are also used by other engine, below has been added to strictly identify the installed engine.
		engineComponents := 0
		for key := range componentDataMap {
			if !strings.Contains(util.NDMComponentNames, key) && strings.Contains(util.CasTypeToComponentNamesMap[casType], key) {
				engineComponents += 1
				if casType == util.JivaCasType && key == util.HostpathComponentNames {
					// Since hostpath component is not a unique engine component for jiva
					engineComponents -= 1
				}
			}
		}
		if engineComponents == 0 {
			return nil, fmt.Errorf("components for %s engine are not installed", casType)
		}

		// The below is to fill in the expected components, for example if 5 out of 7 cstor components are there
		// in the cluster, we would not be able what was the expected number of components, the below would ensure cstor
		// needs 7 component always to work.
		for _, item := range strings.Split(componentNames, ",") {
			if _, ok := componentDataMap[item]; !ok {
				componentDataMap[item] = util.ComponentData{}
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
	} else if healthyComponents < totalComponents && healthyComponents != 0 {
		return "Degraded", fmt.Sprintf("%d/%d", healthyComponents, totalComponents)
	} else {
		return "Unhealthy", fmt.Sprintf("%d/%d", 0, totalComponents)
	}
}

func getLocalPVDeviceStatus(componentDataMap map[string]util.ComponentData) (string, string, error) {
	if ndmData, ok := componentDataMap["ndm"]; ok {
		if localPVData, ok := componentDataMap["openebs-localpv-provisioner"]; ok {
			if ndmData.Namespace == localPVData.Namespace && localPVData.Namespace != "" && localPVData.CasType != "" {
				status, working := getStatus(componentDataMap)
				return status, working, nil
			}
		}
	}
	return "", "", fmt.Errorf("installed NDM is not for Device LocalPV")
}

func getVersion(componentDataMap map[string]util.ComponentData) string {
	for key, val := range componentDataMap {
		if !strings.Contains(util.NDMComponentNames, key) && val.Version != "" && val.Status == "Running" {
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
