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

package storage

import (
	"fmt"
	"reflect"
	"testing"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	fakecstor "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var goodNode = corev1.Node{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Node",
		APIVersion: "core/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "node1"},
	Status: corev1.NodeStatus{Phase: corev1.NodeRunning,
		Conditions: []corev1.NodeCondition{
			{Type: corev1.NodeReady, Status: corev1.ConditionTrue},
			{Type: corev1.NodeMemoryPressure, Status: corev1.ConditionFalse},
			{Type: corev1.NodeDiskPressure, Status: corev1.ConditionFalse},
			{Type: corev1.NodePIDPressure, Status: corev1.ConditionFalse},
			{Type: corev1.NodeNetworkUnavailable, Status: corev1.ConditionFalse},
		},
	}}

var node2 = corev1.Node{
	TypeMeta: metav1.TypeMeta{
		Kind:       "Node",
		APIVersion: "core/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "node2"},
	Status: corev1.NodeStatus{Phase: corev1.NodePending,
		Conditions: []corev1.NodeCondition{
			{Type: corev1.NodeReady, Status: corev1.ConditionFalse},
			{Type: corev1.NodeMemoryPressure, Status: corev1.ConditionTrue},
			{Type: corev1.NodeDiskPressure, Status: corev1.ConditionFalse},
			{Type: corev1.NodePIDPressure, Status: corev1.ConditionFalse},
			{Type: corev1.NodeNetworkUnavailable, Status: corev1.ConditionFalse},
		}},
}

func TestCSPCnodeChange(t *testing.T) {
	type args struct {
		k        *client.K8sClient
		poolName string
		oldNode  string
		newNode  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"No cspc found", args{
			k: &client.K8sClient{
				Ns:        "openebs",
				OpenebsCS: fakecstor.NewSimpleClientset(),
			},
			poolName: "fake-pool",
			oldNode:  "node1",
			newNode:  "node2",
		},
			true,
		},
		{"CSPC found, but newNode is not ready", args{
			k: &client.K8sClient{
				Ns:        "openebs",
				K8sCS:     fake.NewSimpleClientset(&node2),
				OpenebsCS: fakecstor.NewSimpleClientset(&cspc),
			},
			poolName: "cassandra-pool", oldNode: "node1", newNode: "node2"}, true},
		{"CSPC found, but newNode does not exist", args{
			k: &client.K8sClient{
				Ns:        "openebs",
				K8sCS:     fake.NewSimpleClientset(&goodNode),
				OpenebsCS: fakecstor.NewSimpleClientset(&cspc),
			},
			poolName: "cassandra-pool", oldNode: "node3", newNode: "node-456"}, true},
		{
			"CSPC found, newNode exists and is ready but old-node name does not match", args{
				k: &client.K8sClient{
					Ns:        "openebs",
					K8sCS:     fake.NewSimpleClientset(&goodNode),
					OpenebsCS: fakecstor.NewSimpleClientset(&cspc),
				},
				poolName: "cassandra-pool", oldNode: "bad-node", newNode: "node1"}, false,
		},
		{
			"CSPC found, newNode exists and is ready, old-node name matches", args{
				k: &client.K8sClient{
					Ns:        "openebs",
					K8sCS:     fake.NewSimpleClientset(&goodNode),
					OpenebsCS: fakecstor.NewSimpleClientset(&cspc),
				},
				poolName: "cassandra-pool", oldNode: "bad-node", newNode: "node1"}, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CSPCnodeChange(tt.args.k, tt.args.poolName, tt.args.oldNode, tt.args.newNode); (err != nil) != tt.wantErr {
				t.Errorf("CSPCnodeChange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}


var cspc1 = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "cspc1", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{{
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1n1"}, {BlockDeviceName: "bd2n1"}}}},
		PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
			DataRaidGroups: []cstorv1.RaidGroup{{
				CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1n2"}, {BlockDeviceName: "bd2n2"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}}}},
}

var cspc2 = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "cspc1", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{{
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1n1"}, {BlockDeviceName: "bd2n1"}}}},
		PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node3"},
			DataRaidGroups: []cstorv1.RaidGroup{{
				CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1n2"}, {BlockDeviceName: "bd2n2"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}}}},
}

var goodBD1N1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1n1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd1n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var goodBD2N1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd2n1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd1n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var goodBD1N2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1n2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd1n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var goodBD2N2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd2n2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd1n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

var cspc3 = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "cspc1", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{{
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1n1"}, {BlockDeviceName: "bd2n1"}}}},
		PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
			DataRaidGroups: []cstorv1.RaidGroup{{
				CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1n2"}, {BlockDeviceName: "bd2n2"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}}}},
}

func TestDebugCSPCNode(t *testing.T) {
	type args struct {
		k    *client.K8sClient
		cspc string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
		err     error
	}{
		{"unable to find CSPC", args{k: &client.K8sClient{
			Ns: "random", OpenebsCS: fakecstor.NewSimpleClientset()}, cspc: "cspc1"}, nil, true,
			fmt.Errorf(`unable to get cspc cspc1, Error while getting cspc: cstorpoolclusters.cstor.openebs.io "cspc1" not found`)},
		{"CSPC exists, not in the guessed namespace so unable to find it", args{
			k: &client.K8sClient{Ns: "wrongNS", OpenebsCS: fakecstor.NewSimpleClientset(&cspc1)}, cspc: "cspc1"}, nil, true,
			fmt.Errorf(`unable to get cspc cspc1, Error while getting cspc: cstorpoolclusters.cstor.openebs.io "cspc1" not found`)},
		{"CSPC exists, blockdevices are in the same node as expected", args{
			k: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspc1, &goodBD1N1, &goodBD2N1,
				&goodBD1N2, &goodBD2N2)}, cspc: "cspc1"}, nil, true,
			fmt.Errorf(`no change in the storage node`)},
		// it'd make sense to evaluate the error in the above Test suite somehow
		{"CSPC exists, two BD's loc changed to other node", args{
			k: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspc2, &goodBD1N1, &goodBD2N1,
				&goodBD1N2, &goodBD2N2)}, cspc: "cspc1"},
			map[string]string{"node3": "node2"}, false, nil},
		{"CSPC exists, BD nodes swapped", args{
			k: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspc3, &goodBD1N1, &goodBD2N1,
				&goodBD1N2, &goodBD2N2)}, cspc: "cspc1"}, nil, true,
			fmt.Errorf(`more than one node change in the storage instance`)},
		{"CSPC exists, all 4 BDs are missing", args{
			k: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspc3)}, cspc: "cspc1"}, nil, true,
			fmt.Errorf("%d blockdevices are missing from the cluster", 4)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DebugCSPCNode(tt.args.k, tt.args.cspc)
			if (err != nil) != tt.wantErr {
				t.Errorf("DebugCSPCNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && !reflect.DeepEqual(err.Error(), tt.err.Error()) {
				t.Logf("DebugCSPCNode() Got %v error, wanted %v error", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DebugCSPCNode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
