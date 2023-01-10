This test sequence will show how to create a datamodel object from Org Chart Datamodel, in your system.

## Prerequisites

* [Org Chart Datamodel is initialized and added to your application](AppLocalDatamodel.md)
* [Application workspace that is initialized by Nexus SDK](AppLocalDatamodel.md#setup-application-workspace)
* [Application has an operator watch for Org Chart Datamodel objects](WorkingWithDatamodel.md)

## Workflow

* Monitor Test Application Logs
* Create Leader object in Org Chart Datamodel

## Setup Workspace
If you're working through a complete workflow, you may find the following snippet useful to set up the required variables. 

```
export APP_NAME=${APP_NAME:-test-app-local}
export NAMESPACE=${NAMESPACE:-default}
```

<!-- nexus-specific exports
```
# store the current directory before we `cd` into the app dir
export DOCS_INTERNAL_DIR=$PWD/docs/_internal
```
-->

## Monitor Test Application Logs

The operator we created for this datamodel node will be notified of the create event.

* Stream the logs from your app to see the events being received by the operator
```shell
kubectl logs -f $(kubectl get pods -n "${NAMESPACE}" -l control-plane="${APP_NAME}" --no-headers | awk '{print $1}') -n "${NAMESPACE}"
```

<!--
```
kubectl port-forward svc/nexus-proxy-container 45192:80 -n ${NAMESPACE}
# create a leader object
kubectl -s localhost:45192 apply -f $DOCS_INTERNAL_DIR/leader_obj.yaml
# verify that the notification is consumed by the controller and a file is created with name 'Kind_name'
kubectl exec deploy/$APP_NAME -c $APP_NAME -- cat /tmp/Leader_default
```
-->

##  Create Leader object in Org Chart Datamodel

```shell
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: orgchart.vmware.org/v1
kind: Root
metadata:
  name: default' | kubectl -s nexus-api-gw:80 apply -f -"
```

```shell
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: management.vmware.org/v1
kind: Leader
metadata:
  name: default
  labels:
    roots.orgchart.vmware.org: default
spec:
  designation: Chief Executive Officer
  name: Raghu
  employeeID: 1' | kubectl -s nexus-api-gw:80 apply -f -"
```

##  Create Leader Object via REST API

```shell
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- curl -X PUT http://nexus-api-gw/root/default -H 'Content-Type: application/json'
```

```shell
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- curl -X PUT http://nexus-api-gw/root/default/leader/default -H 'Content-Type: application/json' -d '{"designation": "Chief Executive Officer","employeeID":1, "name":"Raghu"}
```

*NOTE: NAMESPACE is where Nexus Runtime microservices are running*

## Create Mgr Object via REST API

```shell
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- curl -X PUT http://nexus-api-gw/root/default/leader/default/mgr/mgr1 -H 'Content-Type: application/json' -d '{"employeeID": 10, "name": "John"}'
```

*NOTE: NAMESPACE is where Nexus Runtime microservices are running*

## Create Dev Object via REST API

```shell
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- curl -X PUT http://nexus-api-gw/root/default/dev/clark?management.Mgr=mgr1 -H 'Content-Type: application/json' -d '{"employeeID": 15, "name": "Clark"}' 
```

*NOTE: NAMESPACE is where Nexus Runtime microservices are running*

## List Dev objects via REST API

```shell
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- curl -X GET http://nexus-api-gw/devs?management.Mgr=mgr1 -H 'Content-Type: application/json'
```

*NOTE: NAMESPACE is where Nexus Runtime microservices are running*

### Next Steps

#### [Configure Admin Runtime to create routes to your application](AdminRuntimeConfiguration.md)

#### [Configure Runtime for enabling CORS headers ](AddingCorsHeaders.md)
