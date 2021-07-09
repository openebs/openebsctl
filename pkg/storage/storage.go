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
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"
)

// Get manages various implementations of Storage listing
func Get(pools []string, openebsNS string, casType string) error {
	// TODO: Change the below implementation once more castype pool listing is in place
	if casType != util.CstorCasType && casType != "" {
		return errors.Errorf("storage listing feature not available for %s", casType)
	}
	// 1. Create the clientset
	k, _ := client.NewK8sClient("")
	// 2. Get the namespaces
	nsMap, _ := k.GetOpenEBSNamespaceMap()
	if openebsNS == "" {
		if val, ok := nsMap[casType]; ok {
			k.Ns = val
		}
	}
	// TODO: Change this line, currently this overwriting empty flag values as cstor
	if list, ok := CasListMap()[util.CstorCasType]; ok {
		err := list(k, pools)
		if err != nil {
			return err
		}
	}
	return nil
}

// Describe manages various implementations of Storage Describing
func Describe(storages []string, openebsNs string) error {
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
		// TODO: Change this line, currently this overwriting empty flag values as cstor
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
	}
}

// CasDescribeMap returns a map cas-types to functions for Storage describing
func CasDescribeMap() map[string]func(*client.K8sClient, string) error {
	// a good hack to implement immutable maps in Golang & also write tests for it
	return map[string]func(*client.K8sClient, string) error{
		util.CstorCasType: DescribeCstorPool,
	}
}
