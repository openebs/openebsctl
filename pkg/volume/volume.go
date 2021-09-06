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

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// Get manages various implementations of Volume listing
func Get(vols []string, openebsNS, casType string) error {
	if casType != "" && !util.IsValidCasType(casType) {
		return fmt.Errorf("cas-type %s is not supported", casType)
	}
	// TODO: Prefer passing the client from outside
	k, _ := client.NewK8sClient("")
	// 1. Get a list of required PersistentVolumes
	var pvList *corev1.PersistentVolumeList
	var err error
	if vols == nil {
		pvList, err = k.GetPVs(nil, "")
	} else {
		pvList, err = k.GetPVs(vols, "")
	}
	if err != nil {
		// stop if no PVs found
		return err
	}
	// TODO: (improvisation) Only call specific cas-functions for a
	// list-obj-by-name & if only 2-3 cas-exist
	var rows []metav1.TableRow
	// 2. Get more information about pvList volumes
	if work, ok := CasListMap()[casType]; ok {
		var err error
		if rows, err = work(k, pvList, openebsNS); err != nil {
			return err
		}
	} else {
		for _, t := range CasList() {
			if jr, err := t(k, pvList, openebsNS); err == nil {
				rows = append(rows, jr...)
			}
		}
	}
	// 3. Print the volumes from rows
	util.TablePrinter(util.VolumeListColumnDefinations, rows, printers.PrintOptions{Wide: true})
	return nil
}

// Describe manages various implementations of Volume Describing
func Describe(vols []string, openebsNs string) error {
	if vols == nil {
		return errors.New("please provide atleast one pv name to describe")
	}
	// Clienset creation
	k, _ := client.NewK8sClient(openebsNs)

	// 1. Get a list of required PersistentVolumes
	var pvList *corev1.PersistentVolumeList
	pvList, err := k.GetPVs(vols, "")
	if err != nil {
		return errors.New("no volumes found corresponding to the names")
	}
	// 2. Get the namespaces
	nsMap, _ := k.GetOpenEBSNamespaceMap()
	// 3. Range over the list of PVs
	for _, pv := range pvList.Items {
		// 4. Fetch the storage class, used to get the cas-type
		//TODO: Add cas-type label in every storage engine pv
		sc, err := k.GetSC(pv.Spec.StorageClassName)
		// 5. Get cas type
		casType := ""
		if err != nil {
			casType = util.GetCasTypeFromPV(&pv)
		} else {
			casType = util.GetCasType(&pv, sc)
		}
		// 6. Assign a namespace corresponding to the engine
		if openebsNs == "" {
			if val, ok := nsMap[casType]; ok {
				k.Ns = val
			} else if casType != util.ZFSCasType && casType != util.LVMCasType {
				// The reason for above condition is that, newer lvm has cas ty
				return errors.New("could not determine the underlying storage engine ns, please provide using '--openebs-namespace' flag")
			}
		}
		// 7. Describe the volume based on its casType
		if desc, ok := CasDescribeMap()[casType]; ok {
			err = desc(k, &pv)
			if err != nil {
				continue
			}
		}
	}
	return nil
}

// CasList returns a list of functions by cas-types for volume listing
func CasList() []func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error) {
	// a good hack to implement immutable lists in Golang & also write tests for it
	return []func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error){GetJiva, GetCStor, GetZFSLocalPVs, GetLVMLocalPV}
}

// CasListMap returns a map cas-types to functions for volume listing
func CasListMap() map[string]func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error) {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error){
		util.JivaCasType:  GetJiva,
		util.CstorCasType: GetCStor,
		util.ZFSCasType:   GetZFSLocalPVs,
		util.LVMCasType:   GetLVMLocalPV,
	}
}

// CasDescribeMap returns a map cas-types to functions for volume describing
func CasDescribeMap() map[string]func(*client.K8sClient, *corev1.PersistentVolume) error {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, *corev1.PersistentVolume) error{
		util.JivaCasType:  DescribeJivaVolume,
		util.CstorCasType: DescribeCstorVolume,
		util.ZFSCasType:   DescribeZFSLocalPVs,
		util.LVMCasType:   DescribeLVMLocalPVs,
	}
}
