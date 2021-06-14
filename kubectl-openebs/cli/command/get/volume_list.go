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
			openebsNs, _ := cmd.Flags().GetString("openebs-namespace")
			casType, _ := cmd.Flags().GetString("cas-type")
			casType = strings.ToLower(casType)
			util.CheckErr(RunVolumesList(openebsNs, casType, args), util.Fatal)
		},
	}
	return cmd
}

// RunVolumesList lists the volumes
func RunVolumesList(openebsNs, casType string, vols []string) error {
	k8sClient, err := client.NewK8sClient(openebsNs)
	if err != nil {
		return err
	}
	var nsMap map[string]string
	// 0. Figure out the openebsNs & casType mess!!
	if openebsNs == "" && casType != "" {
		openebsNs, err := k8sClient.GetOpenEBSNamespace(casType)
		if err == nil {
			// TODO: Verbose log for this estimated namespace
			k8sClient.Ns = openebsNs
		} else {
			return fmt.Errorf("couldn't figure out openebs-namespace for cas-type=%s\nuse \"--openebs-namespace\" flag to provide the namespace", casType)
		}
	} else if casType == "" {
		// show all volumes
		nsMap, err = k8sClient.GetOpenEBSNamespaceMap()
		if err != nil {
			return fmt.Errorf("couldn't figure out the namespace for jiva & cstor")
		}
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
	var rows []metav1.TableRow
	for _, pv := range pvList.Items {
		name := pv.Name
		capacity := pv.Spec.Capacity.Storage()
		sc := pv.Spec.StorageClassName
		attached := pv.Status.Phase
		attachedNode := "N/A"
		storageVersion := "N/A"
		customStatus := "N/A"
		ns := "N/A"
		// TODO: Estimate the cas-type and decide to print it out
		// Should all AccessModes be shown in a csv format, or the highest be displayed ROO < RWO < RWX?
		if pv.Spec.CSI != nil {
			// 2. For eligible PVs fetch the custom-resource to add more info
			if pv.Spec.CSI.Driver == util.CStorCSIDriver && (casType == util.CstorCasType || casType == "") {
				if openebsNs == "" && nsMap != nil {
					if val, ok := nsMap[util.CstorCasType]; ok {
						k8sClient.Ns = val
					} else {
						// cstor CSI pod doesn't exist
						continue
					}
				}
				cv, err := k8sClient.GetcStorVolume(pv.Name)
				if err == nil {
					ns = cv.Namespace
					customStatus = string(cv.Status.Phase)
					storageVersion = cv.VersionDetails.Status.Current
					cva, err := k8sClient.GetCVA(pv.Name)
					if err == nil {
						attachedNode = cva.Labels["nodeID"]
					}
				}
			} else if pv.Spec.CSI.Driver == util.JivaCSIDriver && (casType == util.JivaCasType || casType == "") {
				if openebsNs == "" && nsMap != nil {
					if val, ok := nsMap[util.JivaCasType]; ok {
						k8sClient.Ns = val
					} else {
						// jiva CSI pod doesn't exist
						continue
					}
				}
				jv, err := k8sClient.GetJivaVolume(pv.Name)
				if err == nil {
					ns = jv.Namespace
					customStatus = jv.Status.Status // RW, RO, etc
					attachedNode = jv.Labels["nodeID"]
					storageVersion = jv.VersionDetails.Status.Current
				}
			} else {
				// Skip non-CStor & non-Jiva options
				continue
			}
		} else {
			// Skip non CSI provisoned volumes
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
