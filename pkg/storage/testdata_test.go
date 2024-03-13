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
	lvm "github.com/openebs/lvm-localpv/pkg/apis/openebs.io/lvm/v1alpha1"
	zfs "github.com/openebs/zfs-localpv/pkg/apis/openebs.io/zfs/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	fourGigiByte = resource.MustParse("4Gi")
	fiveGigiByte = resource.MustParse("5Gi")
)
var lvmNode1 = lvm.LVMNode{
	TypeMeta: metav1.TypeMeta{
		Kind:       "LVMNode",
		APIVersion: "local.openebs.io/v1alpha1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "node1",
		Namespace: "lvm",
	},
	VolumeGroups: []lvm.VolumeGroup{{Name: "lvmvg", UUID: "ed6fko-Lf33-AW2d-Vblk-cUJZ-sQt5-Gr4rcH", Size: fiveGigiByte,
		Free: fourGigiByte, LVCount: 1, PVCount: 1},
		{Name: "lvmvg2", UUID: "ed6fko-Lf33-AW2d-Vblk-cUJZ-sQt5-Hr4rcI", Size: fiveGigiByte,
			Free: fourGigiByte, LVCount: 1, PVCount: 1}},
}
var lvmNode2 = lvm.LVMNode{
	TypeMeta: metav1.TypeMeta{
		Kind:       "LVMNode",
		APIVersion: "local.openebs.io/v1alpha1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "node2",
		Namespace: "lvm",
	},
	VolumeGroups: []lvm.VolumeGroup{{Name: "lvmvg", UUID: "ed6fko-Lf33-AW2d-Vblk-cUJZ-sQt5-Gr4rcH", Size: fiveGigiByte,
		Free: fourGigiByte, LVCount: 2, PVCount: 2},
		{Name: "lvmvg2", UUID: "ed6fko-Lf33-AW2d-Vblk-cUJZ-sQt5-Hr4rcI", Size: fiveGigiByte,
			Free: fourGigiByte, LVCount: 1, PVCount: 1}},
}

var zfsNode1 = zfs.ZFSNode{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ZFSNode",
		APIVersion: "zfs.openebs.io/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "node1",
		Namespace: "zfs",
		// OwnerReference: refers to the K8s-node where the zfs volume is created
	},
	Pools: []zfs.Pool{
		{
			Name: "zfs-pool1",
			UUID: "15423895941648453427",
			Free: resource.MustParse("33285828Ki"),
		},
	},
}

var zfsNode2 = zfs.ZFSNode{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ZFSNode",
		APIVersion: "zfs.openebs.io/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "node2",
		Namespace: "zfs",
		// OwnerReference: refers to the K8s-node where the zfs volume is created
	},
	Pools: []zfs.Pool{{Name: "zfs-pool2", UUID: "15423895941648453428", Free: resource.MustParse("33285828Ki")},
		{Name: "zfs-pool3", UUID: "15423895941648453426", Free: resource.MustParse("33285828Ki")}},
}

var zfsNode3 = zfs.ZFSNode{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ZFSNode",
		APIVersion: "zfs.openebs.io/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "node3",
		Namespace: "zfs",
		// OwnerReference: refers to the K8s-node where the zfs volume is created
	},
	Pools: []zfs.Pool{{Name: "zfs-pool1", UUID: "15423895941648453428", Free: resource.MustParse("33285828")},
		{Name: "zfs-poolX", UUID: "15423895941648453426", Free: resource.MustParse("33285828Ki")}},
}
