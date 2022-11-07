## Datamodel Workflow

Nexus Datamodel provides the framework, language, toolkit and runtime to implement state/data required by application, that is:

* hierarchical
* distributed
* consistent
* customizable

Nexus datamodel framework supports formalization of objects/nodes and their relationship in a spec that is expressed using Golang syntax and is fully compliant with the Golang compiler.

**NOTE: If you would prefer to follow a pre-cooked Datamodel, as you try this workflow, execute commands from the code sections below** 

This guided workflow will walk you through setting up a datamodel that is local to your application or Importable datamodel

* [Nexus Install](DatamodelWorkflow.md#nexus-install)

* [Nexus Pre-req verify](DatamodelWorkflow.md#nexus-pre-req-verify)

* [Setup Datamodel Workspace](DatamodelWorkflow.md#setup-datamodel-workspace)

* [Datamodel Build](DatamodelWorkflow.md#datamodel-build)
  
* [Datamodel Install](DatamodelWorkflow.md#datamodel-install)
  
* [Datamodel Playground](DatamodelWorkflow.md#datamodel-playground)


## Nexus Install

Install Nexus CLI

```
go install github.com/vmware-tanzu/graph-framework-for-microservices/cli/cmd/plugin/nexus@NPT-604-Migrate-CLI-Repo
```

## Nexus Pre-req Verify

Verify nexus sdk pre-requisites are satisfied

    nexus prereq verify

## Setup Datamodel Workspace

1. Create and `cd` to your workspace directory to create, compile and install datamodel
    ```
    mkdir -p $HOME/test-datamodel/orgchart && cd $HOME/test-datamodel/orgchart
    ```
     
1. Initialize datamodel workspace
    ```
    nexus datamodel init --name orgchart --group orgchart.org
    ```

1. Start writing datamodel specification for your application.
   
   **To understand the workflow we can use the below example datmodel. To write your own datamodel please refer** [here](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/docs/-/blob/master/Datamodel/DSL/README.md)
   
**Example Orgchart DSL**

   The Orgchart Application has 3 levels. 1. Leader, 2. Manager and 3. Engineer

```shell
echo 'package root

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

var LeaderRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/leader/{root.Leader}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri:     "/leaders",
			Methods: nexus.HTTPListResponse,
		},
	},
}

// nexus-rest-api-gen:LeaderRestAPISpec
type Leader struct {

	// Tags "Root" as a node in datamodel graph
	nexus.Node

	Name          string
	Designation   string
	DirectReports Manager `nexus:"children"`
}

type Manager struct {

	// Tags "Root" as a node in datamodel graph
	nexus.Node

	Name          string
	Designation   string
	DirectReports Engineer `nexus:"children"`
}

type Engineer struct {

	// Tags "Root" as a node in datamodel graph
	nexus.Node

	Name        string
	Designation string
}
' > $HOME/test-datamodel/orgchart/root.go
```

## Datamodel Build

   ```
   nexus datamodel build --name orgchart
   ```

This generates libraries, types, runtime and metadata required to implement the datamodel at runtime.


## Datamodel Install

<details><summary>Pre-requisites</summary>
To install datamodel we need to install nexus runtime as a pre-requisite
  
  [Install Runtime](RuntimeWorkflow.md)
</details>

   ```
   DOCKER_REPO=orgchart VERSION=latest make docker_build
   ```

   ```
   kind load docker-image orgchart:latest --name kind
   ```

   ```
   nexus datamodel install image orgchart:latest
   ```

## Datamodel Playground

#### [App Datamodel Usage Workflow](AppDatamodelUsageWorkflow.md)
