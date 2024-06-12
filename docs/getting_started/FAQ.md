# FAQ

**[Build](FAQ.md#build)**

**[API Gateway](FAQ.md#api-gateway)**

# Build

## How to modifying Datamodel spec, build and install it in the Nexus runtime ?

###  App Local Datamodel

App local datamodels are versioned along with the App. As such can be changed natively as part of App.

Step 1: Build Datamodel

```
nexus datamodel build --name $DATAMODEL_NAME 
```

Step 2: Install updated Datamodel

```
nexus datamodel install name $DATAMODEL_NAME --namespace $NAMESPACE --debug
```

CAUTION :skull:

Validation layer in the Nexus Runtime prevents any updates to annotations on the CRD's / API's.
Updating spec/annotations is a dangerous action in production environments and WILL impact all your tenants. 

Nexus validation layer will throw error like the following:

> Name: "roots.orgchart.vmware.org", Namespace: ""
> for: "vmware/build/crds/root_root.yaml": admission webhook "nexus-validation-crd-type.webhook.svc" denied the request: You are not allowed to change nexus annotation for this CRD

To workaround the validation layer, you can delete the CRD/API on which the error is reported and retry datamodel install command.
```
export CRD_NAME="orgcharts.orgchart.vmware.org"
export NAMESPACE="default"
kubectl exec -it $(kubectl get pods  --no-headers -l app=nexus-proxy-container -n "$NAMESPACE" | awk '{print $1}') -n "$NAMESPACE" -- kubectl -s nexus-apiserver:8080 delete crds "$CRD_NAME"
nexus datamodel install name $DATAMODEL_NAME --namespace $NAMESPACE --debug
```

## API Gateway

### How can I access the OpenApi Spec ?

Step 1: Port-forward to API-GW Pod.

This is required as the API-GW does not have a Public domain/IP.

```
NAMESPACE="default" kubectl port-forward $(kubectl get pods -n "${NAMESPACE}" -l control-plane=api-gw --no-headers | awk '{print $1}') -n "${NAMESPACE}" 8080:80
```

Step 2: Access the following URL

```
http://localhost:8080/openapi.json
```
