package crd_generator

import (
	"strings"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/pkg/util"
	"k8s.io/gengo/namer"
)

func getGroupName(pkgName, baseGroupName string) string {
	return pkgName + "." + baseGroupName // eg root.helloworld.com
}

func getImportPath(pkgName, baseGroupName, version string) string {
	return getGroupName(pkgName, baseGroupName) + "/" + version // eg root.helloworld.com/v1
}

func getBaseImportName(pkgName, baseGroupName, version string) string {
	return "base" + util.RemoveSpecialChars(getImportPath(pkgName, baseGroupName, version)) // eg baseroothelloworldv1
}

func getGroupGoName(baseGroupName string) string {
	return namer.IC(strings.Split(baseGroupName, ".")[0]) // eg Helloworld
}

func getGroupResourceName(nodeName string) string {
	return strings.ToLower(util.ToPlural(nodeName)) // eg roots
}

func getGroupVarName(pkgName, baseGroupName, version string) string {
	return pkgName + getGroupGoName(baseGroupName) + strings.Title(version) // eg rootHelloworldV1
}

func getGroupTypeName(pkgName, baseGroupName, version string) string {
	return strings.Title(pkgName) + getGroupGoName(baseGroupName) + strings.Title(version) // eg RootHelloworldV1
}

func getGroupResourceNameTitle(nodeName string) string {
	return strings.Title(util.ToPlural(nodeName)) // eg Roots
}

func getGroupResourceType(baseNodeName, pkgName, baseGroupName, version string) string {
	return strings.ToLower(baseNodeName) + getGroupTypeName(pkgName, baseGroupName, version) // eg rootRootHelloworld
}
