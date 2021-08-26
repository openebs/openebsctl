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
