# API Gateway to Custom backend services

API Gateway exposes REST API to interact to the UserDefined Datamodel objects installed as CRDs, `Route` adds the support for exposing any custom microservice running on tenant runtime namespace as a API Gateway REST endpoint.

![CustomAPIEndpoints](.content/images/customapi.jpg)

## Pre-Requisites for the microservices

Currently to expose a microservice using `Route` under API Gateway, it should fulfill below requirements

* The microservice should be deployed on the nexus runtime namespace

* The microservice should have `GET` `/` endpoint with HTTP 200 response code


## Workflow for Exposing a microservice

### Deploying and exposing a microservice to APIGateway

* Deploying a microservice with a backend service reachable in the tenant namespace

For this workflow using this testapp as a example

```shell
git clone https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/nexus-testapp-external-svc.git
cd nexus-testapp-external-svc
make build
make deploy KIND_CLUSTER_NAME=kind
```

* Creating RootObjects before creating `Route`:
```shell
kubectl port-forward $(kubectl get pods -n $NAMESPACE | grep api-gw | awk '{print $1}') -n $NAMESPACE 5000:80 & echo $! > tmp-pf.pid
```

### Create `Route` object to expose under APIGateway:

* Specifiction Defintion
```yaml
Spec:
    service:
       Name: Name of the microservice the route should be forwarded
       Port: Port in which the microservice is listening
       Scheme: Scheme of the service to be deployed ( Http/Https)
    resource:
        Name: User Defined name for a particular Set of Service URLs the route should be registered to
    uri: the URI which should be exposed to the outer world from a service (/* for wildcard support)
```

* Creating `Route` object to expose microservice `testappserver` under APIGateway
```shell
echo 'apiVersion: route.nexus.vmware.com/v1
kind: Route
metadata:
  name: custom
spec:
  service:
    name: testappserver
    port: 80
    scheme: Http
  resource:
    name: custom
  uri: "/*"' | kubectl -s localhost:5000 apply -f -
```
### Accessing the customAPI added through APIGateway

* Verify the APIGateway added routes to expose the microservice

Access `http://localhost:5000/` and verify the route is added

![CustomAPIEndpointsVerify](.content/images/VerifyCustomAPI.png)

* Accessing a REST endpoint under the service

URLScheme: `/apis/<resource.Name>/v1/<uri>` the service will be exposed in this prefix

For Example: to reach `teststr` endpoint under `testappserver` exposed by the `Route` created with `resource.Name: custom` , we can use the below URL

```shell
curl localhost:5000/apis/custom/v1/teststr
```

