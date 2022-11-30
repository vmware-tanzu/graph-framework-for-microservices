# Playground

Welcome to TL;DR tutorial on the Nexus framework.

This tutorial will walk you through the fundamental aspects of Nexus.

The goal is to give you a taste on the most interesting and impactful aspects of the framework in the shortest possible time.

  * [Install Nexus CLI](#install-nexus-cli)
  * [Build a Datamodel](#build-a-datamodel)
  * [Install Datamodel](#install-datamodel)
  * [Access datamodel](#access-datamodel)
  * [Play with datamodel](#play-with-datamodel)

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

Lets define a datamodel to implement well known facet in our work: Organization Chart.

1. Create a workspace directory.
    ```
    mkdir -p $HOME/test-datamodel/orgchart && cd $HOME/test-datamodel/orgchart       
    ```

2. Initialize workspace to specify datamodel.
    ```
    nexus datamodel init --name orgchart --group orgchart.org
    ```

3. Write datamodel to implement an organization chart.

   You can choose to write a datamodel from scratch. Detailed grammar and notes [here](../../compiler/DSL.md)

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

4. Compile datamodel.

   ```
   nexus datamodel build --name orgchart
   ```

## Install Datamodel

### Pre-requisites 

**The following steps requires a running Kubernetes cluster >= version 1.19**

1. Install Nexus Runtime.

    ```
    nexus runtime install --namespace default
    ```

2. Install datamodel.

   <details><summary>If your Kubernetes cluster is running on Kind, execute the following </summary>

   ```
   kind load docker-image orgchart:latest --name <kind cluster name>
   ```
   </details>

   ```
   nexus datamodel install image orgchart:latest --namespace default
   ```

## Access datamodel

1. Enable connectivity to Nexus API Gateway.

    ```
    kubectl port-forward svc/nexus-api-gw 5000:80 -n default &
    ```

2. GraphQL dashboard is available [here](http://localhost:5000/apis/graphql/v1)

3. REST API Explorer is available [here](http://localhost:5000/orgchart.org/docs#/)

## Play with datamodel

1. Create a Leader in your Organization.

    ```shell
    curl -X PUT -H 'Content-Type: application/json' -d '{"designation": "CTO","name":"foo"}' http://localhost:5000/leader/MyLeader
    ```

2. Create a Manager reporting to the Leader.

    ```shell
    curl -X POST -H 'Content-Type: application/json' -d '{"apiVersion":"root.orgchart.org/v1","kind":"Manager","metadata":{"labels":{"leaders.root.orgchart.org":"MyLeader"},"name":"Manager1"},"spec":{"designation":"Manager","name":"bar"}}'  http://localhost:5000/apis/root.orgchart.org/v1/managers 
    ```

3. Hire an Engineer reporting to the Manager.

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

5. Access your organization chart through REST API Explorer [here](http://localhost:5000/orgchart.org/docs#/)

## Replicate datamodel from Nexus api-server to base K8s api-server

1. Create NexusEndpoint configuration with destination host and port details. This deploys one instance of nexus-connector that syncs objects to the desired destination endpoint.

    <details><summary>If your Kubernetes cluster is running on Kind, execute the following to get the destination IP and use https://[IP]:6443</summary>

    ```
    docker inspect <cluster-name>-control-plane | jq '.[].NetworkSettings.Networks["kind"].IPAddress'
    ```
    </details>

    ```shell
    echo 'apiVersion: connect.nexus.org/v1
    kind: NexusEndpoint
    metadata:
      name: default
      labels:
        nexus/is_name_hashed: "false"
        connects.connect.nexus.org: default
    spec:
      host: XXX 
      port: XXX' > $HOME/test-datamodel/orgchart/endpoint.yaml && kubectl -s localhost:5000 apply -f $HOME/test-datamodel/orgchart/endpoint.yaml
    ```

2. Install org-chart CRDs on the destination endpoint (base K8s api-server) and give cluster permissions for the API groups.

    ```
    cd $HOME/test-datamodel/orgchart/build/crds
    kubectl apply -f .
    ```

    ```shell
    echo 'apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
        name: nexus-connector-cr
    rules:
    - apiGroups:
        - "*"
      resources:
        - "*"
      verbs:
        - "*"
    ---
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
        name: nexus-connector-crb
    roleRef:
        apiGroup: rbac.authorization.k8s.io
        kind: ClusterRole
        name: nexus-connector-cr
    subjects:
      - kind: ServiceAccount
        name: default
        namespace: default' > $HOME/test-datamodel/orgchart/permissions.yaml && kubectl apply -f $HOME/test-datamodel/orgchart/permissions.yaml
    ```

3. Create the below-given replication-config to replicate `Manager1` to the destination endpoint (base K8s api-server).

    - **Note**: Fill in the `accessToken` spec field before creating the config.


    ```
    kubectl get secret $(kubectl get sa default -o yaml | yq -r '.secrets[0].name') -oyaml | yq '.data.token' | base64 -d
    ```

    ```shell
    echo 'apiVersion: connect.nexus.org/v1
    kind: ReplicationConfig
    metadata:
      name: one
      labels: 
          nexus/is_name_hashed: "false"
          connects.connect.nexus.org: default 
    spec:
      accessToken: XXXXX
      remoteEndpointGvk:
        group: connect.nexus.org
        kind: NexusEndpoint
        name: 4187f4f8437a5f4b8f4535c26d70443591b56856
      source:
        kind: Object
        object:
          name: Manager1
          objectType:
            group: root.orgchart.org
            kind: Manager
            version: v1
          hierarchical: true
          hierarchy:
            labels:
            - key: "leaders.root.orgchart.org"
              value: "MyLeader"
      destination:
        hierarchical: false' > $HOME/test-datamodel/orgchart/replication-config.yaml && kubectl -s localhost:5000 apply -f $HOME/test-datamodel/orgchart/replication-config.yaml
    ```

The manager object `Manager1` will now appear in base K8s api-server. Also, try update and delete on the manager object `Manager1` on the source and verify if it is reflected on the destination endpoint.

## Replicate datamodel from base K8s api-server to Nexus api-server

Step 1: Clone the latest connector.

git clone https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/connector.git


Step 2: Install standalone connector using the below command:

helm install -g nexus-connector/ --set-string global.namespace=default \
--set-string global.connector.tag=v0.0.14 \
--set-string global.statusReplication=DISABLE \
--set-string global.token=<token> --wait --debug


where:

connector.tag -> Latest working connector tag.
Token -> Tenant namespace token fetched from secret.

Step 3: Create NexusEndpoint CR. Click Sample NexusEndpoint CR

conn
Step 4: Create leader object.

apiVersion: management.vmware.org/v1
kind: Leader
metadata:
  name: default
spec:
    designation: CEO
    employeeID: 1
    name: Alice


Step 5: Create the below-given replication-config to replicate leader object to the destination endpoint.
Note: Refer here for the steps to fetch access token.

 apiVersion: connect.nexus.org/v1
  kind: ReplicationConfig
  metadata:
    name: one
  spec:
    accessToken: <token>
    destination:
      hierarchical: true
      hierarchy:
        labels:
        - key: "roots.orgchart.vmware.org"
          value: "default"
      objectType:
        group: management.vmware.org
        kind: Leader
        version: v1
    remoteEndpointGvk:
      group: connect.nexus.org
      kind: NexusEndpoint
      name: default
    source:
      kind: Type
      type:
        group: management.vmware.org
        kind: Leader
        version: v1


The leader object "default" will now appear on the destination endpoint. Also, try update and delete on the leader object "default" on the source and verify if it is reflected on the destination endpoint.