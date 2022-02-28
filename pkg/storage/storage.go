/*
Copyright 2020-2022 The OpenEBS Authors

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

package storage

import (
	"errors"
	"fmt"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// Get manages various implementations of Storage listing
func Get(pools []string, openebsNS string, casType string) error {
	// 1. Create the clientset
	k := client.NewK8sClient()
	// 2. If casType is specified, call the specific function & exit
	if f, ok := CasListMap()[casType]; ok {
		// if a cas-type is found, run it and return the error
		header, rows, err := f(k, pools)
		if err != nil {
			return err
		}
		if len(rows) == 0 {
			return util.HandleEmptyTableError("Storage", openebsNS, casType)
		}
		util.TablePrinter(header, rows, printers.PrintOptions{Wide: true})
	} else if casType != "" {
		return fmt.Errorf("cas-type %s is not supported", casType)
	} else {
		storageResourcesFound := false
		// 3. Call all functions & exit
		for _, f := range CasList() {
			header, row, err := f(k, pools)
			if err == nil {
				if len(row) > 0 {
					storageResourcesFound = true
				}
				// 4. Find the correct heading & print the rows
				util.TablePrinter(header, row, printers.PrintOptions{Wide: true})
				// A visual separator for different cas-type pools/storage entities
				fmt.Println()
			}
		}

		if !storageResourcesFound {
			return util.HandleEmptyTableError("Storage", openebsNS, casType)
		}
	}
	return nil
}

// CasList has a list of method implementations for different cas-types
func CasList() []func(*client.K8sClient, []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error) {
	return []func(*client.K8sClient, []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error){
		GetCstorPools, GetVolumeGroups, GetZFSPools}
}

// Describe manages various implementations of Storage Describing
func Describe(storages []string, openebsNs, casType string) error {
	if len(storages) == 0 || storages == nil {
		return errors.New("please provide atleast one pv name to describe")
	}
	// 1. Create the clientset
	k := client.NewK8sClient(openebsNs)
	// 2. Get the namespace
	nsMap, _ := k.GetOpenEBSNamespaceMap()
	if openebsNs == "" {
		if casType == util.ZFSCasType {
			// a temporary way to get the zfs-namespace
			zfs, _, err := k.GetZFSNodes(nil, util.List, "", util.MapOptions{})
			if err != nil {
				return fmt.Errorf("please specify --openebs-namespace for ZFS LocalPV")
			}
			if zfs != nil && zfs.Items != nil && len(zfs.Items) > 0 {
				k.Ns = zfs.Items[0].Namespace
			}
		} else if val, ok := nsMap[casType]; ok {
			k.Ns = val
		}
	}
	// 3. Run a specific cas-type function
	if casType != "" {
		if work, ok := CasDescribeMap()[casType]; ok {
			for _, storage := range storages {
				_ = work(k, storage)
			}
			return nil
		}
		return fmt.Errorf("cas-type %s unknown", casType)
	}

	// 4. Brute-force run describe the storage by all cas-type functions
	for _, storageName := range storages {
		for _, work := range CasDescribeList() {
			_ = work(k, storageName)
			// TODO: Should the errors be logged
			// Should we ask the user to specify a cas-type for a useful error
		}
	}
	return nil
}

// CasListMap returns a map cas-types to functions for Storage listing
func CasListMap() map[string]func(*client.K8sClient, []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error) {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error){
		util.CstorCasType: GetCstorPools,
		util.LVMCasType:   GetVolumeGroups,
		util.ZFSCasType:   GetZFSPools,
	}
}

// CasDescribeMap returns a map cas-types to functions for Storage describing
func CasDescribeMap() map[string]func(*client.K8sClient, string) error {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, string) error{
		util.CstorCasType: DescribeCstorPool,
		util.ZFSCasType:   DescribeZFSNode,
		util.LVMCasType:   DescribeLVMvg,
	}
}

// CasDescribeList returns a list of functions which describe a Storage i.e. a pool/volume-group
func CasDescribeList() []func(*client.K8sClient, string) error {
	return []func(*client.K8sClient, string) error{DescribeCstorPool, DescribeZFSNode, DescribeLVMvg}
}
