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

	fakecstor "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/api/v2/pkg/client/clientset/versioned/typed/cstor/v1/fake"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	corefake "k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
)

func TestGetCstorPool(t *testing.T) {
	type args struct {
		c        *client.K8sClient
		poolName []string
	}
	tests := []struct {
		name      string
		args      args
		cstorfunc func(sClient *client.K8sClient)
		want      []metav1.TableRow
		wantErr   bool
	}{
		{
			"no cstor pool found",
			args{c: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset()},
				poolName: nil},
			cspiNotFound,
			nil,
			true,
		},
		{
			"two cstor pool found",
			args{c: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspi1, &cspi2)},
				poolName: nil,
			},
			nil,
			[]metav1.TableRow{
				{Cells: []interface{}{"pool-1", "node1", "174.0GiB", "188.1GiB", false, int32(2), int32(2), "ONLINE"}},
				{Cells: []interface{}{"pool-2", "node2", "174.0GiB", "188.1GiB", false, int32(2), int32(2), "ONLINE"}}},

			false,
		},
		{
			"no pool-3 cstor pool found",
			args{c: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspi1, &cspi2)},
				poolName: []string{"pool-3"},
			},
			cspiNotFound,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cstorfunc != nil {
				tt.cstorfunc(tt.args.c)
			}
			if head, row, err := GetCstorPools(tt.args.c, tt.args.poolName); (err != nil) != tt.wantErr {
				t.Errorf("GetCstorPool() error = %v, wantErr %v", err, tt.wantErr)
			} else if err == nil {
				if len(row) != len(tt.want) {
					t.Errorf("GetCstorPool() returned %d rows, wanted %d elements", len(row), len(tt.want))
				}
				for i, cspi := range row {
					if !reflect.DeepEqual(cspi.Cells[0:8], tt.want[i].Cells) {
						t.Errorf("GetCstorPool() returned %v want = %v", row, tt.want)
					}
				}
				if !reflect.DeepEqual(head, util.CstorPoolListColumnDefinations) {
					t.Errorf("GetCstorPools() returned wrong headers = %v want = %v", head,
						util.CstorPoolListColumnDefinations)
				}
			}
			// TODO: Check all but the last item of want
		})
	}
}

func TestDescribeCstorPool(t *testing.T) {
	type args struct {
		c        *client.K8sClient
		poolName string
	}
	tests := []struct {
		name      string
		args      args
		cstorfunc func(sClient *client.K8sClient)
		wantErr   bool
	}{
		{"no cstor pool exist",
			args{c: &client.K8sClient{Ns: "cstor", OpenebsCS: fakecstor.NewSimpleClientset()},
				poolName: ""},
			// a GET on resource which don't exist, returns an error automatically
			nil,
			true,
		},
		{"cspi-3 does not exist",
			args{c: &client.K8sClient{Ns: "cstor", OpenebsCS: fakecstor.NewSimpleClientset()},
				poolName: "cspi-3"},
			nil,
			true,
		},
		{"cspi-1 exists but Namespace mismatched",
			args{c: &client.K8sClient{Ns: "fake", OpenebsCS: fakecstor.NewSimpleClientset(&cspi1)},
				poolName: "cspi-1"},
			nil,
			true,
		},
		{
			"cspi-1 exists and namespace matches but no BD",
			args{c: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspi1)},
				poolName: "pool-1"},
			nil,
			false,
		},
		{
			"cspi-1 exists and BD exists",
			args{c: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspi1, &bd1)},
				poolName: "pool-1"},
			nil,
			false,
		},
		{
			"cspi-1 exists, BD & CVR exists",
			args{c: &client.K8sClient{Ns: "openebs", OpenebsCS: fakecstor.NewSimpleClientset(&cspi1, &bd1, &bd2,
				&cvr1, &cvr2), K8sCS: corefake.NewSimpleClientset(&pv1)},
				poolName: "pool-1"},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DescribeCstorPool(tt.args.c, tt.args.poolName); (err != nil) != tt.wantErr {
				t.Errorf("DescribeCstorPool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func cspiNotFound(c *client.K8sClient) {
	// NOTE: Set the VERB & Resource correctly & make it work for single resources
	c.OpenebsCS.CstorV1().(*fake.FakeCstorV1).Fake.PrependReactor("*", "*", func(action k8stest.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("failed to list CSPI")
	})
}
