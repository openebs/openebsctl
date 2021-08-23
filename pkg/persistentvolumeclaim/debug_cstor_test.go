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
	"fmt"
	"testing"
	"time"

	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	cstortypes "github.com/openebs/api/v2/pkg/apis/types"
	openebsFakeClientset "github.com/openebs/api/v2/pkg/client/clientset/versioned/fake"
	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	fake2 "k8s.io/client-go/kubernetes/typed/core/v1/fake"
	k8stest "k8s.io/client-go/testing"
)

func TestDebugCstorVolumeClaim(t *testing.T) {
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
			"Test with all valid values",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &cspi1, &cspi2, &bd1, &bd2, &bdc1, &bdc2, &cv1, &cva1, &cvc1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
			},
			false,
		},
		{
			"Test with PV missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &cspi1, &cspi2, &bd1, &bd2, &bdc1, &bdc2, &cva1, &cvc1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  nil,
			},
			false,
		},
		{
			"Test with CV missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &cspi1, &cspi2, &bd1, &bd2, &bdc1, &bdc2, &cva1, &cvc1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
			},
			false,
		},
		{
			"Test with CVC Missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &cspi1, &cspi2, &bd1, &bd2, &bdc1, &bdc2, &cv1, &cva1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
			},
			false,
		},
		{
			"Test with CVA missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &cspi1, &cspi2, &bd1, &bd2, &bdc1, &bdc2, &cv1, &cvc1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
			},
			false,
		},
		{
			"Test with CVRs missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &cspi1, &cspi2, &bd1, &bd2, &bdc1, &bdc2, &cv1, &cva1, &cvc1),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
			},
			false,
		},
		{
			"Test with cspc missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspi1, &cspi2, &bd1, &bd2, &bdc1, &bdc2, &cv1, &cva1, &cvc1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
			},
			false,
		},
		{
			"Test with cspis missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &bd1, &bd2, &bdc1, &bdc2, &cv1, &cva1, &cvc1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
			},
			false,
		},
		{
			"Test with bds missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &cspi1, &cspi2, &bdc1, &bdc2, &cv1, &cva1, &cvc1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
			},
			false,
		},
		{
			"Test with bdcs missing",
			args{
				k: &client.K8sClient{
					Ns:        "cstor",
					K8sCS:     fake.NewSimpleClientset(&cstorPV1, &cstorPVC1, &cstorSc),
					OpenebsCS: openebsFakeClientset.NewSimpleClientset(&cspc, &cspi1, &cspi2, &bd1, &bd2, &cv1, &cva1, &cvc1, &cvr1, &cvr2),
				},
				pvc: &cstorPVC1,
				pv:  &cstorPV1,
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

func Test_displayBDCEvents(t *testing.T) {
	type args struct {
		k       client.K8sClient
		crs     util.CstorVolumeResources
		bdcFunc func(*client.K8sClient, map[string]corev1.EventList)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with valid values and events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(&bdcEvent1, &bdcEvent2),
				},
				crs: util.CstorVolumeResources{
					BDCs: &bdcList,
				},
				bdcFunc: eventFunc,
			},
			false,
		},
		{
			"Test with valid values with no events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					BDCs: &bdcList,
				},
				bdcFunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no BDCs",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					BDCs: nil,
				},
				bdcFunc: nil,
			},
			true,
		},
		{
			"Test with valid values with no BDCList as empty",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					BDCs: &v1alpha1.BlockDeviceClaimList{
						Items: []v1alpha1.BlockDeviceClaim{},
					},
				},
				bdcFunc: nil,
			},
			true,
		},
		{
			"Test with valid values with no events errored out",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					BDCs: &v1alpha1.BlockDeviceClaimList{
						Items: []v1alpha1.BlockDeviceClaim{},
					},
				},
				bdcFunc: noEventFunc,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.bdcFunc != nil {
				tt.args.bdcFunc(&tt.args.k, map[string]corev1.EventList{
					"bdc-1": {Items: []corev1.Event{bdcEvent1}},
					"bdc-2": {Items: []corev1.Event{bdcEvent2}},
				})
			}
			if err := displayBDCEvents(tt.args.k, tt.args.crs); (err != nil) != tt.wantErr {
				t.Errorf("displayBDCEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_displayCSPCEvents(t *testing.T) {
	type args struct {
		k        client.K8sClient
		crs      util.CstorVolumeResources
		cspcFunc func(*client.K8sClient, map[string]corev1.EventList)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with valid values and events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(&cspcEvent),
				},
				crs: util.CstorVolumeResources{
					CSPC: &cspc,
				},
				cspcFunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CSPC: &cspc,
				},
				cspcFunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no CSPC",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CSPC: nil,
				},
				cspcFunc: nil,
			},
			true,
		},
		{
			"Test with valid values with no events errored out",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CSPC: &cspc,
				},
				cspcFunc: noEventFunc,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.cspcFunc != nil {
				tt.args.cspcFunc(&tt.args.k, map[string]corev1.EventList{})
			}
			if err := displayCSPCEvents(tt.args.k, tt.args.crs); (err != nil) != tt.wantErr {
				t.Errorf("displayCSPCEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_displayCSPIEvents(t *testing.T) {
	type args struct {
		k        client.K8sClient
		crs      util.CstorVolumeResources
		cspiFunc func(*client.K8sClient, map[string]corev1.EventList)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with valid values and events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(&cspiEvent1, &cspiEvent2),
				},
				crs: util.CstorVolumeResources{
					CSPIs: &cspiList,
				},
				cspiFunc: eventFunc,
			},
			false,
		},
		{
			"Test with valid values with no events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CSPIs: &cspiList,
				},
				cspiFunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no CSPIs",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CSPIs: nil,
				},
				cspiFunc: nil,
			},
			true,
		},
		{
			"Test with valid values with no CSPIList as empty",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CSPIs: &v1.CStorPoolInstanceList{
						Items: []v1.CStorPoolInstance{},
					},
				},
				cspiFunc: nil,
			},
			true,
		},
		{
			"Test with valid values with no events errored out",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CSPIs: &v1.CStorPoolInstanceList{
						Items: []v1.CStorPoolInstance{},
					},
				},
				cspiFunc: noEventFunc,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.cspiFunc != nil {
				tt.args.cspiFunc(&tt.args.k, map[string]corev1.EventList{
					"cspc-1": {Items: []corev1.Event{cspiEvent1}},
					"cspc-2": {Items: []corev1.Event{cspiEvent2}},
				})
			}
			if err := displayCSPIEvents(tt.args.k, tt.args.crs); (err != nil) != tt.wantErr {
				t.Errorf("displayCSPIEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_displayCVCEvents(t *testing.T) {
	type args struct {
		k       client.K8sClient
		crs     util.CstorVolumeResources
		cvcFunc func(*client.K8sClient, map[string]corev1.EventList)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with valid values and events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(&cvcEvent1, &cvcEvent2),
				},
				crs: util.CstorVolumeResources{
					CVC: &cvc1,
				},
				cvcFunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CVC: &cvc1,
				},
				cvcFunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no CVC",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CVC: nil,
				},
				cvcFunc: nil,
			},
			true,
		},
		{
			"Test with valid values with no events errored out",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CVC: &cvc1,
				},
				cvcFunc: noEventFunc,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.cvcFunc != nil {
				tt.args.cvcFunc(&tt.args.k, map[string]corev1.EventList{})
			}
			if err := displayCVCEvents(tt.args.k, tt.args.crs); (err != nil) != tt.wantErr {
				t.Errorf("displayCVCEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_displayCVREvents(t *testing.T) {
	type args struct {
		k       client.K8sClient
		crs     util.CstorVolumeResources
		cvrfunc func(*client.K8sClient, map[string]corev1.EventList)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with valid values and events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(&cvrEvent1, &cvrEvent2),
				},
				crs: util.CstorVolumeResources{
					CVRs: &cvrList,
				},
				cvrfunc: eventFunc,
			},
			false,
		},
		{
			"Test with valid values with no events",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CVRs: &cvrList,
				},
				cvrfunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no CVRs",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CVRs: nil,
				},
				cvrfunc: nil,
			},
			true,
		},
		{
			"Test with valid values with no CVRIList as empty",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CVRs: &v1.CStorVolumeReplicaList{
						Items: []v1.CStorVolumeReplica{},
					},
				},
				cvrfunc: nil,
			},
			true,
		},
		{
			"Test with valid values with no events errored out",
			args{
				k: client.K8sClient{
					Ns:    "cstor",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					CVRs: &v1.CStorVolumeReplicaList{
						Items: []v1.CStorVolumeReplica{},
					},
				},
				cvrfunc: noEventFunc,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.cvrfunc != nil {
				tt.args.cvrfunc(&tt.args.k, map[string]corev1.EventList{
					"pvc-1-rep-1": {Items: []corev1.Event{cvrEvent1}},
					"pvc-1-rep-2": {Items: []corev1.Event{cvrEvent2}},
				})
			}
			if err := displayCVREvents(tt.args.k, tt.args.crs); (err != nil) != tt.wantErr {
				t.Errorf("displayCVREvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_displayPVCEvents(t *testing.T) {
	type args struct {
		k       client.K8sClient
		crs     util.CstorVolumeResources
		pvcFunc func(*client.K8sClient, map[string]corev1.EventList)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with valid values and events",
			args{
				k: client.K8sClient{
					Ns:    "default",
					K8sCS: fake.NewSimpleClientset(&pvcEvent1, &pvcEvent2),
				},
				crs: util.CstorVolumeResources{
					PVC: &cstorPVC1,
				},
				pvcFunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no events",
			args{
				k: client.K8sClient{
					Ns:    "default",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					PVC: &cstorPVC1,
				},
				pvcFunc: nil,
			},
			false,
		},
		{
			"Test with valid values with no events errored out",
			args{
				k: client.K8sClient{
					Ns:    "default",
					K8sCS: fake.NewSimpleClientset(),
				},
				crs: util.CstorVolumeResources{
					PVC: &cstorPVC1,
				},
				pvcFunc: noEventFunc,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.pvcFunc != nil {
				tt.args.pvcFunc(&tt.args.k, map[string]corev1.EventList{})
			}
			if err := displayPVCEvents(tt.args.k, tt.args.crs); (err != nil) != tt.wantErr {
				t.Errorf("displayPVCEvents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_resourceStatus(t *testing.T) {
	var cvrListWithMoreUsedCapacity = v1.CStorVolumeReplicaList{Items: []v1.CStorVolumeReplica{{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "pvc-1-rep-1",
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels:            map[string]string{cstortypes.PersistentVolumeLabelKey: "pvc-1"},
			Finalizers:        []string{},
			Namespace:         "cstor",
		},
		Status: v1.CStorVolumeReplicaStatus{
			Capacity: v1.CStorVolumeReplicaCapacityDetails{
				Total: "4Gi",
				Used:  "3.923GiB",
			},
			Phase: v1.CVRStatusOnline,
		},
	}}}
	type args struct {
		crs util.CstorVolumeResources
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with all valid values",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        &cvrList,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all PV absent",
			args{crs: util.CstorVolumeResources{
				PV:          nil,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        &cvrList,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all CV absent",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          nil,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        &cvrList,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all CVC absent",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         nil,
				CVA:         &cva1,
				CVRs:        &cvrList,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all CVA absent",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         nil,
				CVRs:        &cvrList,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all CVRs absent",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        nil,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all BDs absent",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        &cvrList,
				PresentBDs:  nil,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all BDCs absent",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        &cvrList,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        nil,
				CSPIs:       &cspiList,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all CSPIs absent",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        &cvrList,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       nil,
				CSPC:        &cspc,
			}},
			false,
		},
		{
			"Test with all CSPC absent",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        &cvrList,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        nil,
			}},
			false,
		},
		{
			"Test with all Used Capacity exceeding 80%",
			args{crs: util.CstorVolumeResources{
				PV:          &cstorPV1,
				PVC:         &cstorPVC1,
				CV:          &cv1,
				CVC:         &cvc1,
				CVA:         &cva1,
				CVRs:        &cvrListWithMoreUsedCapacity,
				PresentBDs:  &bdList,
				ExpectedBDs: expectedBDs,
				BDCs:        &bdcList,
				CSPIs:       &cspiList,
				CSPC:        nil,
			}},
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

func eventFunc(c *client.K8sClient, eventMap map[string]corev1.EventList) {
	c.K8sCS.CoreV1().Events(c.Ns).(*fake2.FakeEvents).Fake.PrependReactor("*", "events", func(action k8stest.Action) (handled bool, ret runtime.Object, err error) {
		listOpts, ok := action.(k8stest.ListActionImpl)
		if ok {
			val, matched := listOpts.ListRestrictions.Fields.RequiresExactMatch("involvedObject.name")
			if matched {
				if events, present := eventMap[val]; present {
					return true, &events, nil
				} else {
					return true, nil, fmt.Errorf("invalid fieldSelector")
				}
			} else {
				return true, nil, fmt.Errorf("invalid fieldSelector")
			}
		} else {
			return true, nil, fmt.Errorf("invalid fieldSelector")
		}
	})
}

func noEventFunc(c *client.K8sClient, eventMap map[string]corev1.EventList) {
	c.K8sCS.CoreV1().Events(c.Ns).(*fake2.FakeEvents).Fake.PrependReactor("*", "events", func(action k8stest.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("failed to list events")
	})
}
