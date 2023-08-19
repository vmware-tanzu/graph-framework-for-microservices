## Access backend/app REST API's on Nexus API GW

This workflow will walk through steps to expose API's from backend services/apps directly on the Nexus API GW.

### Prerequisites

1. Please install nexus runtime in a namespace using nexus CLI

   #### **[Nexus Runtime installation](Playground-SockShop-Install-Datamodel.md#install-nexus-runtime)**

### Step 1: Start a backend microservice that serves REST endpoint

#### NOTE: this step is optional and only needed for quick test / poc workflow

1. For this demo we will start a Nginx pod that serves a "/" URI
    ```
    kubectl run my-nginx --image=nginx --port=80
    ```

2. Expose the Nginx pods as a K8s service
    ```
    kubectl expose pod my-nginx --port=80
    ```

### Step 2: Ensure network connectivity to Nexus API GW

#### NOTE: this step is optional and only needed if Nexus API GW is not exposed on a public domain; as is the case for most POC's

1. Port-forward to Nexus API GW service.
    ```
    kubectl port-forward svc/nexus-api-gw 8080:80 &
    ```

### Step 3: Create a Nexus Route object

Nexus Route object is a declarative API through which custom routes can be configured on the Nexus API GW.

1. Create the about Route object
```
cat <<EOF | kubectl -s localhost:8080 apply -f -
apiVersion: route.nexus.vmware.com/v1
kind: Route
metadata:
  name: test                # name of the route object. Use any name that makes sense to you.
  labels:
    nexuses.api.nexus.vmware.com: default
    configs.config.nexus.vmware.com: default
spec:
  uri: /                    # The URI prefix that will be served by the backend microservice
  service:
    name: my-nginx          # Name of the K8s service for the backend microservice
    port: 80                # Port number of K8s service for the backend microservice
    scheme: http            # Protocol of the K8s service for the backend microservice
  resource:
    name: mygroup           # API identifier/group to uniquely identify this route. Use any name that makes sense.
EOF
```

### Done. Access the exposed REST API of the backend microservice using prefix: /apis/mygroup/v1/

```
kubectl -s localhost:8080 get --raw /apis/mygroup/v1
```

***Example:***
```
➜ kubectl -s localhost:8080 get --raw /apis/mygroup/v1
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
html { color-scheme: light dark; }
body { width: 35em; margin: 0 auto;
font-family: Tahoma, Verdana, Arial, sans-serif; }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```
