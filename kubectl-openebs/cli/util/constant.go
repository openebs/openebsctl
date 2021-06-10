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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
	// Healthy cstor volume status
	Healthy = "Healthy"
	// StorageKey key present in pvc status.capacity
	StorageKey = "storage"
	//NotAttached to show when CVA is not present
	NotAttached = "N/A"
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
	// CstorReplicaColumnDefinations stores the Table headers for CVR Details
	CstorReplicaColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Total", Type: "string"},
		{Name: "Used", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Age", Type: "string"},
	}
	// CstorTargetDetailsColumnDefinations stores the Table headers for Cstor Target Details
	CstorTargetDetailsColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Namespace", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Ready", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Age", Type: "string"},
		{Name: "IP", Type: "string"},
		{Name: "Node", Type: "string"},
	}
	// CstorVolumeListColumnDefinations stores the Table headers for Cstor Volume Details
	CstorVolumeListColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Namespace", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Version", Type: "string"},
		{Name: "Capacity", Type: "string"},
		{Name: "Storage Class", Type: "string"},
		{Name: "Attached", Type: "string"},
		{Name: "Access Mode", Type: "string"},
		{Name: "Attached Node", Type: "string"},
	}
	// CstorPoolListColumnDefinations stores the Table headers for Cstor Pool Details
	CstorPoolListColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "HostName", Type: "string"},
		{Name: "Free", Type: "string"},
		{Name: "Capacity", Type: "string"},
		{Name: "Read Only", Type: "bool"},
		{Name: "Provisioned Replicas", Type: "int"},
		{Name: "Healthy Replicas", Type: "int"},
		{Name: "Status", Type: "string"},
		{Name: "Age", Type: "string"},
	}
	// BDListColumnDefinations stores the Table headers for Block Device Details
	BDListColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Capacity", Type: "string"},
		{Name: "State", Type: "string"},
	}
	// PoolReplicaColumnDefinations stores the Table headers for Pool Replica Details
	PoolReplicaColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "PVC Name", Type: "string"},
		{Name: "Size", Type: "string"},
		{Name: "State", Type: "string"},
	}
	// CstorVolumeBackupColumnDefinations stores the Table headers for Cstor Backup Details
	CstorBackupColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Backup Name", Type: "string"},
		{Name: "Volume Name", Type: "string"},
		{Name: "Backup Destination", Type: "string"},
		{Name: "Snap Name", Type: "string"},
		{Name: "Status", Type: "string"},
	}
	// CstorCompletedBackupColumnDefinations stores the Table headers for Cstor Completed Backup Details
	CstorCompletedBackupColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Backup Name", Type: "string"},
		{Name: "Volume Name", Type: "string"},
		{Name: "Last Snap Name", Type: "string"},
	}
	// CstorRestoreColumnDefinations stores the Table headers for Cstor Restore Details
	CstorRestoreColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Restore Name", Type: "string"},
		{Name: "Volume Name", Type: "string"},
		{Name: "Restore Source", Type: "string"},
		{Name: "Storage Class", Type: "string"},
		{Name: "Status", Type: "string"},
	}
)
