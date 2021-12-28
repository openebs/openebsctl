/*
Copyright 2020-2021 The OpenEBS Authors

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
	TypeMeta: metav1.TypeMeta{Kind: "Node", APIVersion: "v1"}, ObjectMeta: metav1.ObjectMeta{Name: "node1", Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Status: corev1.NodeStatus{Phase: corev1.NodeRunning}}

var node2 = corev1.Node{
	TypeMeta: metav1.TypeMeta{Kind: "Node", APIVersion: "v1"}, ObjectMeta: metav1.ObjectMeta{Name: "node2", Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Status: corev1.NodeStatus{Phase: corev1.NodeRunning}}

var node3 = corev1.Node{
	TypeMeta: metav1.TypeMeta{Kind: "Node", APIVersion: "v1"}, ObjectMeta: metav1.ObjectMeta{Name: "node3", Labels: map[string]string{"kubernetes.io/hostname": "node3"}},
	Status: corev1.NodeStatus{Phase: corev1.NodeRunning}}

var activeBDwEXT4 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec:   v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "ext4", Mountpoint: "/dev/sda"}},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

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

var goodBD3N1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd3-n1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd3n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

var goodBD4N1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd4-n1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd4n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

var goodBD5N1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd5-n1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd5n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

var goodBD6N1 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd6-n1", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node1"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd6n1"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
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

var goodBD3N2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd3-n2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd3n2"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

var goodBD4N2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd4-n2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd4n2"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

var goodBD5N2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd5-n2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd5n2"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
		Path: "/dev/sda"},
	Status: v1alpha1.DeviceStatus{ClaimState: v1alpha1.BlockDeviceUnclaimed, State: v1alpha1.BlockDeviceActive}}

var goodBD6N2 = v1alpha1.BlockDevice{
	TypeMeta: metav1.TypeMeta{Kind: "Blockdevice", APIVersion: "openebs.io/v1alpha1"},
	ObjectMeta: metav1.ObjectMeta{Name: "bd6-n2", Namespace: "openebs",
		Labels: map[string]string{"kubernetes.io/hostname": "node2"}},
	Spec: v1alpha1.DeviceSpec{FileSystem: v1alpha1.FileSystemInfo{Type: "", Mountpoint: "/mnt/bd6n2"}, Capacity: v1alpha1.DeviceCapacity{Storage: 1074000000},
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

var mirrorCSPC = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{GenerateName: "cstor", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
			DataRaidGroups: []cstorv1.RaidGroup{{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{
				BlockDeviceName: "bd1-n1"}, {BlockDeviceName: "bd2-n1"}}}}, PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
			DataRaidGroups: []cstorv1.RaidGroup{{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{
				BlockDeviceName: "bd1-n2"}, {BlockDeviceName: "bd2-n2"}}}}, PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node3"},
			DataRaidGroups: []cstorv1.RaidGroup{{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{
				BlockDeviceName: "bd1-n3"}, {BlockDeviceName: "bd2-n3"}}}}, PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}}}}}

var mirrorCSPCFourBDs = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{GenerateName: "cstor", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
			DataRaidGroups: []cstorv1.RaidGroup{{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n1"}, {BlockDeviceName: "bd2-n1"}}},
				{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd3-n1"}, {BlockDeviceName: "bd4-n1"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
			DataRaidGroups: []cstorv1.RaidGroup{{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n2"}, {BlockDeviceName: "bd2-n2"}}},
				{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd3-n2"}, {BlockDeviceName: "bd4-n2"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolMirrored)}}}}}

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
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n1
    nodeSelector:
      kubernetes.io/hostname: node1
    poolConfig:
      dataRaidGroupType: mirror
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n2
    nodeSelector:
      kubernetes.io/hostname: node2
    poolConfig:
      dataRaidGroupType: mirror
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n3
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n3
    nodeSelector:
      kubernetes.io/hostname: node3
    poolConfig:
      dataRaidGroupType: mirror

`
var raidzCSPCThreeBDTwoNode = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{GenerateName: "cstor", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
			DataRaidGroups: []cstorv1.RaidGroup{
				{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n1"}, {BlockDeviceName: "bd2-n1"}, {BlockDeviceName: "bd3-n1"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolRaidz)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
			DataRaidGroups: []cstorv1.RaidGroup{
				{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n2"}, {BlockDeviceName: "bd2-n2"}, {BlockDeviceName: "bd3-n2"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolRaidz)}}}}}

var raidzCSPCstr = `apiVersion: cstor.openebs.io/v1
kind: CStorPoolCluster
metadata:
  creationTimestamp: null
  generateName: cstor
  namespace: openebs
spec:
  pools:
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd3-n1
    nodeSelector:
      kubernetes.io/hostname: node1
    poolConfig:
      dataRaidGroupType: raidz
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd3-n2
    nodeSelector:
      kubernetes.io/hostname: node2
    poolConfig:
      dataRaidGroupType: raidz

`
var raidz2CSPCSixBDTwoNode = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{GenerateName: "cstor", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
			DataRaidGroups: []cstorv1.RaidGroup{
				{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n1"},
					{BlockDeviceName: "bd2-n1"}, {BlockDeviceName: "bd3-n1"}, {BlockDeviceName: "bd4-n1"}, {BlockDeviceName: "bd5-n1"}, {BlockDeviceName: "bd6-n1"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolRaidz2)}},
		{NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
			DataRaidGroups: []cstorv1.RaidGroup{
				{CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n2"},
					{BlockDeviceName: "bd2-n2"}, {BlockDeviceName: "bd3-n2"}, {BlockDeviceName: "bd4-n2"}, {BlockDeviceName: "bd5-n2"}, {BlockDeviceName: "bd6-n2"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolRaidz2)}}}}}

var raidz2CSPCstr = `apiVersion: cstor.openebs.io/v1
kind: CStorPoolCluster
metadata:
  creationTimestamp: null
  generateName: cstor
  namespace: openebs
spec:
  pools:
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd3-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd4-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd5-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd6-n1
    nodeSelector:
      kubernetes.io/hostname: node1
    poolConfig:
      dataRaidGroupType: raidz2
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd3-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd4-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd5-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd6-n2
    nodeSelector:
      kubernetes.io/hostname: node2
    poolConfig:
      dataRaidGroupType: raidz2

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
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1
    nodeSelector:
      kubernetes.io/hostname: node1
    poolConfig:
      dataRaidGroupType: stripe

`
var StripeThreeNodeTwoDev = `apiVersion: cstor.openebs.io/v1
kind: CStorPoolCluster
metadata:
  creationTimestamp: null
  generateName: cstor
  namespace: openebs
spec:
  pools:
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n1
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n1
    nodeSelector:
      kubernetes.io/hostname: node1
    poolConfig:
      dataRaidGroupType: stripe
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n2
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n2
    nodeSelector:
      kubernetes.io/hostname: node2
    poolConfig:
      dataRaidGroupType: stripe
  - dataRaidGroups:
    - blockDevices:
      # /dev/sda  1.0GiB
      - blockDeviceName: bd1-n3
      # /dev/sda  1.0GiB
      - blockDeviceName: bd2-n3
    nodeSelector:
      kubernetes.io/hostname: node3
    poolConfig:
      dataRaidGroupType: stripe

`
var threeNodeTwoDevCSPC = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{GenerateName: "cstor", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{{
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n1"}, {BlockDeviceName: "bd2-n1"}}}},
		PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolStriped)}},
		{
			NodeSelector: map[string]string{"kubernetes.io/hostname": "node2"},
			DataRaidGroups: []cstorv1.RaidGroup{{
				CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n2"}, {BlockDeviceName: "bd2-n2"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolStriped)}},
		{
			NodeSelector: map[string]string{"kubernetes.io/hostname": "node3"},
			DataRaidGroups: []cstorv1.RaidGroup{{
				CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1-n3"}, {BlockDeviceName: "bd2-n3"}}}},
			PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolStriped)}}}},
}

var cspc1Struct = cstorv1.CStorPoolCluster{
	TypeMeta:   metav1.TypeMeta{Kind: "CStorPoolCluster", APIVersion: "cstor.openebs.io/v1"},
	ObjectMeta: metav1.ObjectMeta{GenerateName: "cstor", Namespace: "openebs"},
	Spec: cstorv1.CStorPoolClusterSpec{Pools: []cstorv1.PoolSpec{{
		NodeSelector: map[string]string{"kubernetes.io/hostname": "node1"},
		DataRaidGroups: []cstorv1.RaidGroup{{
			CStorPoolInstanceBlockDevices: []cstorv1.CStorPoolInstanceBlockDevice{{BlockDeviceName: "bd1"}}}},
		PoolConfig: cstorv1.PoolConfig{DataRaidGroupType: string(cstorv1.PoolStriped)}}}},
}
