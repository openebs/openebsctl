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

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/zfs-localpv/pkg/generated/clientset/internalclientset/fake"
	fakezfs "github.com/openebs/zfs-localpv/pkg/generated/clientset/internalclientset/typed/zfs/v1/fake"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
)

func TestGetZFSLocalPVs(t *testing.T) {
	type args struct {
		c           *client.K8sClient
		zfsReactors func(*client.K8sClient)
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
			name: "no zfs volumes present",
			args: args{
				c: &client.K8sClient{
					Ns:        "random-namespace",
					ZFCS:      fake.NewSimpleClientset(),
					K8sCS:     k8sfake.NewSimpleClientset(),
					OpenebsCS: nil,
				},
				pvList:      &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, pv2, pv3}},
				zfsReactors: zfsVolNotExists,
				openebsNS:   "zfslocalpv",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only one zfs volume present",
			args: args{
				c: &client.K8sClient{
					Ns:    "zfslocalpv",
					K8sCS: k8sfake.NewSimpleClientset(&localpvzfsCSICtrlSTS),
					ZFCS:  fake.NewSimpleClientset(&zfsVol1),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, zfsPV1}},
				openebsNS: "zfslocalpv",
			},
			wantErr: false,
			want: []metav1.TableRow{
				{
					Cells: []interface{}{"zfslocalpv", "pvc-1", "Ready", "1.9.0", "4.0GiB", "zfs-sc-1", corev1.VolumeBound, corev1.ReadWriteOnce, "node1"},
				},
			},
		},
		{
			name: "only one zfs volume present with zfsvol absent",
			args: args{
				c: &client.K8sClient{
					Ns:    "zfslocalpv",
					K8sCS: k8sfake.NewSimpleClientset(&localpvzfsCSICtrlSTS),
					ZFCS:  fake.NewSimpleClientset(),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{zfsPV1}},
				openebsNS: "zfslocalpv",
			},
			wantErr: false,
			want:    nil,
		},
		{
			name: "only one zfs volume present, namespace conflicts",
			args: args{
				c: &client.K8sClient{
					Ns:    "jiva",
					K8sCS: k8sfake.NewSimpleClientset(&localpvzfsCSICtrlSTS),
					ZFCS:  fake.NewSimpleClientset(&zfsVol1),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, zfsPV1}},
				openebsNS: "zfslocalpvXYZ",
			},
			wantErr: false,
			want:    nil,
		},
		{
			name: "controller sts not present",
			args: args{
				c: &client.K8sClient{
					Ns:    "jiva",
					K8sCS: k8sfake.NewSimpleClientset(),
					ZFCS:  fake.NewSimpleClientset(&zfsVol1),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{jivaPV1, zfsPV1}},
				openebsNS: "zfslocalpv",
			},
			wantErr: false,
			want: []metav1.TableRow{
				{
					Cells: []interface{}{"zfslocalpv", "pvc-1", "Ready", "N/A", "4.0GiB", "zfs-sc-1", corev1.VolumeBound, corev1.ReadWriteOnce, "node1"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Before func
			if tt.args.zfsReactors != nil {
				tt.args.zfsReactors(tt.args.c)
			}
			got, err := GetZFSLocalPVs(tt.args.c, tt.args.pvList, tt.args.openebsNS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetZFSLocalPVs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetZFSLocalPVs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDescribeZFSLocalPVs(t *testing.T) {
	type args struct {
		c       *client.K8sClient
		zfsfunc func(sClient *client.K8sClient)
		vol     *corev1.PersistentVolume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no zfs volume present",
			args{c: &client.K8sClient{Ns: "zfs", K8sCS: k8sfake.NewSimpleClientset(), ZFCS: fake.NewSimpleClientset()},
				vol:     nil,
				zfsfunc: zfsVolNotExists,
			},
			true,
		},
		{"one zfs volume present and asked for and zfs-controller absent",
			args{c: &client.K8sClient{Ns: "zfslocalpv", ZFCS: fake.NewSimpleClientset(&zfsVol1), K8sCS: k8sfake.NewSimpleClientset()},
				vol: &zfsPV1},
			false,
		},
		{"one zfs volume present and asked for and zfs-controller present",
			args{c: &client.K8sClient{Ns: "zfslocalpv",
				K8sCS: k8sfake.NewSimpleClientset(&localpvzfsCSICtrlSTS),
				ZFCS:  fake.NewSimpleClientset(&zfsVol1)},
				vol: &zfsPV1},
			false,
		},
		{"one zfs volume present and asked for but namespace wrong",
			args{c: &client.K8sClient{Ns: "zfslocalpv1", ZFCS: fake.NewSimpleClientset(&zfsVol1), K8sCS: k8sfake.NewSimpleClientset()},
				vol: &zfsPV1, zfsfunc: zfsVolNotExists},
			false,
		},
		{"one zfs volume present and some other volume asked for",
			args{c: &client.K8sClient{Ns: "zfs", K8sCS: k8sfake.NewSimpleClientset(&zfsPV1), ZFCS: fake.NewSimpleClientset(&zfsVol1)},
				vol:     &cstorPV2,
				zfsfunc: zfsVolNotExists},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.zfsfunc != nil {
				tt.args.zfsfunc(tt.args.c)
			}
			if err := DescribeZFSLocalPVs(tt.args.c, tt.args.vol); (err != nil) != tt.wantErr {
				t.Errorf("DescribeZFSLocalPVs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// zfsVolNotExists makes fakezfsClientSet return error
func zfsVolNotExists(c *client.K8sClient) {
	// NOTE: Set the VERB & Resource correctly & make it work for single resources
	c.ZFCS.ZfsV1().(*fakezfs.FakeZfsV1).Fake.PrependReactor("*", "*", func(action k8stest.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("failed to list ZFSVolumes")
	})
}
