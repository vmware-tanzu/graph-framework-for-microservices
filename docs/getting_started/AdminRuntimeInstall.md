## Nexus Admin Runtime
An overview of the admin runtime and its design can be found [here](../design/Nexus-Runtime.md#nexus-admin-runtime)

This page aims to provide instructions to install and configure the admin runtime.

### Prepare
```
export ADMIN_NAMESPACE=nexus-admin
```

<!-- nexus-specific exports
```
# store the current directory before we `cd` into the app dir
export DOCS_INTERNAL_DIR=$PWD/docs/_internal
```
-->

### Installing
<!-- enable istio-injection with admin and tenant namespaces
```
# install istio to test runtime with istio-injection 
istioctl install --set profile=demo --set hub=gcr.io/nsx-sm/istio -y
kubectl create namespace $ADMIN_NAMESPACE
kubectl label namespace $ADMIN_NAMESPACE istio-injection=enabled --overwrite
kubectl label namespace default istio-injection=enabled --overwrite
```
-->

To install Nexus Runtime with custom resource(cpu/memory) values, please refer [Nexus Sizing](NexusRuntimeSizing.md)

```
nexus runtime install --namespace $ADMIN_NAMESPACE --admin --skip-bootstrap
```

### Next Steps
Install the tenant runtime and deploy your nexus application. And then, configure the admin runtime to be able to route traffic to this tenant. 

#### [Import your datamodel into you application](WorkingWithCommonDatamodel.md)

#### [Build and Deploy Application](BuildDeployApplication.md)

#### [Test your Application -- Based on OrgChart Datamodel](TestingOrgChartDatamodel.md)

#### [Configure Admin Runtime to create routes to your application](AdminRuntimeConfiguration.md)

### [FAQ](FAQ.md)
