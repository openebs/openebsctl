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
	"testing"

	"github.com/openebs/openebsctl/pkg/client"
	fakezfsclient "github.com/openebs/zfs-localpv/pkg/generated/clientset/internalclientset/fake"
	fakezfs "github.com/openebs/zfs-localpv/pkg/generated/clientset/internalclientset/typed/zfs/v1/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stest "k8s.io/client-go/testing"
)

func TestGetZFSPools(t *testing.T) {
	type args struct {
		c        *client.K8sClient
		zfsnodes []string
	}
	tests := []struct {
		name    string
		args    args
		zfsfunc func(sClient *client.K8sClient)
		want    []metav1.TableRow
		wantErr bool
	}{
		{
			"no zfs pools present",
			args{c: &client.K8sClient{Ns: "random", ZFCS: fakezfsclient.NewSimpleClientset()}, zfsnodes: nil},
			zfsNodeNotFound,
			nil,
			true,
		},
		{
			"zfs pools present",
			args{c: &client.K8sClient{Ns: "random", ZFCS: fakezfsclient.NewSimpleClientset(&zfsNode1, &zfsNode2)},
				zfsnodes: nil},
			nil,
			[]metav1.TableRow{
				{Cells: []interface{}{"node1", ""}},
				{Cells: []interface{}{firstElemPrefix + "zfs-pool1", "32 GiB"}},
				{Cells: []interface{}{"", ""}},
				{Cells: []interface{}{"node2", ""}},
				{Cells: []interface{}{firstElemPrefix + "zfs-pool2", "32 GiB"}},
				{Cells: []interface{}{lastElemPrefix + "zfs-pool3", "32 GiB"}},
				{Cells: []interface{}{"", ""}},
			},
			false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Add LVMfakeclient error reactors
			if tt.zfsfunc != nil {
				tt.zfsfunc(tt.args.c)
			}
			if err := GetZFSPools(tt.args.c, tt.args.zfsnodes); (err != nil) != tt.wantErr {
				t.Errorf("GetZFSPools() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// lvnNodeNotFound makes fakelvmClientSet return error
func zfsNodeNotFound(c *client.K8sClient) {
	// NOTE: Set the VERB & Resource correctly & make it work for single resources
	c.ZFCS.ZfsV1().(*fakezfs.FakeZfsV1).Fake.PrependReactor("*", "*", func(action k8stest.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("failed to list ZFS nodes")
	})
}
