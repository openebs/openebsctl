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
	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
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
