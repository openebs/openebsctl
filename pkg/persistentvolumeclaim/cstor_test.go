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

	openebsFakeClientset "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDescribeCstorVolumeClaim(t *testing.T) {
	type args struct {
		c   *client.K8sClient
		pvc *corev1.PersistentVolumeClaim
		pv  *corev1.PersistentVolume
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "All Valid Values",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor, &cstorTargetPod),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pv:  &cstorPV1,
				pvc: &cstorPVC1,
			},
			wantErr: false,
		},
		{
			name: "PV missing",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pv:  nil,
				pvc: &cstorPVC1,
			},
			wantErr: false,
		},
		{
			name: "CV missing",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv2, &cva1, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pv:  &cstorPV1,
				pvc: &cstorPVC1,
			},
			wantErr: false,
		},
		{
			name: "CVC missing",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pv:  &cstorPV1,
				pvc: &cstorPVC1,
			},
			wantErr: false,
		},
		{
			name: "CVA missing",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pv:  &cstorPV1,
				pvc: &cstorPVC1,
			},
			wantErr: false,
		},
		{
			name: "CVRs missing",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva2, &cvc1, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pv:  &cstorPV1,
				pvc: &cstorPVC1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DescribeCstorVolumeClaim(tt.args.c, tt.args.pvc, tt.args.pv); (err != nil) != tt.wantErr {
				t.Errorf("DescribeCstorVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
