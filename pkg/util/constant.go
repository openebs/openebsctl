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
	// OpenEBSCasTypeKey present in label of PV
	OpenEBSCasTypeKey = "openebs.io/cas-type"
	// Unknown to be retuned when cas type is not known
	Unknown = "unknown"
	// OpenEBSCasTypeKeySc present in parameter of SC
	OpenEBSCasTypeKeySc = "cas-type"
	// CstorCasType cas type name
	CstorCasType = "cstor"
	// ZFSCasType cas type name
	ZFSCasType = "localpv-zfs"
	// JivaCasType is the cas type name for Jiva
	JivaCasType = "jiva"
	// LVMCasType cas type name
	LVMCasType = "localpv-lvm"
	// LocalHostpathCasType cas type name
	LocalHostpathCasType = "localpv-hostpath"
	// LocalDeviceCasType cas type name
	LocalDeviceCasType = "localpv-device"
	// Healthy cstor volume status
	Healthy = "Healthy"
	// StorageKey key present in pvc status.capacity
	StorageKey = "storage"
	//NotAttached to show when CVA is not present
	NotAttached = "N/A"
	// CVAVolnameKey present in label of CVA
	CVAVolnameKey = "Volname"
	// UnicodeCross stores the character representation of U+2718
	UnicodeCross = "✘"
	// UnicodeCheck stores the character representation of U+2714
	UnicodeCheck = "✔"
	// NotFound stores the Not Found Status
	NotFound = "Not Found"
	// CVANotAttached stores CVA Not Attached status
	CVANotAttached = "Not Attached to Application"
	// Attached stores CVA Attached Status
	Attached = "Attached"
)

const (
	// CStorCSIDriver is the name of CStor CSI driver
	CStorCSIDriver = "cstor.csi.openebs.io"
	// JivaCSIDriver is the name of the Jiva CSI driver
	JivaCSIDriver = "jiva.csi.openebs.io"
	// ZFSCSIDriver is the name of the ZFS localpv CSI driver
	ZFSCSIDriver = "zfs.csi.openebs.io"
	// LocalPVLVMCSIDriver is the name of the LVM LocalPV CSI driver
	// NOTE: This might also mean local-hostpath, local-device or zfs-localpv later.
	LocalPVLVMCSIDriver = "local.csi.openebs.io"
)

// Constant CSI component-name label values
const (
	// CStorCSIControllerLabelValue is the label value of CSI controller STS & pod
	CStorCSIControllerLabelValue = "openebs-cstor-csi-controller"
	// JivaCSIControllerLabelValue is the label value of CSI controller STS & pod
	JivaCSIControllerLabelValue = "openebs-jiva-csi-controller"
	// LVMLocalPVcsiControllerLabelValue is the label value of CSI controller STS & pod
	LVMLocalPVcsiControllerLabelValue = "openebs-lvm-controller"
	// ZFSLocalPVcsiControllerLabelValue is the label value of CSI controller STS & pod
	ZFSLocalPVcsiControllerLabelValue = "openebs-zfs-controller"
)

const (
	// CstorComponentNames for the cstor control plane components
	CstorComponentNames = "cspc-operator,cvc-operator,cstor-admission-webhook,openebs-cstor-csi-node,openebs-cstor-csi-controller"
	// NDMComponentNames for the ndm components
	NDMComponentNames = "openebs-ndm-operator,ndm"
	// JivaComponentNames for the jiva control plane components
	JivaComponentNames = "openebs-jiva-csi-node,openebs-jiva-csi-controller,jiva-operator"
	// LVMComponentNames for the lvm control plane components
	LVMComponentNames = "openebs-lvm-controller,openebs-lvm-node"
	// ZFSComponentNames for the zfs control plane components
	ZFSComponentNames = "openebs-zfs-controller,openebs-zfs-node"
	// HostpathComponentNames for the hostpath control plane components
	HostpathComponentNames = "openebs-localpv-provisioner"
)

var (
	// CasTypeAndComponentNameMap stores the component name of the corresponding cas type
	// NOTE: Not including ZFSLocalPV as it'd break existing code
	CasTypeAndComponentNameMap = map[string]string{
		CstorCasType: CStorCSIControllerLabelValue,
		JivaCasType:  JivaCSIControllerLabelValue,
		LVMCasType:   LVMLocalPVcsiControllerLabelValue,
		ZFSCasType:   ZFSLocalPVcsiControllerLabelValue,
	}
	// ComponentNameToCasTypeMap is a reverse map of CasTypeAndComponentNameMap
	// NOTE: Not including ZFSLocalPV as it'd break existing code
	ComponentNameToCasTypeMap = map[string]string{
		CStorCSIControllerLabelValue:      CstorCasType,
		JivaCSIControllerLabelValue:       JivaCasType,
		LVMLocalPVcsiControllerLabelValue: LVMCasType,
		ZFSLocalPVcsiControllerLabelValue: ZFSCasType,
	}
	// ProvsionerAndCasTypeMap stores the cas type name of the corresponding provisioner
	ProvsionerAndCasTypeMap = map[string]string{
		CStorCSIDriver: CstorCasType,
		JivaCSIDriver:  JivaCasType,
		// NOTE: In near future this might mean all local-pv volumes
		LocalPVLVMCSIDriver: LVMCasType,
		ZFSCSIDriver:        ZFSCasType,
	}

	// CasTypeToComponentNamesMap stores the names of the control-plane components of each cas-types.
	// To show statuses of new CasTypes, please update this map.
	CasTypeToComponentNamesMap = map[string]string{
		CstorCasType:         CstorComponentNames + "," + NDMComponentNames,
		JivaCasType:          JivaComponentNames + "," + HostpathComponentNames,
		LocalHostpathCasType: HostpathComponentNames,
		LocalDeviceCasType:   HostpathComponentNames + "," + NDMComponentNames,
		ZFSCasType:           ZFSComponentNames,
		LVMCasType:           LVMComponentNames,
	}

	// CstorReplicaColumnDefinations stores the Table headers for CVR Details
	CstorReplicaColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Total", Type: "string"},
		{Name: "Used", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Age", Type: "string"},
	}
	// PodDetailsColumnDefinations stores the Table headers for Pod Details
	PodDetailsColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Namespace", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Ready", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Age", Type: "string"},
		{Name: "IP", Type: "string"},
		{Name: "Node", Type: "string"},
	}
	// JivaPodDetailsColumnDefinations stores the Table headers for Jiva Pod Details
	JivaPodDetailsColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Namespace", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Mode", Type: "string"},
		{Name: "Node", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "IP", Type: "string"},
		{Name: "Ready", Type: "string"},
		{Name: "Age", Type: "string"},
	}
	// VolumeListColumnDefinations stores the Table headers for Volume Details
	VolumeListColumnDefinations = []metav1.TableColumnDefinition{
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
	// CstorBackupColumnDefinations stores the Table headers for Cstor Backup Details
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
	// BDTreeListColumnDefinations stores the Table headers for Block Device Details, when displayed as tree
	BDTreeListColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Path", Type: "string"},
		{Name: "Size", Type: "string"},
		{Name: "ClaimState", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "FsType", Type: "string"},
		{Name: "MountPoint", Type: "string"},
	}
	// LVMvolgroupListColumnDefinitions stores the table headers for listing lvm vg-group when displayed as tree
	LVMvolgroupListColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "FreeSize", Type: "string"},
		{Name: "TotalSize", Type: "string"},
	}
	// ZFSPoolListColumnDefinitions stores the table headers for listing zfs pools when displayed as tree
	ZFSPoolListColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "FreeSize", Type: "string"},
	}

	// JivaReplicaPVCColumnDefinations stores the Table headers for Jiva Replica PVC details
	JivaReplicaPVCColumnDefinations = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Status", Type: "string"},
		{Name: "Volume", Type: "string"},
		{Name: "Capacity", Type: "string"},
		{Name: "Storageclass", Type: "string"},
		{Name: "Age", Type: "string"},
	}

	// CstorVolumeCRStatusColumnDefinitions stores the Table headers for Cstor CRs status details
	CstorVolumeCRStatusColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Kind", Type: "string"},
		{Name: "Name", Type: "string"},
		{Name: "Status", Type: "string"},
	}

	// VolumeTotalAndUsageDetailColumnDefinitions stores the Table headers for volume usage details
	VolumeTotalAndUsageDetailColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Total Capacity", Type: "string"},
		{Name: "Used Capacity", Type: "string"},
		{Name: "Available Capacity", Type: "string"},
	}
	// EventsColumnDefinitions stores the Table headers for events details
	EventsColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string"},
		{Name: "Action", Type: "string"},
		{Name: "Reason", Type: "string"},
		{Name: "Message", Type: "string"},
		{Name: "Type", Type: "string"},
	}

	VersionColumnDefinition = []metav1.TableColumnDefinition{
		{Name: "Component", Type: "string"},
		{Name: "Version", Type: "string"},
	}
	// ClusterInfoColumnDefinitions stores the Table headers for Cluster-Info details
	ClusterInfoColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Cas-Type", Type: "string"},
		{Name: "Namespace", Type: "string"},
		{Name: "Version", Type: "string"},
		{Name: "Working", Type: "string"},
		{Name: "Status", Type: "string"},
	}
)
