package util

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/storage/v1"
)

// CheckError prints err to stderr and exits with code 1 if err is not nil. Otherwise, it is a
// no-op.
func CheckError(err error) {
	if err != nil {
		if err != context.Canceled {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("An error occurred: %v\n", err))
		}
		os.Exit(1)
	}
}

// CheckErr to handle command errors
func CheckErr(err error, handleErr func(string)) {
	if err == nil {
		return
	}
	handleErr(err.Error())
}

// CheckVolAttachmentError is used to check if the volume is correctly attached
// to the cspc
func CheckVolAttachmentError(attachementStatus v1.VolumeAttachmentStatus) string {

	if attachementStatus.Attached == true {
		return "Attached"
	}

	return attachementStatus.AttachError.Message
}
