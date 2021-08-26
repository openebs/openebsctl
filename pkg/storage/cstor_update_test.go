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
	"testing"

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
