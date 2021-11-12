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
	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

func TestDescribeJivaVolumeClaim(t *testing.T) {
	type args struct {
		c         *client.K8sClient
		pvc       *corev1.PersistentVolumeClaim
		vol       *corev1.PersistentVolume
		mountPods string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DescribeJivaVolumeClaim(tt.args.c, tt.args.pvc, tt.args.vol, tt.args.mountPods); (err != nil) != tt.wantErr {
				t.Errorf("DescribeJivaVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
