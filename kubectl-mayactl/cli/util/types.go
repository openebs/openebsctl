package util

import (
	v1 "github.com/openebs/api/pkg/apis/cstor/v1"
	corev1 "k8s.io/api/core/v1"
)

//Volume struct will have all the details we want to give in the output for
// mayactl commands
type Volume struct {
	AccessMode              []corev1.PersistentVolumeAccessMode
	AttachementStatus       string
	Capacity                string
	CSPC                    string
	CSIVolumeAttachmentName string
	Name                    string
	Namespace               string
	Node                    string
	PVC                     string
	Status                  v1.CStorVolumePhase
	StorageClass            string
	VolType                 string
	Version                 string
}

// PortalInfo keep info about the ISCSI Target Portal.
type PortalInfo struct {
	IQN           string
	VolumeName    string
	Portal        string
	Size          string
	Status        []v1.CStorVolumeCondition
	ReplicaCount  int
	ReplicaStatus []v1.ReplicaStatus
	//ControllerNode string
}
