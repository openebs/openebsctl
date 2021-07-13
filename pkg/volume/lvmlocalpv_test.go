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
package volume

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	lvm "github.com/openebs/lvm-localpv/pkg/apis/openebs.io/lvm/v1alpha1"
	"github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset/fake"
	fake2 "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset/typed/lvm/v1alpha1/fake"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
)

var lvmVol1 = lvm.LVMVolume{
	TypeMeta: metav1.TypeMeta{
		Kind:       "LVMVolume",
		APIVersion: "lvm.openebs.io/v1alpha1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:              "pvc-1",
		Namespace:         "lvmlocalpv",
		CreationTimestamp: metav1.Time{Time: time.Now()},
		Labels:            map[string]string{},
		Annotations:       map[string]string{},
		OwnerReferences:   nil,
		Finalizers:        nil,
	},
	Spec: lvm.VolumeInfo{
		OwnerNodeID:   "node1",
		VolGroup:      "lvmpv",
		VgPattern:     "vg1*",
		Capacity:      "4Gi",
		Shared:        "NotShared",
		ThinProvision: "No",
	},
	Status: lvm.VolStatus{
		State: "Ready",
		Error: nil,
	},
}

var lvmPV1 = corev1.PersistentVolume{
	TypeMeta: metav1.TypeMeta{
		Kind:       "PersistentVolume",
		APIVersion: "core/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:        "pvc-1",
		Labels:      map[string]string{},
		Annotations: map[string]string{},
	},
	Spec: corev1.PersistentVolumeSpec{
		// 4GiB
		Capacity:                      corev1.ResourceList{corev1.ResourceStorage: fourGigiByte},
		PersistentVolumeSource:        corev1.PersistentVolumeSource{CSI: &corev1.CSIPersistentVolumeSource{Driver: util.LocalPVLVMCSIDriver}},
		AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
		ClaimRef:                      nil,
		PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimDelete,
		StorageClassName:              "lvm-sc-1",
		VolumeMode:                    &blockFS,
		NodeAffinity: &corev1.VolumeNodeAffinity{
			Required: &corev1.NodeSelector{NodeSelectorTerms: []corev1.NodeSelectorTerm{
				{MatchExpressions: []corev1.NodeSelectorRequirement{
					{Key: "kubernetes.io/hostname", Operator: corev1.NodeSelectorOpIn, Values: []string{"node2"}},
				}},
			}},
		},
	},
	Status: corev1.PersistentVolumeStatus{
		Phase:   corev1.VolumeBound,
		Message: "Storage class not found",
		Reason:  "K8s API was down",
	},
}

var localpvCSICtrlSTS = appsv1.StatefulSet{
	TypeMeta: metav1.TypeMeta{
		Kind:       "StatefulSet",
		APIVersion: "apps/v1",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name:      "fake-LVM-CSI",
		Namespace: "lvm",
		Labels: map[string]string{
			"openebs.io/version":        "1.9.0",
			"openebs.io/component-name": "openebs-lvm-controller"},
	},
}

func TestGetLVMLocalPV(t *testing.T) {
	type args struct {
		c           *client.K8sClient
		lvmReactors func(*client.K8sClient)
		pvList      *corev1.PersistentVolumeList
		openebsNS   string
	}
	tests := []struct {
		name    string
		args    args
		want    []metav1.TableRow
		wantErr bool
	}{
		{
			name: "no lvm volumes present",
			args: args{
				c: &client.K8sClient{
					Ns:        "random-namespace",
					LVMCS:     fake.NewSimpleClientset(),
					K8sCS:     k8sfake.NewSimpleClientset(),
					OpenebsCS: nil,
				},
				pvList:      &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, pv2, pv3}},
				lvmReactors: lvmVolNotExists,
				openebsNS:   "openebs",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only one lvm volume present",
			args: args{
				c: &client.K8sClient{
					Ns:    "lvmlocalpv",
					K8sCS: k8sfake.NewSimpleClientset(&localpvCSICtrlSTS),
					LVMCS: fake.NewSimpleClientset(&lvmVol1),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, lvmPV1}},
				openebsNS: "lvmlocalpv",
			},
			wantErr: false,
			want: []metav1.TableRow{
				{
					Cells: []interface{}{"lvmlocalpv", "pvc-1", "Ready", "1.9.0", "4Gi", "lvm-sc-1", corev1.VolumeBound, corev1.ReadWriteOnce, "node1"},
				},
			},
		},
		{
			name: "only one lvm volume present, namespace conflicts",
			args: args{
				c: &client.K8sClient{
					Ns:    "jiva",
					K8sCS: k8sfake.NewSimpleClientset(&localpvCSICtrlSTS),
					LVMCS: fake.NewSimpleClientset(&lvmVol1),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, lvmPV1}},
				openebsNS: "lvmlocalpv",
			},
			wantErr: false,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Before func
			if tt.args.lvmReactors != nil {
				tt.args.lvmReactors(tt.args.c)
			}
			// 2. Call the code under test
			got, err := GetLVMLocalPV(tt.args.c, tt.args.pvList, tt.args.openebsNS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLVMLocalPV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 3. Test for TC pass/fail & display
			gotLen := len(got)
			expectedLen := len(tt.want)
			if gotLen != expectedLen {
				t.Errorf("GetLVMLocalPV() returned %d elements, wanted %d elements", gotLen, expectedLen)
			}
			for i, gotLine := range got {
				if len(gotLine.Cells) != len(tt.want[i].Cells) {
					t.Errorf("Line#%d in output had %d elements, wanted %d elements", (i + 1), len(gotLine.Cells), len(tt.want[i].Cells))
				}
				if !reflect.DeepEqual(tt.want[i].Cells, gotLine.Cells) {
					t.Errorf("GetLVMLocalPV() line#%d got = %v, want %v", i+1, got, tt.want)
				}
			}
		})
	}
}

// lvmVolNotExists makes fakelvmClientSet return error
func lvmVolNotExists(c *client.K8sClient) {
	// NOTE: Set the VERB & Resource correctly & make it work for single resources
	c.LVMCS.LocalV1alpha1().(*fake2.FakeLocalV1alpha1).Fake.PrependReactor("*", "*", func(action k8stest.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("failed to list LVMVolumes")
	})
}
