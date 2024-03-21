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

	corev1 "k8s.io/api/core/v1"
)

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
			args{casType: LVMCasType},
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
