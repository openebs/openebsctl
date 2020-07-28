package util

import (
	v1 "github.com/openebs/api/pkg/apis/cstor/v1"

	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
)

// CheckVersion returns a message based on the status of the version
func CheckVersion(versionDetail v1.VersionDetails) string {

	if string(versionDetail.Status.State) == "Reconciled" || string(versionDetail.Status.State) == "" {
		return versionDetail.Status.Current
	}

	return string(versionDetail.Status.State) + ", desired version " + versionDetail.Desired
}

// CheckIfAccessable is used to check if the we can get the spec for volume
func CheckIfAccessable(attachment storagev1.VolumeAttachment) []corev1.PersistentVolumeAccessMode {

	if attachment.Status.Attached == true {
		return attachment.Spec.Source.InlineVolumeSpec.AccessModes
	}
	return make([]corev1.PersistentVolumeAccessMode, 0)
}

// CheckForVol is used to check if the we can get the volume, if no volume attachment
// to SC for the corresponding volume is found display error
func CheckForVol(name string, vols map[string]*Volume) *Volume {

	if vol, found := vols[name]; found {
		return vol
	}

	// create & return an empty object to display details as Not Available
	errVol := &Volume{
		StorageClass:      "N/A",
		Node:              "N/A",
		AttachementStatus: "N/A",
		AccessMode:        "N/A",
	}

	return errVol
}

//AccessModeToString Flattens the arrat of AccessModes and returns a string fit to display in the output
func AccessModeToString(accessModeArray []corev1.PersistentVolumeAccessMode) string {

	accessModes := ""

	for _, mode := range accessModeArray {
		accessModes = accessModes + string(mode) + " "
	}
	return accessModes
}
