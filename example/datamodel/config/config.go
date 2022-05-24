package config

import (
	"net/http"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config/gns"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

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
	MyStr  *gns.MyStr
	MyStr1 []gns.MyStr
	MyStr2 map[string]gns.MyStr
}

type CrossPackageTester struct {
	Test gns.MyStr
}

type EmptyStructTest struct{}
