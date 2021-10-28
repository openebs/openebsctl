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
	"github.com/stretchr/testify/assert"
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
			"cstor present, no suggested nodes present",
			args{
				c:     &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod), OpenebsCS: cstorfake.NewSimpleClientset()},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, nil,
			"", true,
		},
		{
			"cstor present, suggested nodes present, blockdevices absent",
			args{
				c:     &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1), OpenebsCS: cstorfake.NewSimpleClientset()},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, nil,
			"", true,
		},
		{
			"cstor present, suggested nodes present, blockdevices present but incompatible",
			args{
				c: &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1),
					OpenebsCS: cstorfake.NewSimpleClientset(&activeBDwEXT4, &inactiveBDwEXT4)},
				nodes: []string{"node1"}, devs: 1, poolType: ""}, nil,
			"", true,
		},
		{
			"cstor present, suggested nodes present, blockdevices present and compatible",
			args{
				c: &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1),
					OpenebsCS: cstorfake.NewSimpleClientset(&activeUnclaimedUnforattedBD)},
				nodes: []string{"node1"}, devs: 1, poolType: "stripe"}, &cspc1Struct, cspc1, false,
		},
		{
			"all good config, CSTOR_NAMESPACE is correctly identified each time",
			args{
				c: &client.K8sClient{Ns: "randomNamespaceWillGetReplaced", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1),
					OpenebsCS: cstorfake.NewSimpleClientset(&activeUnclaimedUnforattedBD)},
				nodes: []string{"node1"}, devs: 1, poolType: "stripe"}, &cspc1Struct, cspc1, false,
		},
		{
			"all good config, 2 disk stripe pool for 3 nodes",
			args{
				c: &client.K8sClient{Ns: "", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1, &node2, &node3),
					OpenebsCS: cstorfake.NewSimpleClientset(&goodBD1N1, &goodBD1N2, &goodBD1N3, &goodBD2N1, &goodBD2N2, &goodBD2N3)},
				nodes: []string{"node1", "node2", "node3"}, devs: 2, poolType: "stripe"}, &threeNodeTwoDevCSPC, StripeThreeNodeTwoDev, false,
		},
		{
			"good config, no BDs",
			args{
				c: &client.K8sClient{Ns: "randomNamespaceWillGetReplaced", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1),
					OpenebsCS: cstorfake.NewSimpleClientset(&inactiveBDwEXT4)},
				nodes: []string{"node1"}, devs: 5, poolType: "stripe"}, nil, "", true,
		},
		{
			"all good mirror pool gets provisioned, 2 nodes of same size on 3 nodes",
			args{
				c: &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1, &node2, &node3),
					OpenebsCS: cstorfake.NewSimpleClientset(&goodBD1N1, &goodBD2N1, &goodBD1N2, &goodBD2N2, &goodBD1N3, &goodBD2N3)},
				nodes: []string{"node1", "node2", "node3"}, devs: 2, poolType: "mirror"}, &mirrorCSPC, mirrorCSPCstr, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.args.GB,
			got, got1, err := CSPC(tt.args.c, tt.args.nodes, tt.args.devs, tt.args.poolType)
			assert.YAMLEq(t, tt.str, got1, "stringified YAML is not the same as expected")
			assert.Equal(t, got, tt.want, "struct is not same")
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
