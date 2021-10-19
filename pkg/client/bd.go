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

	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// required for auth, see: https://github.com/kubernetes/client-go/tree/v0.17.3/plugin/pkg/client/auth
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// GetBD returns the BlockDevice passed as name with OpenEBS's Client
func (k K8sClient) GetBD(bd string) (*v1alpha1.BlockDevice, error) {
	blockDevice, err := k.OpenebsCS.OpenebsV1alpha1().BlockDevices(k.Ns).Get(context.TODO(), bd, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting block device")
	}
	return blockDevice, nil
}

// GetBDs returns a list of BlockDevices based on the values of bdNames slice.
// bdNames slice if is nil or empty, it returns all the BDs in the cluster.
// bdNames slice if is not nil or not empty, it return the BDs whose names are present in the slice.
// labelselector takes the label(key+value) and makes an api call with this filter applied. Can be empty string if label filtering is not needed.
func (k K8sClient) GetBDs(bdNames []string, labelselector string) (*v1alpha1.BlockDeviceList, error) {
	bds, err := k.OpenebsCS.OpenebsV1alpha1().BlockDevices(k.Ns).List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting block device")
	}
	if len(bdNames) == 0 {
		return bds, nil
	}
	bdNameBDmap := make(map[string]v1alpha1.BlockDevice)
	for _, item := range bds.Items {
		bdNameBDmap[item.Name] = item
	}
	var items = make([]v1alpha1.BlockDevice, 0)
	for _, name := range bdNames {
		if _, ok := bdNameBDmap[name]; ok {
			items = append(items, bdNameBDmap[name])
		}
	}
	return &v1alpha1.BlockDeviceList{
		Items: items,
	}, nil
}

// GetBDCs returns a list of BlockDeviceClaims based on the values of bdcNames slice.
// bdcNames slice if is nil or empty, it returns all the BDCs in the cluster.
// bdcNames slice if is not nil or not empty, it return the BDCs whose names are present in the slice.
// labelselector takes the label(key+value) and makes an api call with this filter applied. Can be empty string if label filtering is not needed.
func (k K8sClient) GetBDCs(bdcNames []string, labelselector string) (*v1alpha1.BlockDeviceClaimList, error) {
	bds, err := k.OpenebsCS.OpenebsV1alpha1().BlockDeviceClaims(k.Ns).List(context.TODO(), metav1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return nil, errors.Wrapf(err, "Error while getting block device")
	}
	if len(bdcNames) == 0 {
		return bds, nil
	}
	bdcNameBDCmap := make(map[string]v1alpha1.BlockDeviceClaim)
	for _, item := range bds.Items {
		bdcNameBDCmap[item.Name] = item
	}
	var items = make([]v1alpha1.BlockDeviceClaim, 0)
	for _, name := range bdcNames {
		if _, ok := bdcNameBDCmap[name]; ok {
			items = append(items, bdcNameBDCmap[name])
		}
	}
	return &v1alpha1.BlockDeviceClaimList{
		Items: items,
	}, nil
}
