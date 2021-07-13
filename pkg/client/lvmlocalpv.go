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

package client

import (
	"context"
	"fmt"

	lvm "github.com/openebs/lvm-localpv/pkg/apis/openebs.io/lvm/v1alpha1"
	lvmclient "github.com/openebs/lvm-localpv/pkg/generated/clientset/internalclientset"
	"github.com/openebs/openebsctl/pkg/util"
	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// getOpenEBSClient returns OpenEBS clientset by taking kubeconfig as an
// argument
func getLVMclient(kubeconfig string) (*lvmclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("could not build config from flags: %v", err)
	}
	client, err := lvmclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not get new config: %v", err)
	}
	return client, nil
}

// GetLVMvol returns a list or a map of LVMVolume depending upon rType & options
func (k K8sClient) GetLVMvol(lVols []string, rType util.ReturnType, labelSelector string, options util.MapOptions) (*lvm.LVMVolumeList, map[string]lvm.LVMVolume, error) {
	// NOTE: The resource name must be plural and the API-group should be present for getting CRs
	lvs, err := k.LVMCS.LocalV1alpha1().LVMVolumes(k.Ns).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	var list []lvm.LVMVolume
	if lVols == nil || len(lVols) == 0 {
		list = lvs.Items
	} else {
		lvsMap := make(map[string]lvm.LVMVolume)
		for _, lv := range lvs.Items {
			lvsMap[lv.Name] = lv
		}
		for _, name := range lVols {
			if lv, ok := lvsMap[name]; ok {
				list = append(list, lv)
			} else {
				fmt.Printf("Error from server (NotFound): lvmvolume %s not found\n", name)
			}
		}
	}
	if rType == util.List {
		return &lvm.LVMVolumeList{
			Items: list,
		}, nil, nil
	}
	if rType == util.Map {
		lvMap := make(map[string]lvm.LVMVolume)
		switch options.Key {
		case util.Label:
			for _, lv := range list {
				if vol, ok := lv.Labels[options.LabelKey]; ok {
					lvMap[vol] = lv
				}
			}
			return nil, lvMap, nil
		case util.Name:
			for _, lv := range list {
				lvMap[lv.Name] = lv
			}
			return nil, lvMap, nil
		default:
			return nil, nil, errors.New("invalid map options")
		}
	}
	return nil, nil, errors.New("invalid return type")
}
