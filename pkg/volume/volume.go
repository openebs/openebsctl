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
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// Get manages various implementations of Volume listing
func Get(vols []string, openebsNS, casType string) error {
	// TODO: Prefer passing the client from outside
	k, _ := client.NewK8sClient("")
	// 1. Get a list of required PersistentVolumes
	var pvList *corev1.PersistentVolumeList
	if vols == nil {
		pvList, _ = k.GetPVs(nil, "")
	} else {
		pvList, _ = k.GetPVs(vols, "")
	}
	// TODO: (improvisation) Only call specific cas-functions for a
	// list-obj-by-name & if only 2-3 cas-exist
	var rows []metav1.TableRow
	// 2. Get more information about pvList volumes
	if work, ok := CasMap()[casType]; ok {
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

// CasList returns a list of functions by cas-types for volume listing
func CasList() []func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error) {
	// a good hack to implement immutable lists in Golang & also write tests for it
	return []func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error){GetJiva, GetCStor, GetZFSLocalPVs, GetLVMLocalPV}
}

// CasMap returns a map cas-types to functions for volume listing
func CasMap() map[string]func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error) {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, *corev1.PersistentVolumeList, string) ([]metav1.TableRow, error){
		util.JivaCasType:  GetJiva,
		util.CstorCasType: GetCStor,
		util.ZFSCasType:   GetZFSLocalPVs,
		util.LVMLocalPV:   GetLVMLocalPV,
	}
}
