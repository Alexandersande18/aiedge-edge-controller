/*
Copyright 2022.

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

// EdgeSpec defines the desired state of Edge
type EdgeSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	EdgeName string   `json:"edgeName,omitempty"`
	BaseArch BaseArch   `json:"baseArch,omitempty"`
	NodePortIP nodePortIP `json:"nodePortIp,omitempty"`
	ImageRegistry imageRegistry `json:"imageRegistry,omitempty"`
	Nodes    []string `json:"nodes,omitempty"`
}

type imageRegistry string

type nodePortIP string 

// +kubebuilder:validation:Enum=amd64;arm64
type BaseArch string

type NodeName string 


// EdgeStatus defines the observed state of Edge
type EdgeStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AvailableNodes int32 `json:"availableNodes"`
	EdgeSize       int32 `json:"edgeSize"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="basearch",type=string,JSONPath=`.spec.baseArch`
//+kubebuilder:printcolumn:name="size",type="integer",JSONPath=`.status.edgeSize`
type Edge struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EdgeSpec   `json:"spec,omitempty"`
	Status EdgeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EdgeList contains a list of Edge
type EdgeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Edge `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Edge{}, &EdgeList{})
}
