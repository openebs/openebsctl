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
	"fmt"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const zfsdesc = `
{{.HostName}} Details :

HOSTNAME        : {{.HostName}}
NAMESPACE       : {{.Namespace}}
NUMBER OF POOLS : {{.NumberOfPools}}
TOTAL FREE      : {{.TotalFree}}
`

// GetZFSPools lists all zfspools by zfsnodes
func GetZFSPools(c *client.K8sClient, zfsnodes []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error) {
	zfsNodes, _, err := c.GetZFSNodes(zfsnodes, util.List, "", util.MapOptions{})
	if err != nil {
		return nil, nil, err
	}
	var rows []metav1.TableRow
	for _, zfsNode := range zfsNodes.Items {
		rows = append(rows, metav1.TableRow{Cells: []interface{}{zfsNode.Name, ""}})
		for i, pool := range zfsNode.Pools {
			var prefix string
			if i < len(zfsNode.Pools)-1 {
				prefix = firstElemPrefix
			} else {
				prefix = lastElemPrefix
			}
			rows = append(rows, metav1.TableRow{Cells: []interface{}{prefix + pool.Name,
				util.ConvertToIBytes(pool.Free.String())}})
		}
		rows = append(rows, metav1.TableRow{Cells: []interface{}{"", ""}})
	}
	// 3. Actually print the table or return an error
	if len(rows) == 0 {
		return nil, nil, util.HandleEmptyTableError("zfs pools", c.Ns, "")
	}
	return util.ZFSPoolListColumnDefinitions, rows, nil
}

// ZfsNodeDesc describes a zfsnode
type ZfsNodeDesc struct {
	HostName      string
	Namespace     string
	NumberOfPools int
	TotalFree     string
}

// DescribeZFSNode describes a ZFS node & the zfspools present in it
func DescribeZFSNode(c *client.K8sClient, sName string) error {
	zfsInfo, _, err := c.GetZFSNodes([]string{sName}, util.List, "", util.MapOptions{})
	if err != nil {
		return err
	}
	if len(zfsInfo.Items) == 0 {
		return fmt.Errorf("zfsnode %s not found", sName)
	}
	zfsN := zfsInfo.Items[0]
	var totalFree resource.Quantity
	for _, pools := range zfsN.Pools {
		// TODO: handle case when size is just represented in numbers of bytes
		totalFree.Add(pools.Free)
	}
	desc := ZfsNodeDesc{
		HostName:      zfsN.Name,
		Namespace:     zfsN.Namespace,
		NumberOfPools: len(zfsN.Pools),
		TotalFree:     util.ConvertToIBytes(totalFree.String()),
	}
	return util.PrintByTemplate("zfsnodes", zfsdesc, desc)
}
