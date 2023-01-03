
Workflow to build and deploy an application in Nexus SDK framework

## Workflow

## Setup Workspace
If you're working through a complete workflow, you may find the following snippet useful to set up the required variables. 

```shell
export GOPATH=${GOPATH:-$(go env GOPATH)}
```

```
export APP_NAME=${APP_NAME:-test-app-local}
export NAMESPACE=${NAMESPACE:-default}

# if you're using a kind cluster
export KIND_CLUSTER_NAME=${KIND_CLUSTER_NAME:-kind}
```

### Go to App Workspace
Create and `cd` to your workspace directory for your test application, **in your $GOPATH/src**
```
mkdir -p $GOPATH/src/$APP_NAME && cd $GOPATH/src/$APP_NAME
```

### Build Application
```
make build
```

### Deploy Application
* For "kind" based Kubernetes cluster
    ```
    make deploy CLUSTER=$KIND_CLUSTER_NAME
    ```
* For other K8s clusters on the cloud, reachable via kubectl
    * Publish the image to docker registry. Ensure you're logged into the docker registry.
        ```shell
        make publish
        ```
    * Deploy the App
        ```shell
        make deploy
        ```
### Next Steps

#### [Test your Application -- Based on OrgChart Datamodel](TestingOrgChartDatamodel.md)

#### [Configure Admin Runtime to create routes to your application](AdminRuntimeConfiguration.md)

#### [Configure Runtime for enabling CORS headers ](AddingCorsHeaders.md)
