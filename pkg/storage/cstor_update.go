package storage

import (
	"context"
	"encoding/json"
	"fmt"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/openebsctl/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// CSPCnodeChange helps patch the CSPC for older nodes
func CSPCnodeChange(k *client.K8sClient, poolName, oldNode, newNode string) error {
	if k == nil {
		k, _ = client.NewK8sClient("openebs")
	}
	cspc, err := k.GetCSPC(poolName)
	if err != nil {
		return fmt.Errorf("CStor pool cluster %s not found", poolName)
	}
	var newPool cstorv1.CStorPoolCluster
	_, err = k.GetNode(newNode)
	if err != nil {
		return fmt.Errorf("node %s not found", newNode)
	}
	// TODO: Find a good way to figure out if the newer node is more suitable for the disk-replacement
	//return fmt.Errorf("node %s not in a good state", newNode)
	cspc.DeepCopyInto(&newPool)
	for _, pi := range cspc.Spec.Pools {
		if pi.NodeSelector["kubernetes.io/hostname"] == oldNode {
			pi.NodeSelector["kubernetes.io/hostname"] = newNode
		}
	}

	// cspis, _ := k.GetCSPIs(nil, "openebs.io/cas-type=cstor,openebs.io/cstor-pool-cluster="+cspc.Name)
	oldCSPC, _ := json.Marshal(cspc)
	newCSPC, _ := json.Marshal(newPool)
	data, err := strategicpatch.CreateTwoWayMergePatch(oldCSPC, newCSPC, cspc)
	_, err = k.OpenebsCS.CstorV1().CStorPoolClusters(k.Ns).Patch(context.TODO(), poolName,
		types.StrategicMergePatchType, data, metav1.PatchOptions{}, []string{}...)
	return err
}
