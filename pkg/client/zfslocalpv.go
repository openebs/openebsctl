package client

import (
	"context"
	"fmt"

	"github.com/openebs/openebsctl/pkg/util"
	zfs "github.com/openebs/zfs-localpv/pkg/apis/openebs.io/zfs/v1"
	zvolclient "github.com/openebs/zfs-localpv/pkg/generated/clientset/internalclientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
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
func (k K8sClient) GetZFSNodes() (*zfs.ZFSNodeList, error) {
	return k.ZFCS.ZfsV1().ZFSNodes("").List(context.TODO(), metav1.ListOptions{})
}
