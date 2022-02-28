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

package util

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestCheckForVol(t *testing.T) {
	type args struct {
		name string
		vols map[string]*Volume
	}
	tests := []struct {
		name string
		args args
		want *Volume
	}{
		{
			"volume_attached_to_storage_class",
			args{name: "cstor_volume", vols: map[string]*Volume{"cstor_volume": {CSIVolumeAttachmentName: "volume_one"}}},
			&Volume{CSIVolumeAttachmentName: "volume_one"},
		},
		{
			"volume_not_attached_to_storage_class",
			args{name: "cstor_volume", vols: map[string]*Volume{"cstor_volume_two": {CSIVolumeAttachmentName: "volume_one"}}},
			&Volume{StorageClass: NotAvailable, Node: NotAvailable, AttachementStatus: NotAvailable, AccessMode: NotAvailable},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckForVol(tt.args.name, tt.args.vols); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckForVol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAccessModeToString(t *testing.T) {
	type args struct {
		accessModeArray []corev1.PersistentVolumeAccessMode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Valid Values",
			args{[]corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce, corev1.ReadOnlyMany}},
			"ReadWriteOnce ReadOnlyMany ",
		},
		{
			"In valid Values",
			args{[]corev1.PersistentVolumeAccessMode{}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AccessModeToString(tt.args.accessModeArray); got != tt.want {
				t.Errorf("AccessModeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
