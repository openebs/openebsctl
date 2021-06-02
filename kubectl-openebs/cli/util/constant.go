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

const (
	// BytesToGB used to convert bytes to GB
	BytesToGB = 1073741824
	// BytesToMB used to convert bytes to MB
	BytesToMB = 1048567
	// BytesToKB used to convert bytes to KB
	BytesToKB = 1024
	// MicSec used to convert to microsec to second
	MicSec = 1000000
	// MinWidth used in tabwriter
	MinWidth = 0
	// MaxWidth used in tabwriter
	MaxWidth = 0
	// Padding used in tabwriter
	Padding = 4
	// OpenEBSCasTypeKey present in label of PV
	OpenEBSCasTypeKey = "openebs.io/cas-type"
	// Unknown to be retuned when cas type is not known
	Unknown = "unknown"
	// OpenEBSCasTypeKeySc present in parameter of SC
	OpenEBSCasTypeKeySc = "cas-type"
	// CstorCasType cas type name
	CstorCasType = "cstor"
	// CstorCasType cas type name
	JivaCasType = "jiva"
	// Healthy cstor volume status
	Healthy = "Healthy"
	// StorageKey key present in pvc status.capacity
	StorageKey = "storage"
)

var (
	// CasTypeAndComponentNameMap stores the component name of the corresponding cas type
	CasTypeAndComponentNameMap = map[string]string{
		"cstor": "openebs-cstor-csi-controller",
	}
	// ProvsionerAndCasTypeMap stores the cas type name of the corresponding provisioner
	ProvsionerAndCasTypeMap = map[string]string{
		"cstor.csi.openebs.io":         "cstor",
		"openebs.io/provisioner-iscsi": "jiva",
		"openebs.io/local":             "local",
		"local.csi.openebs.io ":        "localpv-lvm",
		"zfs.csi.openebs.io":           "localpv-zfs",
	}
)
