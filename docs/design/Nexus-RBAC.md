# Nexus RBAC

Nexus Role Based Aceess Control is a mechanism for declaring and enforcing access permissions on resources and APIs.

Nexus SDK provides:

* a declarative syntax to express RBAC specification
* server side enforcement, hence enforcement is common across all clients
* hierarchical, auto-discovering and elastic authorization enforcement at runtime
* rule based, declarative resource selection
* non-imposing design that lets product teams design their own constraints

In Nexus, RBAC is enforced at the Nexus Tenant Runtime.

![RuntimeRBAC](.content/images/runtime-rbac-flow.png)
# RBAC Flow

![RBACFlow](.content/images/RbacFlow.png)

An Authenticated request consists of the following:

* an string identifying end user or thing associated with the request
* the resource/object/api being accessed by the request
* operation being performed on the resource/object/api


The RBAC layer will consume the above information to enforce authorization policies.
# RBAC Constructs

![RBAC](.content/images/nexus-rbac.png)

## User

User is an entity / thing interacting with Nexus Tenant API Gateway.

Nexus RBAC design has the following expectations, when it comes to User info:

1. Product will invoke Add / Remove User API that will be exposed by Nexus SDK. If the user is not added to Nexus runtime, then the user will essentially be unknown.
 2. User is unique at a Tenant or System level.
3. User string is expected to be made avaiable to the RBAC layer, post authentication stage.
4. Product will provide a means to idenity which params in the request carries user info.
5. User can be associated with multiple roles and rolesbindings.

```
package nexus

// AddUser adds a user to Nexus runtime along with a metadata string.
// The metadata string is opaque to nexus runtime and is not interpreted.
// So it can be any string (ex. json encoded).
// AddUser will always upsert the user info.
AddUser(name string, metadata string) error

GetUser(name string) (metadate string, err error)

RemoveUser(name string) error

```
## Resources

Resource in nexus is a first class object in datamodel and as such can be a datamodel node,
custome API etc. 

* A resource is the unit/entity at which RBAC rules are evaluated and enforced.

* A resource may be associated with multiple roles at any time.

![DatamodelExample](.content/images/DatamodelExample.png)

* Datamodel DSL should explicitly callout all possible roles that can be associated with a Resource at runtime.

## Roles

Role specifies "resources" and the "operations" that are possible on them.

![RoleIllustration](.content/images/RoleIllustration.png)

### Resource Selection
When a Role is associated with a Resource, the association is hierarchical: ie. the role is associated with a specific resource and all its child resources in the datamodel. The child resources to be associated with a Role is specified by a rules:

Rules identifying child resources in a Role can be:

    - "*" -- all resources of all types, under this specific resource can be associated with this Role
    - < Nexus Type > -- all resource of specified types can be associated with this Role.
                        Type is idenitifed by its fully qualified Nexus type in DSL.
                        Example: config.MyConfig, root.Root etc 

If no matching rule is specified, then the role is only associable to the Resource on it is specified in the DSL.

### Operations

Operations are a collection of verbs/actions possible on a resouce.
    
    - Get
    - Put
    - Delete
    - List



## Scope

Scope is a logical grouping of a Runtime instance / object of a Resource and all its children, as specified in datamodel.

Scope is always hierarchical: it represents a specific instance of a Resource, along with instances of its children.

![Scope](.content/images/Scope.png)

For example, if the Resource is ClusterType, at runtime there are two object of ClusterType called FooCluster and BarCluster, then there are Scopes for ClusterType:
* FooCluster and all its child objects in the graph
* BarCluster and all its child objects in the graph

## RoleBinding

RoleBinding binds a User to a Scope + Role.

Rolebinding is only possible at runtime, as its needs to know the Scope.

Each runtime object/instance will expose methods that can be used to associate "allowed" Roles on it to a specific User.
### Workflow:

#### Step 1: Define Roles (Compile time)

```
// Create a GNS Admin Role
var GnsConfigAdmin = Role{
	Rules: []ResourceRule {
		{
			"config.Gns": AllVerbs,
		},
	},
}

// Create a GNS Operator Role
var GnsOperator = Role{
	Rules: []ResourceRule {
		{
			"config.Gns": []Verb{ Get, Put, List, },
		},
		{
			"runtime.Gns":[]Verb{ Get, List, },
		},
	},
}
```

#### Step 2: Associate Roles with Resource Type(Compile time)

```

package config

// Gns struct.
// nexus-rest-api-gen:GNSRestAPISpec
// nexus-rbac-roles: GnsConfigAdmin, GnsOperator   <-- Associating roles to Gns Node Type
// specification of GNS.
type Gns struct {
	nexus.Node
	Domain                 string
	UseSharedGateway       bool
	Description            Description
	GnsServiceGroups       map[string]service_group.SvcGroup `nexus:"child"`
	GnsAccessControlPolicy policy.AccessControlPolicy        `nexus:"child"`
	Dns                    Dns                               `nexus:"link"`
	State                  GnsState                          `nexus:"status"`
}
```

```

package runtime

// Gns struct.
// nexus-rest-api-gen:GNSBindingRestAPISpec
// nexus-rbac-roles: GnsOperator   <-- Associating roles to Gns Node Type
// specification of GNS.
type GnsBinding struct {
	nexus.Node

	...
}
```

#### Step 3: Bind an object (and its hierarchy) to a User and Role (Runtime)

```
// Package config

func(g *Gns) GnsAddGnsConfigAdmin(user User)
func(g *Gns) GnsRemoveGnsConfigAdmin(user User)
func(g *Gns) GnsGetGnsConfigAdmin() []User

func(g *Gns) GnsAddGnsOperator(user User)
func(g *Gns) GnsRemoveGnsOperator(user User)
func(g *Gns) GnsGetGnsOperator() []User
```

## Nexus Implementation

### Infra Library
#### Role Type in Nexus SDK
```
package nexus

type Verb string
const (
	Get    Verb = "Get"
	List        = "List"
	Put         = "Put"
	Delete      = "Delete"
)

var (
	AllResources = []string{"*"}
	AllVerbs     = []Verb{"*"}
)
type ResourceRule map[string][]Verb

type Role struct {
	Rules []ResourceRule
}

// Stock Admin Rule
var AdminRole = Role{
	Rules: []ResourceRule {
		{
			AllResources: AllVerbs,
		},
	},
}

// Stock Operator Rule
var OperatorRole = Role{
	Rules: []ResourceRule {
		{
			AllResources: []Verb{
					Get,
					List,
			},
		},
	},
}
```

```

// Create Project Admin Rule
var ProjectAdmin = AdminRole

// Create Project Operator Rule
var ProjectOperator = OperatorRole

```



