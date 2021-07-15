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

// GetZFSNodes lists all zfspools by zfsnodes
func GetZFSNodes(c *client.K8sClient, zfsnodes []string) error {
	zfsNodes, err := c.GetZFSNodes()
	if err != nil {
		return err
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
		return fmt.Errorf("no zfspools found")
	}
	util.TablePrinter(util.ZFSPoolListColumnDefinitions, rows, printers.PrintOptions{Wide: true})
	return nil
}
