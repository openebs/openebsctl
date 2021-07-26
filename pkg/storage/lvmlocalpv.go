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
)

const (
	firstElemPrefix = `├─`
	lastElemPrefix  = `└─`
)

// GetVolumeGroups lists all volume groups by node
func GetVolumeGroups(c *client.K8sClient, vgs []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error) {
	lvmNodes, err := c.GetLVMNodes()
	if err != nil {
		// should this error be white-washed with return fmt.Errorf("no lvm volumegroups found")
		return nil, nil, err
	}
	var rows []metav1.TableRow
	for _, lv := range lvmNodes.Items {
		rows = append(rows, metav1.TableRow{Cells: []interface{}{lv.Name, "", "", ""}})
		for i, vg := range lv.VolumeGroups {
			var prefix string
			if i < len(lv.VolumeGroups)-1 {
				prefix = firstElemPrefix
			} else {
				prefix = lastElemPrefix
			}
			rows = append(rows, metav1.TableRow{Cells: []interface{}{prefix + vg.Name,
				util.ConvertToIBytes(vg.Free.String()), util.ConvertToIBytes(vg.Size.String())}})
		}
		rows = append(rows, metav1.TableRow{Cells: []interface{}{"", "", ""}})
	}
	// 3. Actually print the table or return an error
	if len(rows) == 0 {
		// TODO: Improve this in issue #56
		return nil, nil, fmt.Errorf("no lvm volumegroups found")
	}
	return util.LVMvolgroupListColumnDefinitions, rows, nil
}
