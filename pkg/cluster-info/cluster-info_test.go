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

package cluster_info

import (
	"reflect"
	"testing"

	"github.com/openebs/openebsctl/pkg/client"
	"github.com/openebs/openebsctl/pkg/util"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_getComponentDataByComponents(t *testing.T) {
	type args struct {
		k              *client.K8sClient
		componentNames string
		casType        string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]util.ComponentData
		wantErr bool
	}{
		{
			"All components present and running",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&cspcOperator, &cvcOperator, &cstorAdmissionWebhook, &ndm, &ndmOperator, &openebsCstorCsiController, &openebsCstorCsiNode),
				},
				componentNames: util.CasTypeToComponentNamesMap[util.CstorCasType],
				casType:        util.CstorCasType,
			},
			map[string]util.ComponentData{
				"cspc-operator":                {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"cvc-operator":                 {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"cstor-admission-webhook":      {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"openebs-cstor-csi-node":       {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"openebs-cstor-csi-controller": {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"ndm":                          {Namespace: "openebs", Status: "Running", Version: "1.1", CasType: "cstor"},
				"openebs-ndm-operator":         {Namespace: "openebs", Status: "Running", Version: "1.1", CasType: "cstor"},
			},
			false,
		},
		{
			"Some components present and running",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&cspcOperator, &cvcOperator, &cstorAdmissionWebhook, &ndmOperator, &openebsCstorCsiNode),
				},
				componentNames: util.CasTypeToComponentNamesMap[util.CstorCasType],
				casType:        util.CstorCasType,
			},
			map[string]util.ComponentData{
				"cspc-operator":                {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"cvc-operator":                 {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"cstor-admission-webhook":      {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"openebs-cstor-csi-node":       {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"openebs-cstor-csi-controller": {},
				"ndm":                          {},
				"openebs-ndm-operator":         {Namespace: "openebs", Status: "Running", Version: "1.1", CasType: "cstor"},
			},
			false,
		},
		{
			"All components present and running with some component having evicted pods",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&cspcOperator, &cvcOperator, &cstorAdmissionWebhook, &ndm, &ndmOperator, &openebsCstorCsiController, &openebsCstorCsiNode, &cspcOperatorEvicted, &cvcOperatorEvicted),
				},
				componentNames: util.CasTypeToComponentNamesMap[util.CstorCasType],
				casType:        util.CstorCasType,
			},
			map[string]util.ComponentData{
				"cspc-operator":                {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"cvc-operator":                 {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"cstor-admission-webhook":      {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"openebs-cstor-csi-node":       {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"openebs-cstor-csi-controller": {Namespace: "openebs", Status: "Running", Version: "2.1", CasType: "cstor"},
				"ndm":                          {Namespace: "openebs", Status: "Running", Version: "1.1", CasType: "cstor"},
				"openebs-ndm-operator":         {Namespace: "openebs", Status: "Running", Version: "1.1", CasType: "cstor"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getComponentDataByComponents(tt.args.k, tt.args.componentNames, tt.args.casType)
			if (err != nil) != tt.wantErr {
				t.Errorf("getComponentDataByComponents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getComponentDataByComponents() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLocalPVDeviceStatus(t *testing.T) {
	type args struct {
		componentDataMap map[string]util.ComponentData
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			"ndm and localpv provisioner in same ns",
			args{
				componentDataMap: map[string]util.ComponentData{
					"openebs-localpv-provisioner": {
						Namespace: "openebs",
						Status:    "Running",
						Version:   "1.1",
						CasType:   util.LocalDeviceCasType,
					},
					"ndm": {
						Namespace: "openebs",
						Status:    "Running",
						Version:   "3.1",
						CasType:   util.LocalDeviceCasType,
					},
				},
			},
			"Healthy",
			"2/2",
			false,
		},
		{
			"ndm and localpv provisioner in same ns but ndm down",
			args{
				componentDataMap: map[string]util.ComponentData{
					"openebs-localpv-provisioner": {
						Namespace: "openebs",
						Status:    "Running",
						Version:   "1.1",
						CasType:   util.LocalDeviceCasType,
					},
					"ndm": {
						Namespace: "openebs",
						Status:    "Pending",
						Version:   "3.1",
						CasType:   util.LocalDeviceCasType,
					},
				},
			},
			"Degraded",
			"1/2",
			false,
		},
		{
			"ndm and localpv provisioner in same ns but both down",
			args{
				componentDataMap: map[string]util.ComponentData{
					"openebs-localpv-provisioner": {
						Namespace: "openebs",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   util.LocalDeviceCasType,
					},
					"ndm": {
						Namespace: "openebs",
						Status:    "Pending",
						Version:   "3.1",
						CasType:   util.LocalDeviceCasType,
					},
				},
			},
			"Unhealthy",
			"0/2",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getLocalPVDeviceStatus(tt.args.componentDataMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLocalPVDeviceStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getLocalPVDeviceStatus() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getLocalPVDeviceStatus() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getNamespace(t *testing.T) {
	type args struct {
		componentDataMap map[string]util.ComponentData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"some running components with ndm in same ns",
			args{
				componentDataMap: map[string]util.ComponentData{
					"cstor-csi-controller": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"ndm": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "3.1",
						CasType:   "cstor",
					},
					"cstor-operator": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"cstor-some-xyz-component": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "1.1",
						CasType:   "cstor",
					},
				},
			},
			"cstor",
		},
		{
			"some running components with ndm in different ns",
			args{
				componentDataMap: map[string]util.ComponentData{
					"cstor-csi-controller": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"ndm": {
						Namespace: "openebs",
						Status:    "Running",
						Version:   "3.1",
						CasType:   "cstor",
					},
					"cstor-operator": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"cstor-some-xyz-component": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
				},
			},
			"cstor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNamespace(tt.args.componentDataMap); got != tt.want {
				t.Errorf("getNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStatus(t *testing.T) {
	type args struct {
		componentDataMap map[string]util.ComponentData
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			"some running components",
			args{
				componentDataMap: map[string]util.ComponentData{
					"cstor-csi-controller": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"ndm": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "3.1",
						CasType:   "cstor",
					},
					"cstor-operator": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"cstor-some-xyz-component": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "1.1",
						CasType:   "cstor",
					},
				},
			},
			"Degraded",
			"2/4",
		},
		{
			"No running components",
			args{
				componentDataMap: map[string]util.ComponentData{
					"cstor-csi-controller": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"ndm": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "3.1",
						CasType:   "cstor",
					},
					"cstor-operator": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"cstor-some-xyz-component": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
				},
			},
			"Unhealthy",
			"0/4",
		},
		{
			"All running components",
			args{
				componentDataMap: map[string]util.ComponentData{
					"cstor-csi-controller": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"ndm": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "3.1",
						CasType:   "cstor",
					},
					"cstor-operator": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"cstor-some-xyz-component": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "1.1",
						CasType:   "cstor",
					},
				},
			},
			"Healthy",
			"4/4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getStatus(tt.args.componentDataMap)
			if got != tt.want {
				t.Errorf("getStatus() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getStatus() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getVersion(t *testing.T) {
	type args struct {
		componentDataMap map[string]util.ComponentData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"some running components",
			args{
				componentDataMap: map[string]util.ComponentData{
					"cstor-csi-controller": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"ndm": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "3.1",
						CasType:   "cstor",
					},
					"cstor-operator": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"cstor-some-xyz-component": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "1.1",
						CasType:   "cstor",
					},
				},
			},
			"1.1",
		},
		{
			"No running components except ndm",
			args{
				componentDataMap: map[string]util.ComponentData{
					"cstor-csi-controller": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"ndm": {
						Namespace: "cstor",
						Status:    "Running",
						Version:   "3.1",
						CasType:   "cstor",
					},
					"cstor-operator": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
					"cstor-some-xyz-component": {
						Namespace: "cstor",
						Status:    "Pending",
						Version:   "1.1",
						CasType:   "cstor",
					},
				},
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getVersion(tt.args.componentDataMap); got != tt.want {
				t.Errorf("getVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

//map[cspc-operator:{openebs Running 2.1 cstor} cstor-admission-webhook:{openebs Running 2.1 cstor} cvc-operator:{openebs Running 2.1 cstor} ndm:{openebs Running 1.1 cstor} openebs-cstor-csi-controller:{openebs Running 2.1 cstor} openebs-cstor-csi-node:{openebs Running 2.1 cstor} openebs-ndm-operator:{openebs Running 1.1 cstor}],
//map[cspc-operator:{openebs Running 2.1 cstor} cstor-admission-webhook:{openebs Running 2.1 cstor} cvc-operator:{openebs Running 2.1 cstor} ndm:{openebs Running 1.1 cstor} ndm-operator:{openebs Running 1.1 cstor} openebs-cstor-csi-controller:{openebs Running 2.1 cstor} openebs-cstor-csi-node:{openebs Running 2.1 cstor}]

func Test_compute(t *testing.T) {
	type args struct {
		k *client.K8sClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"All components of cstor present and running",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&cspcOperator, &cvcOperator, &cstorAdmissionWebhook, &ndm, &ndmOperator, &openebsCstorCsiController, &openebsCstorCsiNode),
				},
			},
			false,
		},
		{
			"Some components of cstor present and running",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&cspcOperator, &cvcOperator, &cstorAdmissionWebhook, &ndmOperator, &openebsCstorCsiNode),
				},
			},
			false,
		},
		{
			"All components of cstor present and running with some component having evicted pods",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&cspcOperator, &cvcOperator, &cstorAdmissionWebhook, &ndm, &ndmOperator, &openebsCstorCsiController, &openebsCstorCsiNode, &cspcOperatorEvicted, &cvcOperatorEvicted),
				},
			},
			false,
		},
		{
			"If no components are present",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(),
				},
			},
			true,
		},
		{
			"If ndm and localpv provisioner components are in same ns",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&localpvProvisionerInOpenebs, &ndm, &ndmOperator),
				},
			},
			false,
		},
		{
			"If ndm and localpv provisioner components are in different ns",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&localpvProvisioner, &ndm, &ndmOperator),
				},
			},
			false,
		},
		{
			"If jiva and ndm in same ns",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&jivaOperator, &openebsJivaCsiController, &openebsJivaCsiNode, &ndm, &ndmOperator, &localpvProvisionerInOpenebs),
				},
			},
			false,
		},
		{
			"If jiva and ndm in different ns",
			args{
				k: &client.K8sClient{
					Ns:    "",
					K8sCS: fake.NewSimpleClientset(&jivaOperator, &openebsJivaCsiController, &openebsJivaCsiNode, &ndmXYZ, &ndmOperatorXYZ, &localpvProvisionerInOpenebs),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := compute(tt.args.k); (err != nil) != tt.wantErr {
				t.Errorf("compute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
