/*
Copyright 2020-2022 The OpenEBS Authors

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
	// ZFSCasType cas type name
	ZFSCasType = "localpv-zfs"
	// LVMCasType cas type name
	LVMCasType = "localpv-lvm"
	// LocalPvHostpathCasType cas type name
	LocalPvHostpathCasType = "localpv-hostpath"
	// LocalHostpathCasLabel cas-type label in dynamic-localpv-provisioner
	LocalHostpathCasLabel = "local-hostpath"
	// StorageKey key present in pvc status.capacity
	StorageKey = "storage"
	// NotAvailable shows something is missing, could be a component,
	// unknown version, or some other unknowns
	NotAvailable = "N/A"
)

const (
	// ZFSCSIDriver is the name of the ZFS localpv CSI driver
	ZFSCSIDriver = "zfs.csi.openebs.io"
	// LocalPVLVMCSIDriver is the name of the LVM LocalPV CSI driver
	// NOTE: This might also mean local-hostpath, local-device or zfs-localpv later.
	LocalPVLVMCSIDriver = "local.csi.openebs.io"
)

// Constant CSI component-name label values
const (
	// LVMLocalPVcsiControllerLabelValue is the label value of CSI controller STS & pod
	LVMLocalPVcsiControllerLabelValue = "openebs-lvm-controller"
	// ZFSLocalPVcsiControllerLabelValue is the label value of CSI controller STS & pod
	ZFSLocalPVcsiControllerLabelValue = "openebs-zfs-controller"
)

const (
	// LVMComponentNames for the lvm control plane components
	LVMComponentNames = "openebs-lvm-controller,openebs-lvm-node"
	// ZFSComponentNames for the zfs control plane components
	ZFSComponentNames = "openebs-zfs-controller,openebs-zfs-node"
	// HostpathComponentNames for the hostpath control plane components
	HostpathComponentNames = "openebs-localpv-provisioner"
)

var (
	// CasTypeAndComponentNameMap stores the component name of the corresponding cas type
	CasTypeAndComponentNameMap = map[string]string{
		LVMCasType:             LVMLocalPVcsiControllerLabelValue,
		ZFSCasType:             ZFSLocalPVcsiControllerLabelValue,
		LocalPvHostpathCasType: HostpathComponentNames,
	}
	// ComponentNameToCasTypeMap is a reverse map of CasTypeAndComponentNameMap
	ComponentNameToCasTypeMap = map[string]string{
		LVMLocalPVcsiControllerLabelValue: LVMCasType,
		ZFSLocalPVcsiControllerLabelValue: ZFSCasType,
		HostpathComponentNames:            LocalPvHostpathCasType,
	}
	// ProvsionerAndCasTypeMap stores the cas type name of the corresponding provisioner
	ProvsionerAndCasTypeMap = map[string]string{
		LocalPVLVMCSIDriver: LVMCasType,
		ZFSCSIDriver:        ZFSCasType,
	}
	// CasTypeToCSIProvisionerMap stores the provisioner of corresponding cas-types
	CasTypeToCSIProvisionerMap = map[string]string{
		LVMCasType: LocalPVLVMCSIDriver,
		ZFSCasType: ZFSCSIDriver,
	}

	// CasTypeToComponentNamesMap stores the names of the control-plane components of each cas-types.
	// To show statuses of new CasTypes, please update this map.
	CasTypeToComponentNamesMap = map[string]string{
		LocalPvHostpathCasType: HostpathComponentNames,
		ZFSCasType:             ZFSComponentNames,
		LVMCasType:             LVMComponentNames,
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
