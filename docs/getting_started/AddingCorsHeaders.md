## Configuring the Cors Headers
Ensure the Runtime is installed. 

* For Admin Runtime Instructions [here](AdminRuntimeInstall.md)
* For Tenant Runtime Instruction [here](WorkingWithDatamodel.md#Install-Nexus-Runtime)


### Prepare 

```
export NAMESPACE=${NAMESPACE:-default}
```

<!-- nexus-specific exports
```
# store the current directory before we `cd` into the app dir
export DOCS_INTERNAL_DIR=$PWD/docs/_internal
```
-->

### Enabling CORS Headers for required domains

Create a `CORSConfig` object to configure gateway to enable CORS headers

```
# port-forward the admin api-gw to be able to configure it
kubectl port-forward svc/nexus-api-gw 5001:80 -n $NAMESPACE
```

```shell
#create corsconfig object
kubectl -s localhost:5001 apply -f - <<EOF
apiVersion: domain.nexus.vmware.com/v1
kind: CORSConfig
metadata:
  name: default
  labels:
    nexuses.api.nexus.vmware.com: default
    configs.config.nexus.vmware.com: default
    apigateways.apigateway.nexus.vmware.com: default
spec:
  origins: 
    - http://domain # --> list of domain names where the CORS headers should be enabled
  headers:
    - X-Origin # --> list of headers allowed for the Origins configured ( optional argument)
EOF
```
<!-- nexus-specific exports
```
sleep 5
kubectl -s localhost:5001 apply -f $DOCS_INTERNAL_DIR/cors_config.yaml 
sleep 5
```
-->

### Verify CORS configuration takes effect

<!-- nexus-specific exports
```
sleep 5
bash $DOCS_INTERNAL_DIR/cors_headers_verify.sh
sleep 5
```
-->

```shell
curl -X OPTIONS 'http://localhost:5001/api/v1/namespaces/' -H 'Origin: http://domain' -v
```

