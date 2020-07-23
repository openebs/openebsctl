package util

import (
	v1 "github.com/openebs/api/pkg/apis/cstor/v1"
	corev1 "k8s.io/api/core/v1"
)

//Volume struct will have all the details we want to give in the output for
// mayactl commands
type Volume struct {
	// AccessMode of the underlying PV
	AccessMode string
	// Attachment status of the PV and it's claim
	AttachementStatus string
	//Size of PV
	Capacity string
	//CStorPoolCluster that this volume belongs to
	CSPC string
	// The unique volume name returned by the CSI volume plugin to
	// refer to the volume on all subsequent calls.
	CSIVolumeAttachmentName string
	Name                    string
	Namespace               string
	Node                    string
	PVC                     string
	// Status of the CStor Volume
	Status       v1.CStorVolumePhase
	StorageClass string
	// will be cStorVolume for all cStor volumes
	VolType string
	// version of the spec used to create the volumes
	Version string
}

//VolumeInfo struct will have all the details we want to give in the output for
// mayactl command volume describe
type VolumeInfo struct {
	AccessMode string
	Capacity   string
	CSPC       string
	//cStor Instance Driver
	CSIDriver               string
	CSIVolumeAttachmentName string
	Name                    string
	Namespace               string
	PVC                     string
	//Number of replicas user has specified for thw cStorVolume
	ReplicaCount int
	VolumePhase  corev1.PersistentVolumePhase
	StorageClass string
	Version      string
	Size         string
	// Status of the CStor volume
	Status v1.CStorVolumePhase
}

// PortalInfo keep info about the ISCSI Target Portal.
type PortalInfo struct {
	//iSCSI qualified name to configure the target
	IQN        string
	VolumeName string
	Portal     string
	TargetIP   string
	//Node Name on which the application pod is running
	TargetNodeName string
}

// CStorReplicaInfo holds information about the cstor replicas
type CStorReplicaInfo struct {
	//replica name
	Name string
	//Node on ehoch it is presemt
	NodeName string
	ID       v1.ReplicaID
	// replica status
	Status string
}
