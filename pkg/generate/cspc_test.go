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

package generate

import (
	"testing"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	cstorfake "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/resource"
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
				// Stripe pools can have only one RaidGroup per instance, i.e.
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
		{
			"all good raidz pool gets provisioned, 3 nodes of same size on 2 nodes",
			args{
				c: &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1, &node2),
					OpenebsCS: cstorfake.NewSimpleClientset(&goodBD1N1, &goodBD2N1, &goodBD3N1, &goodBD1N2, &goodBD2N2, &goodBD3N2)},
				nodes: []string{"node1", "node2"}, devs: 3, poolType: "raidz"}, &raidzCSPCThreeBDTwoNode, raidzCSPCstr, false,
		},
		{
			"all good raidz2 pool does not gets provisioned, insufficient BDs 3 nodes of same size on 2 nodes",
			args{
				c: &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1, &node2),
					OpenebsCS: cstorfake.NewSimpleClientset(&goodBD1N1,
						&goodBD2N1, &goodBD3N1, &goodBD4N1, &goodBD1N2,
						&goodBD2N2, &goodBD3N2, &goodBD4N2)},
				nodes: []string{"node1", "node2"}, devs: 3, poolType: "raidz2"}, nil, "", true,
		},
		{
			"raidz2 pool provisioned, 2 nodes, 6 BDs",
			args{
				c: &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1, &node2),
					OpenebsCS: cstorfake.NewSimpleClientset(&goodBD1N1,
						&goodBD2N1, &goodBD3N1, &goodBD4N1, &goodBD5N1,
						&goodBD6N1, &goodBD1N2, &goodBD2N2, &goodBD3N2,
						&goodBD4N2, &goodBD5N2, &goodBD6N2)},
				nodes: []string{"node1", "node2"}, devs: 6, poolType: "raidz2"}, &raidz2CSPCSixBDTwoNode, raidz2CSPCstr, false,
		},
		{
			"raidz2 not provisioned, requires 2 more BDs",
			args{
				c: &client.K8sClient{Ns: "openebs", K8sCS: fake.NewSimpleClientset(&cstorCSIpod, &node1, &node2),
					OpenebsCS: cstorfake.NewSimpleClientset(&goodBD1N1,
						&goodBD2N1, &goodBD3N1, &goodBD4N1, &goodBD1N2,
						&goodBD2N2, &goodBD3N2, &goodBD4N2)},
				nodes: []string{"node1", "node2"}, devs: 4, poolType: "raidz2"}, nil, "", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.args.GB,
			got, got1, err := cspc(tt.args.c, tt.args.nodes, tt.args.devs, tt.args.poolType, resource.MustParse("1Gi"))
			if (err != nil) != tt.wantErr {
				t.Errorf("cspc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.YAMLEq(t, tt.str, got1, "stringified YAML is not the same as expected")
			assert.Exactlyf(t, got, tt.want, "struct is not same")
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
		{name: "invalid pools", poolNames: []string{"striped", "mirrored", "raid-z", "raid-z2", "lvm", "raidz1", "raidz0"}, want: false},
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

func Test_makePools(t *testing.T) {
	type args struct {
		poolType string
		nDevices int
		bd       map[string][]v1alpha1.BlockDevice
		nodes    []string
		hosts    []string
		minSize  resource.Quantity
	}
	tests := []struct {
		name    string
		args    args
		want    *[]cstorv1.PoolSpec
		wantErr bool
	}{
		{"stripe, three node, two disks", args{"stripe", 2,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1},
				"node2": {goodBD1N2, goodBD2N2}, "node3": {goodBD1N3, goodBD2N3}},
			[]string{"node1", "node2", "node3"}, []string{"node1", "node2", "node3"}, resource.MustParse("1Gi")}, &threeNodeTwoDevCSPC.Spec.Pools, false},
		{"stripe, three node, two disks, one node lacking disks", args{"stripe", 2,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1},
				"node2": {goodBD1N2, goodBD2N2}},
			[]string{"node1", "node2", "node3"}, []string{"node1", "node2", "node3"}, resource.MustParse("1Gi")}, nil, true},
		{"stripe, three node, two disks, one node lacking required disks", args{"stripe", 2,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1},
				"node2": {goodBD1N2}, "node3": {goodBD1N3, goodBD2N2}}, []string{"node1", "node2", "node3"},
			[]string{"node1", "node2", "node3"}, resource.MustParse("1Gi")}, nil, true},
		{"raidz, three node, three disks but only two disks present in node3", args{"raidz", 3,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1, goodBD3N1},
				"node2": {goodBD1N2, goodBD2N2, goodBD3N2}, "node3": {goodBD1N3, goodBD2N3}},
			[]string{"node1", "node2", "node3"}, []string{"node1", "node2", "node3"}, resource.MustParse("1Gi")}, nil, true},
		{"unknown pool, three node, two disks", args{"randompoolwhichmakesnosense", 2,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1},
				"node2": {goodBD1N2, goodBD2N2}, "node3": {goodBD1N3, goodBD2N3}},
			[]string{"node1", "node2", "node3"}, []string{"node1", "node2", "node3"}, resource.MustParse("1Gi")}, nil, true},
		{"mirror, three node, two disks", args{"mirror", 2,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1},
				"node2": {goodBD1N2, goodBD2N2}, "node3": {goodBD1N3, goodBD2N3}},
			[]string{"node1", "node2", "node3"}, []string{"node1", "node2", "node3"}, resource.MustParse("1Gi")}, &mirrorCSPC.Spec.Pools, false},
		{"mirror, two node, four disks", args{"mirror", 4,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1, goodBD3N1, goodBD4N1},
				"node2": {goodBD1N2, goodBD2N2, goodBD3N2, goodBD4N2}, "node3": {goodBD1N3, goodBD2N3}},
			// in the above example, presence of node3 BDs don't matter
			[]string{"node1", "node2"}, []string{"node1", "node2"}, resource.MustParse("1Gi")}, &mirrorCSPCFourBDs.Spec.Pools, false},
		{"mirror, three node, one disk", args{"mirror", 1,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1},
				"node2": {goodBD1N2, goodBD2N2}, "node3": {goodBD1N3, goodBD2N3}},
			// one cannot create a mirror pool with just one disk per node
			[]string{"node1", "node2", "node3"}, []string{"node1", "node2", "node3"}, resource.MustParse("1Gi")}, nil, true},
		{"raidz, two node, three disk", args{"raidz", 3,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1, goodBD3N1}, "node2": {goodBD1N2, goodBD2N2, goodBD3N2}},
			[]string{"node1", "node2"}, []string{"node1", "node2"}, resource.MustParse("1Gi")}, &raidzCSPCThreeBDTwoNode.Spec.Pools, false},
		{"raidz2, two node, three disk", args{"raidz2", 6,
			map[string][]v1alpha1.BlockDevice{"node1": {goodBD1N1, goodBD2N1, goodBD3N1, goodBD4N1, goodBD5N1, goodBD6N1}, "node2": {goodBD1N2, goodBD2N2, goodBD3N2, goodBD4N2, goodBD5N2, goodBD6N2}},
			[]string{"node1", "node2"}, []string{"node1", "node2"}, resource.MustParse("1Gi")}, &raidz2CSPCSixBDTwoNode.Spec.Pools, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makePools(tt.args.poolType, tt.args.nDevices, tt.args.bd, tt.args.nodes, tt.args.hosts, tt.args.minSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("makePools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "pool specs differ for %s", tt.name)
		})
	}
}
