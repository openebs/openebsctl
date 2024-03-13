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

// Volume struct will have all the details we want to give in the output for
// openebsctl commands
type Volume struct {
	// AccessModes contains all ways the volume can be mounted
	AccessMode string
	// Attachment status of the PV and it's claim
	AttachementStatus string
	// Represents the actual capacity of the underlying volume.
	Capacity string
	Name     string
	//Namespace defines the space within each name must be unique.
	// An empty namespace is equivalent to the "default" namespace
	Namespace string
	Node      string
	// Name of the PVClaim of the underlying Persistent Volume
	PVC string
	// Name of StorageClass to which this persistent volume belongs.
	StorageClass string
	// version of the spec used to create the volumes
	Version string
}

// VolumeInfo struct will have all the details we want to give in the output for
// openebsctl command volume describe
type VolumeInfo struct {
	AccessMode string
	// Capacity of the underlying PV
	Capacity string
	// Name of the volume & Namespace on which it exists
	Name string
	PVC  string
	// Phase indicates if a volume is available, bound to a claim, or released
	// by a claim.
	VolumePhase corev1.PersistentVolumePhase
	// Name of StorageClass to which this persistent volume belongs.
	StorageClass string
	Size         string
}

type LocalHostPathVolInfo struct {
	VolumeInfo
	Path          string
	ReclaimPolicy string
	CasType       string
}

// LVMPVCInfo struct will have all the details we want to give in the output for describe pvc
// details section for lvm pvc
type LVMPVCInfo struct {
	Name             string
	Namespace        string
	CasType          string
	BoundVolume      string
	StorageClassName string
	Size             string
	PVCStatus        corev1.PersistentVolumeClaimPhase
	MountPods        string
}

// ZFSPVCInfo struct will have all the details we want to give in the output for describe pvc
// details section for zfs pvc
type ZFSPVCInfo struct {
	Name             string
	Namespace        string
	CasType          string
	BoundVolume      string
	StorageClassName string
	Size             string
	PVCStatus        corev1.PersistentVolumeClaimPhase
	MountPods        string
}

// PVCInfo struct will have all the details we want to give in the output for describe pvc
// details section for generic pvc
type PVCInfo struct {
	Name             string
	Namespace        string
	CasType          string
	BoundVolume      string
	StorageClassName string
	Size             string
	PVStatus         corev1.PersistentVolumePhase
	MountPods        string
}

// MapOptions struct to get the resources as Map with the provided options
// Key defines what to use as a key, ex:- name, label, currently these two are supported, add more according to need.
// LabelKey defines which Label to use as key.
type MapOptions struct {
	Key      Key
	LabelKey string
}

// ReturnType defines in which format the object needs to be returned i.e. List or Map
type ReturnType string

// Key defines what should be the key if we create a map, i.e. Label or Name
type Key string

const (
	// List If we want the return type as a list
	List ReturnType = "list"
	// Map If we want the return type as a map
	Map ReturnType = "map"
	// Name key if we want the keys to be made on name
	Name Key = "name"
	// Label key if want to make the keys on labels
	Label Key = "label"
)

// ZFSVolDesc is the output helper for ZfsVolDesc
type ZFSVolDesc struct {
	Name         string
	Namespace    string
	AccessMode   string
	CSIDriver    string
	Capacity     string
	PVC          string
	VolumePhase  corev1.PersistentVolumePhase
	StorageClass string
	Version      string
	Status       string
	VolumeType   string
	PoolName     string
	FileSystem   string
	Compression  string
	Dedup        string
	NodeID       string
	Recordsize   string
}

// LVMVolDesc is the output helper for LVMVolDesc
type LVMVolDesc struct {
	Name            string
	Namespace       string
	AccessMode      string
	CSIDriver       string
	Capacity        string
	PVC             string
	VolumePhase     corev1.PersistentVolumePhase
	StorageClass    string
	Version         string
	Status          string
	VolumeGroup     string
	Shared          string
	ThinProvisioned string
	NodeID          string
}

// ComponentData stores the data for each component of an engine
type ComponentData struct {
	Namespace string
	Status    string
	Version   string
	CasType   string
}
