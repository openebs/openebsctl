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

	"github.com/openebs/openebsctl/pkg/util"
	zfs "github.com/openebs/zfs-localpv/pkg/apis/openebs.io/zfs/v1"
	zvolclient "github.com/openebs/zfs-localpv/pkg/generated/clientset/internalclientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// getLVMClient returns OpenEBS clientset by taking kubeconfig as an
// argument
func getZFSclient(kubeconfig string) (*zvolclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("could not build config from flags: %v", err)
	}
	client, err := zvolclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not get new config: %v", err)
	}
	return client, nil
}

// GetZFSVols returns a list or a map of ZFSVolume depending upon rType & options
func (k K8sClient) GetZFSVols(volNames []string, rType util.ReturnType, labelSelector string, options util.MapOptions) (*zfs.ZFSVolumeList, map[string]zfs.ZFSVolume, error) {
	zvols, err := k.ZFCS.ZfsV1().ZFSVolumes("").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, nil, err
	}
	var list []zfs.ZFSVolume
	if len(volNames) == 0 {
		list = zvols.Items
	} else {
		zvsMap := make(map[string]zfs.ZFSVolume)
		for _, zv := range zvols.Items {
			zvsMap[zv.Name] = zv
		}
		for _, name := range volNames {
			if zv, ok := zvsMap[name]; ok {
				list = append(list, zv)
			} else {
				fmt.Printf("Error from server (NotFound): zfsVolume %s not found\n", name)
			}
		}
	}
	if rType == util.List {
		return &zfs.ZFSVolumeList{
			Items: list,
		}, nil, nil
	}
	if rType == util.Map {
		zvMap := make(map[string]zfs.ZFSVolume)
		switch options.Key {
		case util.Label:
			for _, zv := range list {
				if vol, ok := zv.Labels[options.LabelKey]; ok {
					zvMap[vol] = zv
				}
			}
			return nil, zvMap, nil
		case util.Name:
			for _, zv := range list {
				zvMap[zv.Name] = zv
			}
			return nil, zvMap, nil
		default:
			return nil, nil, fmt.Errorf("invalid map options")
		}
	}
	return nil, nil, fmt.Errorf("invalid return type")
}

// GetZFSNodes return a list of ZFSNodes
func (k K8sClient) GetZFSNodes(volNames []string, rType util.ReturnType, labelSelector string, options util.MapOptions) (*zfs.ZFSNodeList, map[string]zfs.ZFSNode, error) {
	zfsNode, err := k.ZFCS.ZfsV1().ZFSNodes("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	var list []zfs.ZFSNode
	if len(volNames) == 0 {
		list = zfsNode.Items
	} else {
		zvsMap := make(map[string]zfs.ZFSNode)
		for _, zn := range zfsNode.Items {
			zvsMap[zn.Name] = zn
		}
		for _, name := range volNames {
			if zv, ok := zvsMap[name]; ok {
				list = append(list, zv)
			} else {
				// This might be omitted
				fmt.Printf("Error from server (NotFound): zfsVolume %s not found\n", name)
			}
		}
	}
	if rType == util.List {
		return &zfs.ZFSNodeList{
			Items: list,
		}, nil, nil
	}
	if rType == util.Map {
		znMap := make(map[string]zfs.ZFSNode)
		switch options.Key {
		case util.Label:
			for _, zn := range list {
				if vol, ok := zn.Labels[options.LabelKey]; ok {
					znMap[vol] = zn
				}
			}
			return nil, znMap, nil
		case util.Name:
			for _, zn := range list {
				znMap[zn.Name] = zn
			}
			return nil, znMap, nil
		default:
			return nil, nil, fmt.Errorf("invalid map options")
		}
	}
	return nil, nil, fmt.Errorf("invalid return type")
}
