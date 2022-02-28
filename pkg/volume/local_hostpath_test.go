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

package volume

import (
	"reflect"
	"testing"

	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestGetLocalHostpath(t *testing.T) {
	type args struct {
		c         *client.K8sClient
		pvList    *corev1.PersistentVolumeList
		openebsNS string
	}
	tests := []struct {
		name    string
		args    args
		want    []metav1.TableRow
		wantErr bool
	}{
		{
			name: "no local hostpath volumes present",
			args: args{
				c: &client.K8sClient{
					Ns:        "random-namespace",
					K8sCS:     k8sfake.NewSimpleClientset(),
					OpenebsCS: nil,
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, pv2}},
				openebsNS: "openebs",
			},
			want:    []metav1.TableRow{},
			wantErr: false,
		},
		{
			name: "only one local hostpath volume present",
			args: args{
				c: &client.K8sClient{
					Ns:    "lvmlocalpv",
					K8sCS: k8sfake.NewSimpleClientset(&localpvHostpathDpl1),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, localHostpathPv1}},
				openebsNS: "localhostpath",
			},
			wantErr: false,
			want: []metav1.TableRow{
				{
					Cells: []interface{}{"openebs", "pvc-1", "", "1.9.0", &fourGigiByte, "pvc-1-local", corev1.VolumeBound, corev1.ReadWriteOnce, "node1"},
				},
			},
		},
		{
			name: "two local hostpath volume present",
			args: args{
				c: &client.K8sClient{
					Ns:    "lvmlocalpv",
					K8sCS: k8sfake.NewSimpleClientset(&localpvHostpathDpl1, &localpvHostpathDpl2),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, localHostpathPv1}},
				openebsNS: "localhostpath",
			},
			wantErr: false,
			want: []metav1.TableRow{
				{
					Cells: []interface{}{"", "pvc-1", "", "N/A", &fourGigiByte, "pvc-1-local", corev1.VolumeBound, corev1.ReadWriteOnce, "node1"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Call the code under test
			got, err := GetLocalHostpath(tt.args.c, tt.args.pvList, tt.args.openebsNS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocalHostpath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 2. Test conditions of pass/fail & display
			gotLen := len(got)
			expectedLen := len(tt.want)
			if gotLen != expectedLen {
				t.Errorf("GetLocalHostpath() returned %d elements, wanted %d elements", gotLen, expectedLen)
			}
			for i, gotLine := range got {
				if len(gotLine.Cells) != len(tt.want[i].Cells) {
					t.Errorf("Line#%d in output had %d elements, wanted %d elements", i+1, len(gotLine.Cells), len(tt.want[i].Cells))
				}
				if !reflect.DeepEqual(tt.want[i].Cells, gotLine.Cells) {
					t.Errorf("GetLocalHostpath() line#%d got = %v, want %v", i+1, got, tt.want)
				}
			}
		})
	}
}

func TestDescribeLocalHostpathVolume(t *testing.T) {
	type args struct {
		c   *client.K8sClient
		vol *corev1.PersistentVolume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"one local hostpath volume present",
			args{c: &client.K8sClient{Ns: "openebsc", K8sCS: k8sfake.NewSimpleClientset()},
				vol: &localHostpathPv1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Call the code under test and check condition.
			if err := DescribeLocalHostpathVolume(tt.args.c, tt.args.vol); (err != nil) != tt.wantErr {
				t.Errorf("DescribeLocalHostpathVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
