# CRUD app
This is example application showing CRUD operation using nexus client which includes shim layer for
automatic graph operations such as child/link resolution, updating parent with newly created child, adding/removing
soflinks, recurisve delete of object and all it's children.

Example is based on the following [datamodel](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples/-/tree/master/org-chart)

To run this application:
 - make sure nexus-apiserver is available, for example using port forward `kubectl port-forward svc/nexus-proxy-container -n default 45192:80`
 - apply CRDs definitions `kubectl -s localhost:45192 apply -f crds/`
 - run app using: `go run main.go -host http://<host>:<port>` for example `go run main.go -host http://localhost:45192`
