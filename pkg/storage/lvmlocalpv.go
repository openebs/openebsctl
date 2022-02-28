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
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	firstElemPrefix = `├─`
	lastElemPrefix  = `└─`
)

// GetVolumeGroups lists all volume groups by node
func GetVolumeGroups(c *client.K8sClient, vgs []string) ([]metav1.TableColumnDefinition, []metav1.TableRow, error) {
	lvmNodes, _, err := c.GetLVMNodes(vgs, util.List, "", util.MapOptions{})
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
		return nil, nil, util.HandleEmptyTableError("lvm Volumegroups", c.Ns, "")
	}
	return util.LVMvolgroupListColumnDefinitions, rows, nil
}

const lvmdesc = `
{{.HostName}} Details :

HOSTNAME        : {{.HostName}}
NAMESPACE       : {{.Namespace}}
NUMBER OF POOLS : {{.NumberOfPools}}
TOTAL CAPACITY  : {{.Size}}
TOTAL FREE      : {{.TotalFree}}
TOTAL LVs       : {{.TotalLogicalVolumes}}
TOTAL PVs       : {{.TotalPVs}}

`

// LVMvgDesc describes an LVM Volume group & the vgs present in it
type LVMvgDesc struct {
	HostName            string
	Namespace           string
	NumberOfPools       int
	Size                string
	TotalFree           string
	TotalLogicalVolumes int32
	TotalPVs            int32
}

// DescribeLVMvg describes an LVM volume group
func DescribeLVMvg(c *client.K8sClient, vg string) error {
	vgs, _, err := c.GetLVMNodes([]string{vg}, util.List, "", util.MapOptions{})
	if err != nil {
		return err
	}
	if len(vgs.Items) == 0 {
		return fmt.Errorf("vg group %s not found", vg)
	}
	volGrp := vgs.Items[0]
	var totalFree, total resource.Quantity
	var totLV, totPV int32
	for _, pools := range volGrp.VolumeGroups {
		totalFree.Add(pools.Free)
		total.Add(pools.Size)
		totLV += pools.LVCount
		totPV += pools.PVCount
	}

	desc := LVMvgDesc{
		HostName:            volGrp.Name,
		Namespace:           volGrp.Namespace,
		NumberOfPools:       len(volGrp.VolumeGroups),
		Size:                util.ConvertToIBytes(total.String()),
		TotalFree:           util.ConvertToIBytes(totalFree.String()),
		TotalLogicalVolumes: totLV,
		TotalPVs:            totPV,
	}

	var r []metav1.TableRow
	for _, k := range volGrp.VolumeGroups {
		usedPercent := util.GetUsedPercentage(k.Size.String(), k.Free.String())
		r = append(r, metav1.TableRow{Cells: []interface{}{k.Name, k.UUID, k.LVCount, k.PVCount, fmt.Sprintf("%0.1f%%", 100-usedPercent)}})
	}
	_ = util.PrintByTemplate("lvmvgs", lvmdesc, desc)
	fmt.Println("Volume group details")
	fmt.Println("---------------------")
	def := []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "UUID", Type: "string"},
		{Name: "LV count", Type: "string"},
		{Name: "PV count", Type: "string"},
		{Name: "Used percentage", Type: "string"},
	}
	util.TablePrinter(def, r, printers.PrintOptions{Wide: true})
	return nil
}
