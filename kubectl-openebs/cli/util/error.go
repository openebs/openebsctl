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

	if attachementStatus.Attached {
		return "Attached"
	}

	return attachementStatus.AttachError.Message
}
