package util

import (
	"strings"
	"unicode"

	"k8s.io/gengo/namer"
)

func GetCrdName(nodeName, pkgName, baseGroupName string) string {
	return GetGroupResourceName(nodeName) + "." + GetGroupName(pkgName, baseGroupName) // eg roots.root.helloworld.com
}

func GetGroupName(pkgName, baseGroupName string) string {
	return pkgName + "." + baseGroupName // eg root.helloworld.com
}

func GetImportPath(pkgName, baseGroupName, version string) string {
	return GetGroupName(pkgName, baseGroupName) + "/" + version // eg root.helloworld.com/v1
}

func GetBaseImportName(pkgName, baseGroupName, version string) string {
	return "base" + RemoveSpecialChars(GetImportPath(pkgName, baseGroupName, version)) // eg baseroothelloworldv1
}

func GetGroupGoName(baseGroupName string) string {
	return namer.IC(strings.Split(baseGroupName, ".")[0]) // eg Helloworld
}

func GetGroupResourceName(nodeName string) string {
	return strings.ToLower(ToPlural(nodeName)) // eg roots
}

func GetGroupVarName(pkgName, baseGroupName, version string) string {
	return pkgName + GetGroupGoName(baseGroupName) + strings.Title(version) // eg rootHelloworldV1
}

func GetGroupTypeName(pkgName, baseGroupName, version string) string {
	return strings.Title(RemoveSpecialChars(pkgName)) + GetGroupGoName(baseGroupName) + strings.Title(version) // eg RootHelloworldV1
}

func GetGroupResourceNameTitle(nodeName string) string {
	return strings.Title(ToPlural(nodeName)) // eg Roots
}

func GetGroupResourceType(baseNodeName, pkgName, baseGroupName, version string) string {
	return strings.ToLower(baseNodeName) + GetGroupTypeName(pkgName, baseGroupName, version) // eg rootRootHelloworld
}

func GetTag(name string) string {
	return string(unicode.ToLower(rune(name[0]))) + name[1:] // eg serviceGroup
}

func GetGvkFieldName(fieldName string) string {
	return fieldName + "Gvk"
}

func GetGvkFieldTagName(fieldName string) string {
	return GetTag(fieldName) + "Gvk"
}

func GetGroupFromCrdName(crdName string) string {
	parts := strings.Split(crdName, ".")
	return strings.Join(parts[1:], ".")
}

func GetPackageNameFromCrdName(crdName string) string {
	parts := strings.Split(crdName, ".")
	return parts[1]
}
