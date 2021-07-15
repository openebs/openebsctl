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
	"testing"

	fakelvmclient "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset/fake"
	"github.com/openebs/openebsctl/pkg/client"
)

func TestGetVolumeGroup(t *testing.T) {
	type args struct {
		c  *client.K8sClient
		vg []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"no LVM volumegroups present",
			args{
				c: &client.K8sClient{
					Ns:        "lvmlocalpv",
					K8sCS:     nil,
					OpenebsCS: nil,
					LVMCS:     fakelvmclient.NewSimpleClientset(),
				},
				vg: nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetVolumeGroups(tt.args.c, tt.args.vg); (err != nil) != tt.wantErr {
				t.Errorf("GetVolumeGroups() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
