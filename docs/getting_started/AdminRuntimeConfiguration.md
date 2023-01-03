## Configuring the Nexus Admin Runtime
Ensure that the Admin Runtime is installed. Instructions [here](AdminRuntimeInstall.md)

### Prepare
```
export ADMIN_NAMESPACE=${ADMIN_NAMESPACE:-nexus-admin}
```

##### Ensure
1. that the tenant runtime has successfully been installed and
2. that the application has successfully been deployed on the tenant runtime

<!-- nexus-specific exports
```
# store the current directory before we `cd` into the app dir
export DOCS_INTERNAL_DIR=$PWD/docs/_internal
```
-->

##### Routing based on a custom HTTP header
Create a `ProxyRule` object to configure the admin runtime to route based on headers

```
# port-forward the admin api-gw to be able to configure it
kubectl port-forward svc/nexus-api-gw 5000:80 -n $ADMIN_NAMESPACE
```

```shell
kubectl -s localhost:5000 apply -f - <<EOF
apiVersion: admin.nexus.vmware.com/v1
kind: ProxyRule
metadata:
  name: header-based
  labels:
    nexuses.api.nexus.vmware.com: default
    configs.config.nexus.vmware.com: default
    apigateways.apigateway.nexus.vmware.com: default
spec:
  matchCondition:
    type: header
    key: x-tenant # the name of the HTTP header to use for routing
    value: "t-1"  # the value of the HTTP header
  upstream:
    scheme: http
    host: nexus-api-gw.default    # use the nexus-api-gw of tenant-1 as the destination 
    port: 80
EOF
```

### Access the tenant REST APIs through the admin runtime
Try accessing the /leaders endpoint on the app we created through the nexus-proxy on the admin runtime
```
# port-forward the nexus-proxy on the admin-runtime namespace
kubectl port-forward -n $ADMIN_NAMESPACE svc/nexus-proxy 10000:10000
```

```shell
curl localhost:10000/leaders -H "x-tenant:t-1"
```
<!--
```
sleep 5
kubectl -s localhost:5000 apply -f $DOCS_INTERNAL_DIR/proxyrule_header_based_routing.yaml
sleep 5
```
-->

<!--
```
bash $DOCS_INTERNAL_DIR/header_based_routing_test.sh 200
```
-->

### Configuring OpenID connect authentication
This step registers the IDP with the nexus admin-runtime.
```shell
kubectl -s localhost:5000 apply -f - <<EOF
apiVersion: authentication.nexus.vmware.com/v1
kind: OIDC
metadata:
  labels:
    apigateways.apigateway.nexus.vmware.com: default
    configs.config.nexus.vmware.com: default
    nexuses.api.nexus.vmware.com: default
  name: csp
spec:
  config:
    clientId: XXX
    clientSecret: XXX
    oAuthIssuerUrl: https://console-stg.cloud.vmware.com/csp/gateway/am/api
    oAuthRedirectUrl: http://localhost:10000/callback
    scopes:
      - openid
      - profile
      - offline_access
  validationProps:
    insecureIssuerURLContext: true
    skipClientAudValidation: true
    skipClientIdValidation: true
    skipIssuerValidation: true
EOF
```

##### Routing based on a JWT claim
```shell
kubectl -s localhost:5000 apply -f - <<EOF
apiVersion: admin.nexus.vmware.com/v1
kind: ProxyRule
metadata:
  name: csp
  labels:
    nexuses.api.nexus.vmware.com: default
    configs.config.nexus.vmware.com: default
    apigateways.apigateway.nexus.vmware.com: default
spec:
  matchCondition:
    type: jwt
    key: foo   # the name of the JWT claim to use for routing
    value: bar # the value of the JWT claim 

  # if the request matches the matchCondition, envoy will proxy the request to this upstream 
  upstream:
    scheme: http
    host: nexus-api-gw.default
    port: 80
EOF
```

### ProxyRule Gotchas
At least one proxy rule must be created for every tenant that's present in the system. This is required because the upstreams registered with the envoy proxy are based on the proxy rules present.

So if a new tenant is added, be sure to create a proxy rule before trying to access any of the tenant's APIs through the admin runtime.

Note: if there's a mix of proxy rule match types (jwt and header), then envoy will be configured to first try to match the JWT match conditions. 

### Access the tenant REST APIs through the admin runtime
Login via the `/login` endpoint
[Login](http://localhost:10000/login)

Access the `/leaders` endpoint
[/leaders](http://localhost:10000/leaders)

### Delete the ProxyRule
```
kubectl port-forward svc/nexus-proxy-container 45193:80 -n ${ADMIN_NAMESPACE}
kubectl -s localhost:45193 delete proxyrules --all
```

<!--
```
bash $DOCS_INTERNAL_DIR/header_based_routing_test.sh 400
```
-->

### Troubleshooting
```shell
# dump the nexus-proxy config
kubectl port-forward svc/nexus-proxy 19000:19000 -n "$ADMIN_NAMESPACE" &
curl localhost:19000/config_dump

# dump the nexus-proxy logs
kubectl logs $(kubectl get pods  --no-headers -l app=nexus-proxy -n "$ADMIN_NAMESPACE" | awk '{print $1}') -n "$ADMIN_NAMESPACE"

# dump the nexus-api-gw logs
kubectl logs $(kubectl get pods  --no-headers -l control-plane=api-gw -n "$ADMIN_NAMESPACE" | awk '{print $1}') -n "$ADMIN_NAMESPACE"
```
