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

package persistentvolumeclaim

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestDescribeGenericVolumeClaim(t *testing.T) {
	type args struct {
		pvc     *corev1.PersistentVolumeClaim
		pv      *corev1.PersistentVolume
		casType string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "All Valid Values",
			args: args{
				pv:      &cstorPV1,
				pvc:     &cstorPVC1,
				casType: "some-cas",
			},
			wantErr: false,
		},
		{
			name: "PV missing",
			args: args{
				pv:      nil,
				pvc:     &cstorPVC1,
				casType: "some-cas",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DescribeGenericVolumeClaim(tt.args.pvc, tt.args.pv, tt.args.casType); (err != nil) != tt.wantErr {
				t.Errorf("DescribeGenericVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
