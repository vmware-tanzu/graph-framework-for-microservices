This workflow will walk you through the steps to interact and work with your imported datamodel.

**NOTE: If you would prefer to follow a pre-cooked Datamodel, as you try this workflow, execute commands from the code sections below**

## Workflow

* [Setup Application Workspace](WorkingWithCommonDatamodel.md#setup-application-workspace)

* [Add datamodel to App](WorkingWithCommonDatamodel.md#add-datamodel-to-app)

* [Create Operator for Datamodel Node](WorkingWithCommonDatamodel.md#create-operator-for-datamodel-node)

* [Add Business Logic](WorkingWithCommonDatamodel.md#add-business-logic)

* [Install Nexus Runtime](WorkingWithCommonDatamodel.md#install-nexus-runtime)

* [Install the datamodel on Nexus Runtime](WorkingWithCommonDatamodel.md#install-the-imported-datamodel-on-nexus-runtime)

## Setup Workspace
If you're working through a complete workflow, you may find the following snippet useful to set up the required variables.

```shell
export GOPATH=${GOPATH:-$(go env GOPATH)}
```

```
export APP_NAME=test-app-imported-dm
export DATAMODEL_NAME=vmware
export DATAMODEL_GROUP=vmware.org
export NAMESPACE=default
export GIT_MODULE=gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples.git/org-chart
export DATAMODEL_NODE=Leader
export PACKAGE=management
export DOCKER_REPO=testdatamodel-import
export VERSION=test
```

<!-- nexus-specific exports
```
# store the current directory before we `cd` into the app dir
export DOCS_INTERNAL_DIR=$PWD/docs/_internal
```
-->

## Setup Application Workspace

1. Create and `cd` to your workspace directory for your test application, **in your $GOPATH/src**
    ```
    mkdir -p $GOPATH/src/$APP_NAME && cd $GOPATH/src/$APP_NAME
    ```
2. Create a fully-functional and self-contained App workspace with a custom datamodel.
    ```
    nexus app init --name $APP_NAME
    ```

## Add datamodel to App

```
nexus app add-datamodel --name $DATAMODEL_NAME --package-name $GIT_MODULE --default=true
```

## Create Operator for Datamodel Node

```
nexus operator create --group $PACKAGE.$DATAMODEL_GROUP --kind $DATAMODEL_NODE --version v1
```

## Add Business Logic
```shell
cat <<< '--- controllers/management.vmware.org/leader_controller.go.orig    2022-05-11 16:15:27.000000000 +0530
+++ controllers/management.vmware.org/leader_controller.go      2022-05-11 16:15:37.000000000 +0530
@@ -18,6 +18,9 @@

 import (
        "context"
+       "fmt"
+       "os"
+       "path/filepath"

        "k8s.io/apimachinery/pkg/runtime"
        ctrl "sigs.k8s.io/controller-runtime"
@@ -33,9 +36,9 @@
        Scheme *runtime.Scheme
 }

-//+kubebuilder:rbac:groups=management.vmware.org.test-app-imported-dm.com,resources=leaders,verbs=get;list;watch;create;update;patch;delete
-//+kubebuilder:rbac:groups=management.vmware.org.test-app-imported-dm.com,resources=leaders/status,verbs=get;update;patch
-//+kubebuilder:rbac:groups=management.vmware.org.test-app-imported-dm.com,resources=leaders/finalizers,verbs=update
+//+kubebuilder:rbac:groups=management.vmware.org.test-app-imported-dm.com,resources=leaders,verbs=get;list;watch;create;update;patch;delete
+//+kubebuilder:rbac:groups=management.vmware.org.test-app-imported-dm.com,resources=leaders/status,verbs=get;update;patch
+//+kubebuilder:rbac:groups=management.vmware.org.test-app-imported-dm.com,resources=leaders/finalizers,verbs=update

 // Reconcile is part of the main kubernetes reconciliation loop which aims to
 // move the current state of the cluster closer to the desired state.
@@ -49,7 +52,19 @@
 func (r *LeaderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
        _ = log.FromContext(ctx)

-       // TODO(user): your logic here
+       // Business Logic: An event has occurred on the Leader Node.
+       var leader managementvmwareorgv1.Leader
+       if err := r.Get(ctx, req.NamespacedName, &leader); err != nil {
+               return ctrl.Result{}, client.IgnoreNotFound(err)
+       }
+       fmt.Printf("Received event for leader node: Name %s Spec %v\n", leader.Name, leader.Spec)
+
+       // create a file type_name
+       filename := leader.Kind + "_" + leader.ObjectMeta.Name
+       err := os.WriteFile(filepath.Join("/tmp", filename), []byte{}, 0644)
+       if err != nil {
+               fmt.Printf("Failed to write to /tmp/%s due to error: %v\n", filename, err)
+       }

        return ctrl.Result{}, nil
 }
' | patch controllers/management.vmware.org/leader_controller.go --ignore-whitespace
```
or, alternatively, replace `controllers/$PACKAGE.$DATAMODEL_GROUP/leader_controller.go` with [leader_controller.go](../_internal/leader_controller.go.patched_imported)

<!--
```
cp $DOCS_INTERNAL_DIR/leader_controller.go.patched_imported controllers/$PACKAGE.$DATAMODEL_GROUP/leader_controller.go
```
-->

## Install Nexus Runtime

To install Nexus Runtime with custom resource(cpu/memory) values, please refer [Nexus Sizing](NexusRuntimeSizing.md)

For kind clusters
```
nexus runtime install --namespace $NAMESPACE
```

## Install the imported datamodel on Nexus Runtime
<!-- TODO replace this with `nexus datamodel install` once we start supporting that -->

```
git clone git@gitlab.eng.vmware.com:nsx-allspark_users/nexus-sdk/datamodel-examples.git

cd datamodel-examples/org-chart && make docker_build

kind load docker-image $DOCKER_REPO:$VERSION  --name=$KIND_CLUSTER_NAME
nexus datamodel install image $DOCKER_REPO:$VERSION --namespace $NAMESPACE
```


### Next Steps

#### [Build and Deploy Application](BuildDeployApplication.md)

#### [Test your Application -- Based on OrgChart Datamodel](TestingOrgChartDatamodel.md)

#### [Configure Admin Runtime to create routes to your application](AdminRuntimeConfiguration.md)

####  [Configure Runtime for enabling CORS headers ](AddingCorsHeaders.md)
