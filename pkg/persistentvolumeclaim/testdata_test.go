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

package persistentvolumeclaim

import (
	"time"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	fourGigiByte = resource.MustParse("4Gi")
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

var cvrList = v1.CStorVolumeReplicaList{Items: []v1.CStorVolumeReplica{cvr1, cvr2}}

var cstorSc = v12.StorageClass{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cstor-sc",
		CreationTimestamp: metav1.Time{Time: time.Now()},
	},
	Provisioner: "cstor.csi.openebs.io",
	Parameters:  map[string]string{"cstorPoolCluster": "cspc"},
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

var cstorTargetPod = corev1.Pod{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "restore-name",
		Namespace:         "cstor",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{"openebs.io/persistent-volume-claim": "cstor-pvc-1", "openebs.io/persistent-volume": "pvc-1", "openebs.io/target": "cstor-target"},
		Finalizers:        []string{},
	},
	Spec:   corev1.PodSpec{NodeName: "node-1"},
	Status: corev1.PodStatus{ContainerStatuses: []corev1.ContainerStatus{{Ready: true}, {Ready: true}, {Ready: true}}, PodIP: "10.2.2.2", Phase: "Running"},
}

var cspc = v1.CStorPoolCluster{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cspc",
		Namespace:         "cstor",
		CreationTimestamp: metav1.Time{Time: time.Now()},
	},
	Spec: v1.CStorPoolClusterSpec{
		Pools: []v1.PoolSpec{{
			DataRaidGroups: []v1.RaidGroup{
				{CStorPoolInstanceBlockDevices: []v1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd-1"}}},
				{CStorPoolInstanceBlockDevices: []v1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd-2"}}},
				{CStorPoolInstanceBlockDevices: []v1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd-3"}}},
			},
		}},
	},
}

var cspi1 = v1.CStorPoolInstance{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cspc-1",
		Namespace:         "cstor",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels: map[string]string{
			"openebs.io/cstor-pool-cluster": "cspc",
			"openebs.io/cas-type":           "cstor",
		},
	},
	Spec: v1.CStorPoolInstanceSpec{
		DataRaidGroups: []v1.RaidGroup{
			{CStorPoolInstanceBlockDevices: []v1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd-1"}}},
		},
	},
	Status: v1.CStorPoolInstanceStatus{Phase: "ONLINE"},
}

var cspi2 = v1.CStorPoolInstance{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cspc-2",
		Namespace:         "cstor",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels: map[string]string{
			"openebs.io/cstor-pool-cluster": "cspc",
			"openebs.io/cas-type":           "cstor",
		},
	},
	Spec: v1.CStorPoolInstanceSpec{
		DataRaidGroups: []v1.RaidGroup{
			{CStorPoolInstanceBlockDevices: []v1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd-2"}}},
		},
	},
	Status: v1.CStorPoolInstanceStatus{Phase: "ONLINE"},
}

var cspiList = v1.CStorPoolInstanceList{Items: []v1.CStorPoolInstance{cspi1, cspi2}}

/****************
* BDC & BDCs
 ****************/

var bd1 = v1alpha1.BlockDevice{
	TypeMeta:   metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{Name: "bd-1", Namespace: "cstor"},
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
var bd2 = v1alpha1.BlockDevice{
	TypeMeta:   metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{Name: "bd-2", Namespace: "cstor"},
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

var bdList = v1alpha1.BlockDeviceList{
	Items: []v1alpha1.BlockDevice{bd1, bd2},
}

var bdc1 = v1alpha1.BlockDeviceClaim{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "bdc-1",
		Namespace:         "cstor",
		CreationTimestamp: metav1.Time{Time: time.Now()},
	},
	Status: v1alpha1.DeviceClaimStatus{Phase: "Bound"},
}

var bdc2 = v1alpha1.BlockDeviceClaim{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "bdc-2",
		Namespace:         "cstor",
		CreationTimestamp: metav1.Time{Time: time.Now()},
	},
	Status: v1alpha1.DeviceClaimStatus{Phase: "Bound"},
}

var bdcList = v1alpha1.BlockDeviceClaimList{Items: []v1alpha1.BlockDeviceClaim{bdc1, bdc2}}

var expectedBDs = map[string]bool{
	"bdc-1": true,
	"bdc-2": true,
	"bdc-3": false,
}

/****************
* EVENTS
****************/

var pvcEvent1 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cstor-pvc-1.time1",
		Namespace: "default",
		UID:       "some-random-event-uuid-1",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "PersistentVolumeClaim",
		Namespace: "default",
		Name:      "cstor-pvc-1",
		UID:       "some-random-pvc-uuid-1",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var pvcEvent2 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cstor-pvc-1.time2",
		Namespace: "default",
		UID:       "some-random-event-uuid-2",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "PersistentVolumeClaim",
		Namespace: "default",
		Name:      "cstor-pvc-1",
		UID:       "some-random-pvc-uuid-1",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var cvcEvent1 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1.time1",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-3",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "CStorVolumeConfig",
		Namespace: "cstor",
		Name:      "pvc-1",
		UID:       "some-random-cvc-uuid-1",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var cvcEvent2 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1.time2",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-4",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "CStorVolumeConfig",
		Namespace: "cstor",
		Name:      "pvc-1",
		UID:       "some-random-cvc-uuid-2",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var bdcEvent1 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "bdc-1.time1",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-5",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "BlockDeviceClaim",
		Namespace: "cstor",
		Name:      "bdc-1",
		UID:       "some-random-bdc-uuid-1",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var bdcEvent2 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "bdc-2.time1",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-6",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "BlockDeviceClaim",
		Namespace: "cstor",
		Name:      "bdc-2",
		UID:       "some-random-bdc-uuid-1",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var cspiEvent1 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cspc-1.time1",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-7",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "CStorPoolInstance",
		Namespace: "cstor",
		Name:      "cspc-1",
		UID:       "some-random-cspi-uuid-1",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var cspiEvent2 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cspc-2.time1",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-8",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "CStorPoolInstance",
		Namespace: "cstor",
		Name:      "cspc-2",
		UID:       "some-random-cspi-uuid-2",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var cspcEvent = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "cspc.time1",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-9",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "CStorPoolCluster",
		Namespace: "cstor",
		Name:      "cspc",
		UID:       "some-random-cspc-uuid-1",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var cvrEvent1 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1-rep-1.time1",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-10",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "CStorVolumeReplica",
		Namespace: "cstor",
		Name:      "pvc-1-rep-1",
		UID:       "some-random-cvr-uuid-1",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}

var cvrEvent2 = corev1.Event{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1-rep-2.time1",
		Namespace: "cstor",
		UID:       "some-random-event-uuid-11",
	},
	InvolvedObject: corev1.ObjectReference{
		Kind:      "CStorVolumeReplica",
		Namespace: "cstor",
		Name:      "pvc-1-rep-2",
		UID:       "some-random-cvr-uuid-2",
	},
	Reason:  "some-fake-reason",
	Message: "some-fake-message",
	Count:   1,
	Type:    "Warning",
	Action:  "some-fake-action",
}
