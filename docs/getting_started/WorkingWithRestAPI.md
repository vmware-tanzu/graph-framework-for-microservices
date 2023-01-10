# API Gateway - Nexus Runtime

API GW provides a unified point of access to the Nexus Datamodel and Runtime.

![NexusRuntimeAPIGW](.content/images/NexusRuntimeAPIGW.png)

It provides the following interfaces:

* REST Endpoint
* Kubectl Endpoint
* IDP / Oauth Endpoint (In Progress)
* Graphql Endpoint (Future release)

## REST Endpoint

The REST endpoint provides a language agnostic interface to Datamodel through RESTfult APIs. 

### Workflow to expose Datamodel nodes as RESTful API
#### 1. Define a variable of type RestAPISpec

To expose a datamodel node as RESTful API, define a variable of type nexus.RestAPISpec.

```go
// NOTE: In this example we will use `Leader` node from the orgchart example.

var LeaderRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{ // List of URIs
            {
                Uri:     "/root/{orgchart.Root}/leader/{management.Leader}", // Endpoint URI
                Methods: nexus.DefaultHTTPMethodsResponses, // Exposed methods
            },
            {
                Uri:     "/leader", // Endpoint URI
                QueryParams: []string{ // QueryParams field is used to indicate Query Params for endpoint. 
                  "orgchart.Root",
                },
                Methods: nexus.DefaultHTTPMethodsResponses, // Exposed methods
            },
            {
                Uri: "/leaders", // Endpoint URI
                Methods: nexus.HTTPListResponse // HTTPListResponse type indicates that this endpoint will return a list of objects,
            },
	},
}
```
**URI**

URI field holds endpoint where you can declare parent hierarchy for a specified node. To add a new parent you need to add segment in the following format: 
```
{<package name>.<Datamodel node name>}
```
In example above we used `{orgchart.Root}` which is `Root` struct in the `orgchart` package. If you don't want to add these segments to URI you can use query params instead. 

For example in `/leaders` URI we can do a `GET` request with the same hierarchy as URI above. 
```
GET http://localhost:8080/leaders?orgchart.Root=exampleRoot&management.Leader=exampleLeader
```
**Methods**

Methods field holds HTTP methods which will be exposed for the given URI.

Right now, we provide default methods which you can use for the URI.

`nexus.DefaultHTTPMethodsResponses` which contains GET, PUT and DELETE methods:
* DefaultHTTPGETResponses (status 200, 404 and default 501)
* DefaultHTTPPUTResponses (status 200, 201 and default 501)
* DefaultHTTPDELETEResponses (status 200, default 501)
* HTTPListResponse (status 200, 404 and default 501) used to indicate that request will return a list

If you don't want to use default methods, you can define your own methods and responses following the example below.
```go
var LeaderCustomMethodsResponses = nexus.HTTPMethodsResponses{
	http.MethodDelete: nexus.HTTPCodesResponse{
		http.StatusOK:              nexus.HTTPResponse{Description: "ok"},
		http.StatusNotFound:        nexus.HTTPResponse{Description: http.StatusText(http.StatusNotFound)},
		nexus.DefaultHTTPErrorCode: nexus.DefaultHTTPError,
	},
}

var LeaderRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{ // List of URIs
		{
			Uri:     "/root/{Root.orgchart}/leader/{Leader.management}", // Endpoint URI
			Methods: LeaderCustomMethodsResponses, // Exposed methods
		},
	},
}

// nexus-rest-api-gen:LeaderRestAPISpec
type Leader struct {
	EngManagers Mgr               `nexus:"children"` <---multiple child
	HRs         hr.HumanResources `nexus:"links"`    <---multiple links
	Role        role.Executive    `nexus:"link"`     <---single link
	Status      LeaderState       `nexus:"status"`   <--- status annotation
}

type LeaderState struct {
  IsOnVacations            bool
  DaysLeftToEndOfVacations int
}
```

**QueryParams**

QueryParams field is used to indicate Query Params for an endpoint.

In the case where we don't want to put every parameter in URI param like `/root/{Root.orgchart}/leader/{Leader.management}`, we can move params to QueryParams field.
```
Uri: "/leader/{management.Leader}", // Endpoint URI 
QueryParams: []string{ // QueryParams field is used to indicate Query Params for endpoint. 
  "orgchart.Root",
},
```
and then if you want to make a request to that endpoint you can build URL like below:

`/leader/example-leader-name?orgchart.Root=example-root`

#### 2. Associate RestAPISpec with Datamodel Node

Once the RestAPISpec variable is define, associate it with the desired Datamodel node. Association of REST API definition to a datamodel node is specified by annotating the datamodel node with:
```
	key: nexus-rest-api-gen
	value: < variable of type nexus.RestAPISpec >
```
Example

```go
// nexus-rest-api-gen:LeaderRestAPISpec
type Leader struct {
```

#### 3. Build Datamodel
- Based on the nexus annotations **nexus:link** described, when we build the above datamodel, CRD annotations are populated with children, links and URIs as shown below.

```
nexus datamodel build --name $DATAMODEL_NAME 
```

Example:
```
annotations:
    nexus: |
      ...
      "name":"management.Leader",
      "hierarchy":["roots.orgchart.vmware.org"],
      "children":{"EngManagers":{"fieldName":"EngManagers","fieldNameGvk":"engManagersGvk","isNamed":true},
      "links":{"humanresourceses.hr.vmware.org":{"fieldName":"HR","fieldNameGvk":"hRGvk","isNamed":true},"Role":{"fieldName":"Role","fieldNameGvk":"roleGvk","isNamed":false}},
      "nexus-rest-api-gen":{"uris":[{"uri":"/root/{orgchart.Root}/leader/{management.Leader}","methods":{"GET":{"200":{"description":"OK"}
      ...
```

#### 4. Install Datamodel

Pre-req: Nexus runtime is installed and running in $NAMESPACE

```
nexus datamodel install name $DATAMODEL_NAME --namespace $NAMESPACE
```
### Access REST API
#### Setup access to the API GW

```
kubectl port-forward -n $NAMESPACE deployment/api-gw 8080:80
```
#### Access Swagger-UI

API-GW provides Swagger-UI where you can see the list of all declared endpoints.

It's available on `http://localhost:8080/docs`.

![SwaggerUI](.content/images/NexusRestAPISwaggerUI.png)

### Access through Kubectl
#### Setup access to the API GW

```
kubectl port-forward -n $NAMESPACE deployment/api-gw 8080:80
```
#### Access through kubectl

In addition to using REST API, you can use kubectl to retrieve, create and delete objects.

Difference between normal access to api-server and this is that api-gw supports parent hierarchy and name mangling.
You can use standard kubectl syntax to access your data. 

##### Get list of objects
```
$ kubectl -s localhost:8080 get roots
NAME                                       AGE
bf37a17aa378fa6ebf64c530a1ea4c0fff48e8e4   8m34s
```

##### Get specific object
```
$ kubectl -s localhost:5000 get roots bf37a17aa378fa6ebf64c530a1ea4c0fff48e8e4 -o yaml
apiVersion: orgchart.vmware.org/v1
kind: Root
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"orgchart.vmware.org/v1","kind":"Root","metadata":{"annotations":{},"name":"default"}}
  creationTimestamp: "2022-05-25T07:35:19Z"
  generation: 1
  labels:
    nexus/display_name: default
    nexus/is_name_hashed: "true"
  name: bf37a17aa378fa6ebf64c530a1ea4c0fff48e8e4
  resourceVersion: "927"
  selfLink: /apis/orgchart.vmware.org/v1/roots/bf37a17aa378fa6ebf64c530a1ea4c0fff48e8e4
  uid: 8a60791a-6545-4491-bdff-ae479f3330f1

```
##### Create
API gateway will automatically hash the name of an object.

In addition, it also supports parent linking, so when you create a child it will be automatically linked to a corresponding parent.
```
$ echo 'apiVersion: management.vmware.org/v1
kind: Leader
metadata:
  name: leader1
  labels:
    roots.orgchart.vmware.org: root1
  namespace: default
spec:
  designation: test
  employeeID: 0
  name: test' | kubectl -s localhost:5000 apply -f -
leader.management.vmware.org/ace02433b4e96f1367569cb618a45b26e418639a created
```

##### Delete
**kubectl based delete is NOT RECOMMENDED, you can easily remove all of your objects stored in the system by an accident**

**To safely delete an object use `nexus` CLI**


API gateway supports recursive deletion of children, so the whole tree will be deleted when you remove top level parent.
```
$ kubectl -s localhost:5000 delete leaders.management.vmware.org ace02433b4e96f1367569cb618a45b26e418639a
leader.management.vmware.org "ace02433b4e96f1367569cb618a45b26e418639a" deleted
```

To delete an object without using hashed name, you can use label selector.
```
kubectl -s localhost:5000 delete leaders.management.vmware.org -l nexus/display_name=leader1,roots.orgchart.vmware.org=root1
leader.management.vmware.org "ace02433b4e96f1367569cb618a45b26e418639a" deleted
```
`nexus/display_name` label is required for this to work. This label stores the real name of the object.

If you don't specify `nexus/display_name` label, then **ALL** objects which contains those parents will be deleted.

## REST API support for GET children, links and status

We can access it in the following way,
 1. Access the APIs through the Swagger UI.
 2. Access the APIs through the shim layer.

### API Gateway
- Based on the above CRD annotations, API GW is responsible for automatically creating endpoints for children, links and status.
- Only explicitly provided REST APIs will have an implicit URIs to query children, links and status.
- Currently, the following APIs are created from the  **parent URI** and generated annotations in `#### 3. Build Datamodel`.

For example, when `/leaders` URI specified in API spec, api-gw create an API with children, link and status fieldName as

- `/leader/EngManagers`
- `/leader/HR`
- `/leader/Role`
- `/leader/Status`

### Access the APIs through the Swagger UI.

##### GET multiple child API `/leader/EngManagers`
- GET `http://localhost:8080/leader/EngManagers'
- response body resemble the example below.

```
[
  {
    "group": "management.vmware.org/v1",
    "kind": "Mgr",
    "name": "Manager1",
    "hierarchy": [
      "roots.orgchart.vmware.org:default"
    ]
  },
  {
    "group": "management.vmware.org/v1",
    "kind": "Mgr",
    "name": "Manager2",
    "hierarchy": [
      "roots.orgchart.vmware.org:default"
    ]
  }
]
```

##### GET single softlink API `/leader/Role`
- GET `http://localhost:8080/leader/Role'

Example graph:
```
                  Root
                 /     \
                /       \
               /          \
            Role <-------  Leader
                softlink

```
- response body resemble the example below.

```
  {
    "group": "role.vmware.org/v1",
    "kind": "Executive",
    "name": "Role1",
    "hierarchy": [
      "roots.orgchart.vmware.org:default"
    ]
  }
```

##### GET multiple link API `/leader/HRs`
- GET `http://localhost:8080/leader/HRs'

Example graph:
```
                  Root
                 /     \
                /       \
               /         \
           ____________   \
          |  HR        |    Leader
          |  /\        |        |
          | /  \       |        / softlink
          | |   |      |        /
          | HR1 HR2 .. | <-----/
          |            |
          |____________|
```

- response body resemble the example below.

```
[
  {
    "group": "hr.vmware.org/v1",
    "kind": "HumanResources",
    "name": "HR1",
    "hierarchy": [
      "roots.orgchart.vmware.org:default"
    ]
  },
  {
    "group": "hr.vmware.org/v1",
    "kind": "HumanResources",
    "name": "HR2",
    "hierarchy": [
      "roots.orgchart.vmware.org:default"
    ]
  }
]
```

##### PUT Leader status API `/leader/Status`
- PUT `http://localhost:8080/leader/Status'

```
  {
    "status": {
     "IsOnVacations": true,
     "DaysLeftToEndOfVacations":1
    }
  }
```

##### GET Leader status API `/leader/Status`
- GET `http://localhost:8080/leader/Status'

```
  {
    "status": {
     "IsOnVacations": true,
     "DaysLeftToEndOfVacations":1
    }
  }
```

### Access the APIs through the shim layer.

For the full set of APIs: getting_started/WorkingWithShimLayer.md.

Example APIs for above Leader spec:

**Multiple child**
1. To add a EngManager to Leader `/leader/EngManagers` use the following pattern,

   ```_, err = ceo.AddEngManagers(context.TODO(), <managerObject>)```

2. To GET all the EngManagers `/leader/EngManagers` that is added to leader object.

   ```engManagers, err := ceo.GetAllEngManagers(context.TODO())```

**Link**
1. To add a softlink between two `Leader -> Role` nodes, the API looks like the following.

   ```err = ceo.LinkRole(context.TODO(), <roleObject>)```

2. For GET role object from which is linked above

   ```getExecRole, err := ceo.GetRole(context.TODO())```

**Status**
1. To update the status of CEO(Leader) object `/leader/Status`

   ```err = ceo.SetStatus(context.TODO(), <statusObject>)```

2. To GET the updated status of Leader object, use the API below

   ```updatedStatus, err := ceo.GetStatus(context.TODO())```



FAQs
1. Can we also display the spec of a linked object in the response body?
    - No. We can only view the metadata and parent information of the linked object.
2. Can we add linked object as query param to a parent GET API URI?
    - No. We don't support query param. We support only querying directional relationships.
