# Example operator app

Goal of this application is to show how to use operator together with nexus shim layer. In
controllers/management.vmware.org/leader_controller.go we can see Leader objects reconciler. Business logic is to link
default role to each Leader object for which event has occurred and which doesn't have role applied.

# Steps to run application

1. Setup env
```
export NAMESPACE=${NAMESPACE:-default}
export DATAMODEL_NAME=${DATAMODEL_NAME:-vmware}

nexus runtime install --namespace $NAMESPACE
nexus datamodel install name $DATAMODEL_NAME --namespace $NAMESPACE --debug

kubectl port-forward svc/nexus-proxy-container -n default 45192:80
kubectl port-forward deployment/api-gw   5000:80
```

2. Create root and default executive role (you can also use REST api
at http://localhost:5000/docs and create Root "default" and Executive "default-executive-role")

```
echo 'apiVersion: orgchart.vmware.org/v1
kind: Root
metadata:
   name: default' | kubectl -s localhost:5000 apply -f -
```

```
echo 'apiVersion: role.vmware.org/v1
kind: Executive
metadata:
  name: default-executive-role
  labels:
    roots.orgchart.vmware.org: default' | kubectl -s localhost:5000 apply -f -
```

3. Run example application
`go run main.go -host http://localhost:45192`

4. Add leader object (you can also use REST api at
http://localhost:5000/docs and create a leader there)
```
echo 'apiVersion: management.vmware.org/v1
kind: Leader
metadata:
  name: default
  labels:
    roots.orgchart.vmware.org: default
spec:
  designation: Chief Executive Officer
  name: Raghu
  employeeID: 1' | kubectl -s localhost:5000 apply -f -
```

5. Expected behaviour is to observe in application logs that role is linked to new Leader:
```
New event for Leader reconciler occured
Received event for leader node: default
Current role is nil, updating role to default
Leader's ceo role is: default-executive-role
```
