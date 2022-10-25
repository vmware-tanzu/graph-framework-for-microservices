# Nexus DSL

In Nexus, you can specify your data model using Domain Specific Language, which is Golang.

Why Go? Go is easy to read, easy to type and easy to parse. Golang's structs are
natural data types definitions and Golang's annotations and comments can be used as additional parameters.

## How to define your data model?

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
