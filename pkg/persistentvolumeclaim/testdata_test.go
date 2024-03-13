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

package persistentvolumeclaim

import (
	"time"

	lvm "github.com/openebs/lvm-localpv/pkg/apis/openebs.io/lvm/v1alpha1"
	"github.com/openebs/openebsctl/pkg/util"
	zfs "github.com/openebs/zfs-localpv/pkg/apis/openebs.io/zfs/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	fourGigiByte = resource.MustParse("4Gi")
	blockFS      = corev1.PersistentVolumeBlock
)

/****************
* LVM LOCAL PV
****************/

var lvmScName = "lvmsc"

var lvmVol1 = lvm.LVMVolume{
	TypeMeta: metav1.TypeMeta{
		Kind:       "LVMVolume",
		APIVersion: "lvm.openebs.io/v1alpha1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-1",
		Namespace:         "lvmlocalpv",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{},
		Annotations:       map[string]string{},
		OwnerReferences:   nil,
		Finalizers:        nil,
	},
	Spec: lvm.VolumeInfo{
		OwnerNodeID:   "node1",
		VolGroup:      "lvmpv",
		VgPattern:     "vg1*",
		Capacity:      "4Gi",
		Shared:        "NotShared",
		ThinProvision: "No",
	},
	Status: lvm.VolStatus{
		State: "Ready",
		Error: nil,
	},
}

var lvmPV1 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{
		Kind:       "PersistentVolume",
		APIVersion: "core/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "pvc-1",
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	},
	Spec: corev1.PersistentVolumeSpec{
		// 4GiB
		Capacity:               corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
		PersistentVolumeSource: corev1.PersistentVolumeSource{CSI: &corev1.CSIPersistentVolumeSource{Driver: util.LocalPVLVMCSIDriver}},
		AccessModes:            []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		ClaimRef: &corev1.ObjectReference{Kind: "PersistentVolumeClaim", Namespace: "pvc-namespace",
			Name: "lvm-pvc-1", APIVersion: "v1"},
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
		StorageClassName:              "lvm-sc-1",
		VolumeMode:                    &blockFS,
		NodeAffinity: &corev1.VolumeNodeAffinity{
			Required: &corev1.NodeSelector{NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{MatchExpressions: []corev1.NodeSelectorRequirement{
					{Key: "kubernetes.io/hostname", Operator: corev1.NodeSelectorOpIn, Values: []string{"node2"}},
				}},
			}},
		},
	},
	Status: corev1.PersistentVolumeStatus{
		Phase:   corev1.VolumeBound,
		Message: "Storage class not found",
		Reason:  "K8s API was down",
	},
}

var lvmPVC1 = corev1.PersistentVolumeClaim{
	TypeMeta: metav1.TypeMeta{
		Kind:       "PersistentVolumeClaim",
		APIVersion: "core/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "lvm-pvc",
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	},
	Spec: corev1.PersistentVolumeClaimSpec{
		AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Resources:        corev1.ResourceRequirements{Requests: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: fourGigiByte}},
		VolumeName:       "pvc-1",
		StorageClassName: &lvmScName,
		VolumeMode:       &blockFS,
	},
	Status: corev1.PersistentVolumeClaimStatus{
		Phase: corev1.PersistentVolumeClaimPhase(corev1.VolumeBound),
	},
}

/****************
* ZFS LOCAL PV
****************/

var zfsVol1 = zfs.ZFSVolume{
	TypeMeta: metav1.TypeMeta{
		Kind:       "ZFSVolume",
		APIVersion: "zfs.openebs.io/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-1",
		Namespace:         "zfslocalpv",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"kubernetes.io/nodename": "node1"},
		Annotations:       map[string]string{},
		OwnerReferences:   nil,
		Finalizers:        nil,
	},
	Spec: zfs.VolumeInfo{
		OwnerNodeID:   "node1",
		PoolName:      "zfspv",
		Capacity:      "4Gi",
		RecordSize:    "4k",
		Compression:   "off",
		Dedup:         "off",
		ThinProvision: "No",
		VolumeType:    "DATASET",
		FsType:        "zfs",
		Shared:        "NotShared",
	},
	Status: zfs.VolStatus{State: "Ready"},
}

var zfsPV1 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{
		Kind:       "PersistentVolume",
		APIVersion: "core/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "pvc-1",
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	},
	Spec: corev1.PersistentVolumeSpec{
		// 4GiB
		Capacity:               corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
		PersistentVolumeSource: corev1.PersistentVolumeSource{CSI: &corev1.CSIPersistentVolumeSource{Driver: util.ZFSCSIDriver}},
		AccessModes:            []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		ClaimRef: &corev1.ObjectReference{Kind: "PersistentVolumeClaim", Namespace: "pvc-namespace",
			Name: "zfs-pvc-1", APIVersion: "v1"},
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
		StorageClassName:              "zfs-sc-1",
		VolumeMode:                    &blockFS,
		NodeAffinity: &corev1.VolumeNodeAffinity{
			Required: &corev1.NodeSelector{NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{MatchExpressions: []corev1.NodeSelectorRequirement{
					{Key: "kubernetes.io/hostname", Operator: corev1.NodeSelectorOpIn, Values: []string{"node2"}},
				}},
			}},
		},
	},
	Status: corev1.PersistentVolumeStatus{
		Phase:   corev1.VolumeBound,
		Message: "Storage class not found",
		Reason:  "K8s API was down",
	},
}

var zfsPVC1 = corev1.PersistentVolumeClaim{
	TypeMeta: metav1.TypeMeta{
		Kind:       "PersistentVolumeClaim",
		APIVersion: "core/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "zfs-pvc",
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	},
	Spec: corev1.PersistentVolumeClaimSpec{
		AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		Resources:        corev1.ResourceRequirements{Requests: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: fourGigiByte}},
		VolumeName:       "pvc-1",
		StorageClassName: &lvmScName,
		VolumeMode:       &blockFS,
	},
	Status: corev1.PersistentVolumeClaimStatus{
		Phase: corev1.PersistentVolumeClaimPhase(corev1.VolumeBound),
	},
}
