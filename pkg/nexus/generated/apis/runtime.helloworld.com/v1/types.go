// Code generated by nexus. DO NOT EDIT.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:openapi-gen=true
type Child struct {
	Group string `json:"group" yaml:"group"`
	Kind  string `json:"kind" yaml:"kind"`
	Name  string `json:"name" yaml:"name"`
}

// +k8s:openapi-gen=true
type Link struct {
	Group string `json:"group" yaml:"group"`
	Kind  string `json:"kind" yaml:"kind"`
	Name  string `json:"name" yaml:"name"`
}

/* ------------------- CRDs definitions ------------------- */

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Runtime struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              RuntimeSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// +k8s:openapi-gen=true
type RuntimeSpec struct {
	MyRuntimeInt int `json:"myRuntimeInt" yaml:"myRuntimeInt"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RuntimeList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Runtime `json:"items" yaml:"items"`
}
