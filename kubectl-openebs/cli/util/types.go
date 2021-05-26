/*
Copyright 2020 The OpenEBS Authors

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
	v1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	corev1 "k8s.io/api/core/v1"
)

//Volume struct will have all the details we want to give in the output for
// openebsctl commands
type Volume struct {
	// AccessModes contains all ways the volume can be mounted
	AccessMode string
	// Attachment status of the PV and it's claim
	AttachementStatus string
	// Represents the actual capacity of the underlying volume.
	Capacity string
	// CStorPoolCluster that this volume belongs to
	CSPC string
	// The unique volume name returned by the CSI volume plugin to
	// refer to the volume on all subsequent calls.
	CSIVolumeAttachmentName string
	Name                    string
	//Namespace defines the space within each name must be unique.
	// An empty namespace is equivalent to the "default" namespace
	Namespace string
	Node      string
	// Name of the PVClaim of the underlying Persistent Volume
	PVC string
	// Status of the CStor Volume
	Status v1.CStorVolumePhase
	// Name of StorageClass to which this persistent volume belongs.
	StorageClass string
	// will be cStorVolume for all cStor volumes
	VolType string
	// version of the spec used to create the volumes
	Version string
}

//VolumeInfo struct will have all the details we want to give in the output for
// openebsctl command volume describe
type VolumeInfo struct {
	AccessMode string
	// Capacity of the underlying PV
	Capacity string
	// CStorPoolCluster that the volume belongs to
	CSPC string
	// cStor Instance Driver
	CSIDriver               string
	CSIVolumeAttachmentName string
	// Name of the volume & Namespace on which it exists
	Name      string
	Namespace string
	// Name of the underlying PVC
	PVC string
	// ReplicationFactor represents number of volume replica created during
	// volume provisioning connect to the target
	ReplicaCount int
	// Phase indicates if a volume is available, bound to a claim, or released
	// by a claim.
	VolumePhase corev1.PersistentVolumePhase
	// Name of StorageClass to which this persistent volume belongs.
	StorageClass string
	// Version of the OpenEBS resource definition being used
	Version string
	Size    string
	// Status of the CStor volume
	Status v1.CStorVolumePhase
}

// PortalInfo keep info about the ISCSI Target Portal.
type PortalInfo struct {
	// Target iSCSI Qualified Name.combination of nodeBase
	IQN        string
	VolumeName string
	// iSCSI Target Portal. The Portal is combination of IP:port
	// (typically TCP ports 3260)
	Portal string
	// TargetIP IP of the iSCSI target service
	TargetIP string
	//Node Name on which the application pod is running
	TargetNodeName string
}

// CStorReplicaInfo holds information about the cStor replicas
type CStorReplicaInfo struct {
	// Replica name present on ObjectMetadata
	Name string
	// Node on which the replica is present
	NodeName string
	ID       v1.ReplicaID
	//Replica Status reflects the phase, i.e hold result of last action.
	// ec. Healthy, Offline ,Degraded etc.
	Status string
}

// CstorPVCInfo struct will have all the details we want to give in the output for describe pvc
// details section for cstor pvc
type CstorPVCInfo struct {
	Name             string
	Namespace        string
	CasType          string
	BoundVolume      string
	AttachedToNode   string
	Pool             string
	StorageClassName string
	Size             string
	Used             string
	PVStatus         v1.CStorVolumePhase
}

// CVCInfo struct will have all the details we want to give in the output for describe pvc
// cvc section for cstor pvc
type CVCInfo struct {
	Name         string
	ReplicaCount int
	PoolInfo     []string
	Version      string
	Upgrading    bool
}

// PVCInfo struct will have all the details we want to give in the output for describe pvc
// details section for non-cstor pvc
type PVCInfo struct {
	Name             string
	Namespace        string
	CasType          string
	BoundVolume      string
	StorageClassName string
	Size             string
	PVStatus         corev1.PersistentVolumePhase
}
