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
)

// Get manages various implementations of Storage listing
func Get(pools []string, openebsNS string, casType string) error {
	// 1. Create the clientset
	k, _ := client.NewK8sClient("")
	// 2. If casType is specified, call the specific function & exit
	if f, ok := CasListMap()[casType]; ok {
		// if a cas-type is found, run it and return the error
		return f(k, pools)
	}
	// 3. Call all functions & exit
	var separator bool
	for _, f := range CasList() {
		if separator {
			fmt.Println()
		}
		err := f(k, pools)
		if err == nil {
			// A visual separator for different cas-type pools/storage entities
			separator = true
		} else {
			separator = false
		}
	}
	return nil
}

// CasList has a list of method implementations for different cas-types
func CasList() []func(*client.K8sClient, []string) error {
	return []func(*client.K8sClient, []string) error{
		GetCstorPools, GetVolumeGroups, GetZFSNodes}
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
func CasListMap() map[string]func(*client.K8sClient, []string) error {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, []string) error{
		util.CstorCasType: GetCstorPools,
		util.LVMLocalPV:   GetVolumeGroups,
		util.ZFSCasType:   GetZFSNodes,
	}
}

// CasDescribeMap returns a map cas-types to functions for Storage describing
func CasDescribeMap() map[string]func(*client.K8sClient, string) error {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, string) error{
		util.CstorCasType: DescribeCstorPool,
	}
}
