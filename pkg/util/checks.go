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

package util

import (
	corev1 "k8s.io/api/core/v1"
)

// CheckForVol is used to check if the we can get the volume, if no volume attachment
// to SC for the corresponding volume is found display error
func CheckForVol(name string, vols map[string]*Volume) *Volume {
	if vol, found := vols[name]; found {
		return vol
	}
	// create & return an empty object to display details as Not Available
	errVol := &Volume{
		StorageClass:      NotAvailable,
		Node:              NotAvailable,
		AttachementStatus: NotAvailable,
		AccessMode:        NotAvailable,
	}
	return errVol
}

// AccessModeToString Flattens the arrat of AccessModes and returns a string fit to display in the output
func AccessModeToString(accessModeArray []corev1.PersistentVolumeAccessMode) string {
	accessModes := ""
	for _, mode := range accessModeArray {
		accessModes = accessModes + string(mode) + " "
	}
	return accessModes
}
