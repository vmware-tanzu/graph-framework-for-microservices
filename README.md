# Nexus SDK

Playground - Early Edition :)

Nexus SDK is an next-gen platform, framework and runtime. Its mission is to provide an extensible, distributed platform that will accelerate application development, simplify consumption, providing a highly functional, stable and feature rich platform, for all products in the AllSpark family.

This Dev/Tech preview implements Nexus Datamodel as K8s CRD's. Applications can consume Nexus datamodel using standard K8s libraries, just as they consume any K8s CRD.

We will demonstrate:

- Simple / intuitive DSL to specify Nexus Datamodel
- Nexus compiler that understands Nexus DSL and generates datamodel spec and libraries
- Nexus Datamodel implemention using K8s CRD's
- Nexus runtime to host Nexus Datamodel
- Ability to consume Nexus Datamodel using standard / opensource K8s libraries

This wiki will setup you up with your own playgound, to explore Nexus SDK in action.
The playground also provides a Sample App to interact with the datamodel.

**2 - ways to get started:**

- **[TLDR version -- Helloworld Datamodel](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#tldr-version-helloworld-datamodel)**

- **[Expert version -- Build-your-own Datamodel](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#build-your-own-datamodel)**

- **[Expert version -- Build centralized datamodel](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#create-your-own-datamodel-and-push-it-to-repo)**

## FAQ

Refert to FAQ's [here](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#faqs)

## Support

Reach out to Platform Team on [#nexus-sdk](https://vmware.slack.com/archives/C017KTHQ10X) slack channel for additional info and support.

## Pre-requisites

#### Go is installed in the environment and $GOPATH/src is valid. Recommended(and tested with) go 1.17+
```
    ➜ go version
    go version go1.17 darwin/amd64
```
#### $GOPATH/bin is in the unix PATH

```
export PATH=$PATH:$GOPATH/bin
```

#### Docker daemon is installed and running on the machine
```
    ➜ docker info
    Client:
    Context:    default
    Debug Mode: false
    Plugins:
      app: Docker App (Docker Inc., v0.9.1-beta3)
      buildx: Build with BuildKit (Docker Inc., v0.5.1-docker)
      scan: Docker Scan (Docker Inc., v0.6.0)
```
#### GoImport to work

1. Ensure you have environment variables set for pulling private golang repositories
```
go env -w GOPRIVATE=gitlab.eng.vmware.com
```
or
```
export GOPRIVATE=gitlab.eng.vmware.com
```
2. (optional) In rare cases if it fails due to go get using https instead of ssh://
```
git config --global --add url."git@gitlab.eng.vmware.com:".insteadOf "https://gitlab.eng.vmware.com/"
```

#### Permission to pull docker image from 284299419820.dkr.ecr.us-west-2.amazonaws.com repository

1. Ensure you are logged into "shared" AWS account

```
  aws sts get-caller-identity
  {
    "UserId": "**************************",
    "Account": "284299419820",
    "Arn": "arn:aws:iam::284299419820:user/xxxxxxxxxx@vmware.com"
}
```

2. Setup Docker credentials

```
docker login --username AWS -p $(aws ecr get-login-password --region us-west-2) 284299419820.dkr.ecr.us-west-2.amazonaws.com
```
or
```
flash repo login .
```

Note: if you are using any cloudgate based dev group , please run
```
clearaws
```
and retry step 2 once again.

#### Kubebuilder installed (optional)
```
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
```

#### The K8s cluster should be reachable from the terminal via "kubectl"

#### Compiled to work in OSX, Linux

## TLDR version -- Helloworld Datamodel

Playaround with a pre-built datamodel that has a Root Node with 3 children: Config, Runtime and Inventory

### Download Nexus CLI

```
 GOPRIVATE="gitlab.eng.vmware.com" go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/cmd/plugin/nexus@nexus-cli-dev
 ```
### Create your workspace/playground in $GOPATH/src directory

```
mkdir $WORKSPACE/helloworldapp   <-- $WORKSPACE is a directory under $GOPATH/src
cd $WORKSPACE/helloworldapp
```
## Application workspace init: Sets up files required to successfully create a Go microservice
Pre-requisites:
 - [LatestCLI](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#check-if-latest-nexus-cli-is-installed)
```
nexus app init --name helloworldapp --datamodel-init helloworld
```

This will prep the application to support consumers from datamodel.

### Build Datamodel

```
nexus datamodel build
```

## Install Nexus Runtime

```
nexus runtime install
```

## Install Datamodel

```
nexus datamodel install
```

### Add controllers to consume

```
nexus operator create --group root.helloworld.com --kind Root --version v1 --datamodel helloworld
```

edit controllers/root.helloworld.com/root_controller.go businesslogic

```
func (r *RootReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var config roothelloworldcomv1.Root
	if err := r.Get(ctx, req.NamespacedName, &config); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Printf("Received config node: Name %s Spec %v\n", config.Name, config.Spec)

	return ctrl.Result{}, nil
}
```

### Deploy sample application in running cluster

Pre-requisites:

- [DockerLogin](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#permission-to-pull-docker-image-from-284299419820dkrecrus-west-2amazonawscom-repository)

```
make build
make publish
make deploy
```

For Kind clusters
```
make build
make deploy CLUSTER=kind
```

After deploying you will see a pod called helloworldapp
```
kubectl get pods | grep helloworldapp
kubectl logs <pod-name>
```

## Create Datamodel Objects

**Root Node Object**
```
NAMESPACE=default kubectl exec -it $(kubectl get pods  --no-headers -l app=proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: root.helloworld.com/v1
kind: Root
metadata:
  name: myroot3
spec:
  myInt: 100
  config:
    group: config.helloworld.com
    kind: Config
    name: myConfig
  runtime:
    group: runtime.helloworld.com
    kind: Runtime
    name: myruntime
  inventory:
    group: inventory.helloworld.com
    kind: Inventory
    name: myinventory' | kubectl -s apiserver:8080 apply -f -"
```
**Config Node Object**

```
NAMESPACE=default kubectl exec -it $(kubectl get pods  --no-headers -l app=proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: "config.helloworld.com/v1"
kind: Config
metadata:
  name: myconfig
  namespace: default
spec:
  exampleStr: "foo"' | kubectl -s apiserver:8080 apply -f -"
```

**Runtime Node Object**

```
NAMESPACE=default kubectl exec -it $(kubectl get pods  --no-headers -l app=proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: "runtime.helloworld.com/v1"
kind: Runtime
metadata:
  name: myroot
  namespace: default
spec:
  myRuntimeInt: 20' | kubectl -s apiserver:8080 apply -f -"
```
**Inventory Node Object**

```
NAMESPACE=default kubectl exec -it $(kubectl get pods  --no-headers -l app=proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: "inventory.helloworld.com/v1"
kind: Inventory
metadata:
  name: myinventory
  namespace: default
spec:
  inventoryId: 1000' | kubectl -s apiserver:8080 apply -f -"
```
## Expected Results

Application subscribes to and watches for changes to Datamodel objects that are implemented in K8s CRDs.

## Build Your Own Datamodel

Playaround with templatized files that will allow you to define your own datamodel and see it in action.

### Download Nexus CLI

```
  GOPRIVATE="gitlab.eng.vmware.com" go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/cmd/plugin/nexus@nexus-cli-dev
```

### Create your workspace/playground in $GOPATH/src directory

```
mkdir $WORKSPACE/myapp   <-- $WORKSPACE is a directory under $GOPATH/src
cd $WORKSPACE/myapp
```

### Init Datamodel

Pre-requisites:
 - [LatestCLI](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#check-if-latest-nexus-cli-is-installed)

This will prep the playgound with required components and structure.

```
nexus app init --name myapp --datamodel-init customdatamodel
```

### Edit root.go Datamodel file with your desired fields in the Root Node

Refer to the [Nexus DSL Spec](https://confluence.eng.vmware.com/display/NSBU/Nexus+Platform#NexusPlatform-TL;DR;) for the grammar and structure of the Nexus DSL

```
vi nexus/customdatamodel/root.go
```

**Example** (Provided for your convenience)

Defines Root with a child node called ClusterStatus.

ClusterStatus has custom fields in the spece and also has a child node called ClusterStatusMetadata.

ClusterStatusMetadata is a leaf node with custom fields.

```
package root

import (
        "customdatamodel/nexus"
)

//Add your openapispec fields and crd fields here
type Root struct {
        nexus.Node

        //Add your child nodes and fields here

        ClusterStatus ClusterStatus `nexus:"child"`

}

type ClusterState string

type ClusterStatus struct {
        nexus.Node

        State           ClusterState
        Code            int32
        Message         string
        UpdateTimestamp string

        Metadata        ClusterStatusMetadata `nexus:"child"`
}

type ClusterStatusMetadata struct {
        nexus.Node

        Substate string
        Progress int32
}
```

### Build Datamodel

```
nexus datamodel build
```

## Install Nexus Runtime

```
nexus runtime install
```

## Install Datamodel

```
nexus datamodel install
```

### Add controllers to consume

```
nexus operator create --group root.customdatamodel.com --kind Root --version v1 --datamodel customdatamodel
```

#### Edit controllers/root.customdatamodel.com/root_controller.go file to include reconcile section.

```
func (r *RootReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var root rootcustomdatamodelcomv1.Root
	if err := r.Get(ctx, req.NamespacedName, &root); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Printf("Received root node: Name %s Spec %v\n", root.Name, root.Spec)

	return ctrl.Result{}, nil
}
```

### Deploy sample application in running cluster

- [DockerLogin](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#permission-to-pull-docker-image-from-284299419820dkrecrus-west-2amazonawscom-repository)

For Kind clusters
```
make build
make deploy CLUSTER=kind
```

For EKS clusters

```
make build
make publish
make deploy
```

After deploying you will see a pod called myapp
```
kubectl get pods | grep myapp
kubectl logs <pod-name>
```
## Create Datamodel Objects


**Root Node Object**

```
NAMESPACE=default kubectl exec -it $(kubectl get pods  --no-headers -l app=proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: "root.customdatamodel.com/v1"
kind: Root
metadata:
  name: myroot
  namespace: default
spec:
  clusterStatus:
    group: root.customdatamodel.com
    kind: Clusterstatus
    name: status' | kubectl -s apiserver:8080 apply -f -"
```
**Cluster Status Object**

```
NAMESPACE=default kubectl exec -it $(kubectl get pods  --no-headers -l app=proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: "root.customdatamodel.com/v1"
kind: Clusterstatus
metadata:
  name: status
  namespace: default
spec:
  state: "test"
  code: 10
  message: "mymessage"
  updateTimestamp: "1.2.3"
  metadata:
    group: root.customdatamodel.com
    kind: Clusterstatusmetadata
    name: metadata' | kubectl -s apiserver:8080 apply -f -"
```

**Cluster Status Metadata Object**

```
NAMESPACE=default kubectl exec -it $(kubectl get pods  --no-headers -l app=proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- bash -c "echo 'apiVersion: "root.customdatamodel.com/v1"
kind: Clusterstatusmetadata
metadata:
  name: metadata
  namespace: default
spec:
  substate: "test"
  progress: 10' | kubectl -s apiserver:8080 apply -f -"
```
## Expected Results

Application subscribes to and watches for changes to Datamodel objects that are implemented in K8s CRDs.

## Create your Own datamodel and push it to repo
**Step 0**: Create a Gitlab repository for storing datamodel and to import it as a _go_ _module_ from applications

For our example, we create a repo called _customdatamodel_ under the _nexus_ group in gitlab.eng.vmware.com

   * Repo name would be : `gitlab.eng.vmware.com/nexus/customdatamodel`

Cloning the repository to a folder in your workspace (under GOPATH)
  ```
  git clone git@gitlab.eng.vmware.com:nexus/customdatamodel.git
  ```

would provide us with intialized git repository under customdatamodel folder

**Step 1**: Create bare-datamodel for the user to edit and build to generate CRD and OpenAPISpecs:

```
cd customdatamodel
nexus datamodel init --repo gitlab.eng.vmware.com/nexus/customdatamodel --group customgroup.com
```

Note: use the correct name to avoid import errors while creating the application.

The above command would provide basic structure
```
├── Makefile
├── go.mod
├── nexus
│   └── nexus.go
├── nexus.yaml
└── root.go
```

**Step 2**:  Create and Build your Datamodel

  Step 2a (optional): Add your own nodes in root.go and create subdirectories to include your schema

  Step 2b: Build your datamodel to generate CRDs and APIs.

```
 nexus datamodel build
```
  Step 2c: Push the generated artifacts to the git repo so that applications may import them.

### Consume the Generated APIs and CRDs from the datamodel

**Step 3**: Create an app called _myapp_ under the $GOPATH/src directory to avoid datamodel resolution conflicts

```
$GOPATH/src --> $WORKSPACE
mkdir $WORKSPACE/src
cd $WORKSPACE/myapp
nexus app init --name myapp
```

**Step 4**:
 Step 4a (optional): If you want to associate the app with a default datamodel, you may use the following command to do so. Having a default DM set helps the nexus CLI decide the datamodel to use if a DM wasn't explicitly specified (as we'll see in the following command)
 ```
 # Add a datamodel
 nexus app add-datamodel --name customdatamodel --location gitlab.eng.vmware.com/nexus/customdatamodel --default=true
 ```

 Pre-requisites:

 * [GoImports](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#goimport-to-work)

 Step 4b: Create operators to consume the datamodel apis

```
nexus operator create --group root.customgroup.com --kind Root --version v1 --datamodel gitlab.eng.vmware.com/nexus/customdatamodel

# If the default datamodel had been set, we can safely skip the --datamodel argument above. The operator will assume that the GroupVersionKind (GVK) comes from the default DM
nexus operator create --group root.customgroup.com --kind Root --version v1

```

Note: The above command creates a folder controllers/root.customgroup.com/... , as each object created in the datamodel add a operator with your GroupVersionKind and --datamodel variable should point it to git repository name of the datamdoel to consume.

 Step 4c:
 Edit controllers/root.customgroup.com/root_controller.go file to include reconcile section.

```
func (r *RootReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var root rootcustomgroupcomv1.Root
	if err := r.Get(ctx, req.NamespacedName, &root); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	fmt.Printf("Received root node: Name %s Spec %v\n", root.Name, root.Spec)

	return ctrl.Result{}, nil
}
```

Note: rootcustomgroupcomv1 denotes the import group , please get the name from import appropriate to your module in the same file you are editing.

**Step 5**: Preparing the runtime for installing the application(microservice) to test the consumption of datamodel works

```
cd $WORKSPACE/myapp
nexus runtime install
```

**Step 6**: Installing the CRDs on the runtime for reconcilers on the application to start.

Please go to the repo cloned in **Step 1**.
```
cd customdatamodel
kubectl port-forward svc/proxy-container 45192:80 & echo $! > pid
kubectl -s localhost:45192 apply -f build/crds
kill $(cat pid) && rm -rf pid
```

**Step 7**: Build the application and deploy on the runtime
   ```
   make build
   ```

**Step 8**: Deploy application
    on Kind Cluster
```
make deploy CLUSTER=kind
```
Note: Please verify if you are current context and Steps 5 and 6 should be run successfully before calling deploy.

on EkS Cluster
Pre-req:
- [DockerLogin](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus/-/blob/nexus-sdk-dev/Nexus-SDK-README.md#permission-to-pull-docker-image-from-284299419820dkrecrus-west-2amazonawscom-repository)
```
make publish
make deploy
```
 Note: Please verify the docker login was set to appropriate shared group if you get 403 errors

## FAQ's

**1. How can I get additional help ?**

Reach out to Platform Team on [#nexus-sdk](https://vmware.slack.com/archives/C017KTHQ10X) slack channel for additional info and support.


**2. I am running into an issue with above steps. What info will the team need to debug my issue ?**

Please provide the following details, for us to give the best support:

```
go env; go version
```
```
which nexus; nexus version
```
```
aws sts get-caller-identity
```
```
docker ps
```
For Publish related issues
```
cat ~/.docker/config.json
<<<check if below line present>>>
"auths": {
		"284299419820.dkr.ecr.us-west-2.amazonaws.com"
```

For installation related issues

```
kubectl version; kubectl get pods -A
```


### Update latest tags for each components

To Update tags of each component
```
<edit> nexus-runtime/values.yaml
<section>
api_gateway:
    tag: v0.0.17
  validation:
    tag: v0.0.8
  connector:
    tag: v0.0.1
  controller:
    tag: v0.0.1
  api:
    tag: v0.0.2 <--- replace tag for your component
```
To update chart changes
```
git submodule update --init --remote
```

To use dev branch for testing chart changes
```
<edit> .gitmodules
<replace> branch = master with branch = <intended> of your repo
```

And run
```
git submodule update --init --remote
```

Please commit and push subcharts folder too along with other changes

### Helm Publish to Harbor Repo

To publish helm chart to nexus-runtime repository in ECR.
```
helm repo add "harbor-vmware" "https://harbor-repo.vmware.com/chartrepo/nexus" --username <vmware_username> --password <vmware_password>
```

Create Helm chart
```
rm -rf nexus-runtime/charts/*
helm dependency update
helm package nexus-runtime --version <version>
```

Publish helm chart to ECR
```
helm cm-push nexus-runtime-<version>.tgz harbor-vmware
```

### Install nexus-runtime from HELM using harbor chart

Login to Harbor registry
```
 helm repo add "harbor-vmware" "https://harbor-repo.vmware.com/chartrepo/nexus"
```

For deploying tenant
```
kubectl create ns <namespace> --label name=<namespace> --dry-run -o yaml | kubectl apply -f  -
helm install --wait nexus-runtime harbor-vmware/nexus-runtime --version <version> \
            --set-string global.namespace=<namespace>\
            --set-string global.registry=harbor-repo.vmware.com/nexus
```

* Note: global.repository=284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus if you are using EKS cluster

Please refer below section for available helm variables


For deploying admin namespace
```
kubectl create ns <namespace> --label name=<namespace> --nexus=admin --dry-run -o yaml | kubectl apply -f  -
helm install --wait nexus-runtime nexus-runtime harbor-vmware/nexus-runtime  --version <version> \
            --set-string global.namespace=newtestv\
            --set-string global.repository=harbor-repo.vmware.com/nexus\
            --set global.nexusAdmin=true
```

#### HELM Variables defintion

| Variable | Definition | Example  | Default |
| :---:   | :-: | :-: | :---: |
| global.namespace | Namespace where you want to deploy | --set global.namespace=default | default |
global.nexusAdmin | To deploy admin namespace Currently tsm can deployed in 2 modes (Tenant, Admin) | --set global.nexusAdmin=true | false
global.registry | To pull images from particular private registry | --set registry=harbor-repo.vmware.com | 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus |
global.connector.tag | To override connector image tag | --set global.connector.tag=<> |
global.api_gateway.tag | To override api_gateway image tag | --set global.api_gateway.tag=<> |
global.validation.tag | To override validation image tag | --set global.validation.tag |
global.controller.tag | To override connect-controller image tag | --set global.controller.tag |
global.runtimeEnabled | To enable all components of nexus runtime | --set global.runtimeEnabled=true | true
global.CertEnabled | To enable cert creation for api-gateway to use SSL Termination | --set global.certEnabled=true | false
global.imagepullsecret | To pull images from private registry setup used with --set global.registry=<> | --set global.imagepullsecret="xx" | ""







