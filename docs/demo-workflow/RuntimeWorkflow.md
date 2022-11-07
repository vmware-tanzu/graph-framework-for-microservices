## Runtime Workflow
## Prerequisites
   Cluster(kind or EKS) should be deployed and available to the nexus runtime deployment

* [Nexus Install](RuntimeWorkflow.md#nexus-install)

* [Nexus Pre-req verify](RuntimeWorkflow.md#nexus-pre-req-verify)

* [Setup Environment](RuntimeWorkflow.md#setup-environment)

* [Install Nexus Admin Runtime](RuntimeWorkflow.md#install-nexus-admin-runtime)
  
* [Install Nexus Tenant Runtime](RuntimeWorkflow.md#install-nexus-tenant-runtime)

## Nexus Install

Install Nexus CLI

```
GOPRIVATE="gitlab.eng.vmware.com" go install gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/cli.git/cmd/plugin/nexus@master
```

   <details><summary>FAQ</summary>
   The above commmand shows unable to connect to gitlab.eng.vmware.com ?

   Verify that you have permissions to the repo

       git ls-remote git@gitlab.eng.vmware.com:nsx-allspark_users/nexus-sdk/cli.git

   Update gitconfig to use ssh instead of https

       git config --global url.git@gitlab.eng.vmware.com:.insteadOf https://gitlab.eng.vmware.com

   </details>

## Nexus Pre-req Verify

Verify nexus sdk pre-requisites are satisfied

    nexus prereq verify

<!-- nexus-specific exports
```
# store the current directory before we `cd` into the app dir
export DOCS_INTERNAL_DIR=$PWD/docs/_internal
```
-->

## Install Nexus Admin Runtime

:bulb: This step is optional, based on the user's requirement, this step shall be skipped or executed

An overview of the admin runtime and its design can be found [here](../design/Nexus-Runtime.md#nexus-admin-runtime)

This page aims to provide instructions to install and configure the admin runtime.

<!-- enable istio-injection with admin and tenant namespaces
```
# install istio to test runtime with istio-injection 
istioctl install --set profile=demo --set hub=gcr.io/nsx-sm/istio -y
kubectl create namespace $ADMIN_NAMESPACE
kubectl label namespace $ADMIN_NAMESPACE istio-injection=enabled --overwrite
kubectl label namespace default istio-injection=enabled --overwrite
```
-->

```
nexus runtime install --namespace nexus-admin --admin --skip-bootstrap
```

:bulb: To override the default value, *nexus runtime install --namespace \<admin-namespace\> --admin --skip-bootstrap*

For EKS cluster in the TSM devgroups(https://confluence.eng.vmware.com/pages/viewpage.action?spaceKey=NSBU&title=AWS+Shared+Accounts)
```shell
nexus runtime install --admin --registry 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus --skip-bootstrap
```
:bulb: To override the default value, *nexus runtime install --namespace \<admin-namespace\> --admin --registry 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus --skip-bootstrap*

## Install Nexus Tenant Runtime

For kind clusters

```
nexus runtime install --namespace default
```
:bulb: To override the default value, *nexus runtime install --namespace \<namespace\>*

For EKS cluster in the TSM devgroups(https://confluence.eng.vmware.com/pages/viewpage.action?spaceKey=NSBU&title=AWS+Shared+Accounts)
```shell
nexus runtime install --namespace default --registry 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus
```
:bulb: To override the default value, *nexus runtime install --namespace \<namespace\> --registry 284299419820.dkr.ecr.us-west-2.amazonaws.com/nexus*

This will install Nexus Runtime microservices to your kubernetes cluster.

*NOTE: Optionally override namespace "default" with a preferred namespace to install nexus runtime*


