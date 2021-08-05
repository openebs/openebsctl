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

package persistentvolumeclaim

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// DebugCstorVolumeClaim is used to debug a cstor volume by calling various modules
func DebugCstorVolumeClaim(k *client.K8sClient, pvc *corev1.PersistentVolumeClaim, pv *corev1.PersistentVolume) error {
	var cstorResources util.CstorVolumeResources
	cstorResources.PVC, _ = k.GetPVC(pvc.Name, "default")
	cstorResources.PV = pv
	if pv != nil {
		cv, _ := k.GetCV(pv.Name)
		cstorResources.CV = cv
		cvc, _ := k.GetCVC(pv.Name)
		cstorResources.CVC = cvc
		cva, _ := k.GetCVA(util.CVAVolnameKey + "=" + pv.Name)
		cstorResources.CVA = cva
		cvrs, _ := k.GetCVRs(cstortypes.PersistentVolumeLabelKey + "=" + pv.Name)
		cstorResources.CVRs = cvrs
	}
	sc, _ := k.GetSC(*pvc.Spec.StorageClassName)
	if sc != nil {
		cspc, _ := k.GetCSPC(sc.Parameters["cstorPoolCluster"])
		cstorResources.CSPC = cspc
		if cspc != nil {
			cspis, _ := k.GetCSPIs(nil, "openebs.io/cas-type=cstor,openebs.io/cstor-pool-cluster="+cspc.Name)
			cstorResources.CSPIs = cspis
			expectedBlockDevicesInPool := make(map[string]bool)

			for _, pool := range cspc.Spec.Pools {
				dataRaidGroups := pool.DataRaidGroups
				for _, dataRaidGroup := range dataRaidGroups {
					for _, bdName := range dataRaidGroup.GetBlockDevices() {
						expectedBlockDevicesInPool[bdName] = false
					}
				}
			}

			var presentBlockDevicesInPool []string
			for _, pool := range cspis.Items {
				raidGroupsInPool := pool.GetAllRaidGroups()
				for _, item := range raidGroupsInPool {
					presentBlockDevicesInPool = append(presentBlockDevicesInPool, item.GetBlockDevices()...)
				}
			}

			cstorResources.PresentBDs, _ = k.GetBDs(presentBlockDevicesInPool, "")
			for _, item := range cstorResources.PresentBDs.Items {
				if _, ok := expectedBlockDevicesInPool[item.Name]; ok {
					expectedBlockDevicesInPool[item.Name] = true
				}
			}
			cstorResources.ExpectedBDs = expectedBlockDevicesInPool
			cstorResources.BDCs, _ = k.GetBDCs(nil, "openebs.io/cstor-pool-cluster="+cspc.Name)

		}
	}
	_ = resourceStatus(cstorResources)
	return nil
}

func resourceStatus(crs util.CstorVolumeResources) error {
	var totalCapacity, usedCapacity, availableCapacity string
	totalCapacity = util.ConvertToIBytes(crs.PVC.Spec.Resources.Requests.Storage().String())
	usedCapacity = util.ConvertToIBytes(util.GetUsedCapacityFromCVR(crs.CVRs))
	if usedCapacity != "" {
		availableCapacity = util.GetAvailableCapacity(totalCapacity, usedCapacity)
		percentage := util.GetUsedPercentage(totalCapacity, usedCapacity)
		if percentage >= 80.00 {
			availableCapacity = color.HiRedString(availableCapacity)
		} else {
			availableCapacity = color.HiGreenString(availableCapacity)
		}
	}
	_, _ = fmt.Fprint(os.Stdout, "Volume Usage Stats:\n-------------------\n")

	util.TablePrinter([]metav1.TableColumnDefinition{
		{Name: "Total Capacity", Type: "string"},
		{Name: "Used Capacity", Type: "string"},
		{Name: "Available Capacity", Type: "string"},
	}, []metav1.TableRow{{Cells: []interface{}{totalCapacity, usedCapacity, availableCapacity}}}, printers.PrintOptions{})

	_, _ = fmt.Fprint(os.Stdout, "\nRelated CR Statuses:\n-------------------\n")

	var crStatusRows []metav1.TableRow
	if crs.PV != nil {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"PersistentVolume", crs.PV.Name, util.ColorStringOnStatus(string(crs.PV.Status.Phase))}})
	} else {
		crStatusRows = append(
			crStatusRows,

			metav1.TableRow{Cells: []interface{}{"PersistentVolume", "", util.ColorStringOnStatus("Not Found")}},
		)
	}
	if crs.CV != nil {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolume", crs.CV.Name, util.ColorStringOnStatus(string(crs.CV.Status.Phase))}})
	} else {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolume", "", util.ColorStringOnStatus("Not Found")}})
	}
	if crs.CVC != nil {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolumeConfig", crs.CVC.Name, util.ColorStringOnStatus(string(crs.CVC.Status.Phase))}})
	} else {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolumeConfig", "", util.ColorStringOnStatus("Not Found")}})
	}
	if crs.CVA != nil {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolumeAttachment", crs.CVA.Name, util.ColorStringOnStatus("Attached")}})
	} else {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolumeAttachment", "", util.ColorStringOnStatus("Volume Not Attached to Application")}})
	}

	util.TablePrinter([]metav1.TableColumnDefinition{
		{Name: "Kind", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Status", Type: "string"},
	}, crStatusRows, printers.PrintOptions{})

	_, _ = fmt.Fprint(os.Stdout, "\nReplica Statuses:\n-------------------\n")
	crStatusRows = []metav1.TableRow{}
	if crs.CVRs != nil {
		for _, item := range crs.CVRs.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{item.Kind, item.Name, util.ColorStringOnStatus(string(item.Status.Phase))}})
		}
	}

	util.TablePrinter([]metav1.TableColumnDefinition{
		{Name: "Kind", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Status", Type: "string"},
	}, crStatusRows, printers.PrintOptions{})

	_, _ = fmt.Fprint(os.Stdout, "\nBlockDevice and BlockDeviceClaim Statuses:\n-------------------\n")
	crStatusRows = []metav1.TableRow{}
	if crs.PresentBDs != nil {
		for _, item := range crs.PresentBDs.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{item.Kind, item.Name, util.ColorStringOnStatus(string(item.Status.State))}})
		}
		for key, val := range crs.ExpectedBDs {
			if !val {
				crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"BlockDevice", key, util.ColorStringOnStatus("Not Found")}})
			}
		}
	}
	crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"", "", ""}})
	if crs.BDCs != nil {
		for _, item := range crs.BDCs.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{item.Kind, item.Name, util.ColorStringOnStatus(string(item.Status.Phase))}})
		}
	}

	util.TablePrinter([]metav1.TableColumnDefinition{
		{Name: "Kind", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Status", Type: "string"},
	}, crStatusRows, printers.PrintOptions{})

	_, _ = fmt.Fprint(os.Stdout, "\nPool Instance Statuses:\n-------------------\n")
	crStatusRows = []metav1.TableRow{}
	if crs.CSPIs != nil {
		for _, item := range crs.CSPIs.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{item.Kind, item.Name, util.ColorStringOnStatus(string(item.Status.Phase))}})
		}
	}

	util.TablePrinter([]metav1.TableColumnDefinition{
		{Name: "Kind", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Status", Type: "string"},
	}, crStatusRows, printers.PrintOptions{})

	return nil
}
