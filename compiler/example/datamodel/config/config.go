package config

import (
	"net/http"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns/service-group"
	py "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/policy"
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var nonNexusValue = 1
var nonValue int

var BarCustomCodesResponses = nexus.HTTPCodesResponse{
	http.StatusBadRequest: nexus.HTTPResponse{Description: "Bad Request"},
}

var BarCustomMethodsResponses = nexus.HTTPMethodsResponses{
	http.MethodPatch: BarCustomCodesResponses,
}

// nexus-graphql-query:root.GeneralGraphQLQuerySpec
type Config struct {
	nexus.Node
	GNS         gns.Gns                `nexus:"child"`
	DNS         gns.Dns                `nexus:"child"`
	VMPPolicies py.VMpolicy            `nexus:"child"`
	ACPPolicies py.AccessControlPolicy `nexus:"links"`
	Domain      Domain                 `nexus:"child"`
	// Examples for cross-package import.
	MyStr0 *gns.MyStr
	MyStr1 []gns.MyStr
	MyStr2 map[string]gns.MyStr

	XYZPort           gns.Port
	ABCHost           []gns.Host
	ClusterNamespaces []ClusterNamespace

	TestValMarkers TestValMarkers `json:"testValMarkers" yaml:"testValMarkers"`
	FooExample     FooTypeABC     `nexus:"children"`
	Instance       float32
	CuOption       string                        `json:"option_cu"`
	SvcGrpInfo     servicegroup.SvcGroupLinkInfo `nexus:"child"`
}

type FooTypeABC struct {
	nexus.Node
	FooA AMap
	FooB BArray
	FooC CInt   `nexus-graphql:"ignore:true"`
	FooD DFloat `nexus-graphql:"type:string"`
	FooE CInt   `json:"foo_e" nexus-graphql:"ignore:true"`
	FooF DFloat `json:"foo_f" yaml:"c_int" nexus-graphql:"type:string"`
}

type Domain struct {
	nexus.Node
	PointPort        *gns.Port
	PointString      *string
	PointInt         *int
	PointMap         *map[string]string
	PointSlice       *[]string
	SliceOfPoints    []*string
	SliceOfArrPoints []*BArray
	MapOfArrsPoints  map[string]*BArray
	PointStruct      *Cluster
}

type ClusterNamespace struct {
	Cluster   MatchCondition
	Namespace MatchCondition
}

type MatchCondition struct {
	Name string
	Type gns.Host
}

type Cluster struct {
	Name string
	MyID int
}

type AMap map[string]string

type BArray []string
type CInt uint8
type DFloat float32

type CrossPackageTester struct {
	Test gns.MyStr
}

type EmptyStructTest struct{}

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

type SomeStruct struct{}

type AliasedField gns.MyStr
type AliasedFieldMap map[string]gns.MyStr
type AliasedFieldList []gns.MyStr

type StructWithEmbeddedField struct {
	SomeStruct
	gns.MyStr
	ExplicitField    gns.MyStr
	AliasedField     AliasedField
	AliasedFieldMap  AliasedFieldMap
	AliasedFieldList AliasedFieldList
}
