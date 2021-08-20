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
	"errors"
	"fmt"

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
	_ = resourceStatus(cstorResources)
	_ = displayPVCEvents(*k, cstorResources)
	_ = displayCVCEvents(*k, cstorResources)
	_ = displayCVREvents(*k, cstorResources)
	_ = displayCSPIEvents(*k, cstorResources)
	_ = displayCSPCEvents(*k, cstorResources)
	_ = displayBDCEvents(*k, cstorResources)
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
			availableCapacity = util.ColorText(availableCapacity, util.Red)
		} else {
			availableCapacity = util.ColorText(availableCapacity, util.Green)
		}
	}
	// 3. Display the usage status
	fmt.Println("Volume Usage Stats:")
	fmt.Println("-------------------")

	util.TablePrinter(util.VolumeTotalAndUsageDetailColumnDefinitions, []metav1.TableRow{{Cells: []interface{}{totalCapacity, usedCapacity, availableCapacity}}}, printers.PrintOptions{})

	fmt.Println()
	fmt.Println("Related CR Statuses:")
	fmt.Println("--------------------")

	var crStatusRows []metav1.TableRow
	if crs.PV != nil {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"PersistentVolume", crs.PV.Name, util.ColorStringOnStatus(string(crs.PV.Status.Phase))}})
	} else {
		crStatusRows = append(
			crStatusRows,
			metav1.TableRow{Cells: []interface{}{"PersistentVolume", "", util.ColorStringOnStatus(util.NotFound)}},
		)
	}
	if crs.CV != nil {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolume", crs.CV.Name, util.ColorStringOnStatus(string(crs.CV.Status.Phase))}})
	} else {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolume", "", util.ColorStringOnStatus(util.NotFound)}})
	}
	if crs.CVC != nil {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolumeConfig", crs.CVC.Name, util.ColorStringOnStatus(string(crs.CVC.Status.Phase))}})
	} else {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolumeConfig", "", util.ColorStringOnStatus(util.NotFound)}})
	}
	if crs.CVA != nil {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolumeAttachment", crs.CVA.Name, util.ColorStringOnStatus(util.Attached)}})
	} else {
		crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"CstorVolumeAttachment", "", util.ColorStringOnStatus(util.CVANotAttached)}})
	}
	// 4. Display the CRs statuses
	util.TablePrinter(util.CstorVolumeCRStatusColumnDefinitions, crStatusRows, printers.PrintOptions{})

	fmt.Println()
	fmt.Println("Replica Statuses:")
	fmt.Println("-----------------")
	crStatusRows = []metav1.TableRow{}
	if crs.CVRs != nil {
		for _, item := range crs.CVRs.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{item.Kind, item.Name, util.ColorStringOnStatus(string(item.Status.Phase))}})
		}
	}
	// 5. Display the CRs statuses
	util.TablePrinter(util.CstorVolumeCRStatusColumnDefinitions, crStatusRows, printers.PrintOptions{})

	fmt.Println()
	fmt.Println("BlockDevice and BlockDeviceClaim Statuses:")
	fmt.Println("------------------------------------------")
	crStatusRows = []metav1.TableRow{}
	if crs.PresentBDs != nil {
		for _, item := range crs.PresentBDs.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{item.Kind, item.Name, util.ColorStringOnStatus(string(item.Status.State))}})
		}
		for key, val := range crs.ExpectedBDs {
			if !val {
				crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{"BlockDevice", key, util.ColorStringOnStatus(util.NotFound)}})
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

	fmt.Println()
	fmt.Println("Pool Instance Statuses:")
	fmt.Println("-----------------------")
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
	// 2. Fetch the events of the concerned PVC.
	// The PVCs donot have the Kind filled, thus we have hardcoded here.
	events, err := k.GetEvents(fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=PersistentVolumeClaim", crs.PVC.Name))
	// 3. Display the events
	fmt.Println()
	if err == nil && len(events.Items) != 0 {
		fmt.Println("Checking PVC Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCross, len(events.Items)), util.Red), "-------->")
		var crStatusRows []metav1.TableRow
		for _, event := range events.Items {
			crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{event.InvolvedObject.Name, event.Action, event.Reason, event.Message, util.ColorStringOnStatus(event.Type)}})
		}
		util.TablePrinter(util.EventsColumnDefinitions, crStatusRows, printers.PrintOptions{})
		return nil
	} else if err == nil && len(events.Items) == 0 {
		fmt.Println("Checking PVC Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCheck, len(events.Items)), util.Green), "-------->")
		return nil
	} else {
		return err
	}
}

func displayBDCEvents(k client.K8sClient, crs util.CstorVolumeResources) error {
	if crs.BDCs != nil && len(crs.BDCs.Items) != 0 {
		// 1. Set the namespace of the resource to the client
		k.Ns = crs.BDCs.Items[0].Namespace
		// 2. Fetch the events of the concerned BDC
		fmt.Println()
		var crStatusRows []metav1.TableRow
		for _, BDC := range crs.BDCs.Items {
			events, err := k.GetEvents(fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=BlockDeviceClaim", BDC.Name))
			// 3. Display the events
			if err == nil && len(events.Items) != 0 {
				for _, event := range events.Items {
					crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{event.InvolvedObject.Name, event.Action, event.Reason, event.Message, util.ColorStringOnStatus(event.Type)}})
				}
			}
		}
		if len(crStatusRows) == 0 {
			fmt.Println("Checking BDC Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCheck, len(crStatusRows)), util.Green), "-------->")
			return nil
		} else {
			fmt.Println("Checking BDC Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCross, len(crStatusRows)), util.Red), "-------->")
			util.TablePrinter(util.EventsColumnDefinitions, crStatusRows, printers.PrintOptions{})
			return nil
		}
	}
	return errors.New("no BDC present to display events")
}

func displayCVCEvents(k client.K8sClient, crs util.CstorVolumeResources) error {
	if crs.CVC != nil {
		// 1. Set the namespace of the resource to the client
		k.Ns = crs.CVC.Namespace
		// 2. Fetch the events of the concerned CVC
		events, err := k.GetEvents(fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=CStorVolumeConfig", crs.CVC.Name))
		// 3. Display the events
		fmt.Println()
		if err == nil && len(events.Items) != 0 {
			fmt.Println("Checking CVC Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCross, len(events.Items)), util.Red), "-------->")
			var crStatusRows []metav1.TableRow
			for _, event := range events.Items {
				crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{event.InvolvedObject.Name, event.Action, event.Reason, event.Message, util.ColorStringOnStatus(event.Type)}})
			}
			defer util.TablePrinter(util.EventsColumnDefinitions, crStatusRows, printers.PrintOptions{})
			return nil
		} else if err == nil && len(events.Items) == 0 {
			fmt.Println("Checking CVC Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCheck, len(events.Items)), util.Green), "-------->")
			return nil
		} else {
			return err
		}
	} else {
		return errors.New("no CVC present to display events")
	}
}

func displayCSPCEvents(k client.K8sClient, crs util.CstorVolumeResources) error {
	if crs.CSPC != nil {
		// 1. Set the namespace of the resource to the client
		k.Ns = crs.PVC.Namespace
		// 2. Fetch the events of the concerned PVC.
		// The PVCs donot have the Kind filled, thus we have hardcoded here.
		events, err := k.GetEvents(fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=CStorPoolCluster", crs.PVC.Name))
		// 3. Display the events
		fmt.Println()
		if err == nil && len(events.Items) != 0 {
			fmt.Println("Checking CSPC Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCross, len(events.Items)), util.Red), "-------->")
			var crStatusRows []metav1.TableRow
			for _, event := range events.Items {
				crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{event.InvolvedObject.Name, event.Action, event.Reason, event.Message, util.ColorStringOnStatus(event.Type)}})
			}
			util.TablePrinter(util.EventsColumnDefinitions, crStatusRows, printers.PrintOptions{})
			return nil
		} else if err == nil && len(events.Items) == 0 {
			fmt.Println("Checking CSPC Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCheck, len(events.Items)), util.Green), "-------->")
			return nil
		} else {
			return err
		}
	} else {
		return errors.New("no CSPC present to display events")
	}
}

func displayCSPIEvents(k client.K8sClient, crs util.CstorVolumeResources) error {
	if crs.CSPIs != nil && len(crs.CSPIs.Items) != 0 {
		// 1. Set the namespace of the resource to the client
		k.Ns = crs.CSPIs.Items[0].Namespace
		// 2. Fetch the events of the concerned CSPIs
		fmt.Println()
		var crStatusRows []metav1.TableRow
		for _, CSPI := range crs.CSPIs.Items {
			events, err := k.GetEvents(fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=CStorPoolInstance", CSPI.Name))
			// 3. Display the events
			if err == nil && len(events.Items) != 0 {
				for _, event := range events.Items {
					crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{event.InvolvedObject.Name, event.Action, event.Reason, event.Message, util.ColorStringOnStatus(event.Type)}})
				}
			}
		}
		if len(crStatusRows) == 0 {
			fmt.Println("Checking CSPI Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCheck, len(crStatusRows)), util.Green), "-------->")
			return nil
		} else {
			fmt.Println("Checking CSPI Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCross, len(crStatusRows)), util.Red), "-------->")
			util.TablePrinter(util.EventsColumnDefinitions, crStatusRows, printers.PrintOptions{})
			return nil
		}
	}
	return errors.New("no CSPIs present to display events")
}

func displayCVREvents(k client.K8sClient, crs util.CstorVolumeResources) error {
	if crs.CVRs != nil && len(crs.CVRs.Items) != 0 {
		// 1. Set the namespace of the resource to the client
		k.Ns = crs.CVRs.Items[0].Namespace
		// 2. Fetch the events of the concerned CVRs
		fmt.Println()
		var crStatusRows []metav1.TableRow
		for _, CVR := range crs.CVRs.Items {
			events, err := k.GetEvents(fmt.Sprintf("involvedObject.name=%s,involvedObject.kind=CStorVolumeReplica", CVR.Name))
			// 3. Display the events
			if err == nil && len(events.Items) != 0 {
				for _, event := range events.Items {
					crStatusRows = append(crStatusRows, metav1.TableRow{Cells: []interface{}{event.InvolvedObject.Name, event.Action, event.Reason, event.Message, util.ColorStringOnStatus(event.Type)}})
				}
			}
		}
		if len(crStatusRows) == 0 {
			fmt.Println("Checking CVR Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCheck, len(crStatusRows)), util.Green), "-------->")
			return nil
		} else {
			fmt.Println("Checking CVR Events:", util.ColorText(fmt.Sprintf(" %s %d! ", util.UnicodeCross, len(crStatusRows)), util.Red), "-------->")
			util.TablePrinter(util.EventsColumnDefinitions, crStatusRows, printers.PrintOptions{})
			return nil
		}
	}
	return errors.New("no CVRs present to display events")
}
