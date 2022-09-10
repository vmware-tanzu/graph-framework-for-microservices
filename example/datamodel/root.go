package root

import (
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/config"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler.git/example/datamodel/nexus"
)

type Foo string

type Bar struct {
	Foo float32
}
type Root struct {
	nexus.SingletonNode
	Name   int
	Config config.Config `nexus:"child"`
	Foochildren config.Config `nexus:"children"`
	Foolink config.Config `nexus:"link"`
	Foolinks config.Config `nexus:"links"`
	CustomBar Bar
	NonStructFoo Foo
	CustomBarMap map[string]Bar
	StatusBarMap Bar `nexus:"status"`
	ArrayBar []Bar
}

type NonNexusType struct {
	nexus.Node
	Test int
	Foo  Foo
	Bar  Bar
	StatusBar Bar `nexus:"status"`
}


// {
// 	Node{
// 		IsSingletonNode
// 		Child
// 		Children
// 		Link
// 		Links
// 		MapFields
// 		CustomFields
// 		NonStructFields
// 		PkgName
// 		NodeName

// 	}
// }