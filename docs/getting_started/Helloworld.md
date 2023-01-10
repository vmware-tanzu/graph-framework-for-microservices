Helloworld Datamodel is a use and throw playground that will give you a self-guided tour of the Nexus SDK, its capabilities and overall experience.

**[Pre-Requisites](README.md#before-getting-started)**

## Workflow

* [Setup Application Workspace](Helloworld.md#setup-workspace)

* [Install Nexus Runtime](Helloworld.md#install-runtime)

* [Build & Run Test Application](Helloworld.md#build-run-test-application)

## Setup Workspace
If you're working through a complete workflow, you may find the following snippet useful to set up the required variables. 

```shell
export GOPATH=${GOPATH:-$(go env GOPATH)}
```

```
export APP_NAME=test-app
export DATAMODEL_NAME=helloworld
export DATAMODEL_GROUP=helloworld.com
export NAMESPACE=default

# if you're using a kind cluster
export KIND_CLUSTER_NAME=${KIND_CLUSTER_NAME:-kind}
```

<!-- nexus-specific exports
```
# store the docs/_internal directory before we `cd` into the app dir
export DOCS_INTERNAL_DIR=$PWD/docs/_internal
```
-->

### Setup Test Application Workspace

1. Create and `cd` to your workspace directory for your test application, **in your $GOPATH/src**
    ```
    mkdir -p $GOPATH/src/$APP_NAME && cd $GOPATH/src/$APP_NAME
    ```

2. Create a fully functional and self-contained App workspace with a Helloworld datamodel.

    ```
    nexus app init --name $APP_NAME
    ```
    <details><summary>Commentary</summary>
    App-init creates and bootstraps an application workspace for Golang microservice, with following inbuilt attributes:

    * It is cloud native
    * Buildable
    * Deployable to K8s cluster
    * CI Ready
    * Incorporates best practices of cloud native app development, like:
        * Image scan
        * Race / Lint enabled
        * Coverage scan enabled

    Refer: [Microservice boostrapped by App Init](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api-gw)

    App Init will create files and directories in the directly in which the command is invoked.
    </details>
3. Setup Datamodel Workspace
    ```
    nexus datamodel init --name $DATAMODEL_NAME --group $DATAMODEL_GROUP --local
    ```
    <details><summary>Commentary</summary>
    Datamodel Init initializes a datamodel instance.
    A Datamodel instance is a self-contained directory that will host:

    * Hosts the datamodel spec expressed using Nexus DSL
    * Build artifacts and targets
    * Generated artifacts
        * libraries for App development
        * runtime installation artifacts
        * test / simulation artifacts 

    Refer: [An example datamodel instance](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples/-/tree/master/org-chart)

    Datamodel instance is a self-contained and complete Go package.

    * As such it can be exported as its own source versioned / git repo.
      It can then be imported by any application across the product as a Go package.

    * It can be local to an application, as such does not need to be exported.
      Local datamodel instance should be commited and source version along with the application.


    "--local" flag indicates the intent to create a App Local datamodel instance.

    </details>

### Build Helloworld Datamodel

```
nexus datamodel build --name $DATAMODEL_NAME
```

<details><summary>Commentary</summary>
Datamodel Build generates library, runtime, test and other relevant artifacts from the defined datamodel spec.

The artifacts are store in the nexus/helloworld directory.

The artifacts include, but is not limited to following directories:

* crds/  -  K8s CRDs generated for each datamodel node. [Refer](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples/-/tree/master/org-chart/build/crds)
* apis/  -  API library code generated for datamodel nodes. This library is built on top of K8s client / runtime libraries. [Refer](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples/-/tree/master/org-chart/build/apis)
* nexus-client/ - A convenience / shim library that presents a hierarchical view of the datamodel. [Refer](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples/-/tree/master/org-chart/build/nexus-client)

</details>

## Add datamodel to App
```
nexus app add-datamodel --name $DATAMODEL_NAME
```
<details><summary>Commentary</summary>

Add a built or imported datamodel to the local application.

This will setup and resolve build / module dependencies, in go.mod etc.

When this is done, App is ready to consume datamodel through the Nexus libraries.
</details>

## Install Runtime

### Install Nexus Runtime

For kind clusters
```
nexus runtime install --namespace $NAMESPACE
```


<details><summary>Commentary</summary>

Nexus runtime install will deploy nexus runtime microservices to your kubernetes cluster.

The goal fo the nexus runtime is to provide a distributed, eventually consistent, scalable runtime for application datamodel / stage.

The nexus runtime microservices, among many things, will deploy the following:

* a persistent datastore (etcd)
* an api server (k8s api server)
* proxy microserver for api service

Refer: [Design](design/Nexus-Runtime.md)
</details>

*NOTE: Optionally override namespace "default" with a preferred namespace to install nexus runtime*

### Install Helloworld Datamodel on Nexus Runtime
```
nexus datamodel install name $DATAMODEL_NAME --namespace $NAMESPACE
```

<details><summary>Commentary</summary>

Datamodel install will deploy CRD's and other manifests generated from datamodel spec.

This usually introduces now datamodel nodes, new API's, RBAC rules etc.

</details>

## Build & Run Test Application

### Create / init an **operator** for a Datamodel node

```
nexus operator create --group root.$DATAMODEL_GROUP --kind Root --version v1 --datamodel $DATAMODEL_NAME
```

<details><summary>Commentary</summary>

Operator Create creates an k8s operator for the requested node in the datamodel.

This operator is built on the well-known kubebuilder framewok and as such, provides a out-of-the-box / ready-to-use operators
for K8s CRD's, which is how datamodel nodes are represented in the runtime.

Refer: [Example Operator that watches for K8s CRD Types](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api-gw/-/blob/master/controllers/customresourcedefinition_controller.go)

</details>

### Add the operator business logic

##### Patch the Root controller
```shell
cat <<< '--- root_controller.go.orig    2022-05-11 15:53:02.000000000 +0530
+++ root_controller.go  2022-05-10 12:09:48.000000000 +0530
@@ -18,8 +18,10 @@

 import (
        "context"
-
+       "fmt"
        "k8s.io/apimachinery/pkg/runtime"
+       "os"
+       "path/filepath"
        ctrl "sigs.k8s.io/controller-runtime"
        "sigs.k8s.io/controller-runtime/pkg/client"
        "sigs.k8s.io/controller-runtime/pkg/log"
@@ -49,7 +51,18 @@
 func (r *RootReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
        _ = log.FromContext(ctx)

-       // TODO(user): your logic here
+       var root roothelloworldcomv1.Root
+       if err := r.Get(ctx, req.NamespacedName, &root); err != nil {
+               return ctrl.Result{}, client.IgnoreNotFound(err)
+       }
+       fmt.Printf("Received root node: Name %s Spec %v\n", root.Name, root.Spec)
+
+       // create a file type_name
+       filename := root.Kind + "_" + root.ObjectMeta.Name
+       err := os.WriteFile(filepath.Join("/tmp", filename), []byte{}, 0644)
+       if err != nil {
+               fmt.Printf("Failed to write to /tmp/%s due to error: %v\n", filename, err)
+       }

        return ctrl.Result{}, nil
 }
 ' | patch controllers/root.helloworld.com/root_controller.go --ignore-whitespace
```
or, alternatively, replace `controllers/root.$DATAMODEL_GROUP/root_controller.go` with [root_controller.go](../_internal/root_controller.go.patched)

<!--
```
cp $DOCS_INTERNAL_DIR/root_controller.go.patched controllers/root.$DATAMODEL_GROUP/root_controller.go
```
-->

### Build Application
```
make build
```

### Deploy Application
* For "kind" based Kubernetes cluster
    ```
    make deploy CLUSTER=$KIND_CLUSTER_NAME
    ```

  <details><summary>Commentary</summary>

  Deploys the $APP_NAME as a K8s Deployment in the Kind based K8s cluster, in namespace "default".

  Alternatively override the namespace to the desired K8s namespace.

  </details>

* For other K8s clusters on the cloud, reachable via kubectl

    <details><summary>Commentary</summary>

    Deploys the $APP_NAME as a K8s Deployment in the K8s cluster, in namespace "default".

    Alternatively override the namespace to the desired K8s namespace.

    NOTE: The docker image of the application has to published to a publically reachable image repository, before it the microservice can be deployed.

    </details>

    * Publish the image to docker registry. Ensure you're logged into the docker registry. 
        ```shell
        make publish
        ```
    * Deploy the App
        ```shell
        make deploy
        ```
NOTE: Optionally override namespace "default" with a preferred namespace to install nexus runtime

## Test your Application

* Stream the logs from your application
```shell
kubectl logs -f $(kubectl get pods -l control-plane="${APP_NAME}" --no-headers | awk '{print $1}')
```

* CRUD datamodel node in K8s cluster to trigger App business logic

```shell
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: root.helloworld.com/v1
kind: Root
metadata:
  name: myroot3
spec:
  myInt: 100
  configGvk:
    group: config.helloworld.com
    kind: Config
    name: myConfig
  runtimeGvk:
    group: runtime.helloworld.com
    kind: Runtime
    name: myruntime
  inventoryGvk:
    group: inventory.helloworld.com
    kind: Inventory
    name: myinventory' | kubectl -s nexus-apiserver:8080 apply -f -"
```

<!-- For CI purposes: verify the controller consumes the notification
```
kubectl port-forward svc/nexus-proxy-container 45192:80 -n ${NAMESPACE}
# create a root object
kubectl -s localhost:45192 apply -f $DOCS_INTERNAL_DIR/root_obj.yaml
# verify that the notification is consumed by the controller and a file is created with name 'Kind_name'
kubectl exec deploy/$APP_NAME -- cat /tmp/Root_myroot3
```
-->
*NOTE: NAMESPACE is where Nexus Runtime microservices are running*
