// Code generated by nexus. DO NOT EDIT.

package v1

import (
	gnstsmtanzuvmwarecomv1 "nexustempmodule/apis/gns.tsm.tanzu.vmware.com/v1"

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
type Config struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              ConfigSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            ConfigNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type ConfigNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Config) CRDName() string {
	return "configs.config.tsm.tanzu.vmware.com"
}

func (c *Config) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type ConfigSpec struct {
	MyStr0            *gnstsmtanzuvmwarecomv1.MyStr           `json:"myStr0" yaml:"myStr0"`
	MyStr1            []gnstsmtanzuvmwarecomv1.MyStr          `json:"myStr1" yaml:"myStr1"`
	MyStr2            map[string]gnstsmtanzuvmwarecomv1.MyStr `json:"myStr2" yaml:"myStr2"`
	XYZPort           gnstsmtanzuvmwarecomv1.Port             `json:"xYZPort" yaml:"xYZPort"`
	ABCHost           []gnstsmtanzuvmwarecomv1.Host           `json:"aBCHost" yaml:"aBCHost"`
	ClusterNamespaces []ClusterNamespace                      `json:"clusterNamespaces" yaml:"clusterNamespaces"`
	TestValMarkers    TestValMarkers                          `json:"testValMarkers" yaml:"testValMarkers"`
	Instance          float32                                 `json:"instance" yaml:"instance"`
	CuOption          string                                  `json:"option_cu"`
	GNSGvk            *Child                                  `json:"gNSGvk,omitempty" yaml:"gNSGvk,omitempty" nexus:"child"`
	DNSGvk            *Child                                  `json:"dNSGvk,omitempty" yaml:"dNSGvk,omitempty" nexus:"child"`
	VMPPoliciesGvk    *Child                                  `json:"vMPPoliciesGvk,omitempty" yaml:"vMPPoliciesGvk,omitempty" nexus:"child"`
	DomainGvk         *Child                                  `json:"domainGvk,omitempty" yaml:"domainGvk,omitempty" nexus:"child"`
	FooExampleGvk     map[string]Child                        `json:"fooExampleGvk,omitempty" yaml:"fooExampleGvk,omitempty" nexus:"children"`
	SvcGrpInfoGvk     *Child                                  `json:"svcGrpInfoGvk,omitempty" yaml:"svcGrpInfoGvk,omitempty" nexus:"child"`
	ACPPoliciesGvk    map[string]Link                         `json:"aCPPoliciesGvk,omitempty" yaml:"aCPPoliciesGvk,omitempty" nexus:"links"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ConfigList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Config `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type FooTypeABC struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              FooTypeABCSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            FooTypeABCNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type FooTypeABCNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *FooTypeABC) CRDName() string {
	return "footypeabcs.config.tsm.tanzu.vmware.com"
}

func (c *FooTypeABC) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type FooTypeABCSpec struct {
	FooA AMap   `json:"fooA" yaml:"fooA"`
	FooB BArray `json:"fooB" yaml:"fooB"`
	FooC CInt   `nexus-graphql:"ignore:true"`
	FooD DFloat `nexus-graphql:"type:string"`
	FooE CInt   `json:"foo_e" nexus-graphql:"ignore:true"`
	FooF DFloat `json:"foo_f" yaml:"c_int" nexus-graphql:"type:string"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FooTypeABCList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []FooTypeABC `json:"items" yaml:"items"`
}

// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
type Domain struct {
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec              DomainSpec        `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status            DomainNexusStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +k8s:openapi-gen=true
type DomainNexusStatus struct {
	Nexus NexusStatus `json:"nexus,omitempty" yaml:"nexus,omitempty"`
}

func (c *Domain) CRDName() string {
	return "domains.config.tsm.tanzu.vmware.com"
}

func (c *Domain) DisplayName() string {
	if c.GetLabels() != nil {
		return c.GetLabels()[common.DISPLAY_NAME_LABEL]
	}
	return ""
}

// +k8s:openapi-gen=true
type DomainSpec struct {
	PointPort        *gnstsmtanzuvmwarecomv1.Port `json:"pointPort" yaml:"pointPort"`
	PointString      *string                      `json:"pointString" yaml:"pointString"`
	PointInt         *int                         `json:"pointInt" yaml:"pointInt"`
	PointMap         *map[string]string           `json:"pointMap" yaml:"pointMap"`
	PointSlice       *[]string                    `json:"pointSlice" yaml:"pointSlice"`
	SliceOfPoints    []*string                    `json:"sliceOfPoints" yaml:"sliceOfPoints"`
	SliceOfArrPoints []*BArray                    `json:"sliceOfArrPoints" yaml:"sliceOfArrPoints"`
	MapOfArrsPoints  map[string]*BArray           `json:"mapOfArrsPoints" yaml:"mapOfArrsPoints"`
	PointStruct      *Cluster                     `json:"pointStruct" yaml:"pointStruct"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DomainList struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`
	metav1.ListMeta `json:"metadata" yaml:"metadata"`
	Items           []Domain `json:"items" yaml:"items"`
}

// +k8s:openapi-gen=true
type ClusterNamespace struct {
	Cluster   MatchCondition `json:"cluster" yaml:"cluster"`
	Namespace MatchCondition `json:"namespace" yaml:"namespace"`
}

// +k8s:openapi-gen=true
type MatchCondition struct {
	Name string                      `json:"name" yaml:"name"`
	Type gnstsmtanzuvmwarecomv1.Host `json:"type" yaml:"type"`
}

// +k8s:openapi-gen=true
type Cluster struct {
	Name string `json:"name" yaml:"name"`
	MyID int    `json:"myID" yaml:"myID"`
}

// +k8s:openapi-gen=true
type CrossPackageTester struct {
	Test gnstsmtanzuvmwarecomv1.MyStr `json:"test" yaml:"test"`
}

// +k8s:openapi-gen=true
type EmptyStructTest struct {
}

// +k8s:openapi-gen=true
type TestValMarkers struct {
	//nexus-validation: MaxLength=8, MinLength=2, Pattern=ab
	MyStr string `json:"myStr" yaml:"myStr"`
	//nexus-validation: Maximum=8, Minimum=2
	//nexus-validation: ExclusiveMaximum=true
	MyInt int `json:"myInt" yaml:"myInt"`
	//nexus-validation: MaxItems=3, MinItems=2
	//nexus-validation: UniqueItems=true
	MySlice []string `json:"mySlice" yaml:"mySlice"`
}

// +k8s:openapi-gen=true
type SomeStruct struct {
}

// +k8s:openapi-gen=true
type StructWithEmbeddedField struct {
	SomeStruct
	gnstsmtanzuvmwarecomv1.MyStr
	ExplicitField    gnstsmtanzuvmwarecomv1.MyStr `json:"explicitField" yaml:"explicitField"`
	AliasedField     AliasedField                 `json:"aliasedField" yaml:"aliasedField"`
	AliasedFieldMap  AliasedFieldMap              `json:"aliasedFieldMap" yaml:"aliasedFieldMap"`
	AliasedFieldList AliasedFieldList             `json:"aliasedFieldList" yaml:"aliasedFieldList"`
}

type AMap map[string]string
type BArray []string
type CInt uint8
type DFloat float32
type AliasedField gnstsmtanzuvmwarecomv1.MyStr
type AliasedFieldMap map[string]gnstsmtanzuvmwarecomv1.MyStr
type AliasedFieldList []gnstsmtanzuvmwarecomv1.MyStr
