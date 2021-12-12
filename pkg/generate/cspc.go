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

package generate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ghodss/yaml"
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
	}
	return false
}

// CSPC calls the generate routine for different cas-types
func CSPC(nodes []string, devs int, raid string) error {
	c := client.NewK8sClient()
	if !isPoolTypeValid(strings.ToLower(raid)) {
		// TODO: Use the well defined pool constant types from openebs/api when added there
		return fmt.Errorf("invalid pool type %s", raid)
	}
	_, str, err := cspc(c, nodes, devs, strings.ToLower(raid))
	if err != nil {
		return err
	}
	fmt.Println(str)
	return nil
}

// cspc takes eligible nodes, number of devices and poolType to create a pool cluster template
func cspc(c *client.K8sClient, nodes []string, devs int, poolType string) (*cstorv1.CStorPoolCluster, string, error) {
	// 0. Figure out the OPENEBS_NAMESPACE for CStor
	cstorNS, err := c.GetOpenEBSNamespace(util.CstorCasType)
	// assume CSTOR's OPENEBS_NAMESPACE has all the relevant blockdevices
	c.Ns = cstorNS
	if err != nil {
		return nil, "", fmt.Errorf("unable to determine the cStor namespace error: %v", err)
	}
	// 1. Validate nodes & poolType, fetch disks
	nodeList, err := c.GetNodes(nodes, "", "")
	if err != nil {
		return nil, "", fmt.Errorf("(server error) unable to fetch node information %s", err)
	}
	if len(nodeList.Items) != len(nodes) {
		return nil, "", fmt.Errorf("not all worker nodes are available for provisioning a cspc")
	}
	// 1.1 Translate nodeNames to node's hostNames to fetch disks
	// while they might seem equivalent, they aren't equal, this quirk is
	// visible clearly for EKS clusters
	var hostnames []string
	for _, node := range nodeList.Items {
		// I hope it is unlikely for a K8s node to have an empty hostname
		hostnames = append(hostnames, node.Labels["kubernetes.io/hostname"])
	}
	// 2. Fetch BD's from the eligible/valid nodes by hostname labels
	bds, err := c.GetBDs(nil, "kubernetes.io/hostname in ("+strings.Join(hostnames, ",")+")")
	if err != nil || len(bds.Items) == 0 {
		return nil, "", fmt.Errorf("no blockdevices found in nodes with %v hostnames", hostnames)
	}
	_, err = filterCStorCompatible(bds)
	if err != nil {
		return nil, "", fmt.Errorf("(server error) unable to fetch bds from %v nodes", nodes)
	}
	// 3. Choose devices at the valid BDs by hostname
	hostToBD := make(map[string][]v1alpha1.BlockDevice)
	for _, bd := range bds.Items {
		hostToBD[bd.Labels["kubernetes.io/hostname"]] = append(hostToBD[bd.Labels["kubernetes.io/hostname"]], bd)
	}
	// 4. Select disks and create the PoolSpec
	p, err := makePools(poolType, devs, hostToBD, nodes, hostnames)
	if err != nil {
		return nil, "", err
	}

	// 5. Write the cspc object with a dummy name
	cspc := cstorv1.CStorPoolCluster{
		TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
		ObjectMeta: metav1.ObjectMeta{Namespace: cstorNS, GenerateName: "cstor"},
		Spec: cstorv1.CStorPoolClusterSpec{
			Pools: *p,
		},
	}
	// 6. Unmarshall it into a string
	y, err := yaml.Marshal(cspc)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return nil, "", err
	}
	specStr := string(y)
	// 7. removing status and versionDetails field
	specStr = specStr[:strings.Index(specStr, "status: {}")]
	// 8. Split the string by the newlines/carriage returns and insert the BD's link
	specStr = addBDDetailComments(specStr, bds)
	return &cspc, specStr, nil
}

// addBDDetailComments adds more information about the blockdevice in a cspc YAML string
func addBDDetailComments(yaml string, bdList *v1alpha1.BlockDeviceList) string {
	finalYaml := ""
	for _, l := range strings.Split(yaml, "\n") {
		if strings.Contains(l, "- blockDeviceName:") {
			name := strings.Trim(strings.Split(l, ":")[1], " ")
			finalYaml = finalYaml + getBDComment(name, bdList) + "\n"
		}
		finalYaml = finalYaml + l + "\n"
	}
	return finalYaml
}

// getBDComment returns information about a blockdevice, with fixed whitespace
// to match the identation level
func getBDComment(name string, bdList *v1alpha1.BlockDeviceList) string {
	for _, bd := range bdList.Items {
		if bd.Name == name {
			return "      # " + bd.Spec.Path + "  " + util.ConvertToIBytes(strconv.FormatUint(bd.Spec.Capacity.Storage, 10))
		}
	}
	return ""
}

// makePools creates a poolSpec based on the poolType, number of devices per
// pool instance and a collection of blockdevices by nodes
func makePools(poolType string, nDevices int, bd map[string][]v1alpha1.BlockDevice, nodes []string, hosts []string) (*[]cstorv1.PoolSpec, error) {
	// IMPORTANT: User is more likely to see the nodeNames, so the errors
	// should preferably be shown in terms of nodeNames and not hostNames
	var spec []cstorv1.PoolSpec
	if poolType == string(cstorv1.PoolStriped) {
		// always single RAID-group with nDevices patched together, cannot disk replace,
		// no redundancy in a pool, redundancy possible across pool instances

		// for each eligible set of BDs from each eligible nodes with hostname
		// "host", take nDevices number of BDs
		for i, host := range hosts {
			bds, ok := bd[host]
			if !ok {
				// DOUBT: Do 0 or lesser number of BDs demand a separate error string?
				// I can ask to create a stripe pool with 1 disk and my
				// choice of node might not have eligible BDs
				return nil, fmt.Errorf("no eligible blockdevices found in node %s", nodes[i])
			}
			if len(bds) < nDevices {
				// the node might have lesser number of BDs
				return nil, fmt.Errorf("not enough blockdevices found on node %s, want %d, found %d", nodes[i], nDevices, len(bds))
			}
			var raids []cstorv1.CStorPoolInstanceBlockDevice

			for d := 0; d < nDevices; d++ {
				raids = append(raids, cstorv1.CStorPoolInstanceBlockDevice{BlockDeviceName: bds[d].Name})
			}
			spec = append(spec, cstorv1.PoolSpec{
				NodeSelector:   map[string]string{"kubernetes.io/hostname": host},
				DataRaidGroups: []cstorv1.RaidGroup{{CStorPoolInstanceBlockDevices: raids}},
				PoolConfig: cstorv1.PoolConfig{
					DataRaidGroupType: string(cstorv1.PoolStriped),
				},
			})
		}
		return &spec, nil
	} else if poolType == string(cstorv1.PoolMirrored) {
		if nDevices%2 != 0 {
			return nil, fmt.Errorf("mirrored pool requires multiples of two block device")
		}
		for i, host := range hosts {
			var raidGroups []cstorv1.RaidGroup
			// add all BDs to a CSPCs CSPI spec
			bds := bd[host]
			if len(bds) < nDevices {
				return nil, fmt.Errorf("not enough eligible blockdevices found on node %s, want %d, found %d", nodes[i], nDevices, len(bds))
			}
			index := 0
			for d := 0; d < nDevices/2; d++ {
				raids := []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: bds[index].Name},
					{BlockDeviceName: bds[index+1].Name}}
				raidGroups = append(raidGroups, cstorv1.RaidGroup{CStorPoolInstanceBlockDevices: raids})
				index += 2
			}
			// add the CSPI BD spec inside cspc to a PoolSpec
			spec = append(spec, cstorv1.PoolSpec{
				NodeSelector:   map[string]string{"kubernetes.io/hostname": host},
				DataRaidGroups: raidGroups,
				PoolConfig: cstorv1.PoolConfig{
					DataRaidGroupType: string(cstorv1.PoolMirrored),
				},
			})
		}
		return &spec, nil
		// 2ⁿ devices per RaidGroup, (confirm) not more than 2 devices per RaidGroup
		// DOUBT: Should this throw an error if nDevices isn't 2ⁿ?
	} else if poolType == string(cstorv1.PoolRaidz) {
		return nil, fmt.Errorf("%s is not supported yet", poolType)
		// 2ⁿ+1 devices per RaidGroup
	} else if poolType == string(cstorv1.PoolRaidz2) {
		return nil, fmt.Errorf("%s is not supported yet", poolType)
		// 2ⁿ+2 devices per RaidGroup
	}
	return nil, fmt.Errorf("unknown pool-type")
}

// filterCStorCompatible takes a list of BDs and gives out a list of BDs which can be used to provision a pool
func filterCStorCompatible(bds *v1alpha1.BlockDeviceList) (*v1alpha1.BlockDeviceList, error) {
	// TODO: Optionally reject sparse-disks depending on configs
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
