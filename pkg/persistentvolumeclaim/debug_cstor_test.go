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

	"github.com/openebs/openebsctl/pkg/util"

	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This Currently is for debugging the code and doesnot involve mocking.
func TestDebugCstorVolumeClaim(t *testing.T) {
	k, _ := client.NewK8sClient("openebs")
	sc := "common-storageclass"
	type args struct {
		k   *client.K8sClient
		pvc *corev1.PersistentVolumeClaim
		pv  *corev1.PersistentVolume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test 1",
			args{
				k: k,
				pvc: &corev1.PersistentVolumeClaim{
					ObjectMeta: metav1.ObjectMeta{
						Name: "mongo",
					},
					Spec:   corev1.PersistentVolumeClaimSpec{StorageClassName: &sc},
					Status: corev1.PersistentVolumeClaimStatus{},
				},
				pv: &corev1.PersistentVolume{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pvc-b84f60ae-3f26-4110-a85d-bce7ec00dacc",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DebugCstorVolumeClaim(tt.args.k, tt.args.pvc, tt.args.pv); (err != nil) != tt.wantErr {
				t.Errorf("DebugCstorVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_displayPVCEvents(t *testing.T) {
	k, _ := client.NewK8sClient("openebs")
	type args struct {
		k   *client.K8sClient
		crs util.CstorVolumeResources
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test 1",
			args{k: k, crs: util.CstorVolumeResources{}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := displayPVCEvents(*tt.args.k, tt.args.crs); (err != nil) != tt.wantErr {
				t.Errorf("displayPVCEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceStatus(t *testing.T) {
	type args struct {
		crs util.CstorVolumeResources
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"test 1",
			args{crs: util.CstorVolumeResources{}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resourceStatus(tt.args.crs); (err != nil) != tt.wantErr {
				t.Errorf("resourceStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}