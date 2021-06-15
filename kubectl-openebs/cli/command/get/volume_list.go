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
	"strings"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/jiva-operator/pkg/apis/openebs/v1alpha1"
	"github.com/openebs/openebsctl/client"
	"github.com/openebs/openebsctl/kubectl-openebs/cli/util"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	volumesListCommandHelpText = `
This command displays status of available zfs Volumes.
If no volume ID is given, a list of all known volumes will be displayed.

Usage: kubectl openebs get volume [options]
`
)

// NewCmdGetVolume displays status of OpenEBS Volume(s)
func NewCmdGetVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "volume",
		Aliases: []string{"vol", "v", "volumes"},
		Short:   "Displays status information about Volume(s)",
		Long:    volumesListCommandHelpText,
		Run: func(cmd *cobra.Command, args []string) {
			casType, _ := cmd.Flags().GetString("cas-type")
			casType = strings.ToLower(casType)
			util.CheckErr(RunVolumesList(casType, args), util.Fatal)
		},
	}
	return cmd
}

// RunVolumesList lists the volumes
func RunVolumesList(casType string, vols []string) error {
	k8sClient, err := client.NewK8sClient("")
	if err != nil {
		return err
	}
	// 1. Fetch all or required PVs
	var pvList *corev1.PersistentVolumeList
	if len(vols) == 0 {
		pvList, err = k8sClient.GetPVs()
	} else {
		pvList, err = k8sClient.GetPVbyName(vols)
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
		jvMap, _ = k8sClient.GetJivaVolumeMap()
		cvMap, _ = k8sClient.GetCStorVolumeMap()
		cvaMap, _ = k8sClient.GetCStorVolumeAttachmentMap()
	} else if casType == util.JivaCasType {
		jvMap, _ = k8sClient.GetJivaVolumeMap()
	} else if casType == util.CstorCasType {
		cvMap, _ = k8sClient.GetCStorVolumeMap()
		cvaMap, _ = k8sClient.GetCStorVolumeAttachmentMap()
	}
	var rows []metav1.TableRow
	// 3. Show the required ones
	for _, pv := range pvList.Items {
		name := pv.Name
		capacity := pv.Spec.Capacity.Storage()
		sc := pv.Spec.StorageClassName
		var attached, attachedNode, storageVersion, customStatus, ns string
		// TODO: Estimate the cas-type and decide to print it out
		// Should all AccessModes be shown in a csv format, or the highest be displayed ROO < RWO < RWX?
		if pv.Spec.CSI != nil {
			// 2. For eligible PVs fetch the custom-resource to add more info
			if pv.Spec.CSI.Driver == util.CStorCSIDriver && (casType == util.CstorCasType || casType == "") {
				// For all CSI CStor PV, there exist a CV
				cv, ok := cvMap[pv.Name]
				if !ok {
					// condition not possible
					return fmt.Errorf("couldn't find cv %s", pv.Name)
				}
				ns = cv.Namespace
				customStatus = string(cv.Status.Phase)
				storageVersion = cv.VersionDetails.Status.Current
				cva, cvaOk := cvaMap[pv.Name]
				if cvaOk {
					attachedNode = cva.Labels["nodeID"]
				}
			} else if pv.Spec.CSI.Driver == util.JivaCSIDriver && (casType == util.JivaCasType || casType == "") {
				jv, ok := jvMap[pv.Name]
				if !ok {
					return fmt.Errorf("couldn't find jv %s", pv.Name)
				}
				ns = jv.Namespace
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
