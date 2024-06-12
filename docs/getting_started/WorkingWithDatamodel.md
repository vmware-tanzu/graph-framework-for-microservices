This workflow will walk you through the steps to interact and work with your datamodel.

**NOTE: If you would prefer to follow a pre-cooked Datamodel, as you try this workflow, execute commands from the code sections below** 

## Prerequisites

* Datamodel is initialized and added to your application
## Workflow

* [Create Operator for Datamodel Node](WorkingWithDatamodel.md#create-operator-for-datamodel-node)

* [Add Business Logic](WorkingWithDatamodel.md#add-business-logic)

* [SSL Secrets Configuration](WorkingWithDatamodel.md#ssl-secrets-configuration)

* [Install Nexus Runtime](WorkingWithDatamodel.md#install-nexus-runtime)

* [Install the datamodel on Nexus Runtime](WorkingWithDatamodel.md#install-the-app-datamodel-on-nexus-runtime)

* [Build and Deploy Application](WorkingWithDatamodel.md#build-and-deploy-application)

## Setup Workspace
If you're working through a complete workflow, you may find the following snippet useful to set up the required variables. 

```shell
export GOPATH=${GOPATH:-$(go env GOPATH)}
```

```
export APP_NAME=${APP_NAME:-test-app-local}
export DATAMODEL_NAME=${DATAMODEL_NAME:-vmware}
export DATAMODEL_GROUP=${DATAMODEL_GROUP:-vmware.org}
export NAMESPACE=${NAMESPACE:-default}
export PACKAGE=management
export DATAMODEL_NODE=Leader
```

<!-- nexus-specific exports
```
# store the current directory before we `cd` into the app dir
export DOCS_INTERNAL_DIR=$PWD/docs/_internal
```
-->

## Create Operator for Datamodel Node

1. Create and `cd` to your workspace directory for your test application, **in your $GOPATH/src**
    ```
    mkdir -p $GOPATH/src/$APP_NAME && cd $GOPATH/src/$APP_NAME
    ```

2. Create an operator for your datamodel node of interest
    ```
    nexus operator create --group $PACKAGE.$DATAMODEL_GROUP --kind $DATAMODEL_NODE --version v1 --datamodel $DATAMODEL_NAME
    ```

## Add Business Logic

```shell
cat <<< '--- controllers/management.vmware.org/leader_controller.go.orig  2022-05-11 16:07:39.000000000 +0530
+++ controllers/management.vmware.org/leader_controller.go      2022-05-11 16:08:00.000000000 +0530
@@ -18,6 +18,9 @@

 import (
        "context"
+       "fmt"
+       "os"
+       "path/filepath"

        "k8s.io/apimachinery/pkg/runtime"
        ctrl "sigs.k8s.io/controller-runtime"
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
or, alternatively, replace `controllers/$PACKAGE.$DATAMODEL_GROUP/leader_controller.go` with [leader_controller.go](../_internal/leader_controller.go.patched_local)

<!--
```
cp $DOCS_INTERNAL_DIR/leader_controller.go.patched_local controllers/$PACKAGE.$DATAMODEL_GROUP/leader_controller.go
```
-->

## SSL Secrets Configuration

This will be an optional step. Please execute this, if you want to start nexus-api-gateway with secure port.
```shell
  bash create_webhook_signed_cert.sh --service api-gw --namespace $NAMESPACE --secret  api-gw-server-cert
```
*NOTE: Please download the script from here:  gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/docs/example/create_webhook_signed_cert.sh

## Install Nexus Runtime

To install Nexus Runtime with custom resource(cpu/memory) values, please refer [Nexus Sizing](NexusRuntimeSizing.md)

For kind clusters

```
nexus runtime install --namespace $NAMESPACE
```


This will install Nexus Runtime microservices to your kubernetes cluster.

*NOTE: Optionally override namespace "default" with a preferred namespace to install nexus runtime*

## Install the datamodel on Nexus Runtime
```
nexus datamodel install name $DATAMODEL_NAME --namespace $NAMESPACE
```

### Next Steps

#### 1. [Build and Deploy Application](BuildDeployApplication.md)

#### 2. [Test your Application -- Based on OrgChart Datamodel](TestingOrgChartDatamodel.md)
