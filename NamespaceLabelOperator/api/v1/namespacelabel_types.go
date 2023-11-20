/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type LabelName string

type NamespacePhase string

type NamespaceConditionType string

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

const (
	// NamespaceDeletionDiscoveryFailure contains information about namespace deleter errors during resource discovery.
	NamespaceDeletionDiscoveryFailure NamespaceConditionType = "NamespaceDeletionDiscoveryFailure"
	// NamespaceDeletionContentFailure contains information about namespace deleter errors during deletion of resources.
	NamespaceDeletionContentFailure NamespaceConditionType = "NamespaceDeletionContentFailure"
	// NamespaceDeletionGVParsingFailure contains information about namespace deleter errors parsing GV for legacy types.
	NamespaceDeletionGVParsingFailure NamespaceConditionType = "NamespaceDeletionGroupVersionParsingFailure"
	// NamespaceContentRemaining contains information about resources remaining in a namespace.
	NamespaceContentRemaining NamespaceConditionType = "NamespaceContentRemaining"
	// NamespaceFinalizersRemaining contains information about which finalizers are on resources remaining in a namespace.
	NamespaceFinalizersRemaining NamespaceConditionType = "NamespaceFinalizersRemaining"
)

const (
	// NamespaceActive means the namespace is available for use in the system
	NamespaceActive NamespacePhase = "Active"
	// NamespaceTerminating means the namespace is undergoing graceful termination
	NamespaceTerminating NamespacePhase = "Terminating"
)

const (
	FinalizerKubernetes LabelName = "kubernetes"
)

type NamespaceCondition struct {
	// Type of namespace controller condition.
	Type NamespaceConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=NamespaceConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// NamespaceLabelSpec defines the desired state of NamespaceLabel
type NamespaceLabelSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of NamespaceLabel. Edit namespacelabel_types.go to remove/update
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,1,rep,name=labels,casttype=FinalizerName"`
}

// NamespaceLabelStatus defines the observed state of NamespaceLabel
type NamespaceLabelStatus struct {
	// Phase is the current lifecycle phase of the namespace.
	// More info: https://kubernetes.io/docs/tasks/administer-cluster/namespaces/
	// +optional
	Phase NamespacePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=NamespacePhase"`

	// Represents the latest available observations of a namespace's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []NamespaceCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// NamespaceLabel is the Schema for the namespacelabels API
type NamespaceLabel struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec defines the behavior of the Namespace.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec NamespaceLabelSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status describes the current status of a Namespace.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status NamespaceLabelStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

//+kubebuilder:object:root=true

// NamespaceLabelList contains a list of NamespaceLabel
type NamespaceLabelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceLabel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NamespaceLabel{}, &NamespaceLabelList{})
}
