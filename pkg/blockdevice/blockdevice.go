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

package blockdevice

import (
	"fmt"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	"github.com/openebs/openebsctl/pkg/client"
)

// Get manages various implementations of blockdevice listing
func Get(bds []string, openebsNS string) error {
	// TODO: Prefer passing the client from outside
	k, _ := client.NewK8sClient("")
	// 1. Get a list of all BlockDevices
	var bdList *v1alpha1.BlockDeviceList
	bdList, err := k.GetBDs(bds,"")
	if err != nil {
		return err
	}
	var nodeBDlistMap = map[string][]v1alpha1.BlockDevice{}
	for _, bd := range bdList.Items {
		if _, ok := nodeBDlistMap[bd.Spec.NodeAttributes.NodeName]; ok {
			nodeBDlistMap[bd.Spec.NodeAttributes.NodeName] = append(nodeBDlistMap[bd.Spec.NodeAttributes.NodeName], bd)
		}else{
			nodeBDlistMap[bd.Spec.NodeAttributes.NodeName] = []v1alpha1.BlockDevice{bd}
		}
	}
	fmt.Println(nodeBDlistMap)
	return nil
}
