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
	// CasType key in label of PV
	OPENEBS_CAS_TYPE_KEY = "openebs.io/cas-type"
	// Unknown to be retuned when cas type is not known
	UNKNOWN = "unknown"
	// CasType key in parameter of SC
	OPENBEBS_CAS_TYPE_KEY_SC = "cas-type"
	// cstor cas type
	CSTOR_CAS_TYPE = "cstor"
	// cstor volume, replica healthy status
	Healthy = "Healthy"
	// Total Storage key in pvc status.capacity
	STORAGE = "storage"
)
