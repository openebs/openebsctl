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
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// CSPCnodeChange helps patch the CSPC for older nodes
func CSPCnodeChange(k *client.K8sClient, poolName, oldNode, newNode string) error {
	cspc, err := k.GetCSPC(poolName)
	if err != nil {
		return fmt.Errorf("CStor pool cluster %s not found", poolName)
	}
	node, err := k.GetNode(newNode)
	if err != nil {
		return fmt.Errorf("node %s not found", newNode)
	} else if !util.IsNodeReady(node) {
		return fmt.Errorf("node %s is not ready", newNode)
	}
	// TODO: Find a good way to figure out if the newer node is more suitable
	// for the disk-replacement, i.e. doesn't have PID pressure, scheduling is
	// not possible, etc
	//return fmt.Errorf("node %s not in a good state", newNode)
	newPool := cspc.DeepCopy()
	for _, pi := range newPool.Spec.Pools {
		if pi.NodeSelector["kubernetes.io/hostname"] == oldNode {
			pi.NodeSelector["kubernetes.io/hostname"] = newNode
		}
	}
	// Patch the CSPC
	oldCSPC, _ := json.Marshal(cspc)
	newCSPC, _ := json.Marshal(newPool)
	data, err := strategicpatch.CreateTwoWayMergePatch(oldCSPC, newCSPC, cspc)
	if err != nil {
		return err
	}
	_, err = k.OpenebsCS.CstorV1().CStorPoolClusters(k.Ns).Patch(context.TODO(), poolName,
		types.MergePatchType, data, metav1.PatchOptions{}, []string{}...)
	return err
}

// DebugCSPCNode returns the appropriate storage node for the cspc's nodes in a
// map[InitialNode->FinalNode] manner and/or an error
func DebugCSPCNode(k *client.K8sClient, cspc string) (map[string]string, error) {
	// 1. Get the CSPC
	pool, err := k.GetCSPC(cspc)
	if err != nil {
		return nil, fmt.Errorf("unable to get cspc %s, %v", cspc, err)
	}
	// 2. Get the hostnames with the BDs
	expectedBDToHost := make(map[string]string)
	var devices []string
	for _, specs := range pool.Spec.Pools {
		host := specs.NodeSelector["kubernetes.io/hostname"]
		for _, rgs := range specs.DataRaidGroups {
			for _, bds := range rgs.CStorPoolInstanceBlockDevices {
				expectedBDToHost[bds.BlockDeviceName] = host
				devices = append(devices, bds.BlockDeviceName)
			}
		}
	}
	// 3. Fetch all the BDs
	bds, _ := k.GetBDs(devices, "")
	actualBDToHost := make(map[string]string)
	for _, bd := range bds.Items {
		actualBDToHost[bd.Name] = bd.Labels["kubernetes.io/hostname"]
	}

	if len(expectedBDToHost) != len(actualBDToHost) {
		diff := len(expectedBDToHost) - len(actualBDToHost)
		return nil, fmt.Errorf("%d blockdevices are missing from the cluster", diff)
	}
	// 4. Get the BD map difference
	if reflect.DeepEqual(expectedBDToHost, actualBDToHost) {
		return nil, fmt.Errorf("no change in the storage node")
	}
	// AWESOME: This piece of code can also handle circular node changes,
	// the patching function would need some more work to handle this case,
	// which is why this has been restricted to single node swaps

	// 5. Calculate the difference, map expectedHost to the current host to
	// help the existing function for swapping the node-selector values for
	// the CSPCnodeChange function

	diff := make(map[string]string)
	for bd, presentHostForBD := range actualBDToHost {
		if expectedHost, ok := expectedBDToHost[bd]; ok {
			if expectedHost != presentHostForBD {
				// difference spotted
				if _, found := diff[expectedHost]; !found {
					// add the difference to the map
					diff[expectedHost] = presentHostForBD
				} else if diff[expectedHost] != presentHostForBD {
					// beyond the scope of this project
					return nil, fmt.Errorf("multiple bds have switched to different nodes")
				}
			}
			// do nothing if the host is same in expectation & present reality
		}
	}
	// 5.1. Currently, the caller isn't supposed to be dealing with multiple
	// node changes

	// CHALLENGE: This logic will not work for if a multi-BD
	// RAIDGroup's devices have moved to different nodes, in this case,
	// the BlockDevices need to be moved from the current node to another
	// node via an appropriate external tool, later this tool might offer
	// suggestions for the same
	if len(diff) > 1 {
		return nil, fmt.Errorf("more than one node change in the storage instance")
	}
	// 6. Return the diff
	// NOTE: While this is a legit error that the nodes in CSPC and in reality
	// aren't in sync, it is a bad practice to return a usable value with a non-nil error
	return diff, nil
}
