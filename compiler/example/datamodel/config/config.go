package config

import (
	"net/http"

	py "github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/policy"

	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/config/gns"
	"github.com/vmware-tanzu/graph-framework-for-microservices/compiler/example/datamodel/nexus"
)

var nonNexusValue = 1
var nonValue int

var BarCustomCodesResponses = nexus.HTTPCodesResponse{
	http.StatusBadRequest: nexus.HTTPResponse{Description: "Bad Request"},
}

var BarCustomMethodsResponses = nexus.HTTPMethodsResponses{
	http.MethodPatch: BarCustomCodesResponses,
}

type Config struct {
	nexus.Node
	GNS         gns.Gns                `nexus:"child"`
	DNS         gns.Dns                `nexus:"child"`
	VMPPolicies py.VMpolicy            `nexus:"child"`
	ACPPolicies py.AccessControlPolicy `nexus:"links"`

	// Examples for cross-package import.
	MyStr0 *gns.MyStr
	MyStr1 []gns.MyStr
	MyStr2 map[string]gns.MyStr

	TestValMarkers TestValMarkers `json:"testValMarkers" yaml:"testValMarkers"`
}

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

type StructWithEmbeddedField struct {
	SomeStruct
	gns.MyStr
}
