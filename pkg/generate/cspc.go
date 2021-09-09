package generate

import (
	"fmt"
	"strings"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// isPoolValid checks if a CStor pool is valid
func isPoolTypeValid(raid string) bool {
	if raid == "stripe" || raid == "mirror" || raid == "raidz" || raid == "raidz2" {
		return true
	} else {
		return false
	}
}

// Pool calls the generate routine for different cas-types
func Pool(nodes []string, devs int, raid string) error {
	c, _ := client.NewK8sClient("")
	_, _, err := CSPC(c, nodes, devs, strings.ToLower(raid))
	if err != nil {
		return err
	}
	return nil
}

// CSPC takes eligible nodes, number of devices and poolType to create a pool cluster template
func CSPC(c *client.K8sClient, nodes []string, devs int, poolType string) (*cstorv1.CStorPoolCluster, string, error) {
	// 0. Figure out the OPENEBS_NAMESPACE for CStor
	cstorNS, err := c.GetOpenEBSNamespace(util.CstorCasType)
	// assume CSTOR's OPENEBS_NAMESPACE has all the relevant blockdevices
	c.Ns = cstorNS
	if err != nil {
		return nil, "", fmt.Errorf("cannot find an active cstor installation")
	}
	// 1. Validate nodes & poolType, fetch disks
	nodeList, err := c.GetNodes(nodes, "", "")
	if err != nil {
		return nil, "", fmt.Errorf("(server error) unable to fetch nodes %s", err)
	}
	if len(nodeList.Items) != len(nodes) {
		return nil, "", fmt.Errorf("not all worker nodes are available for provisioning a CSPC")
	}
	bds, err := c.GetBDs(nil, "kubernetes.io/hostname in ("+strings.Join(nodes, ",")+")")
	if err != nil || len(bds.Items) == 0 {
		return nil, "", fmt.Errorf("no blockdevices found in %s nodes", nodes)
	}
	_, err = filterCStorCompatible(bds)
	if err != nil {
		return nil, "", fmt.Errorf("(server error) unable to fetch bds from %s nodes", nodes)
	}
	// 2. Choose devices at the valid nodes
	// 3. Write the CSPC object with a dummy name
	p := make([]cstorv1.PoolSpec, 0, len(nodes))
	cspc := cstorv1.CStorPoolCluster{
		TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Name: "", Namespace: cstorNS, GenerateName: "cstor"},
		Spec: cstorv1.CStorPoolClusterSpec{
			Pools: p,
		}}
	// 4. Unmarshall it into a string
	// 5. Split the string by the newlines/carriage returns and insert the BD's link
	return &cspc, "", nil
}

// filterCStorCompatible takes a list of BDs and gives out a list of BDs which can be used to provision a pool
func filterCStorCompatible(bds *v1alpha1.BlockDeviceList) (*v1alpha1.BlockDeviceList, error) {
	var final []v1alpha1.BlockDevice
	for _, bd := range bds.Items {
		// an eligible blockdevice is in active+unclaimed state and lacks a file-system
		if bd.Status.State == v1alpha1.BlockDeviceActive &&
			bd.Status.ClaimState == v1alpha1.BlockDeviceUnclaimed &&
			bd.Spec.FileSystem.Type == "" &&
			// BD's capacity >=64 MiB
			bd.Spec.Capacity.Storage >= 67110000 {
			final = append(final, bd)
		}
	}
	bds.Items = final
	if len(final) == 0 {
		return nil, fmt.Errorf("found no eligble blockdevices")
	}
	return bds, nil
}
