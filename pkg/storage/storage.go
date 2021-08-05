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

package storage

import (
	"fmt"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

// Get manages various implementations of Storage listing
func Get(pools []string, openebsNS string, casType string) error {
	// 1. Create the clientset
	k, _ := client.NewK8sClient("")
	// 2. If casType is specified, call the specific function & exit
	if f, ok := CasListMap()[casType]; ok {
		// if a cas-type is found, run it and return the error
		header, rows, err := f(k, pools)
		if err != nil {
			return err
		} else {
			util.TablePrinter(header, rows, printers.PrintOptions{Wide: true})
		}
	} else if casType != "" {
		return fmt.Errorf("cas-type %s is not supported", casType)
	} else {
		// 3. Call all functions & exit
		for _, f := range CasList() {
			header, row, err := f(k, pools)
			if err == nil {
				// 4. Find the correct heading & print the rows
				util.TablePrinter(header, row, printers.PrintOptions{Wide: true})
				// A visual separator for different cas-type pools/storage entities
				fmt.Println()
			}
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
	// 1. Create the clientset
	k, _ := client.NewK8sClient(openebsNs)
	// 2. Get the namespaces
	nsMap, _ := k.GetOpenEBSNamespaceMap()
	if openebsNs == "" {
		// TODO: Change this line, currently this overwriting empty flag values as cstor
		if val, ok := nsMap[util.CstorCasType]; ok {
			k.Ns = val
		}
	}
	for _, storageName := range storages {
		// 3. Describe the storage
		if list, ok := CasDescribeMap()[util.CstorCasType]; ok {
			err := list(k, storageName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CasListMap returns a map cas-types to functions for Storage listing
func CasListMap() map[string]func(*client.K8sClient, []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error) {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error){
		util.CstorCasType: GetCstorPools,
		util.LVMLocalPV:   GetVolumeGroups,
		util.ZFSCasType:   GetZFSPools,
	}
}

// CasDescribeMap returns a map cas-types to functions for Storage describing
func CasDescribeMap() map[string]func(*client.K8sClient, string) error {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, string) error{
		util.CstorCasType: DescribeCstorPool,
	}
}
