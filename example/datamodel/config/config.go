package config

import (
	"net/http"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
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
	GNS gns.Gns `nexus:"child"`
	DNS gns.Dns `nexus:"child"`

	// Examples for cross-package import.
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

// type StructWithEmbeddedField struct {
// 	SomeStruct
// 	gns.MyStr
// }
