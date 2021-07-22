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
	"time"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	fakecstor "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/api/v2/pkg/client/clientset/versioned/typed/cstor/v1/fake"
	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	corefake "k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
)

var cspi1 = cstorv1.CStorPoolInstance{
	TypeMeta: metav1.TypeMeta{Kind: "CStorPoolInstance", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "pool-1", Namespace: "openebs",
		Finalizers: []string{"cstorpoolcluster.openebs.io/finalizer", "openebs.io/pool-protection"},
		Labels: map[string]string{
			"kubernetes.io/hostname":        "node1",
			"openebs.io/cas-type":           "cstor",
			"openebs.io/cstor-pool-cluster": "cassandra-pool",
			"openebs.io/version":            "2.11"},
		// OwnerReference links to the CSPC
	},
	Spec: cstorv1.CStorPoolInstanceSpec{
		HostName:     "node1",
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
		PoolConfig:   cstorv1.PoolConfig{DataRaidGroupType: "stripe", WriteCacheGroupType: "", Compression: "off"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd-1", Capacity: 1234567, DevLink: "/dev/disk/by-id/abcd/def"}}}},
		WriteCacheRaidGroups: nil,
	},
	Status: cstorv1.CStorPoolInstanceStatus{
		Conditions: []cstorv1.CStorPoolInstanceCondition{{
			Type:               cstorv1.CSPIPoolLost,
			Status:             "True",
			LastUpdateTime:     metav1.Time{Time: time.Now()},
			LastTransitionTime: metav1.Time{Time: time.Now()},
			Reason:             "PoolLost",
			Message:            "failed to importcstor-xyzabcd",
		}},
		Phase: cstorv1.CStorPoolStatusOnline,
		Capacity: cstorv1.CStorPoolInstanceCapacity{
			Used:  resource.MustParse("18600Mi"),
			Free:  resource.MustParse("174Gi"),
			Total: resource.MustParse("192600Mi"),
			ZFS:   cstorv1.ZFSCapacityAttributes{},
		},
		ReadOnly: false, ProvisionedReplicas: 2, HealthyReplicas: 2,
	},
	VersionDetails: cstorv1.VersionDetails{Desired: "2.11",
		Status: cstorv1.VersionStatus{Current: "2.11", State: cstorv1.ReconcileComplete, LastUpdateTime: metav1.Time{Time: time.Now()}},
	},
}

var bd1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "BlockDevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd-1", Namespace: "openebs",
		Annotations: map[string]string{
			"internal.openebs.io/partition-uuid": "49473bca-97c3-f340-beaf-dae9b2ce99bc",
			"internal.openebs.io/uuid-scheme":    "legacy"}},
	Spec: v1alpha1.DeviceSpec{Capacity: v1alpha1.DeviceCapacity{
		Storage:            123456789,
		PhysicalSectorSize: 123456789,
		LogicalSectorSize:  123456789,
	}},
	Status: v1alpha1.DeviceStatus{},
}

var bd2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{
		Kind:       "BlockDevice",
		APIVersion: "openebs.io/v1alpha1",
	},
	ObjectMeta: metav1.ObjectMeta{Name: "bd-2", Namespace: "openebs"},
	Spec: v1alpha1.DeviceSpec{Capacity: v1alpha1.DeviceCapacity{
		Storage:            123456789,
		PhysicalSectorSize: 123456789,
		LogicalSectorSize:  123456789,
	},
		FileSystem: v1alpha1.FileSystemInfo{Type: "zfs_member", Mountpoint: "/home/kubernetes/volume-abcd"}},
	Status: v1alpha1.DeviceStatus{
		ClaimState: "Claimed",
		State:      "Active",
	},
}

var cvr1 = v1.CStorVolumeReplica{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1-rep-1",
		Labels:    map[string]string{cstortypes.CStorPoolInstanceNameLabelKey: "pool-1", "openebs.io/persistent-volume": "pv1"},
		Namespace: "openebs",
	},
	Status: v1.CStorVolumeReplicaStatus{
		Capacity: v1.CStorVolumeReplicaCapacityDetails{
			Total: "4Gi",
			Used:  "70Mi",
		},
		Phase: v1.CVRStatusOnline,
	},
}

var cvr2 = v1.CStorVolumeReplica{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "pvc-1-rep-2",
		Labels:    map[string]string{cstortypes.CStorPoolInstanceNameLabelKey: "pool-1", "openebs.io/persistent-volume": "pv1"},
		Namespace: "openebs",
	},
	Status: v1.CStorVolumeReplicaStatus{
		Capacity: v1.CStorVolumeReplicaCapacityDetails{
			Total: "40Gi",
			Used:  "70Mi",
		},
		Phase: v1.CVRStatusOnline,
	},
}
var pv1 = corev1.PersistentVolume{
	TypeMeta:   metav1.TypeMeta{Kind: "PersistentVolume", APIVersion: "core/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "pv1"},
	Spec:       corev1.PersistentVolumeSpec{ClaimRef: &corev1.ObjectReference{Name: "mongopv1"}},
	Status:     corev1.PersistentVolumeStatus{},
}

var cspi2 = cstorv1.CStorPoolInstance{
	TypeMeta: metav1.TypeMeta{Kind: "CStorPoolInstance", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "pool-2", Namespace: "openebs",
		Finalizers: []string{"cstorpoolcluster.openebs.io/finalizer", "openebs.io/pool-protection"},
		Labels: map[string]string{
			"kubernetes.io/hostname":        "node2",
			"openebs.io/cas-type":           "cstor",
			"openebs.io/cstor-pool-cluster": "cassandra-pool",
			"openebs.io/version":            "2.11"},
		// OwnerReference links to the CSPC
	},
	Spec: cstorv1.CStorPoolInstanceSpec{
		HostName:     "node2",
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
		PoolConfig:   cstorv1.PoolConfig{DataRaidGroupType: "stripe", WriteCacheGroupType: "", Compression: "off"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd2", Capacity: 1234567, DevLink: "/dev/disk/by-id/abcd/def"}}}},
		WriteCacheRaidGroups: nil,
	},
	Status: cstorv1.CStorPoolInstanceStatus{
		Conditions: []cstorv1.CStorPoolInstanceCondition{{
			Type:               cstorv1.CSPIPoolLost,
			Status:             "True",
			LastUpdateTime:     metav1.Time{Time: time.Now()},
			LastTransitionTime: metav1.Time{Time: time.Now()},
			Reason:             "PoolLost",
			Message:            "failed to importcstor-xyzabcd",
		}},
		Phase: cstorv1.CStorPoolStatusOnline,
		Capacity: cstorv1.CStorPoolInstanceCapacity{
			Used:  resource.MustParse("18600Mi"),
			Free:  resource.MustParse("174Gi"),
			Total: resource.MustParse("192600Mi"),
			ZFS:   cstorv1.ZFSCapacityAttributes{},
		},
		ReadOnly: false, ProvisionedReplicas: 2, HealthyReplicas: 2,
	},
	VersionDetails: cstorv1.VersionDetails{Desired: "2.11",
		Status: cstorv1.VersionStatus{Current: "2.11", State: cstorv1.ReconcileComplete, LastUpdateTime: metav1.Time{Time: time.Now()}},
	},
}

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
				{Cells: []interface{}{"pool-1", "node1", "174 GiB", "188 GiB", false, 2, 2, "ONLINE"}},
				{Cells: []interface{}{"pool-2", "node2", "174 GiB", "188 GiB", false, 2, 2, "ONLINE"}}},

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
			if err := GetCstorPools(tt.args.c, tt.args.poolName); (err != nil) != tt.wantErr {
				t.Errorf("DescribeCstorPool() error = %v, wantErr %v", err, tt.wantErr)
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
