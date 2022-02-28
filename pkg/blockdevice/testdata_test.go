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

package blockdevice

import (
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	bd1 = v1alpha1.BlockDevice{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: "some-fake-bd-1"},
		Spec: v1alpha1.DeviceSpec{
			Path:     "/dev/sdb",
			Capacity: v1alpha1.DeviceCapacity{Storage: uint64(132131321)},
			FileSystem: v1alpha1.FileSystemInfo{
				Type:       "zfs_member",
				Mountpoint: "/var/some-fake-point",
			},
			NodeAttributes: v1alpha1.NodeAttribute{
				NodeName: "fake-node-1",
			},
		},
		Status: v1alpha1.DeviceStatus{
			ClaimState: "Claimed",
			State:      "Active",
		},
	}
	bd2 = v1alpha1.BlockDevice{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: "some-fake-bd-2"},
		Spec: v1alpha1.DeviceSpec{
			Path:     "/dev/sdb",
			Capacity: v1alpha1.DeviceCapacity{Storage: uint64(132131321)},
			FileSystem: v1alpha1.FileSystemInfo{
				Type:       "zfs_member",
				Mountpoint: "/var/some-fake-point",
			},
			NodeAttributes: v1alpha1.NodeAttribute{
				NodeName: "fake-node-1",
			},
		},
		Status: v1alpha1.DeviceStatus{
			ClaimState: "Claimed",
			State:      "Active",
		},
	}
	bd3 = v1alpha1.BlockDevice{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{Name: "some-fake-bd-3", Namespace: "fake-ns"},
		Spec: v1alpha1.DeviceSpec{
			Path:     "/dev/sdb",
			Capacity: v1alpha1.DeviceCapacity{Storage: uint64(132131321)},
			FileSystem: v1alpha1.FileSystemInfo{
				Type:       "lvm_member",
				Mountpoint: "/var/some-fake-point",
			},
			NodeAttributes: v1alpha1.NodeAttribute{
				NodeName: "fake-node-2",
			},
		},
		Status: v1alpha1.DeviceStatus{
			ClaimState: "Claimed",
			State:      "Active",
		},
	}
)
