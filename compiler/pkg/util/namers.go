package util

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"

	"golang.org/x/text/language"
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

func GetInformerImportName(pkgName, baseGroupName, version string) string {
	return "informer" + RemoveSpecialChars(GetImportPath(pkgName, baseGroupName, version)) // eg informerroothelloworldv1
}

func GetGroupGoName(baseGroupName string) string {
	baseGroupName = strings.Replace(baseGroupName, "-", ".", 1)
	return namer.IC(strings.Split(baseGroupName, ".")[0]) // eg Helloworld
}

func GetGroupResourceName(nodeName string) string {
	return strings.ToLower(ToPlural(nodeName)) // eg roots
}

func GetNodeNameTitle(nodeName string) string {
	return cases.Title(language.Und, cases.NoLower).String(nodeName) // eg Root
}

func GetGroupVarName(pkgName, baseGroupName, version string) string {
	return pkgName + GetGroupGoName(baseGroupName) + cases.Title(language.Und, cases.NoLower).String(version) // eg rootHelloworldV1
}

func GetGroupTypeName(pkgName, baseGroupName, version string) string {
	return cases.Title(language.Und, cases.NoLower).String(RemoveSpecialChars(pkgName)) + GetGroupGoName(baseGroupName) + cases.Title(language.Und, cases.NoLower).String(version) // eg RootHelloworldV1
}

func GetSimpleGroupTypeName(pkgName string) string {
	return cases.Title(language.Und, cases.NoLower).String(RemoveSpecialChars(pkgName)) // eg Root
}

func GetGroupResourceNameTitle(nodeName string) string {
	return cases.Title(language.Und, cases.NoLower).String(ToPlural(nodeName)) // eg Roots
}

func GetGroupResourceType(baseNodeName, pkgName, baseGroupName, version string) string {
	return strings.ToLower(baseNodeName) + GetGroupTypeName(pkgName, baseGroupName, version) // eg rootRootHelloworldV1
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

func GetInternalImport(crdModulePath string, packageName string) string {
	return `"` + crdModulePath + packageName + `"`
}

func GetPackageName(groupName string) string {
	replacer := strings.NewReplacer(".", "_", "-", "_")
	return replacer.Replace(groupName)
}
