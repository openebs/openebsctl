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

package volume

import (
	"fmt"
	"os"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetCStor returns a list of CStor volumes
func GetCStor(c *client.K8sClient, pvList *corev1.PersistentVolumeList, openebsNS string) ([]metav1.TableRow, error) {
	var (
		cvMap  map[string]v1.CStorVolume
		cvaMap map[string]v1.CStorVolumeAttachment
	)
	// no need to proceed if CVs/CVAs don't exist
	var err error
	_, cvMap, err = c.GetCVs(nil, util.Map, "", util.MapOptions{
		Key: util.Name})
	if err != nil {
		return nil, fmt.Errorf("failed to list CStorVolumes")
	}
	_, cvaMap, err = c.GetCVAs(util.Map, "", util.MapOptions{
		Key:      util.Label,
		LabelKey: "Volname"})
	if err != nil {
		return nil, fmt.Errorf("failed to list CStorVolumeAttachments")
	}
	var rows []metav1.TableRow
	// 3. Show the required ones
	for _, pv := range pvList.Items {
		var attachedNode, storageVersion, customStatus, ns string
		// TODO: Estimate the cas-type and decide to print it out
		// Should all AccessModes be shown in a csv format, or the highest be displayed ROO < RWO < RWX?
		// 2. For eligible PVs fetch the custom-resource to add more info
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == util.CStorCSIDriver {
			// For all CSI CStor PV, there exist a CV
			cv, ok := cvMap[pv.Name]
			if !ok {
				// condition not possible
				_, _ = fmt.Fprintf(os.Stderr, "couldn't find cv "+pv.Name)
			}
			ns = cv.Namespace
			if openebsNS != "" && openebsNS != ns {
				continue
			}
			customStatus = string(cv.Status.Phase)
			storageVersion = cv.VersionDetails.Status.Current
			cva, cvaOk := cvaMap[pv.Name]
			if cvaOk {
				attachedNode = cva.Labels["nodeID"]
			}
			// TODO: What should be done for multiple AccessModes
			accessMode := pv.Spec.AccessModes[0]
			rows = append(rows, metav1.TableRow{
				Cells: []interface{}{
					ns, pv.Name, customStatus, storageVersion, pv.Spec.Capacity.Storage(), pv.Spec.StorageClassName, pv.Status.Phase,
					accessMode, attachedNode}})
		}
	}
	return rows, nil
}
