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
	"github.com/fatih/color"
	"os"

	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// DebugCstorVolumeClaim is used to debug a cstor volume by calling various modules
func DebugCstorVolumeClaim(k *client.K8sClient, pvc *corev1.PersistentVolumeClaim, pv *corev1.PersistentVolume) error {
	// 1. Main Struture Creation which contains all cstor CRs, this structure will be passed accross all modules.
	var cstorResources util.CstorVolumeResources
	cstorResources.PVC = pvc
	cstorResources.PV = pv
	// 2. Fill in the available CRs
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
	// 3. Fill in the Pool and Blockdevice Details
	if sc != nil {
		cspc, _ := k.GetCSPC(sc.Parameters["cstorPoolCluster"])
		cstorResources.CSPC = cspc
		if cspc != nil {
			cspis, _ := k.GetCSPIs(nil, "openebs.io/cas-type=cstor,openebs.io/cstor-pool-cluster="+cspc.Name)
			cstorResources.CSPIs = cspis
			expectedBlockDevicesInPool := make(map[string]bool)
			// This map contains the list of BDs we specified at the time of Pool Creation
			for _, pool := range cspc.Spec.Pools {
				dataRaidGroups := pool.DataRaidGroups
				for _, dataRaidGroup := range dataRaidGroups {
					for _, bdName := range dataRaidGroup.GetBlockDevices() {
						expectedBlockDevicesInPool[bdName] = false
					}
				}
			}
			// This list contains the list of BDs which are actually present in the system.
			var presentBlockDevicesInPool []string
			for _, pool := range cspis.Items {
				raidGroupsInPool := pool.GetAllRaidGroups()
				for _, item := range raidGroupsInPool {
					presentBlockDevicesInPool = append(presentBlockDevicesInPool, item.GetBlockDevices()...)
				}
			}

			// Mark the present BDs are true.
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
	// 4. Call the resource showing module
	// TODO: Integration of all modules
	return nil
}

func resourceStatus(crs util.CstorVolumeResources) error {
	// 1. Fetch the total and usage details and humanize them
	var totalCapacity, usedCapacity, availableCapacity string
	totalCapacity = util.ConvertToIBytes(crs.PVC.Spec.Resources.Requests.Storage().String())
	usedCapacity = util.ConvertToIBytes(util.GetUsedCapacityFromCVR(crs.CVRs))
	// 2. Calculate the available capacity and usage percentage is used capacity is available
	if usedCapacity != "" {
		availableCapacity = util.GetAvailableCapacity(totalCapacity, usedCapacity)
		percentage := util.GetUsedPercentage(totalCapacity, usedCapacity)
		if percentage >= 80.00 {
			availableCapacity = color.HiRedString(availableCapacity)
		} else {
			availableCapacity = color.HiGreenString(availableCapacity)
		}
	}
	// 3. Display the usage status
	_, _ = fmt.Fprint(os.Stdout, "Volume Usage Stats:\n-------------------\n")

	util.TablePrinter(util.VolumeTotalAndUsageDetailColumnDefinitions, []metav1.TableRow{{Cells: []interface{}{totalCapacity, usedCapacity, availableCapacity}}}, printers.PrintOptions{})

	_, _ = fmt.Fprint(os.Stdout, "\nRelated CR Statuses:\n--------------------\n")

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
	// 4. Display the CRs statuses
	util.TablePrinter(util.CstorVolumeCRStatusColumnDefinitions, crStatusRows, printers.PrintOptions{})

	_, _ = fmt.Fprint(os.Stdout, "\nReplica Statuses:\n-----------------\n")
	crStatusRows = []metav1.TableRow{}
	if crs.CVRs != nil {
		for _, item := range crs.CVRs.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{item.Kind, item.Name, util.ColorStringOnStatus(string(item.Status.Phase))}})
		}
	}
	// 5. Display the CRs statuses
	util.TablePrinter(util.CstorVolumeCRStatusColumnDefinitions, crStatusRows, printers.PrintOptions{})

	_, _ = fmt.Fprint(os.Stdout, "\nBlockDevice and BlockDeviceClaim Statuses:\n------------------------------------------\n")
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
	// 6. Display the BDs and BDCs statuses
	util.TablePrinter(util.CstorVolumeCRStatusColumnDefinitions, crStatusRows, printers.PrintOptions{})

	_, _ = fmt.Fprint(os.Stdout, "\nPool Instance Statuses:\n-----------------------\n")
	crStatusRows = []metav1.TableRow{}
	if crs.CSPIs != nil {
		for _, item := range crs.CSPIs.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{item.Kind, item.Name, util.ColorStringOnStatus(string(item.Status.Phase))}})
		}
	}
	// 7. Display the Pool statuses
	util.TablePrinter(util.CstorVolumeCRStatusColumnDefinitions, crStatusRows, printers.PrintOptions{})

	return nil
}

func displayPVCEvents(k client.K8sClient, crs util.CstorVolumeResources) error {
	// 1. Set the namespace of the resource to the client
	k.Ns = crs.PVC.Namespace
	// 2. Fetch the events of the concerned PVC
	events, err := k.GetEvents(fmt.Sprintf("regarding.name=%s,regarding.kind=PersistentVolumeClaim", crs.PVC.Name))
	// 3. Display the events
	if len(events.Items) != 0 && err == nil {
		_, _ = fmt.Fprint(os.Stdout, "\nChecking PVC Events:", color.HiRedString(fmt.Sprintf(" %s %d! ", util.UnicodeCross, len(events.Items))), "-------->\n")
		var crStatusRows []metav1.TableRow
		for _, event := range events.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{event.Action, event.Reason, event.Note, util.ColorStringOnStatus(event.Type)}})
		}
		defer util.TablePrinter(util.EventsColumnDefinitions, crStatusRows, printers.PrintOptions{})
		return nil
	} else if len(events.Items) == 0 && err == nil {
		_, _ = fmt.Fprint(os.Stdout, "\nChecking PVC Events:", color.HiGreenString(fmt.Sprintf(" %s %d! ", util.UnicodeCheck, len(events.Items))), "-------->\n")
		return nil
	} else {
		return err
	}

}
