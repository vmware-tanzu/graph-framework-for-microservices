# Nexus DSL

Nexus DSL empowers applications to specify its datamodel in Go (also called Golang or Go language).

Using Go, applications can express its datamodel as a graph, the specification of a node in the graph, its hierarchy, its relationships etc.

Why Go? Go is easy to read, easy to type and easy to parse. Golang's structs are
natural data types definitions and Golang's annotations and comments can be used as additional parameters.

## Datamodel as a Graph

A datamodel allows an application to structure its business data for simplified organization, storage and consumption. Cloud native applications are rarely static. These applications evolve constantly and often. As these applications evolve, so does its datamodel as well. Relational datamodels are not well suited to handle this agility.
Graphl datamodels fit that bill.

Graph datamodels are agile, flexible and highly performant while allow the datamodel and inturn data to grow along with the application and business needs.

Graph datamodel is composed of two types of elements: nodes and relationships
- nodes representing an entity (a person, place, thing etc)
- relationships representing how any two nodes are associated

Nexus DSL represent your application data as a Graph Datamodel.
## Graph Syntax

 #### TL;DR [here](#Nexus-DSL-syntax-shortcut)

### Nodes: Go Structs

Nexus node is a Go struct annotated as a graph node in Nexus DSL.

It is a Go struct, and so can hold fields of all valid Go types.

A struct can be annotated as a Nexus node by including `nexus.Node` (defined [here](https://github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus))as an embedded field in the struct. A struct without this annotation is not a Nexus node, but just a valid Go struct.

In essence, not all Go structs are Nexus nodes, but all Nexus nodes are Go structs.

Here is a sample Nexus node:

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
### Relationships

Nexus nodes can be associated with other Nexus nodes through Relationships.

Nexus DSL supports two types of relations between Nexus nodes:

- child / children
- link / links

#### Child / Children

A Child relationship provides a way to designate one Nexus node as a "parent" or as a hierarchical root to other Nexus node/s in the graph. The Nexus nodes that are associated with the Parent node are referred to has "children" or "child" nodes.

Parent-child relationships in Nexus DSL have the following attributes:

* a Nexus node cannot be claimed by more than one Nexus node as a child. So, each Nexus node can have atmost one parent
* object for child Nexus node can only be created if the object of its parent Nexus node exists in the graph
* the lifecycle of the child objects are strictly tied to the lifecycle of the parent object. If the parent object is deleted, all children are deleted as well
* lifecycle of the parent object is independent of the lifecycle of the children object. So parent can exist even if the child object does not exist
* circular relationships are prohibited; i.e a parent node cannot be claimed as a child by any of the Nexus nodes in the parent's hierarchy


Child relationship can be created by annotating a field of the parent Nexus node with one of the following:

 * `nexus:"child"` if the parent can only claim a specific object of a Nexus Node, as a child
 * `nexus:"children` if the parent can claim multiple child objects of a Nexus node, as children

Example datamodel with Child relationship:

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

#### Link / Links

A Link relationship provides a way to designate one Nexus node "linked" to other Nexus node/s in the graph. Link relationships are useful to provide a soft or non-hierarchical construct to associate nodes in the graph.

Links can be across hierarchy and so provide a loose coupling between nodes of the graph. As such, Links come with very little riders or restrictions.

Link relationships in Nexus DSL have the following attributes:

* an node can be linked by any other Nexus node in the graph, without restrictions
* lifecycle of the linked nodes are independent of each other

Link relationship can be created by annotating a field of the Nexus node with one of the following:

 * `nexus:"link"` if the node would like to link to a specific object of a Nexus node
 * `nexus:"links` if the node would like to link to multiple objects of a Nexus node

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

TBR: As you can see you can use Go imports for adding nodes from different package. Structs which
don't have `nexus.Node` field can be used for defining spec.

### Node Status

Nexus DSL provides the ability to capture "status" data on a Nexus node.

While the spec fields and relationships of a Nexus node capture its user specified configuration and state, a Status, on the other hand captures runtime state / data that are relevant, as deemed by the user / application.

A custome status can be associated with the Nexus node, by annotating a field in the node struct annotation:

```
 `nexus:"status"`
```

The below example associates status of type "LeaderStatus" with Nexus node "Leader"

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

### Nexus Graph

A graph in Nexus datamodel is a collection of Nexus nodes and their relationships.

Graph built using Nexus DSL will have the following attributes:

* graph will be directed and acyclic
* graph is rooted by a Nexus node. A root node in the graph is a node which is not claimed as child by any other nexus node. Nexus datamodel only allows one root node
* an instance of a graph is identified by a name and domain in Nexus runtime

## API Syntax

Nexus Graph is a collection of Nexus Nodes and Relationships.

Nexus DSL provides the following API schemes to interact with the graph datamodel:

* REST API
* GraphQL

The Nexus node is the unit at which API's are defined and exposed. Through these API's, the user can CRUD the node spec, status, its relationships etc. 

### REST API

Nexus DSL provides the syntax to access a Nexus node through one or more REST API's. 

The syntax for declaration of REST API for a Nexus node has 2 parts in Nexus DSL:

1. Create an instance of type [nexus.RestAPISpec](https://github.com/vmware-tanzu/graph-framework-for-microservices/blob/main/nexus/nexus/nexus.go) that defines one or most REST APIs.

2. Associate the `nexus.RestAPISpec` instance with a Nexus node.

#### nexus.RestAPISpec

Declares one or more REST API's.

Each API spec captures information about the REST API, such as:

* URI
* Allowed http methods
* Desired Response codes

#### Associate nexus.RestAPISpec with Nexus node

The association of an instance/variable of type nexus.RestAPISpec to a Nexus node is achieved by adding a `comment` on the Nexus node in the following format:

`//nexus-rest-api-gen:<variable of type nexus.RestAPISpec>`

The keyword `nexus-rest-api-gen` is an instruction in Nexus DSL to the Nexus Compiler that the referenced REST API's should be associated with the Nexus Node.

As an exanple:

```Go
package role

import (
  "github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var LeaderRestAPISpec = nexus.RestAPISpec{
  Uris: []nexus.RestURIs{
    {
      Uri: "/v1alpha2/root/{root}/leader/{role.Leader}",
      Methods: nexus.HTTPMethodsResponses{
        http.MethodGet: nexus.DefaultHTTPGETResponses,
      },
    },
    {
      Uri: "/v1alpha2/leader",
      QueryParams: []string{
        "root",
        "role.Leader"
      },
      Methods: nexus.HTTPMethodsResponses{
        http.MethodGet: nexus.DefaultHTTPGETResponses,
      },
    },
    {
      Uri:     "/v1alpha2/root/{root}/leader",
      Methods: nexus.HTTPListResponse,
    },
  },
}

// nexus-rest-api-gen:LeaderRestAPISpec
type Leader struct {
  nexus.Node
  EmployeeID int
}
```

The above example, defined a variable LeaderRestAPISpec of type nexus.RestAPISpec. This variable is then referenced in the code comment on Nexus node using the keyword: `nexus-rest-api-gen`

### Custom GraphQL query spec

Custom GraphQl query spec is a way of extending GraphQl server with queries to external GRPC servers. 

To add custom queries you need to use `nexus.GraphQLQuerySpec` struct imported from
[nexus](https://github.com/vmware-tanzu/graph-framework-for-microservices/blob/main/nexus/nexus/nexus.go) package.

This is a collection  of `nexus.GraphQLQuery` structs. A GraphQLQuery specifies a custom query available via GraphQL API.
Each GraphQLQuery is self contained unit of the exposed custom query.
Format is as following:

```Go
type GraphQLQuery struct {
    Name            string               `json:"name,omitempty"`             // query identifier
    ServiceEndpoint GraphQLQueryEndpoint `json:"service_endpoint,omitempty"` // endpoint that serves this query
    Args            interface{}          `json:"args,omitempty"`             // custom graphql filters and arguments
    ApiType         GraphQlApiType       `json:"api_type,omitempty"`         // type of GRPC API endpoint
}
```

Currently there are two API types supported:
- NexusGraphQL Query, [proto file](https://github.com/vmware-tanzu/graph-framework-for-microservices/blob/main/nexus/proto/graphql/query.proto)
  specifies Query requests and responses which server should implement
- GetMetrics, [proto_file](https://github.com/vmware-tanzu/graph-framework-for-microservices/blob/main/nexus/proto/query-manager/server.proto)
  specifies expected requests and responses.

In NexusGraphQL Query you can provide any arguments, they will be translated into map. In GetMetrics arguments must be a
subset of [MetricsArg](https://github.com/vmware-tanzu/graph-framework-for-microservices/blob/main/nexus/generated/query-manager/server.pb.go#L23)
arguments.

GraphQlQuerySpec can be attached to a Nexus Node using comment above a Node.

Example custom query:
```Go
package role

import (
  "github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

// nexus-graphql-query:MyCustomQueries
type Leader struct {
  nexus.Node
  EmployeeID int
}

var MyCustomQueries = nexus.GraphQLQuerySpec{
  Queries: []nexus.GraphQLQuery{
    {
      Name: "queryGns1",
      ServiceEndpoint: nexus.GraphQLQueryEndpoint{
        Domain: "nexus-query-responder",
        Port:   15000,
      },
      Args:    QueryFilters{},
      ApiType: nexus.GraphQLQueryApi,
    },
    {
      Name: "queryMetrics1",
      ServiceEndpoint: nexus.GraphQLQueryEndpoint{
        Domain: "query-manager",
        Port:   15002,
      },
      Args:    nil,
      ApiType: nexus.GetMetricsApi,
    },
    {
      Name: "queryMetrics2",
      ServiceEndpoint: nexus.GraphQLQueryEndpoint{
        Domain: "query-manager",
        Port:   15003,
      },
      Args:    metricsFilers{},
      ApiType: nexus.GetMetricsApi,
    },
  },
}

type QueryFilters struct {
  foo           string
  bar             string
}

type metricsFilers struct {
  StartTime string
  EndTime   string
  TimeInterval  string
}
```

### OpenAPI validation

Spec fields of nexus nodes can be extended with additional validation, for field which should be validated you can add
comments above a field with format `//nexus-validation: Validation pattern`.

Example:
```Go
package role

import (
  "github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Leader struct {
  nexus.Node
  //nexus-validation: MaxLength=8, MinLength=2
  //nexus-validation: Pattern=abc
  Department string
  EmployeeID int
}
```

### Singleton nodes

Singleton Nodes are Nexus Nodes for which we are enforcing that in a given hierarchy there will be only one node of a
given type. To specify node to be a singleton use `nexus.SingletonNode` as a field, instead of `nexus.Node`.

Example:
```Go
package role

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type Leader struct {
	nexus.SingletonNode
	EmployeeID int
}
```

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