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

	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// DeviceList is a LinkedList of BlockDevices
type DeviceList struct {
	item *v1alpha1.BlockDevice
	next *DeviceList
}

// New returns a new DeviceList node with a bd device
func New(bd v1alpha1.BlockDevice) *DeviceList {
	return &DeviceList{&bd, nil}
}

// Generate returns a new initialized *DeviceList(linked list) with the list of Blockdevices
func Generate(list v1alpha1.BlockDeviceList) *DeviceList {
	if len(list.Items) == 0 {
		return nil
	}
	var head *DeviceList
	curr := head
	for _, bd := range list.Items {
		if curr == nil {
			head = New(bd)
			curr = head
		} else {
			curr.next = New(bd)
			curr = curr.next
		}
	}
	return head
}

// Select returns count number of Blockdevices from the DeviceList LinkedList
func (head *DeviceList) Select(size resource.Quantity, count int) (*DeviceList, []v1alpha1.BlockDevice, error) {
	if count == 1 {
		// there's only one way of selecting one disk such that losses are
		// minimized in a single RaidGroup
		curr := head
		head = head.next
		return head, []v1alpha1.BlockDevice{*curr.item}, nil
	}
	curr := head
	fakeHead := &DeviceList{item: &v1alpha1.BlockDevice{}, next: head}
	prev := fakeHead
	results := []v1alpha1.BlockDevice{}
	// ahead is count nodes ahead of curr
	ahead := head
	for i := 1; i < count; i++ {
		if ahead == nil {
			return head, nil, fmt.Errorf("wanted %d blockdevices, have %d to pick", count, i)
		}
		ahead = ahead.next
	}
	for ahead != nil {
		capFirst := resource.MustParse(fmt.Sprintf("%d", curr.item.Spec.Capacity.Storage))
		capLast := resource.MustParse(fmt.Sprintf("%d", ahead.item.Spec.Capacity.Storage))
		if capFirst.Cmp(capLast) == 0 {
			// add all the devices in the same group
			for curr != ahead {
				results = append(results, *curr.item)
				curr = curr.next
			}
			results = append(results, *curr.item)
			// 1. Remove the set of BDs from the LinkedList
			prev.next = ahead.next
			if len(results) == count {
				break
			}
		}
		prev = curr
		curr = curr.next
		ahead = ahead.next
	}
	head = fakeHead.next
	if len(results) != count {
		return head, nil, fmt.Errorf("wanted %d blockdevices, have %d to pick", count, len(results))
	}
	return head, results, nil
}
