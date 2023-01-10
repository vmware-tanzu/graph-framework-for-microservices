Nexus Datamodel provides the framework, language, toolkit and runtime to implement state/data required by application, that is:

* hierarchical
* distributed
* consistent
* customizable

Nexus datamodel framework supports formalization of objects/nodes and their relationship in a spec that is expressed using Golang syntax and is fully compliant with the Golang compiler.

**NOTE: If you would prefer to follow a pre-cooked Datamodel, as you try this workflow, execute commands from the code sections below** 

## Workflow

This guided workflow will walk you through setting up a datamodel that is local to your application.

* [Setup Application Workspace](AppLocalDatamodel.md#setup-application-workspace)

* [Setup Datamodel Workspace](AppLocalDatamodel.md#setup-datamodel-workspace)

* [Define your datamodel spec](AppLocalDatamodel.md#define-your-datamodel-specification)

* [Build datamodel](AppLocalDatamodel.md#build-datamodel)

* [Add datamodel to App](AppLocalDatamodel.md#add-datamodel-to-app)

## Setup Workspace
If you're working through a complete workflow, you may find the following snippet useful to set up the required variables. 

```shell
export GOPATH=${GOPATH:-$(go env GOPATH)}
```

```
export APP_NAME=test-app-local
export DATAMODEL_NAME=vmware
export DATAMODEL_GROUP=vmware.org
```

## Setup Application Workspace

1. Create and `cd` to your workspace directory for your test application, **in your $GOPATH/src**
    ```
    mkdir -p $GOPATH/src/$APP_NAME && cd $GOPATH/src/$APP_NAME
    ```

2. Create a fully functional and self-contained App workspace with a Helloworld datamodel.
    ```
    nexus app init --name $APP_NAME
    ```

## Setup Datamodel Workspace

Initialize a **"local"** datamodel workspace
```
nexus datamodel init --name $DATAMODEL_NAME --group $DATAMODEL_GROUP --local
```

## Define your datamodel specification

Start writing datamodel specification for your application.

```
 git clone https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/datamodel-examples.git
 cp -r datamodel-examples/org-chart-app-local/nexus/vmware/pkg datamodel-examples/org-chart-app-local/nexus/vmware/root.go nexus/vmware/
 rm -rf datamodel-examples
```

## Build datamodel
```
nexus datamodel build --name $DATAMODEL_NAME 
```

This generates libraries, types, runtime and metadata required to implement the datamodel at runtime.
## Add datamodel to App
```
nexus app add-datamodel --name $DATAMODEL_NAME
```

### Next Steps

#### 1. [Work with Datamodel objects](WorkingWithDatamodel.md)

#### 2. [Build and Deploy Application](BuildDeployApplication.md)

#### 3. [Test your Application -- Based on OrgChart Datamodel](TestingOrgChartDatamodel.md)

### [FAQ](FAQ.md)
