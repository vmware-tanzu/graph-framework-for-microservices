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
	nexus.SingletonNode
	ConfigName        string
	GNS               gns.Gns `nexus:"child"`
	Cluster           Cluster
	FooA              AMap
	FooMap            map[string]string
	FooB              BArray
	FooC              CInt   `nexus-graphql:"ignore:true"`
	FooD              DFloat `nexus-graphql:"type:string"`
	XYZPort           []gns.Description
	ABCHost           []gns.Host
	ClusterNamespaces []ClusterNamespace
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
