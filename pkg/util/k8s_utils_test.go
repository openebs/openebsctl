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
	"testing"

	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/storage/v1"
)

var (
	inAnnotationPV, inLabelsPV = cstorPV1, cstorPV1
)

func TestGetCasType(t *testing.T) {
	inAnnotationPV.Annotations = map[string]string{"openebs.io/cas-type": "cstor"}
	inLabelsPV.Labels = map[string]string{"openebs.io/cas-type": "cstor"}
	type args struct {
		v1PV *corev1.PersistentVolume
		v1SC *v1.StorageClass
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"PV and SC both present",
			args{
				v1PV: &cstorPV1,
				v1SC: &cstorSC,
			},
			"cstor",
		},
		{
			"PV present and SC absent",
			args{
				v1PV: &cstorPV1,
				v1SC: nil,
			},
			"cstor",
		},
		{
			"PV present and SC absent",
			args{
				v1PV: &inLabelsPV,
				v1SC: nil,
			},
			"cstor",
		},
		{
			"PV present and SC absent",
			args{
				v1PV: &inAnnotationPV,
				v1SC: nil,
			},
			"cstor",
		},
		{
			"PV absent and SC present",
			args{
				v1PV: nil,
				v1SC: &cstorSC,
			},
			"cstor",
		},
		{
			"PV absent and SC present",
			args{
				v1PV: nil,
				v1SC: &jivaSC,
			},
			"jiva",
		},
		{
			"Both PV and SC absent",
			args{
				v1PV: nil,
				v1SC: nil,
			},
			"unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCasType(tt.args.v1PV, tt.args.v1SC); got != tt.want {
				t.Errorf("GetCasType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCasTypeFromPV(t *testing.T) {
	inAnnotationPV.Annotations = map[string]string{"openebs.io/cas-type": "cstor"}
	inLabelsPV.Labels = map[string]string{"openebs.io/cas-type": "cstor"}
	type args struct {
		v1PV *corev1.PersistentVolume
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"From volume attributes",
			args{
				v1PV: &cstorPV1,
			},
			"cstor",
		},
		{
			"From labels",
			args{
				v1PV: &inLabelsPV,
			},
			"cstor",
		},
		{
			"From annotations",
			args{
				v1PV: &inAnnotationPV,
			},
			"cstor",
		},
		{
			"nil pv",
			args{
				v1PV: nil,
			},
			"unknown",
		},
		{
			"zfs pv, from CSI driver",
			args{v1PV: &zfspv},
			"localpv-zfs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCasTypeFromPV(tt.args.v1PV); got != tt.want {
				t.Errorf("GetCasTypeFromPV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCasTypeFromSC(t *testing.T) {
	type args struct {
		v1SC *v1.StorageClass
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"From provisioner",
			args{v1SC: &cstorSC},
			"cstor",
		},
		{
			"From parameters",
			args{v1SC: &jivaSC},
			"jiva",
		},
		{
			"SC nil",
			args{v1SC: nil},
			"unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCasTypeFromSC(tt.args.v1SC); got != tt.want {
				t.Errorf("GetCasTypeFromSC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetReadyContainers(t *testing.T) {
	type args struct {
		containers []corev1.ContainerStatus
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"valid Values",
			args{containers: []corev1.ContainerStatus{{Ready: true}, {Ready: true}, {Ready: false}}},
			"2/3",
		},
		{
			"Invalid Values",
			args{containers: []corev1.ContainerStatus{}},
			"0/0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetReadyContainers(tt.args.containers); got != tt.want {
				t.Errorf("GetReadyContainers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUsedCapacityFromCVR(t *testing.T) {
	type args struct {
		cvrList *cstorv1.CStorVolumeReplicaList
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Valid values",
			args{cvrList: &cstorv1.CStorVolumeReplicaList{Items: []cstorv1.CStorVolumeReplica{{Status: cstorv1.CStorVolumeReplicaStatus{
				Phase:    "Init",
				Capacity: cstorv1.CStorVolumeReplicaCapacityDetails{Total: "5.0 GiB", Used: "2.1 GiB"},
			}}, {Status: cstorv1.CStorVolumeReplicaStatus{
				Phase:    "Healthy",
				Capacity: cstorv1.CStorVolumeReplicaCapacityDetails{Total: "5.0 GiB", Used: "2.5 GiB"},
			}}}}},
			"2.5 GiB",
		},
		{
			"Valid values",
			args{cvrList: &cstorv1.CStorVolumeReplicaList{Items: []cstorv1.CStorVolumeReplica{{Status: cstorv1.CStorVolumeReplicaStatus{
				Phase:    "Init",
				Capacity: cstorv1.CStorVolumeReplicaCapacityDetails{Total: "5.0 GiB", Used: "2.5 GiB"},
			}}, {Status: cstorv1.CStorVolumeReplicaStatus{
				Phase:    "Init",
				Capacity: cstorv1.CStorVolumeReplicaCapacityDetails{Total: "5.0 GiB", Used: "2.5 GiB"},
			}}}}},
			"",
		},
		{
			"Valid values",
			args{cvrList: &cstorv1.CStorVolumeReplicaList{Items: []cstorv1.CStorVolumeReplica{}}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUsedCapacityFromCVR(tt.args.cvrList); got != tt.want {
				t.Errorf("GetUsedCapacityFromCVR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidCasType(t *testing.T) {
	type args struct {
		casType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Valid Cas Name",
			args{casType: CstorCasType},
			true,
		},
		{
			"Invalid Cas Name",
			args{casType: "some-invalid-cas"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidCasType(tt.args.casType); got != tt.want {
				t.Errorf("IsValidCasType() = %v, want %v", got, tt.want)
			}
		})
	}
}
