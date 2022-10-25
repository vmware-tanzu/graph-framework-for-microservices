# Nexus DSL

In Nexus, you can specify your data model using Domain Specific Language, which is Golang.

Why Go? Go is easy to read, easy to type and easy to parse. Golang's structs are
natural data types definitions and Golang's annotations and comments can be used as additional parameters.

## DSL syntax

In the following points we'll explain Nexus DSL, for a shortcut you can go [here](#Nexus-DSL-syntax-shortcut).

### Go packages as API groups
When you define your data model first of all you need to think about logical grouping of
your data resources.
For example, you can partition model of a company into departments or roles. In Nexus this
can be represented by creating separate Go package for each of your logical group. Nexus compiler
will translate this into separate API groups, which are equivalents to [Kubernetes API
groups](https://kubernetes.io/docs/reference/using-api/#api-groups).

### Go structs as data types

In nexus logical unit of data type is Nexus Node. As Nexus
follows a graph data-model-centric design pattern each Nexus Node is a node in a graph.
Defining Nexus nodes is easy, you just need to define Go struct and as one of embedded fields
use `nexus.Node`. Node definition can be imported from
`https://github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus`
package. As a rest of struct fields you can specify spec of your node.
For a spec you can use standard Go types like `int`, `string` or more complex structs.
So your first data model can look like such Go file:
```Go
package role

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Leader struct {
	nexus.Node
	EmployeeID int
}
```

### Go annotations as relation specifications

In a graph of your data you need to specify relations between nodes. In a Nexus graph there are
two types of relations which you can specify:
- parent-child relation to specify tree-like hierarchy
- soft link relation between nodes from different parts of tree.

To specify relations you can use nexus annotations.
For parent-child relations you can use two types of annotatations:
- `nexus:"child"`
- `nexus:"children`.

`nexus:"child"` means that given parent object can have only one child of this type, `nexus:"children"`
that given parent object can have many children objects of this type.

Example datamodel with parent-child relations:
```Go
package role

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Leader struct {
	nexus.Node
	HR         HR        `nexus:"child"`
	Devs       Developer `nexus:"children"`
	EmployeeID int
}

type HR struct {
	nexus.Node
	EmployeeID int
}

type Developer struct {
	nexus.Node
	EmployeeID int
}
```

So in this example parent node Leader can have one child object HR and multiple child objects Developer.

Similarly, you can add soft links relations. `nexus:"link"` means there can be one linked object,
`nexus:"links"` - multiple objects. Links are unidirectional.

For example to link `Developer` to `Location` you can use follow syntax:
```Go
package geo

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
	"yourdatamodel/role"
)

type Location struct {
	nexus.Node
	Devs    role.Developer `nexus:"links"`
	Address Address
}

type Address struct {
	Country    string
	PostalCode string
	Street     string
	No         int
}
```

As you can see you can use Go imports for adding nodes from different package. Structs which
don't have `nexus.Node` field can be used for defining spec.

### Status of a node

You can add custom status of nexus node by using `nexus:"status"` annotation. In the runtime
you can use this field for specifying current state of the object. Example:

```Go
package role

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Leader struct {
	nexus.Node
	EmployeeID int
	State      LeaderStatus `nexus:"status"`
}

type LeaderStatus struct {
	IsOnVacations            bool
	DaysLeftToEndOfVacations int
}
```


### REST API spec

...

### Custom GraphQL query spec

...

### OpenAPI validation spec

...

## Data model structure

Here are the restrictions which you should follow when defining data model.

- your data model should have go.mod file in main directory (you can add it by running `go mod init`)
- in main directory of your data model should be a package with root node of graph,
  there can be only one root node in data model
- there can be no disjoined node, so each node expect for root should have a parent.

# Nexus DSL syntax shortcut

```Go
/*  Section 1: Group Declaration */

package gns                                                                       <--- API / Node group name


/* Section 2: Attribute Definition (OPTIONAL) */

var GNSRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{                                                       <--- List of REST URL on which the Nexus Node should be exposed
		{
			Uri:     "/v1alpha2/projects/{project}/global-namespace/{gns.Gns}",   <--- REST URL on which the Nexus Node should be exposed
			Methods: nexus.HTTPMethodsResponses{                                  <--- Methods and responses to be enabled on this REST URL
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
		{
			Uri:     "/v1alpha2/global-namespace",                                <--- REST URL on which the Nexus Node should be exposed
			QueryParams: []string{
				"project",                                                        <--- Instead of URI param we are using QueryParams to specify project
			    "gns.Gns"
			},
			Methods: nexus.HTTPMethodsResponses{                                  <--- Methods and responses to be enabled on this REST URL
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
		{
			Uri:     "/v1alpha2/projects/{project}/global-namespaces",            <--- REST URL on which the Nexus Node should be exposed
			Methods: nexus.HTTPListResponse,                                      <--- nexus.HTTPListResponse indicates that this request will return a list of objects
		},
	},
}

/* Section 3: Attribute Binding (OPTIONAL) */
// nexus-rest-api-gen:GNSRestAPISpec                            <--- Binds a REST API attribute with a Nexus Node
// nexus-description: This is a my GNS node description         <--- Adds custom description to Nexus Node. This custom description will be propagated to references to this node in OpenAPI spec.

/* Section 4: Node Definition */
type Gns struct {

    nexus.Node                                                  <--- Declares type "Gns" to be a Nexus Node
                                                                     Alternatively nexus.SingletonNode can be used. SingletonNodes can only have 'default' as a display name
                                                                     For root level nodes there can be only one singleton node present in the system, for non-root objects
                                                                     only one can be present for given parents.


    F1 string                                                   <--- Defines a field of standard type.
                                                                     Supported standard types: bool, int32, float32, string

    F2 string `json:"f2,omitempty"`                             <--- Defines an optional field.
                                                                     To make a field optional, omitempty tag should be added.

    F3 CustomType                                               <--- Defines a field of custom type. The type definition can in the same go package or
                                                                     can be imported from other packages. It should be resolvable by Go compiler.

    C1 ChildType1 `nexus:"child"`                               <--- Declares:
                                                                     * a child of type "ChildType1"
                                                                     * field C1 through which a specific instance/object of type "ChildType1" can be accessed

                                                                     ChildType1 should be another nexus node in the graph.
                                                                     ChildType1 should be resolvable by Go compiler either in the local package or through Go import

                                                                     C1 is the handle through which a specific object of type ChildType1 can be accessed.
                                                                     While there can be multiple objects of type ChildType1 in the system, C1 can only hold single object.

    C2 ChildType2 `nexus:"children"`                            <--- Declares:
                                                                     * a child of type "ChildType2"
                                                                     * field C2 through which multiple objects/instances of type "ChildType2" can be accessed, with each object queryable by name

                                                                     ChildType2 should be another nexus node in the graph.
                                                                     ChildType2 should be resolvable by Go compiler either in the local package or through Go import.

                                                                     C2 is the handle through which multiple objects of type ChildType2 can be accessed.
                                                                     Objects in C2 are queryable by name.

    L1 LinkType1 `nexus:"link"`                                 <--- Declares:
                                                                     * a link of type "LinkType1"
                                                                     * field L1 through which a specific instance/object of type "LinkType1" can be accessed

                                                                     LinkType1 should be another nexus node in the graph.
                                                                     LinkType1 should be resolvable by Go compiler either in the local package or through Go import.

                                                                     L1 is the handle through which a specific object of type LinkType1 can be accessed.
                                                                     While there can be multiple objects of type LinkType1 in the system, L1 can only hold single object.

    L2 LinkType2 `nexus:"links"`                                <--- Declares:
                                                                    * a link of type "LinkType2"
                                                                    * field L2 through which multiple objects/instances of type "LinkType2" can be accessed, with each object queryable by name

                                                                     LinkType2 should be another nexus node in the graph.
                                                                     LinkType2 should be resolvable by Go compiler either in the local package or through Go import.

                                                                     L2 is the handle through which multiple objects of type LinkType2 can be accessed.
                                                                     Objects in L2 are queryable by name.

    S1 StatusType1 `nexus:"status"`                             <--- Declares a status field of type "StatusType".
}

```
