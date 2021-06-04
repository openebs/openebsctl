/*
Copyright © 2020-2021 The OpenEBS Authors

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

package client

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ISCSISpec struct {
	TargetIP   string `json:"targetIP,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
	Iqn        string `json:"iqn,omitempty"`
}

type MountInfo struct {
	// StagingPath is the path provided by K8s during NodeStageVolume
	// rpc call, where volume is mounted globally.
	StagingPath string `json:"stagingPath,omitempty"`
	// TargetPath is the path provided by K8s during NodePublishVolume
	// rpc call where bind mount happens.
	TargetPath string `json:"targetPath,omitempty"`
	FSType     string `json:"fsType,omitempty"`
	DevicePath string `json:"devicePath,omitempty"`
}

// JivaVolumeSpec defines the desired state of JivaVolume
// +k8s:openapi-gen=true
type JivaVolumeSpec struct {
	PV       string `json:"pv"`
	Capacity string `json:"capacity"`
	// AccessType can be specified as Block or Mount type
	AccessType string `json:"accessType"`
	// +nullable
	ISCSISpec ISCSISpec `json:"iscsiSpec,omitempty"`
	// +nullable
	MountInfo MountInfo `json:"mountInfo,omitempty"`
	// Policy is the configuration used for creating target
	// and replica pods during volume provisioning
	// +nullable
	Policy                   JivaVolumePolicySpec `json:"policy,omitempty"`
	DesiredReplicationFactor int                  `json:"desiredReplicationFactor,omitempty"`
}

// JivaVolumeStatus defines the observed state of JivaVolume
// +k8s:openapi-gen=true
type JivaVolumeStatus struct {
	Status       string `json:"status,omitempty"`
	ReplicaCount int    `json:"replicaCount,omitempty"`
	// +nullable
	ReplicaStatuses []ReplicaStatus `json:"replicaStatus,omitempty"`
	// Phase represents the current phase of JivaVolume.
	Phase JivaVolumePhase `json:"phase,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JivaVolume is the Schema for the jivavolumes API
// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced,shortName=jv
// +kubebuilder:printcolumn:name="ReplicaCount",type="string",JSONPath=`.status.replicaCount`
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=`.status.status`
type JivaVolume struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec           JivaVolumeSpec   `json:"spec,omitempty"`
	Status         JivaVolumeStatus `json:"status,omitempty"`
	VersionDetails VersionDetails   `json:"versionDetails,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// JivaVolumeList contains a list of JivaVolume
type JivaVolumeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JivaVolume `json:"items"`
}

// ReplicaStatus stores the status of replicas
type ReplicaStatus struct {
	Address string `json:"address,omitempty"`
	Mode    string `json:"mode,omitempty"`
}

// JivaVolumePhase represents the current phase of JivaVolume.
type JivaVolumePhase string

const (
	// JivaVolumePhasePending indicates that the jivavolume is still waiting for
	// the jivavolume to be created and bound
	JivaVolumePhasePending JivaVolumePhase = "Pending"

	// JivaVolumePhaseSyncing indicates that the jivavolume has been
	// provisioned and replicas are syncing
	JivaVolumePhaseSyncing JivaVolumePhase = "Syncing"

	// JivaVolumePhaseFailed indicates that the jivavolume provisioning
	// has failed
	JivaVolumePhaseFailed JivaVolumePhase = "Failed"

	// JivaVolumePhaseUnkown indicates that the jivavolume status get
	// failed as controller was not reachable
	JivaVolumePhaseUnkown JivaVolumePhase = "Unknown"

	// JivaVolumePhaseReady indicates that the jivavolume provisioning
	// has Created
	JivaVolumePhaseReady JivaVolumePhase = "Ready"

	// JivaVolumePhaseDeleting indicates the the jivavolume is deprovisioned
	JivaVolumePhaseDeleting JivaVolumePhase = "Deleting"
)

// JivaVolumePolicySpec defines the desired state of JivaVolumePolicy
type JivaVolumePolicySpec struct {
	// ReplicaSC represents the storage class used for
	// creating the pvc for the replicas (provisioned by localpv provisioner)
	ReplicaSC string `json:"replicaSC,omitempty"`
	// EnableBufio ...
	EnableBufio bool `json:"enableBufio"`
	// AutoScaling ...
	AutoScaling bool `json:"autoScaling"`
	// ServiceAccountName can be provided to enable PSP
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// PriorityClassName if specified applies to the pod
	// If left empty, no priority class is applied.
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// TargetSpec represents configuration related to jiva target and its resources
	// +nullable
	Target TargetSpec `json:"target,omitempty"`
	// ReplicaSpec represents configuration related to replicas resources
	// +nullable
	Replica ReplicaSpec `json:"replica,omitempty"`
}

// TargetSpec represents configuration related to jiva target deployment
type TargetSpec struct {
	// Monitor enables or disables the target exporter sidecar
	Monitor bool `json:"monitor,omitempty"`

	// ReplicationFactor represents maximum number of replicas
	// that are allowed to connect to the target
	ReplicationFactor int `json:"replicationFactor,omitempty"`

	// PodTemplateResources represents the configuration for target deployment.
	PodTemplateResources `json:",inline"`

	// AuxResources are the compute resources required by the jiva-target pod
	// side car containers.
	AuxResources *corev1.ResourceRequirements `json:"auxResources,omitempty"`
}

// ReplicaSpec represents configuration related to jiva replica sts
type ReplicaSpec struct {
	// PodTemplateResources represents the configuration for replica sts.
	PodTemplateResources `json:",inline"`
}

// PodTemplateResources represents the common configuration field for
// jiva target deployment and jiva replica sts.
type PodTemplateResources struct {
	// Resources are the compute resources required by the jiva
	// container.
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Tolerations, if specified, are the pod's tolerations
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// Affinity if specified, are the pod's affinities
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// NodeSelector is the labels that will be used to select
	// a node for pod scheduleing
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// JivaVolumePolicyStatus is for handling status of JivaVolumePolicy
type JivaVolumePolicyStatus struct {
	Phase string `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JivaVolumePolicy is the Schema for the jivavolumes API
// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Namespaced,shortName=jvp
type JivaVolumePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JivaVolumePolicySpec   `json:"spec,omitempty"`
	Status JivaVolumePolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JivaVolumePolicyList contains a list of JivaVolumePolicy
type JivaVolumePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JivaVolumePolicy `json:"items"`
}

// VersionDetails provides the details for upgrade
type VersionDetails struct {
	// If AutoUpgrade is set to true then the resource is
	// upgraded automatically without any manual steps
	AutoUpgrade bool `json:"autoUpgrade,omitempty"`
	// Desired is the version that we want to
	// upgrade or the control plane version
	Desired string `json:"desired,omitempty"`
	// Status gives the status of reconciliation triggered
	// when the desired and current version are not same
	Status VersionStatus `json:"status,omitempty"`
}

// VersionStatus is the status of the reconciliation of versions
type VersionStatus struct {
	// DependentsUpgraded gives the details whether all children
	// of a resource are upgraded to desired version or not
	DependentsUpgraded bool `json:"dependentsUpgraded,omitempty"`
	// Current is the version of resource
	Current string `json:"current,omitempty"`
	// State is the state of reconciliation
	State VersionState `json:"state,omitempty"`
	// Message is a human readable message if some error occurs
	Message string `json:"message,omitempty"`
	// Reason is the actual reason for the error state
	Reason string `json:"reason,omitempty"`
	// LastUpdateTime is the time the status was last  updated
	// +nullable
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
}

// VersionState is the state of reconciliation
type VersionState string

const (
	// ReconcileComplete is the state when desired and current version are equal.
	ReconcileComplete VersionState = "Reconciled"
	// ReconcileInProgress is the state when desired and current version are
	// not same and the reconcile functions is retrying to make them same.
	ReconcileInProgress VersionState = "ReconcileInProgress"
	// ReconcilePending is the state the reconciliation is still not started yet
	ReconcilePending VersionState = "ReconcilePending"
)

// SetErrorStatus sets the message and reason for the error
func (vs *VersionStatus) SetErrorStatus(msg string, err error) {
	vs.Message = msg
	vs.Reason = err.Error()
	vs.LastUpdateTime = metav1.Now()
}

// SetInProgressStatus sets the state as ReconcileInProgress
func (vs *VersionStatus) SetInProgressStatus() {
	vs.State = ReconcileInProgress
	vs.LastUpdateTime = metav1.Now()
}

// SetSuccessStatus resets the message and reason and sets the state as
// Reconciled
func (vd *VersionDetails) SetSuccessStatus() {
	vd.Status.Current = vd.Desired
	vd.Status.Message = ""
	vd.Status.Reason = ""
	vd.Status.State = ReconcileComplete
	vd.Status.LastUpdateTime = metav1.Now()
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ISCSISpec) DeepCopyInto(out *ISCSISpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ISCSISpec.
func (in *ISCSISpec) DeepCopy() *ISCSISpec {
	if in == nil {
		return nil
	}
	out := new(ISCSISpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JivaVolume) DeepCopyInto(out *JivaVolume) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	in.VersionDetails.DeepCopyInto(&out.VersionDetails)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JivaVolume.
func (in *JivaVolume) DeepCopy() *JivaVolume {
	if in == nil {
		return nil
	}
	out := new(JivaVolume)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JivaVolume) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JivaVolumeList) DeepCopyInto(out *JivaVolumeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]JivaVolume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JivaVolumeList.
func (in *JivaVolumeList) DeepCopy() *JivaVolumeList {
	if in == nil {
		return nil
	}
	out := new(JivaVolumeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JivaVolumeList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JivaVolumePolicy) DeepCopyInto(out *JivaVolumePolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JivaVolumePolicy.
func (in *JivaVolumePolicy) DeepCopy() *JivaVolumePolicy {
	if in == nil {
		return nil
	}
	out := new(JivaVolumePolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JivaVolumePolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JivaVolumePolicyList) DeepCopyInto(out *JivaVolumePolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]JivaVolumePolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JivaVolumePolicyList.
func (in *JivaVolumePolicyList) DeepCopy() *JivaVolumePolicyList {
	if in == nil {
		return nil
	}
	out := new(JivaVolumePolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *JivaVolumePolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JivaVolumePolicySpec) DeepCopyInto(out *JivaVolumePolicySpec) {
	*out = *in
	in.Target.DeepCopyInto(&out.Target)
	in.Replica.DeepCopyInto(&out.Replica)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JivaVolumePolicySpec.
func (in *JivaVolumePolicySpec) DeepCopy() *JivaVolumePolicySpec {
	if in == nil {
		return nil
	}
	out := new(JivaVolumePolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JivaVolumePolicyStatus) DeepCopyInto(out *JivaVolumePolicyStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JivaVolumePolicyStatus.
func (in *JivaVolumePolicyStatus) DeepCopy() *JivaVolumePolicyStatus {
	if in == nil {
		return nil
	}
	out := new(JivaVolumePolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JivaVolumeSpec) DeepCopyInto(out *JivaVolumeSpec) {
	*out = *in
	out.ISCSISpec = in.ISCSISpec
	out.MountInfo = in.MountInfo
	in.Policy.DeepCopyInto(&out.Policy)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JivaVolumeSpec.
func (in *JivaVolumeSpec) DeepCopy() *JivaVolumeSpec {
	if in == nil {
		return nil
	}
	out := new(JivaVolumeSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JivaVolumeStatus) DeepCopyInto(out *JivaVolumeStatus) {
	*out = *in
	if in.ReplicaStatuses != nil {
		in, out := &in.ReplicaStatuses, &out.ReplicaStatuses
		*out = make([]ReplicaStatus, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JivaVolumeStatus.
func (in *JivaVolumeStatus) DeepCopy() *JivaVolumeStatus {
	if in == nil {
		return nil
	}
	out := new(JivaVolumeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MountInfo) DeepCopyInto(out *MountInfo) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MountInfo.
func (in *MountInfo) DeepCopy() *MountInfo {
	if in == nil {
		return nil
	}
	out := new(MountInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodTemplateResources) DeepCopyInto(out *PodTemplateResources) {
	*out = *in
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodTemplateResources.
func (in *PodTemplateResources) DeepCopy() *PodTemplateResources {
	if in == nil {
		return nil
	}
	out := new(PodTemplateResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicaSpec) DeepCopyInto(out *ReplicaSpec) {
	*out = *in
	in.PodTemplateResources.DeepCopyInto(&out.PodTemplateResources)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicaSpec.
func (in *ReplicaSpec) DeepCopy() *ReplicaSpec {
	if in == nil {
		return nil
	}
	out := new(ReplicaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicaStatus) DeepCopyInto(out *ReplicaStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicaStatus.
func (in *ReplicaStatus) DeepCopy() *ReplicaStatus {
	if in == nil {
		return nil
	}
	out := new(ReplicaStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TargetSpec) DeepCopyInto(out *TargetSpec) {
	*out = *in
	in.PodTemplateResources.DeepCopyInto(&out.PodTemplateResources)
	if in.AuxResources != nil {
		in, out := &in.AuxResources, &out.AuxResources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TargetSpec.
func (in *TargetSpec) DeepCopy() *TargetSpec {
	if in == nil {
		return nil
	}
	out := new(TargetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VersionDetails) DeepCopyInto(out *VersionDetails) {
	*out = *in
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VersionDetails.
func (in *VersionDetails) DeepCopy() *VersionDetails {
	if in == nil {
		return nil
	}
	out := new(VersionDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VersionStatus) DeepCopyInto(out *VersionStatus) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VersionStatus.
func (in *VersionStatus) DeepCopy() *VersionStatus {
	if in == nil {
		return nil
	}
	out := new(VersionStatus)
	in.DeepCopyInto(out)
	return out
}
