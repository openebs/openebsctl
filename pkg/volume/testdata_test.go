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

package volume

import (
	"github.com/openebs/jiva-operator/pkg/apis/openebs/v1alpha1"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)
// Some storage sizes for PVs
var (
	fourGigiByte = resource.MustParse("4Gi")
	fiveGigaByte = resource.MustParse("5G")
	fiveGigaBit  = resource.MustParse("5G")
	fiveGigiBit  = resource.MustParse("5Gi")
	blockFS      = corev1.PersistentVolumeBlock
)
var nsJiva = corev1.Namespace{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "jiva",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{},
		Finalizers:        []string{},
	},
	Spec: corev1.NamespaceSpec{Finalizers: []corev1.FinalizerName{corev1.FinalizerKubernetes}},
}
var nsCstor = corev1.Namespace{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "cstor",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{},
		Finalizers:        []string{},
	},
	Spec: corev1.NamespaceSpec{Finalizers: []corev1.FinalizerName{corev1.FinalizerKubernetes}},
}
var nsLocalPV = corev1.Namespace{
	ObjectMeta: metav1.ObjectMeta{
		Name:              "localpv",
		CreationTimestamp: metav1.Time{time.Now()},
		Labels:            map[string]string{},
		Finalizers:        []string{},
	},
	Spec: corev1.NamespaceSpec{Finalizers: []corev1.FinalizerName{corev1.FinalizerKubernetes}},
}
// pvc-1 JivaVolume from jiva namespace attached on worker-node-1 & 1-replica & 2.10.0
var jv1 = v1alpha1.JivaVolume{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1",
		Namespace: "jiva",
		Labels:    map[string]string{"nodeID": "worker-node-1"},
	},
	Spec: v1alpha1.JivaVolumeSpec{},
	Status: v1alpha1.JivaVolumeStatus{
		Status:          "RW",
		ReplicaCount:    1,
		ReplicaStatuses: nil, // TODO
		Phase:           "Attached",
	},
	VersionDetails: v1alpha1.VersionDetails{
		AutoUpgrade: false,
		Desired:     "2.10.0",
		Status: v1alpha1.VersionStatus{
			DependentsUpgraded: false,
			Current:            "2.10.0",
		},
	},
}
// pvc-2 JivaVolume from jiva namespace attached on worker-node-2, two replicas & 2.10.0
var jv2 = v1alpha1.JivaVolume{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-2",
		Namespace: "jiva",
		Labels:    map[string]string{"nodeID": "worker-node-2"},
	},
	Spec: v1alpha1.JivaVolumeSpec{
		PV:         "pvc-2",
		Capacity:   "4Gi",
		AccessType: "",
		ISCSISpec: v1alpha1.ISCSISpec{
			TargetIP:   "1.2.3.4",
			TargetPort: 8080,
			Iqn:        "nice-iqn",
		},
		MountInfo: v1alpha1.MountInfo{
			StagingPath: "/home/staging/",
			TargetPath:  "/home/target",
			FSType:      "ext4",
			DevicePath:  "",
		},
		Policy:                   v1alpha1.JivaVolumePolicySpec{},
		DesiredReplicationFactor: 0,
	},
	Status: v1alpha1.JivaVolumeStatus{
		Status:       "RO",
		ReplicaCount: 2,
		ReplicaStatuses: []v1alpha1.ReplicaStatus{
			{Address: "tcp://192.168.2.7:9502", Mode: "RW"},
			{Address: "tcp://192.168.2.8:9502", Mode: "RO"},
		},
		Phase: "Ready",
	},
	VersionDetails: v1alpha1.VersionDetails{
		AutoUpgrade: false,
		Desired:     "2.10.0",
		Status: v1alpha1.VersionStatus{
			DependentsUpgraded: false,
			Current:            "2.10.0",
		},
	},
}
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
var jivaPV2 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "pvc-2",
		Namespace:   "jiva",
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	},
	Spec: corev1.PersistentVolumeSpec{
		// 4GiB
		Capacity:                      corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
		PersistentVolumeSource:        corev1.PersistentVolumeSource{CSI: &corev1.CSIPersistentVolumeSource{Driver: util.JivaCSIDriver}},
		AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		ClaimRef:                      nil,
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
		StorageClassName:              "pvc-2-sc",
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
		Phase:   corev1.VolumePending,
		Message: "Storage class not found",
		Reason:  "K8s API was down",
	},
}
var pv2 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Name: "pvc-1",
	},
	Spec:   corev1.PersistentVolumeSpec{
		Capacity: corev1.ResourceList{corev1.ResourceStorage:resource.Quantity{}},
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