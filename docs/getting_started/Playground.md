# Playground

Welcome to TL;DR tutorial on the Nexus framework.

This tutorial will walk you through the fundamental aspects of Nexus.

The goal is to give you a taste on the most interesting and impactful aspects of the framework in the shortest possible time.

[Install Nexus CLI](#install-nexus-cli)

[Build a datamodel](#build-a-datamodel)

[Install datamodel](#install-datamodel)

[Access datamodel](#access-datamodel)

[Play with datamodel](#play-with-datamodel)

## Install Nexus CLI 

1. Download Nexus CLI

    ```
    curl -fsSL https://raw.githubusercontent.com/vmware-tanzu/graph-framework-for-microservices/main/cli/get-nexus-cli.sh -o get-nexus-cli.sh
    bash get-nexus-cli.sh
    ```
    <details><summary>FAQs</summary>
      
    To install the specific version
    ```
    bash get-nexus-cli.sh --version <version-tag> 
    ``` 
    
    To install the specific version and the user given destination directory
    ```
    bash get-nexus-cli.sh --version <version-tag> --dst_dir <destination-directoy-path>
    ``` 
	
    </details>
    


2. Verify your environment meets the expected pre-requisites

   ```
   nexus prereq verify
   ```

    <details><summary>See Pre-requisites</summary>

    a. To list all relevant pre-requisites:

        nexus prereq list

    </details>

## Build a Datamodel

Lets define a datamodel to implement well known facet in our work: Organization Chart

1. Create a workspace directory
    ```
    mkdir -p $HOME/test-datamodel/orgchart && cd $HOME/test-datamodel/orgchart       
    ```

2. Initialize workspace to specify datamodel
    ```
    nexus datamodel init --name orgchart --group orgchart.org
    ```

3. Write datamodel to implement an organization chart

   You can chose to write a datamodel from scratch. Detailed grammar and notes [here](../../compiler/DSL.md)

   Since we are in a playground workflow, feel free to use this pre-cooked datamodel implementing an organization chart. 

```shell
echo 'package root

import (
	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
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

4. Compile datamodel

   ```
   nexus datamodel build --name orgchart
   ```

## Install Datamodel

### Pre-requisites

**The following steps requires a running Kubernetes cluster >= version 1.19**

1. Install Nexus Runtime

```
nexus runtime install --namespace default
```

2. Install datamodel

   <details><summary>If your Kubernetes cluster is running on Kind, execute the following </summary>

   ```
   kind load docker-image orgchart:latest --name <kind cluster name>
   ```
   </details>


   ```
   nexus datamodel install image orgchart:latest --namespace default
   ```

## Access datamodel

1. Enable connectivity to Nexus API Gateway

```
kubectl port-forward svc/nexus-api-gw 5000:80 -n default &
```

2. GraphQL dashboard is available [here](http://localhost:5000/apis/graphql/v1)


3. REST API Explorer is available [here](http://localhost:5000/orgchart.org/docs#/)


## Play with datamodel

1. Create a Leader in your Organization

```shell
curl -X PUT -H 'Content-Type: application/json' -d '{"designation": "CTO","name":"foo"}' http://localhost:5000/leader/MyLeader
```

2. Create a Manager reporting to the Leader

```shell
curl -X POST -H 'Content-Type: application/json' -d '{"apiVersion":"root.orgchart.org/v1","kind":"Manager","metadata":{"labels":{"leaders.root.orgchart.org":"MyLeader"},"name":"Manager1"},"spec":{"designation":"Manager","name":"bar"}}'  http://localhost:5000/apis/root.orgchart.org/v1/managers 
```

3. Hire an Engineer reporting to the Manager

```shell
echo 'apiVersion: root.orgchart.org/v1
apiVersion: root.orgchart.org/v1
kind: Engineer
metadata:
  name: Engineer1
  labels:
    leaders.root.orgchart.org: MyLeader
    managers.root.orgchart.org: Manager1
spec:
  designation: Engineer
  name: zoo' > engineers.yaml
```

```shell
nexus login -s localhost:5000 --in-secure
nexus apply -f engineers.yaml
```

4. Access your organization chart through GraphQL [here](http://localhost:5000/apis/graphql/v1)


3. Access your organization chart through REST API Explorer [here](http://localhost:5000/orgchart.org/docs#/)
