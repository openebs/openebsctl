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
	"time"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	lvm "github.com/openebs/lvm-localpv/pkg/apis/openebs.io/lvm/v1alpha1"
	zfs "github.com/openebs/zfs-localpv/pkg/apis/openebs.io/zfs/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cspi1 = cstorv1.CStorPoolInstance{
	TypeMeta: metav1.TypeMeta{Kind: "CStorPoolInstance", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "pool-1", Namespace: "openebs",
		Finalizers: []string{"cstorpoolcluster.openebs.io/finalizer", "openebs.io/pool-protection"},
		Labels: map[string]string{
			"kubernetes.io/hostname":        "node1",
			"openebs.io/cas-type":           "cstor",
			"openebs.io/cstor-pool-cluster": "cassandra-pool",
			"openebs.io/version":            "2.11"},
		// OwnerReference links to the CSPC
	},
	Spec: cstorv1.CStorPoolInstanceSpec{
		HostName:     "node1",
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
		PoolConfig:   cstorv1.PoolConfig{DataRaidGroupType: "stripe", WriteCacheGroupType: "", Compression: "off"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd-1", Capacity: 1234567, DevLink: "/dev/disk/by-id/abcd/def"}}}},
		WriteCacheRaidGroups: nil,
	},
	Status: cstorv1.CStorPoolInstanceStatus{
		Conditions: []cstorv1.CStorPoolInstanceCondition{{
			Type:               cstorv1.CSPIPoolLost,
			Status:             "True",
			LastUpdateTime:     metav1.Time{Time: time.Now()},
			LastTransitionTime: metav1.Time{Time: time.Now()},
			Reason:             "PoolLost",
			Message:            "failed to importcstor-xyzabcd",
		}},
		Phase: cstorv1.CStorPoolStatusOnline,
		Capacity: cstorv1.CStorPoolInstanceCapacity{
			Used:  resource.MustParse("18600Mi"),
			Free:  resource.MustParse("174Gi"),
			Total: resource.MustParse("192600Mi"),
			ZFS:   cstorv1.ZFSCapacityAttributes{},
		},
		ReadOnly: false, ProvisionedReplicas: 2, HealthyReplicas: 2,
	},
	VersionDetails: cstorv1.VersionDetails{Desired: "2.11",
		Status: cstorv1.VersionStatus{Current: "2.11", State: cstorv1.ReconcileComplete, LastUpdateTime: metav1.Time{Time: time.Now()}},
	},
}

var bd1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "BlockDevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd-1", Namespace: "openebs",
		Annotations: map[string]string{
			"internal.openebs.io/partition-uuid": "49473bca-97c3-f340-beaf-dae9b2ce99bc",
			"internal.openebs.io/uuid-scheme":    "legacy"}},
	Spec: v1alpha1.DeviceSpec{Capacity: v1alpha1.DeviceCapacity{
		Storage:            123456789,
		PhysicalSectorSize: 123456789,
		LogicalSectorSize:  123456789,
	}},
	Status: v1alpha1.DeviceStatus{},
}

var bd2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{
		Kind:       "BlockDevice",
		APIVersion: "openebs.io/v1alpha1",
	},
	ObjectMeta: metav1.ObjectMeta{Name: "bd-2", Namespace: "openebs"},
	Spec: v1alpha1.DeviceSpec{Capacity: v1alpha1.DeviceCapacity{
		Storage:            123456789,
		PhysicalSectorSize: 123456789,
		LogicalSectorSize:  123456789,
	},
		FileSystem: v1alpha1.FileSystemInfo{Type: "zfs_member", Mountpoint: "/home/kubernetes/volume-abcd"}},
	Status: v1alpha1.DeviceStatus{
		ClaimState: "Claimed",
		State:      "Active",
	},
}

var cvr1 = v1.CStorVolumeReplica{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1-rep-1",
		Labels:    map[string]string{cstortypes.CStorPoolInstanceNameLabelKey: "pool-1", "openebs.io/persistent-volume": "pv1"},
		Namespace: "openebs",
	},
	Status: v1.CStorVolumeReplicaStatus{
		Capacity: v1.CStorVolumeReplicaCapacityDetails{
			Total: "4Gi",
			Used:  "70Mi",
		},
		Phase: v1.CVRStatusOnline,
	},
}

var cvr2 = v1.CStorVolumeReplica{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1-rep-2",
		Labels:    map[string]string{cstortypes.CStorPoolInstanceNameLabelKey: "pool-1", "openebs.io/persistent-volume": "pv1"},
		Namespace: "openebs",
	},
	Status: v1.CStorVolumeReplicaStatus{
		Capacity: v1.CStorVolumeReplicaCapacityDetails{
			Total: "40Gi",
			Used:  "70Mi",
		},
		Phase: v1.CVRStatusOnline,
	},
}
var pv1 = corev1.PersistentVolume{
	TypeMeta:   metav1.TypeMeta{Kind: "PersistentVolume", APIVersion: "core/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "pv1"},
	Spec:       corev1.PersistentVolumeSpec{ClaimRef: &corev1.ObjectReference{Name: "mongopv1"}},
	Status:     corev1.PersistentVolumeStatus{},
}

var cspi2 = cstorv1.CStorPoolInstance{
	TypeMeta: metav1.TypeMeta{Kind: "CStorPoolInstance", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "pool-2", Namespace: "openebs",
		Finalizers: []string{"cstorpoolcluster.openebs.io/finalizer", "openebs.io/pool-protection"},
		Labels: map[string]string{
			"kubernetes.io/hostname":        "node2",
			"openebs.io/cas-type":           "cstor",
			"openebs.io/cstor-pool-cluster": "cassandra-pool",
			"openebs.io/version":            "2.11"},
		// OwnerReference links to the CSPC
	},
	Spec: cstorv1.CStorPoolInstanceSpec{
		HostName:     "node2",
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
		PoolConfig:   cstorv1.PoolConfig{DataRaidGroupType: "stripe", WriteCacheGroupType: "", Compression: "off"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd2", Capacity: 1234567, DevLink: "/dev/disk/by-id/abcd/def"}}}},
		WriteCacheRaidGroups: nil,
	},
	Status: cstorv1.CStorPoolInstanceStatus{
		Conditions: []cstorv1.CStorPoolInstanceCondition{{
			Type:               cstorv1.CSPIPoolLost,
			Status:             "True",
			LastUpdateTime:     metav1.Time{Time: time.Now()},
			LastTransitionTime: metav1.Time{Time: time.Now()},
			Reason:             "PoolLost",
			Message:            "failed to importcstor-xyzabcd",
		}},
		Phase: cstorv1.CStorPoolStatusOnline,
		Capacity: cstorv1.CStorPoolInstanceCapacity{
			Used:  resource.MustParse("18600Mi"),
			Free:  resource.MustParse("174Gi"),
			Total: resource.MustParse("192600Mi"),
			ZFS:   cstorv1.ZFSCapacityAttributes{},
		},
		ReadOnly: false, ProvisionedReplicas: 2, HealthyReplicas: 2,
	},
	VersionDetails: cstorv1.VersionDetails{Desired: "2.11",
		Status: cstorv1.VersionStatus{Current: "2.11", State: cstorv1.ReconcileComplete, LastUpdateTime: metav1.Time{Time: time.Now()}},
	},
}

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
