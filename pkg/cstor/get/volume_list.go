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

package get

import (
	"fmt"
	"os"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/jiva-operator/pkg/apis/openebs/v1alpha1"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/cli-runtime/pkg/printers"
)

// RunVolumesList lists the volumes
func RunVolumesList(openebsNS, casType string, vols []string) error {
	k8sClient, err := client.NewK8sClient("")
	if err != nil {
		return err
	}
	// 1. Fetch all or required PVs
	var pvList *corev1.PersistentVolumeList
	if len(vols) == 0 {
		pvList, err = k8sClient.GetPVs(nil, "")
	} else {
		pvList, err = k8sClient.GetPVs(vols, "")
	}
	if err != nil {
		return err
	}
	// 2. Fetch all relevant volume CRs without worrying about openebsNS
	var (
		jvMap  map[string]v1alpha1.JivaVolume
		cvMap  map[string]v1.CStorVolume
		cvaMap map[string]v1.CStorVolumeAttachment
	)
	if casType == "" {
		// fetch all
		_, jvMap, _ = k8sClient.GetJVs(nil, util.Map, "", util.MapOptions{Key: util.Name})
		_, cvMap, _ = k8sClient.GetCVs(nil, util.Map, "", util.MapOptions{Key: util.Name})
		_, cvaMap, _ = k8sClient.GetCVAs(util.Map, "", util.MapOptions{Key: util.Label, LabelKey: "Volname"})
	} else if casType == util.JivaCasType {
		_, jvMap, _ = k8sClient.GetJVs(nil, util.Map, "", util.MapOptions{Key: util.Name})
	} else if casType == util.CstorCasType {
		_, cvMap, _ = k8sClient.GetCVs(nil, util.Map, "", util.MapOptions{Key: util.Name})
		_, cvaMap, _ = k8sClient.GetCVAs(util.Map, "", util.MapOptions{Key: util.Label, LabelKey: "Volname"})
	}
	var rows []metav1.TableRow
	// 3. Show the required ones
	for _, pv := range pvList.Items {
		name := pv.Name
		capacity := pv.Spec.Capacity.Storage()
		sc := pv.Spec.StorageClassName
		attached := pv.Status.Phase
		var attachedNode, storageVersion, customStatus, ns string
		// TODO: Estimate the cas-type and decide to print it out
		// Should all AccessModes be shown in a csv format, or the highest be displayed ROO < RWO < RWX?
		if pv.Spec.CSI != nil {
			// 2. For eligible PVs fetch the custom-resource to add more info
			if pv.Spec.CSI.Driver == util.CStorCSIDriver && (casType == util.CstorCasType || casType == "") {
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
			} else if pv.Spec.CSI.Driver == util.JivaCSIDriver && (casType == util.JivaCasType || casType == "") {
				jv, ok := jvMap[pv.Name]
				if !ok {
					_, _ = fmt.Fprintln(os.Stderr, "couldn't find jv "+pv.Name)
				}
				ns = jv.Namespace
				if openebsNS != "" && openebsNS != ns {
					continue
				}
				customStatus = jv.Status.Status // RW, RO, etc
				attachedNode = jv.Labels["nodeID"]
				storageVersion = jv.VersionDetails.Status.Current
			} else {
				// Skip non-CStor & non-Jiva options
				continue
			}
		} else {
			// Skip non CSI provisioned volumes
			continue
		}
		accessMode := pv.Spec.AccessModes[0]
		rows = append(rows, metav1.TableRow{
			Cells: []interface{}{
				ns, name, customStatus, storageVersion, capacity, sc, attached,
				accessMode, attachedNode},
		})
	}
	if len(rows) == 0 {
		if casType == "" {
			return fmt.Errorf("no cstor and/or jiva volumes found")
		}
	}
	util.TablePrinter(util.VolumeListColumnDefinations, rows, printers.PrintOptions{Wide: true})
	return nil
}
