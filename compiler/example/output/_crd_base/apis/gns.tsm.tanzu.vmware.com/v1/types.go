// Code generated by nexus. DO NOT EDIT.

package v1

import (
	cartv1 "github.com/vmware-tanzu/cartographer/pkg/apis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"nexustempmodule/common"
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

// +k8s:openapi-gen=true
type NexusStatus struct {
	SourceGeneration int64 `json:"sourceGeneration" yaml:"sourceGeneration"`
	RemoteGeneration int64 `json:"remoteGeneration" yaml:"remoteGeneration"`
}

/* ------------------- CRDs definitions ------------------- */

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type RandomGnsData struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              RandomGnsDataSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            RandomGnsDataNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type RandomGnsDataNexusStatus struct {
	Status RandomStatus `json:"status,omitempty" yaml:"status,omitempty"`
	Nexus  NexusStatus  `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *RandomGnsData) CRDName() string {
	return "randomgnsdatas.gns.tsm.tanzu.vmware.com"
}

func (c *RandomGnsData) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type RandomGnsDataSpec struct {
	Description RandomDescription `json:"description" yaml:"description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RandomGnsDataList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []RandomGnsData `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Gns struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              GnsSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            GnsNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type GnsNexusStatus struct {
	State GnsState    `json:"state,omitempty" yaml:"state,omitempty"`
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Gns) CRDName() string {
	return "gnses.gns.tsm.tanzu.vmware.com"
}

func (c *Gns) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type GnsSpec struct {
	//nexus-validation: MaxLength=8, MinLength=2
	//nexus-validation: Pattern=abc
	Domain                    string                       `json:"domain" yaml:"domain"`
	UseSharedGateway          bool                         `json:"useSharedGateway" yaml:"useSharedGateway"`
	Description               Description                  `json:"description" yaml:"description"`
	Meta                      string                       `json:"meta" yaml:"meta"`
	Port                      *int                         `json:"port" yaml:"port"`
	OtherDescription          *Description                 `json:"otherDescription" yaml:"otherDescription"`
	MapPointer                *map[string]string           `json:"mapPointer" yaml:"mapPointer"`
	SlicePointer              *[]string                    `json:"slicePointer" yaml:"slicePointer"`
	WorkloadSpec              cartv1.WorkloadSpec          `json:"workloadSpec" yaml:"workloadSpec"`
	DifferentSpec             *cartv1.WorkloadSpec         `json:"differentSpec" yaml:"differentSpec"`
	ServiceSegmentRef         ServiceSegmentRef            `json:"serviceSegmentRef,omitempty"`
	ServiceSegmentRefPointer  *ServiceSegmentRef           `json:"serviceSegmentRefPointer,omitempty"`
	ServiceSegmentRefs        []ServiceSegmentRef          `json:"serviceSegmentRefs,omitempty"`
	ServiceSegmentRefMap      map[string]ServiceSegmentRef `json:"serviceSegmentRefMap,omitempty"`
	GnsServiceGroupsGvk       map[string]Child             `json:"gnsServiceGroupsGvk,omitempty" yaml:"gnsServiceGroupsGvk,omitempty" nexus:"children"`
	GnsAccessControlPolicyGvk *Child                       `json:"gnsAccessControlPolicyGvk,omitempty" yaml:"gnsAccessControlPolicyGvk,omitempty" nexus:"child"`
	FooChildGvk               *Child                       `nexus:"child" nexus-graphql:"type:string"`
	IgnoreChildGvk            *Child                       `nexus:"child" nexus-graphql:"ignore:true"`
	DnsGvk                    *Link                        `json:"dnsGvk,omitempty" yaml:"dnsGvk,omitempty" nexus:"link"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GnsList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Gns `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type BarChild struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              BarChildSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            BarChildNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type BarChildNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *BarChild) CRDName() string {
	return "barchilds.gns.tsm.tanzu.vmware.com"
}

func (c *BarChild) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type BarChildSpec struct {
	Name string `json:"name" yaml:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type BarChildList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []BarChild `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type IgnoreChild struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              IgnoreChildSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            IgnoreChildNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type IgnoreChildNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *IgnoreChild) CRDName() string {
	return "ignorechilds.gns.tsm.tanzu.vmware.com"
}

func (c *IgnoreChild) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type IgnoreChildSpec struct {
	Name string `json:"name" yaml:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IgnoreChildList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []IgnoreChild `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Dns struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`

	Status DnsNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type DnsNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Dns) CRDName() string {
	return "dnses.gns.tsm.tanzu.vmware.com"
}

func (c *Dns) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DnsList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Dns `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type AdditionalGnsData struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              AdditionalGnsDataSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            AdditionalGnsDataNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type AdditionalGnsDataNexusStatus struct {
	Status AdditionalStatus `json:"status,omitempty" yaml:"status,omitempty"`
	Nexus  NexusStatus      `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *AdditionalGnsData) CRDName() string {
	return "additionalgnsdatas.gns.tsm.tanzu.vmware.com"
}

func (c *AdditionalGnsData) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type AdditionalGnsDataSpec struct {
	Description AdditionalDescription `json:"description" yaml:"description"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AdditionalGnsDataList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []AdditionalGnsData `json:"items" yaml:"items"`
}

// +k8s:openapi-gen=true
type RandomDescription struct {
	DiscriptionA string
	DiscriptionB string
	DiscriptionC string
	DiscriptionD string
}

// +k8s:openapi-gen=true
type RandomStatus struct {
	StatusX int
	StatusY int
}

// +k8s:openapi-gen=true
type HostPort struct {
	Host Host
	Port Port
}

// +k8s:openapi-gen=true
type ReplicationSource struct {
	SourceKind SourceKind
}

// +k8s:openapi-gen=true
type ServiceSegmentRef struct {
	Field1 string
	Field2 string
}

// +k8s:openapi-gen=true
type Description struct {
	Color     string
	Version   string
	ProjectId string
	TestAns   []Answer
	Instance  Instance
	HostPort  HostPort
}

// +k8s:openapi-gen=true
type Answer struct {
	Name string
}

// +k8s:openapi-gen=true
type GnsState struct {
	Working     bool
	Temperature int
}

// +k8s:openapi-gen=true
type AdditionalDescription struct {
	DiscriptionA string
	DiscriptionB string
	DiscriptionC string
	DiscriptionD string
}

// +k8s:openapi-gen=true
type AdditionalStatus struct {
	StatusX int
	StatusY int
}

type RandomConst1 string
type RandomConst2 string
type RandomConst3 string
type MyConst string
type SourceKind string
type Port uint16
type Host string
type Instance float32
type AliasArr []int
type MyStr string
type TempConst1 string
type TempConst2 string
type TempConst3 string

const (
	MyConst3 RandomConst3 = "Const3"
	MyConst2 RandomConst2 = "Const2"
	MyConst1 RandomConst1 = "Const1"
	Object   SourceKind   = "Object"
	Type     SourceKind   = "Type"
	XYZ      MyConst      = "xyz"
	Const3   TempConst3   = "Const3"
	Const2   TempConst2   = "Const2"
	Const1   TempConst1   = "Const1"
)
