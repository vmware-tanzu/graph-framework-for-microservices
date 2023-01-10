Nexus Datamodel provides the framework, language, toolkit and runtime to implement state/data required by application, that is:

* hierarchical
* distributed
* consistent
* customizable

Nexus datamodel framework supports formalization of objects/nodes and their relationship in a spec that is expressed using Golang syntax and is fully compliant with the Golang compiler.

**NOTE: If you would prefer to follow a pre-cooked Datamodel, as you try this workflow, execute commands from the code sections below** 

## Workflow

* [Setup Datamodel Workspace](ImportableDatamodel.md#setup-datamodel-workspace)

* [Define your datamodel spec](ImportableDatamodel.md#define-your-datamodel-specification)

* [Build datamodel](ImportableDatamodel.md#build-datamodel)

* [Publish datamodel](ImportableDatamodel.md#publish-datamodel)

## Setup Workspace
Export the required environment variables

```shell
export GOPATH=${GOPATH:-$(go env GOPATH)}
```

```
export DATAMODEL_WORKSPACE_DIRECTORY=$GOPATH/src/test-app-imported-dm
export DATAMODEL_NAME=gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples.git/org-chart
export DATAMODEL_GROUP=vmware.org
```

## Setup Datamodel Workspace

* Create and `cd` to a workspace for datamodel, in your $GOPATH/src
    ```
    mkdir -p $DATAMODEL_WORKSPACE_DIRECTORY && cd $DATAMODEL_WORKSPACE_DIRECTORY
    ```

* Initialize the datamodel workspace
    ```
    nexus datamodel init --name $DATAMODEL_NAME --group $DATAMODEL_GROUP
    ```
## Define your datamodel specification

## Build datamodel
```
nexus datamodel build
```

This generates libraries, types, runtime and metadata required to implement the datamodel at runtime.
## Publish datamodel

For the datamodel to be useful, it has can be published to a Git repo for use by other applications.
It can also be installed directly on a K8s cluster.
### Publish datamodel to Git/Gitlab/Github

Create a git repo and get its URL. Then:

```
git init
git remote add origin git@<git-repo-url>
git add .
git commit -m "Initial commit"
git push -u origin master
```

### Next Steps

#### [Install Admin Runtime](AdminRuntimeInstall.md)

#### [Import your datamodel into you application](WorkingWithCommonDatamodel.md)

#### [Build and Deploy Application](BuildDeployApplication.md)

#### [Test your Application -- Based on OrgChart Datamodel](TestingOrgChartDatamodel.md)

#### [Configure Admin Runtime to create routes to your application](AdminRuntimeConfiguration.md)

#### [Configure Runtime for enabling CORS headers ](AddingCorsHeaders.md)

### [FAQ](FAQ.md)
