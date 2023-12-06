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

// PodFinderSpec defines the desired state of PodFinder
type PodFinderSpec struct {
	// Name of the pod, PodFinder is looking for
	Name string `json:"name,omitempty"`
}

// PodFinderStatus defines the observed state of PodFinder
type PodFinderStatus struct {
	Found bool `json:"found,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodFinder is the Schema for the podfinders API
type PodFinder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodFinderSpec   `json:"spec,omitempty"`
	Status PodFinderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodFinderList contains a list of PodFinder
type PodFinderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodFinder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodFinder{}, &PodFinderList{})
}
