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

package blockdevice

import (
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	openebsFakeClientset "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/openebsctl/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
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

func Test_createTreeByNode(t *testing.T) {
	k8sCS := fake.NewSimpleClientset()
	type args struct {
		k   *client.K8sClient
		bds []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with valid bd inputs and across all namespaces",
			args{
				k: &client.K8sClient{
					Ns:        "",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: nil,
			},
			false,
		},
		{
			"Test with valid bd inputs and in some valid ns",
			args{
				k: &client.K8sClient{
					Ns:        "fake-ns",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: nil,
			},
			false,
		},
		{
			"Test with valid bd inputs and in some invalid ns",
			args{
				k: &client.K8sClient{
					Ns:        "fake-invalid-ns",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: nil,
			},
			true,
		},
		{
			"Test with invalid bd inputs and in some valid ns",
			args{
				k: &client.K8sClient{
					Ns:        "fake-ns",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(),
				},
				bds: nil,
			},
			true,
		},
		{
			"Test with invalid bd inputs across all namespaces",
			args{
				k: &client.K8sClient{
					Ns:        "",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(),
				},
				bds: nil,
			},
			true,
		},
		{
			"Test with valid bd inputs across all namespaces with some valid bd name passed as args",
			args{
				k: &client.K8sClient{
					Ns:        "",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: []string{"some-fake-bd-3"},
			},
			false,
		},
		{
			"Test with valid bd inputs across all namespaces with multiple valid bd names passed as args",
			args{
				k: &client.K8sClient{
					Ns:        "",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: []string{"some-fake-bd-3", "some-fake-bd-2"},
			},
			false,
		},
		{
			"Test with valid bd inputs across all namespaces with some valid and some invalid bd names passed as args",
			args{
				k: &client.K8sClient{
					Ns:        "",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: []string{"some-fake-bd-365", "some-fake-bd-2"},
			},
			false,
		},
		{
			"Test with valid bd inputs across all namespaces with some invalid bd name passed as args",
			args{
				k: &client.K8sClient{
					Ns:        "",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: []string{"some-fake-bd-365"},
			},
			true,
		},
		{
			"Test with valid bd inputs in a namespace with some valid bd name passed as args",
			args{
				k: &client.K8sClient{
					Ns:        "fake-ns",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: []string{"some-fake-bd-3"},
			},
			false,
		},
		{
			"Test with valid bd inputs in an invalid namespace with some valid bd name passed as args",
			args{
				k: &client.K8sClient{
					Ns:        "fake-invalid-ns",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: []string{"some-fake-bd-3"},
			},
			true,
		},
		{
			"Test with valid bd inputs in a valid namespace with some valid bd name passed as args",
			args{
				k: &client.K8sClient{
					Ns:        "fake-ns",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: []string{"some-fake-bd-3"},
			},
			false,
		},
		{
			"Test with valid bd inputs in a valid namespace with some invalid bd name passed as args",
			args{
				k: &client.K8sClient{
					Ns:        "fake-ns",
					K8sCS:     k8sCS,
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&bd1, &bd2, &bd3),
				},
				bds: []string{"some-fake-bd-365"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := createTreeByNode(tt.args.k, tt.args.bds); (err != nil) != tt.wantErr {
				t.Errorf("createTreeByNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
