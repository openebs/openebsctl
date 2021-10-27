package generate

import (
	cstorv1 "github.com/openebs/api/v2/pkg/apis/cstor/v1"
	"github.com/openebs/api/v2/pkg/apis/openebs.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var cstorCSIpod = corev1.Pod{
	TypeMeta: metav1.TypeMeta{Kind: "Pod", APIVersion: "v1"},
	ObjectMeta: metav1.ObjectMeta{Name: "fake-cstor-CSI", Namespace: "openebs",
		Labels: map[string]string{"openebs.io/version": "1.9.0", "openebs.io/component-name": "openebs-cstor-csi-controller"}},
	Status: corev1.PodStatus{Phase: corev1.PodRunning},
}

var node1 = corev1.Node{
	TypeMeta: metav1.TypeMeta{Kind: "Node", APIVersion: "v1"}, ObjectMeta: metav1.ObjectMeta{Name: "node1"},
	Status: corev1.NodeStatus{Phase: corev1.NodeRunning},
}

var node2 = corev1.Node{
	TypeMeta: metav1.TypeMeta{Kind: "Node", APIVersion: "v1"}, ObjectMeta: metav1.ObjectMeta{Name: "node2"},
	Status: corev1.NodeStatus{Phase: corev1.NodeRunning},
}

var node3 = corev1.Node{
	TypeMeta: metav1.TypeMeta{Kind: "Node", APIVersion: "v1"}, ObjectMeta: metav1.ObjectMeta{Name: "node3"},
	Status: corev1.NodeStatus{Phase: corev1.NodeRunning}}

var activeBDwEXT4 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec:   v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "ext4", Mountpoint: "/dev/sda"}},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive},
}

var inactiveBDwEXT4 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1-inactive", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec:   v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "ext4", Mountpoint: "/dev/sda"}, Capacity: v1alpha1.DeviceCapacity{Storage: 6711000}},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceInactive}}

var activeUnclaimedUnforattedBD = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/dev/sda"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

var goodBD1N1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1-n1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd1n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var goodBD2N1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd2-n1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd2n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var goodBD1N2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1-n2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/dev/sda"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var goodBD2N2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd2-n2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/dev/sda"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var goodBD1N3 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1-n3", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node3"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/dev/sdc"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var goodBD2N3 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd2-n3", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node3"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/dev/sdc"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}
var _ = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec:   v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/dev/sda"}, Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceInactive},
}

var mirrorCSPC = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{GenerateName: "cstor", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
			DataRaidGroups: []cstorv1.RaidGroup{{[]cstorv1.CStorPoolInstanceBlockDevice{{
				BlockDeviceName: "bd1-n1"}, {BlockDeviceName: "bd2-n1"}}}}, PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
			DataRaidGroups: []cstorv1.RaidGroup{{[]cstorv1.CStorPoolInstanceBlockDevice{{
				BlockDeviceName: "bd1-n2"}, {BlockDeviceName: "bd2-n2"}}}}, PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node3"},
			DataRaidGroups: []cstorv1.RaidGroup{{[]cstorv1.CStorPoolInstanceBlockDevice{{
				BlockDeviceName: "bd1-n3"}, {BlockDeviceName: "bd2-n3"}}}}, PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}}}}}

var mirrorCSPCstr = `apiVersion: cstor.openebs.io/v1
kind: CStorPoolCluster
metadata:
  creationTimestamp: null
  generateName: cstor
  namespace: openebs
spec:
  pools:
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1GB
      - blockDeviceName: bd1-n2
      # /dev/sda  1GB
      - blockDeviceName: bd2-n2
    nodeSelector:
      kubernetes.io/hostname: node2
    poolConfig:
      dataRaidGroupType: mirror
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1GB
      - blockDeviceName: bd1-n3
      # /dev/sda  1GB
      - blockDeviceName: bd2-n3
    nodeSelector:
      kubernetes.io/hostname: node3
    poolConfig:
      dataRaidGroupType: mirror
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1GB
      - blockDeviceName: bd1-n1
      # /dev/sda  1GB
      - blockDeviceName: bd2-n1
    nodeSelector:
      kubernetes.io/hostname: node1
    poolConfig:
      dataRaidGroupType: mirror
`

var cspc1 = `apiVersion: cstor.openebs.io/v1
kind: CStorPoolCluster
metadata:
  creationTimestamp: null
  generateName: cstor
  namespace: openebs
spec:
  pools:
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1GB
      - blockDeviceName: bd1
    nodeSelector:
      kubernetes.io/hostname: node1
    poolConfig:
      dataRaidGroupType: stripe

`

var cspc1Struct = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{GenerateName: "cstor", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{{
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			[]cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1"}}}},
		PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolStriped)}}}},
}
