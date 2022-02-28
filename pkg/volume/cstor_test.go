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

	openebsFakeClientset "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestDescribeCstorVolume(t *testing.T) {
	cvRep := cv1
	cvRep.Spec.ReplicationFactor = 0
	cvRep.Status.ReplicaStatuses = nil
	type args struct {
		c   *client.K8sClient
		vol *corev1.PersistentVolume
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
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				vol: &cstorPV1,
			},
			wantErr: false,
		},
		{
			name: "All Valid Values with CV absent",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv2, &cva1, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				vol: &cstorPV1,
			},
			wantErr: true,
		},
		{
			name: "All Valid Values with CVC absent",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				vol: &cstorPV1,
			},
			wantErr: true,
		},
		{
			name: "All Valid Values with CVA absent",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				vol: &cstorPV1,
			},
			wantErr: false,
		},
		{
			name: "All Valid Values with CVRs absent",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva2, &cvc1, &cvc2, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				vol: &cstorPV1,
			},
			wantErr: false,
		},
		{
			name: "All Valid Values with CVR count as 0",
			args: args{
				c: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cvRep, &cv2, &cva2, &cvc1, &cvc2, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				vol: &cstorPV1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DescribeCstorVolume(tt.args.c, tt.args.vol); (err != nil) != tt.wantErr {
				t.Errorf("DescribeCstorVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCStor(t *testing.T) {
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
			name: "Test with all valid resources present.",
			args: args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{cstorPV1, cstorPV2}},
				openebsNS: "cstor",
			},
			want: []metav1.TableRow{{Cells: []interface{}{
				"cstor", "pvc-1", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-1"},
			}, {Cells: []interface{}{
				"cstor", "pvc-2", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-2"},
			}},
			wantErr: false,
		},
		{
			name: "Test with one of the required cv not present",
			args: args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cva1, &cva2, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{cstorPV1, cstorPV2}},
				openebsNS: "cstor",
			},
			want: []metav1.TableRow{{Cells: []interface{}{
				"cstor", "pvc-1", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-1"},
			}},
			wantErr: false,
		},
		{
			name: "Test with one of the required cva not present, i.e node cannot be determined",
			args: args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cvc1, &cvc2, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{cstorPV1, cstorPV2}},
				openebsNS: "cstor",
			},
			want: []metav1.TableRow{{Cells: []interface{}{
				"cstor", "pvc-1", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-1"},
			}, {Cells: []interface{}{
				"cstor", "pvc-2", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, ""},
			}},
			wantErr: false,
		},
		{
			name: "Test with one of the required cvc not present, i.e nothing should break in code",
			args: args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc1, &cvr1, &cvr2, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{cstorPV1, cstorPV2}},
				openebsNS: "cstor",
			},
			want: []metav1.TableRow{{Cells: []interface{}{
				"cstor", "pvc-1", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-1"},
			}, {Cells: []interface{}{
				"cstor", "pvc-2", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-2"},
			}},
			wantErr: false,
		},
		{
			name: "Test with two of the required cvrs not present, i.e nothing should break in code",
			args: args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc1, &cvr3, &cvr4, &cbkp, &ccbkp, &crestore),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{cstorPV1, cstorPV2}},
				openebsNS: "cstor",
			},
			want: []metav1.TableRow{{Cells: []interface{}{
				"cstor", "pvc-1", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-1"},
			}, {Cells: []interface{}{
				"cstor", "pvc-2", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-2"},
			}},
			wantErr: false,
		},
		{
			name: "Test with backup and restore crs not present, i.e nothing should break in code",
			args: args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1, &cv2, &cva1, &cva2, &cvc1, &cvr3, &cvr4),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{cstorPV1, cstorPV2}},
				openebsNS: "cstor",
			},
			want: []metav1.TableRow{{Cells: []interface{}{
				"cstor", "pvc-1", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-1"},
			}, {Cells: []interface{}{
				"cstor", "pvc-2", util.Healthy, "2.11.0", "4.0GiB", "cstor-sc", corev1.VolumeBound, corev1.ReadWriteOnce, "node-2"},
			}},
			wantErr: false,
		},
		{
			name: "Test with none of the underlying cstor crs",
			args: args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{cstorPV1}},
				openebsNS: "cstor",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test with none of the underlying cvas are present",
			args: args{
				c: &client.K8sClient{
					Ns:        "",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPV2, &cstorPVC1, &cstorPVC2, &nsCstor),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cv1),
				},
				pvList:    &corev1.PersistentVolumeList{Items: []corev1.PersistentVolume{cstorPV1}},
				openebsNS: "cstor",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCStor(tt.args.c, tt.args.pvList, tt.args.openebsNS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCStor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCStor() got = %v, want %v", got, tt.want)
			}
		})
	}
}
