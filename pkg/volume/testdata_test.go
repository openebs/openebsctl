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

package volume

import (
	"time"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	lvm "github.com/openebs/lvm-localpv/pkg/apis/openebs.io/lvm/v1alpha1"
	"github.com/openebs/openebsctl/pkg/util"
	zfs "github.com/openebs/zfs-localpv/pkg/apis/openebs.io/zfs/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Some storage sizes for PVs
var (
	fourGigiByte = resource.MustParse("4Gi")
	blockFS      = corev1.PersistentVolumeBlock
)

/****************
* CSTOR
****************/

var nsCstor = corev1.Namespace{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cstor",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{},
		Finalizers:        []string{},
	},
	Spec: corev1.NamespaceSpec{Finalizers: []corev1.FinalizerName{corev1.FinalizerKubernetes}},
}

var cv1 = v1.CStorVolume{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-1",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Spec: v1.CStorVolumeSpec{
		Capacity:                 fourGigiByte,
		TargetIP:                 "10.2.2.2",
		TargetPort:               "3002",
		Iqn:                      "pvc1-some-fake-iqn",
		TargetPortal:             "10.2.2.2:3002",
		ReplicationFactor:        3,
		ConsistencyFactor:        0,
		DesiredReplicationFactor: 0,
		ReplicaDetails: v1.CStorVolumeReplicaDetails{KnownReplicas: map[v1.ReplicaID]string{
			"some-id-1": "pvc-1-rep-1", "some-id-2": "pvc-1-rep-2", "some-id-3": "pvc-1-rep-3"},
		},
	},
	Status: v1.CStorVolumeStatus{
		Phase:           util.Healthy,
		ReplicaStatuses: []v1.ReplicaStatus{{ID: "some-id-1", Mode: "Healthy"}, {ID: "some-id-2", Mode: "Healthy"}, {ID: "some-id-3", Mode: "Healthy"}},
		Capacity:        fourGigiByte,
		ReplicaDetails: v1.CStorVolumeReplicaDetails{KnownReplicas: map[v1.ReplicaID]string{
			"some-id-1": "pvc-1-rep-1", "some-id-2": "pvc-1-rep-2", "some-id-3": "pvc-1-rep-3"},
		},
	},
	VersionDetails: v1.VersionDetails{
		AutoUpgrade: false,
		Desired:     "2.11.0",
		Status: v1.VersionStatus{
			DependentsUpgraded: true,
			Current:            "2.11.0",
			LastUpdateTime:     metav1.Time{},
		},
	},
}

var cv2 = v1.CStorVolume{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-2",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Spec: v1.CStorVolumeSpec{
		Capacity:                 fourGigiByte,
		TargetIP:                 "10.2.2.2",
		TargetPort:               "3002",
		Iqn:                      "pvc1-some-fake-iqn",
		TargetPortal:             "10.2.2.2:3002",
		ReplicationFactor:        3,
		ConsistencyFactor:        0,
		DesiredReplicationFactor: 0,
		ReplicaDetails: v1.CStorVolumeReplicaDetails{KnownReplicas: map[v1.ReplicaID]string{
			"some-id-1": "pvc-2-rep-1"},
		},
	},
	Status: v1.CStorVolumeStatus{
		Phase:           util.Healthy,
		ReplicaStatuses: []v1.ReplicaStatus{{ID: "some-id-1", Mode: "Healthy"}},
		Capacity:        fourGigiByte,
		ReplicaDetails: v1.CStorVolumeReplicaDetails{KnownReplicas: map[v1.ReplicaID]string{
			"some-id-1": "pvc-2-rep-1"},
		},
	},
	VersionDetails: v1.VersionDetails{
		AutoUpgrade: false,
		Desired:     "2.11.0",
		Status: v1.VersionStatus{
			DependentsUpgraded: true,
			Current:            "2.11.0",
			LastUpdateTime:     metav1.Time{},
		},
	},
}

var cvc1 = v1.CStorVolumeConfig{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-1",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Spec: v1.CStorVolumeConfigSpec{Provision: v1.VolumeProvision{
		Capacity:     corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
		ReplicaCount: 3,
	}},
	Publish: v1.CStorVolumeConfigPublish{},
	Status:  v1.CStorVolumeConfigStatus{PoolInfo: []string{"pool-1", "pool-2", "pool-3"}},
	VersionDetails: v1.VersionDetails{
		AutoUpgrade: false,
		Desired:     "2.11.0",
		Status:      v1.VersionStatus{Current: "2.11.0"},
	},
}

var cvc2 = v1.CStorVolumeConfig{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-2",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Spec: v1.CStorVolumeConfigSpec{Provision: v1.VolumeProvision{
		Capacity:     corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
		ReplicaCount: 3,
	}},
	Publish: v1.CStorVolumeConfigPublish{},
	Status:  v1.CStorVolumeConfigStatus{PoolInfo: []string{"pool-1"}},
	VersionDetails: v1.VersionDetails{
		AutoUpgrade: false,
		Desired:     "2.11.0",
		Status:      v1.VersionStatus{Current: "2.11.0"},
	},
}

var cva1 = v1.CStorVolumeAttachment{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-1-cva",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"Volname": "pvc-1", "nodeID": "node-1"},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Spec: v1.CStorVolumeAttachmentSpec{Volume: v1.VolumeInfo{OwnerNodeID: "node-1"}},
}

var cva2 = v1.CStorVolumeAttachment{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-2-cva",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"Volname": "pvc-2", "nodeID": "node-2"},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Spec: v1.CStorVolumeAttachmentSpec{Volume: v1.VolumeInfo{OwnerNodeID: "node-2"}},
}

var cvr1 = v1.CStorVolumeReplica{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-1-rep-1",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-1"},
		Finalizers:        []string{},
		Namespace:         "cstor",
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
		Name:              "pvc-1-rep-2",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-1"},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Status: v1.CStorVolumeReplicaStatus{
		Capacity: v1.CStorVolumeReplicaCapacityDetails{
			Total: "4Gi",
			Used:  "70Mi",
		},
		Phase: v1.CVRStatusOnline,
	},
}

var cvr3 = v1.CStorVolumeReplica{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-1-rep-3",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-1"},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Status: v1.CStorVolumeReplicaStatus{
		Capacity: v1.CStorVolumeReplicaCapacityDetails{
			Total: "4Gi",
			Used:  "70Mi",
		},
		Phase: v1.CVRStatusOnline,
	},
}

var cvr4 = v1.CStorVolumeReplica{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-2-rep-1",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-2"},
		Finalizers:        []string{},
		Namespace:         "cstor",
	},
	Status: v1.CStorVolumeReplicaStatus{
		Capacity: v1.CStorVolumeReplicaCapacityDetails{
			Total: "4Gi",
			Used:  "70Mi",
		},
		Phase: v1.CVRStatusOnline,
	},
}

var (
	cstorScName     = "cstor-sc"
	cstorVolumeMode = corev1.PersistentVolumeFilesystem
	cstorPVC1       = corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "cstor-pvc-1",
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-2"},
			Finalizers:        []string{},
			Namespace:         "default",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources:        corev1.ResourceRequirements{Requests: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: fourGigiByte}},
			VolumeName:       "pvc-1",
			StorageClassName: &cstorScName,
			VolumeMode:       &cstorVolumeMode,
		},
		Status: corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimBound, Capacity: corev1.ResourceList{corev1.ResourceStorage: fourGigiByte}},
	}
)

var (
	cstorPVC2 = corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "cstor-pvc-2",
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-2"},
			Finalizers:        []string{},
			Namespace:         "default",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes:      []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources:        corev1.ResourceRequirements{Requests: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: fourGigiByte}},
			VolumeName:       "pvc-2",
			StorageClassName: &cstorScName,
			VolumeMode:       &cstorVolumeMode,
		},
		Status: corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimBound, Capacity: corev1.ResourceList{corev1.ResourceStorage: fourGigiByte}},
	}
)

var (
	cstorPV1 = corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "pvc-1",
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-1"},
			Finalizers:        []string{},
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity:    corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			ClaimRef: &corev1.ObjectReference{
				Namespace: "default",
				Name:      "cstor-pvc-1",
			},
			PersistentVolumeReclaimPolicy: "Retain",
			StorageClassName:              cstorScName,
			VolumeMode:                    &cstorVolumeMode,
			PersistentVolumeSource: corev1.PersistentVolumeSource{CSI: &corev1.CSIPersistentVolumeSource{
				Driver: "cstor.csi.openebs.io",
			}},
		},
		Status: corev1.PersistentVolumeStatus{Phase: corev1.VolumeBound},
	}
)

var (
	cstorPV2 = corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "pvc-2",
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-2"},
			Finalizers:        []string{},
		},
		Spec: corev1.PersistentVolumeSpec{
			Capacity:    corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
			AccessModes: []corev1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			ClaimRef: &corev1.ObjectReference{
				Namespace: "default",
				Name:      "cstor-pvc-2",
			},
			PersistentVolumeReclaimPolicy: "Retain",
			StorageClassName:              cstorScName,
			VolumeMode:                    &cstorVolumeMode,
			PersistentVolumeSource: corev1.PersistentVolumeSource{CSI: &corev1.CSIPersistentVolumeSource{
				Driver: "cstor.csi.openebs.io",
			}},
		},
		Status: corev1.PersistentVolumeStatus{Phase: corev1.VolumeBound},
	}
)

var cbkp = v1.CStorBackup{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "bkp-name",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-1"},
		Finalizers:        []string{},
	},
	Spec: v1.CStorBackupSpec{
		BackupName:   "bkp-name",
		VolumeName:   "pvc-1",
		SnapName:     "snap-name",
		PrevSnapName: "prev-snap-name",
		BackupDest:   "10.2.2.7",
		LocalSnap:    true,
	},
	Status: v1.BKPCStorStatusDone,
}

var ccbkp = v1.CStorCompletedBackup{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "completed-bkp-name",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-1"},
		Finalizers:        []string{},
	},
	Spec: v1.CStorCompletedBackupSpec{
		BackupName:         "completed-bkp-name",
		VolumeName:         "pvc-1",
		SecondLastSnapName: "secondlast-snapshot-name",
		LastSnapName:       "last-snapshot-name",
	},
}

var crestore = v1.CStorRestore{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "restore-name",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-1"},
		Finalizers:        []string{},
	},
	Spec: v1.CStorRestoreSpec{
		RestoreName:   "restore-name",
		VolumeName:    "pvc-1",
		RestoreSrc:    "10.2.2.7",
		MaxRetryCount: 3,
		RetryCount:    2,
		StorageClass:  "cstor-sc",
		Size:          fourGigiByte,
		Local:         true,
	},
}

/****************
* LVM LOCAL PV
****************/

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

var localpvCSICtrlSTS = appsv1.StatefulSet{
	TypeMeta: metav1.TypeMeta{
		Kind:       "StatefulSet",
		APIVersion: "apps/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fake-LVM-CSI",
		Namespace: "lvm",
		Labels: map[string]string{
			"openebs.io/version":        "1.9.0",
			"openebs.io/component-name": "openebs-lvm-controller"},
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

var localpvzfsCSICtrlSTS = appsv1.StatefulSet{
	TypeMeta: metav1.TypeMeta{
		Kind:       "StatefulSet",
		APIVersion: "apps/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fake-ZFS-CSI",
		Namespace: "zfslocalpv",
		Labels: map[string]string{
			"openebs.io/version":        "1.9.0",
			"openebs.io/component-name": "openebs-zfs-controller"},
	},
}

/****************
* JIVA
****************/

// var nsJiva = corev1.Namespace{
// 	ObjectMeta: metav1.ObjectMeta{
// 		Name:              "jiva",
// 		CreationTimestamp: metav1.Time{Time: time.Now()},
// 		Labels:            map[string]string{},
// 		Finalizers:        []string{},
// 	},
// 	Spec: corev1.NamespaceSpec{Finalizers: []corev1.FinalizerName{corev1.FinalizerKubernetes}},
// }

// pvc-1 JivaVolume from jiva namespace attached on worker-node-1 & 1-replica & 2.10.0
// var jv1 = v1alpha1.JivaVolume{
//	TypeMeta: metav1.TypeMeta{},
//	ObjectMeta: metav1.ObjectMeta{
//		Name:      "pvc-1",
//		Namespace: "jiva",
//		Labels:    map[string]string{"nodeID": "worker-node-1"},
//	},
//	Spec: v1alpha1.JivaVolumeSpec{},
//	Status: v1alpha1.JivaVolumeStatus{
//		Status:          "RW",
//		ReplicaCount:    1,
//		ReplicaStatuses: nil, // TODO
//		Phase:           "Attached",
//	},
//	VersionDetails: v1alpha1.VersionDetails{
//		AutoUpgrade: false,
//		Desired:     "2.10.0",
//		Status: v1alpha1.VersionStatus{
//			DependentsUpgraded: false,
//			Current:            "2.10.0",
//		},
//	},
// }
//// pvc-2 JivaVolume from jiva namespace attached on worker-node-2, two replicas & 2.10.0
// var jv2 = v1alpha1.JivaVolume{
// 	TypeMeta: metav1.TypeMeta{},
// 	ObjectMeta: metav1.ObjectMeta{
// 		Name:      "pvc-2",
// 		Namespace: "jiva",
// 		Labels:    map[string]string{"nodeID": "worker-node-2"},
// 	},
// 	Spec: v1alpha1.JivaVolumeSpec{
// 		PV:         "pvc-2",
// 		Capacity:   "4Gi",
// 		AccessType: "",
// 		ISCSISpec: v1alpha1.ISCSISpec{
// 			TargetIP:   "1.2.3.4",
// 			TargetPort: 8080,
// 			Iqn:        "nice-iqn",
// 		},
// 		MountInfo: v1alpha1.MountInfo{
// 			StagingPath: "/home/staging/",
// 			TargetPath:  "/home/target",
// 			FSType:      "ext4",
// 			DevicePath:  "",
// 		},
// 		Policy:                   v1alpha1.JivaVolumePolicySpec{},
// 		DesiredReplicationFactor: 0,
// 	},
// 	Status: v1alpha1.JivaVolumeStatus{
// 		Status:       "RO",
// 		ReplicaCount: 2,
// 		ReplicaStatuses: []v1alpha1.ReplicaStatus{
// 			{Address: "tcp://192.168.2.7:9502", Mode: "RW"},
// 			{Address: "tcp://192.168.2.8:9502", Mode: "RO"},
// 		},
// 		Phase: "Ready",
// 	},
// 	VersionDetails: v1alpha1.VersionDetails{
// 		AutoUpgrade: false,
// 		Desired:     "2.10.0",
// 		Status: v1alpha1.VersionStatus{
// 			DependentsUpgraded: false,
// 			Current:            "2.10.0",
// 		},
// 	},
// }
var jivaPV1 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "pvc-1",
		Namespace:   "jiva",
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	},
	Spec: corev1.PersistentVolumeSpec{
		// 4GiB
		Capacity:               corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
		PersistentVolumeSource: corev1.PersistentVolumeSource{CSI: &corev1.CSIPersistentVolumeSource{Driver: util.JivaCSIDriver}},
		AccessModes:            []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		ClaimRef: &corev1.ObjectReference{
			Kind:            "PersistentVolumeClaim",
			Namespace:       "jiva-app",
			Name:            "mongo-jiva",
			APIVersion:      "v1",
			ResourceVersion: "123"},
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
		StorageClassName:              "pvc-1-sc",
		VolumeMode:                    &blockFS,
		NodeAffinity: &corev1.VolumeNodeAffinity{
			Required: &corev1.NodeSelector{NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{MatchExpressions: []corev1.NodeSelectorRequirement{
					{Key: "kubernetes.io/hostname", Operator: corev1.NodeSelectorOpIn, Values: []string{"node1"}},
				}},
			}},
		},
	},
	Status: corev1.PersistentVolumeStatus{
		Phase:   corev1.VolumeBound,
		Message: "",
		Reason:  "",
	},
}

//var jivaPV2 = corev1.PersistentVolume{
//	TypeMeta: metav1.TypeMeta{},
//	ObjectMeta: metav1.ObjectMeta{
//		Name:        "pvc-2",
//		Namespace:   "jiva",
//		Labels:      map[string]string{},
//		Annotations: map[string]string{},
//	},
//	Spec: corev1.PersistentVolumeSpec{
//		// 4GiB
//		Capacity:                      corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
//		PersistentVolumeSource:        corev1.PersistentVolumeSource{CSI: &corev1.CSIPersistentVolumeSource{Driver: util.JivaCSIDriver}},
//		AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
//		ClaimRef:                      nil,
//		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
//		StorageClassName:              "pvc-2-sc",
//		VolumeMode:                    &blockFS,
//		NodeAffinity: &corev1.VolumeNodeAffinity{
//			Required: &corev1.NodeSelector{NodeSelectorTerms: []corev1.NodeSelectorTerm{
//				{MatchExpressions: []corev1.NodeSelectorRequirement{
//					{Key: "kubernetes.io/hostname", Operator: corev1.NodeSelectorOpIn, Values: []string{"node2"}},
//				}},
//			}},
//		},
//	},
//	Status: corev1.PersistentVolumeStatus{
//		Phase:   corev1.VolumePending,
//		Message: "Storage class not found",
//		Reason:  "K8s API was down",
//	},
//}
var pv2 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name: "pvc-1",
	},
	Spec: corev1.PersistentVolumeSpec{
		Capacity: corev1.ResourceList{corev1.ResourceStorage: resource.Quantity{}},
	},
	Status: corev1.PersistentVolumeStatus{},
}
var pv3 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name: "pvc-1",
	},
	Spec:   corev1.PersistentVolumeSpec{},
	Status: corev1.PersistentVolumeStatus{},
}

/****************
* Local Hostpath
****************/
var localHostpathVolumeCapacity = corev1.ResourceList{corev1.ResourceStorage: fourGigiByte}

var localHostpathPv1 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1",
		Namespace: "localhostpath",
		Labels: map[string]string{
			"openebs.io/component-name": "openebs-localpv-provisioner",
			"openebs.io/cas-type":       "local-hostpath",
			"openebs.io/version":        "1.9.0",
		},
		Annotations: map[string]string{},
	},
	Spec: corev1.PersistentVolumeSpec{
		// 4GiB
		Capacity: localHostpathVolumeCapacity,
		PersistentVolumeSource: corev1.PersistentVolumeSource{
			Local: &corev1.LocalVolumeSource{Path: "/random/path"},
		},
		AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		ClaimRef: &corev1.ObjectReference{
			Kind:            "PersistentVolumeClaim",
			Namespace:       "local-app",
			Name:            "mongo-local",
			APIVersion:      "v1",
			ResourceVersion: "123"},
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
		StorageClassName:              "pvc-1-local",
		VolumeMode:                    &blockFS,
		NodeAffinity: &corev1.VolumeNodeAffinity{
			Required: &corev1.NodeSelector{NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{MatchExpressions: []corev1.NodeSelectorRequirement{
					{Key: "kubernetes.io/hostname", Operator: corev1.NodeSelectorOpIn, Values: []string{"node1"}},
				}},
			}},
		},
	},
	Status: corev1.PersistentVolumeStatus{
		Phase:   corev1.VolumeBound,
		Message: "",
		Reason:  "",
	},
}
var localpvHostpathDpl1 = appsv1.Deployment{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fake-deploy-hostpath-1",
		Namespace: "openebs",
		Labels: map[string]string{
			"openebs.io/version":        "1.9.0",
			"openebs.io/component-name": "openebs-localpv-provisioner"},
	},
}

var localpvHostpathDpl2 = appsv1.Deployment{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fake-deploy-hostpath2",
		Namespace: "openebs",
		Labels: map[string]string{
			"openebs.io/version":        "1.9.0",
			"openebs.io/component-name": "openebs-localpv-provisioner"},
	},
}
