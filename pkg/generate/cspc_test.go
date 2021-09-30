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

package generate

import (
	"reflect"
	"testing"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	cstorfake "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/openebsctl/pkg/client"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCSPC(t *testing.T) {
	type args struct {
		c        *client.K8sClient
		nodes    []string
		devs     int
		GB       int
		poolType string
	}
	tests := []struct {
		name    string
		args    args
		want    *cstorv1.CStorPoolCluster
		str     string
		wantErr bool
	}{
		{
			"no cstor installation present",
			args{
				c:     &client.K8sClient{Ns: "", K8sCS: fake.NewSimpleClientset(), OpenebsCS: cstorfake.NewSimpleClientset()},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, nil,
			"", true,
		},
		{
			"stripe kind CSPC with one block-device",
			args{
				c:     &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(), OpenebsCS: cstorfake.NewSimpleClientset()},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, &cstorv1.CStorPoolCluster{},
			"", false,
		},
		{
			"stripe kind CSPC with two block-device on different nodes",
			args{
				c:     &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(), OpenebsCS: cstorfake.NewSimpleClientset()},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, &cstorv1.CStorPoolCluster{},
			"", false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.args.GB,
			got, got1, err := CSPC(tt.args.c, tt.args.nodes, tt.args.devs, tt.args.poolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CSPC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CSPC() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.str {
				t.Errorf("CSPC() got1 = %v, want %v", got1, tt.str)
			}
		})
	}
}

func Test_isPoolTypeValid(t *testing.T) {
	tests := []struct {
		name      string
		poolNames []string
		want      bool
	}{
		{name: "valid pools", poolNames: []string{"stripe", "mirror", "raidz", "raidz2"}, want: true},
		{name: "invalid pools", poolNames: []string{"striped", "mirrored", "raid-z", "raid-z2", "lvm"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, poolType := range tt.poolNames {
				if got := isPoolTypeValid(poolType); got != tt.want {
					t.Errorf("isPoolTypeValid() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
