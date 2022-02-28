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
	"fmt"
	"reflect"
	"testing"

	fakelvmclient "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset/fake"
	fakelvm "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset/typed/lvm/v1alpha1/fake"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stest "k8s.io/client-go/testing"
)

func TestGetVolumeGroup(t *testing.T) {
	type args struct {
		c       *client.K8sClient
		vg      []string
		lvmfunc func(*client.K8sClient)
	}
	tests := []struct {
		name    string
		args    args
		want    []metav1.TableRow
		wantErr bool
	}{
		{
			"no LVM volumegroups present",
			args{
				c: &client.K8sClient{
					Ns:        "lvmlocalpv",
					K8sCS:     nil,
					OpenebsCS: nil,
					LVMCS:     fakelvmclient.NewSimpleClientset()},
				vg:      nil,
				lvmfunc: lvnNodeNotFound,
			},
			nil,
			true,
		},
		{
			name: "four LVM volumegroup present on two nodes",
			args: args{
				c: &client.K8sClient{
					Ns:    "lvmlocalpv",
					LVMCS: fakelvmclient.NewSimpleClientset(&lvmNode1, &lvmNode2),
				},
				vg: nil,
			},
			want: []metav1.TableRow{
				{Cells: []interface{}{"node1", "", "", ""}},
				{Cells: []interface{}{firstElemPrefix + "lvmvg", "4.0GiB", "5.0GiB"}},
				{Cells: []interface{}{lastElemPrefix + "lvmvg2", "4.0GiB", "5.0GiB"}},
				{Cells: []interface{}{"", "", ""}},
				{Cells: []interface{}{"node2", "", "", ""}},
				{Cells: []interface{}{firstElemPrefix + "lvmvg", "4.0GiB", "5.0GiB"}},
				{Cells: []interface{}{lastElemPrefix + "lvmvg2", "4.0GiB", "5.0GiB"}},
				{Cells: []interface{}{"", "", ""}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Add LVMfakeclient error reactors
			if tt.args.lvmfunc != nil {
				tt.args.lvmfunc(tt.args.c)
			}
			// 2. Run the test & assert the result
			if head, row, err := GetVolumeGroups(tt.args.c, tt.args.vg); (err != nil) != tt.wantErr {
				t.Errorf("GetVolumeGroups() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				if !reflect.DeepEqual(row, tt.want) {
					t.Errorf("GetVolumeGroups() returned %v want = %v", row, tt.want)
				}
				if !reflect.DeepEqual(head, util.LVMvolgroupListColumnDefinitions) {
					t.Errorf("GetVolumeGroups() returned wrong headers = %v want = %v", head,
						util.LVMvolgroupListColumnDefinitions)
				}
			}
		})
	}
}

// lvnNodeNotFound makes fakelvmClientSet return error
func lvnNodeNotFound(c *client.K8sClient) {
	// NOTE: Set the VERB & Resource correctly & make it work for single resources
	c.LVMCS.LocalV1alpha1().(*fakelvm.FakeLocalV1alpha1).Fake.PrependReactor("*", "*", func(action k8stest.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("failed to list LVMVolumes")
	})
}

func TestDescribeLVMvg(t *testing.T) {
	type args struct {
		c       *client.K8sClient
		lvmFunc func(sClient *client.K8sClient)
		vg      string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"no LVM vgs exist",
			args{c: &client.K8sClient{Ns: "", LVMCS: fakelvmclient.NewSimpleClientset()}, lvmFunc: lvnNodeNotFound, vg: "cstor-pv1"},
			true,
		},
		{
			"one LVM node exist and asked for",
			args{c: &client.K8sClient{Ns: "lvm", LVMCS: fakelvmclient.NewSimpleClientset(&lvmNode1)}, vg: "node1"},
			false,
		},
		{
			"one ZFS node exist with differing namespace",
			args{c: &client.K8sClient{Ns: "zfs", LVMCS: fakelvmclient.NewSimpleClientset(&lvmNode1)}, vg: "node1"},
			false,
		},
		{
			"two ZFS node exist, none asked for",
			args{c: &client.K8sClient{Ns: "zfs", LVMCS: fakelvmclient.NewSimpleClientset(&lvmNode1, &lvmNode2)}, vg: ""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.lvmFunc != nil {
				tt.args.lvmFunc(tt.args.c)
			}
			if err := DescribeLVMvg(tt.args.c, tt.args.vg); (err != nil) != tt.wantErr {
				t.Errorf("DescribeLVMvg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
