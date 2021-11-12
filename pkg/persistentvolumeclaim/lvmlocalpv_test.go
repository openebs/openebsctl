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

	"github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset/fake"
	fakelvm "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset/typed/lvm/v1alpha1/fake"
	"github.com/openebs/openebsctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stest "k8s.io/client-go/testing"
)

func TestDescribeLVMVolumeClaim(t *testing.T) {
	type args struct {
		c         *client.K8sClient
		pvc       *corev1.PersistentVolumeClaim
		pv        *corev1.PersistentVolume
		mountPods string
		lvmfunc   func(sClient *client.K8sClient)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test with all valid values",
			args{c: &client.K8sClient{Ns: "lvmlocalpv", LVMCS: fake.NewSimpleClientset(&lvmVol1), K8sCS: k8sfake.NewSimpleClientset()},
				pv:        &lvmPV1,
				pvc:       &lvmPVC1,
				mountPods: "",
			},
			false,
		},
		{
			"Test with PV missing",
			args{c: &client.K8sClient{Ns: "lvmlocalpv", LVMCS: fake.NewSimpleClientset(&lvmVol1), K8sCS: k8sfake.NewSimpleClientset()},
				pv:        nil,
				pvc:       &lvmPVC1,
				mountPods: "",
			},
			false,
		},
		{
			"Test with LVM Vol missing",
			args{c: &client.K8sClient{Ns: "lvmlocalpv", LVMCS: fake.NewSimpleClientset(), K8sCS: k8sfake.NewSimpleClientset()},
				pv:        &lvmPV1,
				pvc:       &lvmPVC1,
				mountPods: "",
				lvmfunc:   lvmVolNotExists,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.lvmfunc != nil {
				tt.args.lvmfunc(tt.args.c)
			}
			if err := DescribeLVMVolumeClaim(tt.args.c, tt.args.pvc, tt.args.pv, tt.args.mountPods); (err != nil) != tt.wantErr {
				t.Errorf("DescribeLVMVolumeClaim() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// lvmVolNotExists makes fakelvmClientSet return error
func lvmVolNotExists(c *client.K8sClient) {
	// NOTE: Set the VERB & Resource correctly & make it work for single resources
	c.LVMCS.LocalV1alpha1().(*fakelvm.FakeLocalV1alpha1).Fake.PrependReactor("*", "*", func(action k8stest.Action) (handled bool, ret runtime.Object, err error) {
		return true, nil, fmt.Errorf("failed to list LVMVolumes")
	})
}
