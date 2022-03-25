<h1>Tanzu TSM CLI</h1>


<h2>1. Install </h2>


Before building the TSM plugins, make sure you download and install tanzu CLI framework.

<h3>MacOS</h3>

```
curl -o tanzu https://storage.googleapis.com/tanzu-cli/artifacts/core/latest/tanzu-core-darwin_amd64 && \
    mv tanzu /usr/local/bin/tanzu && \
    chmod +x /usr/local/bin/tanzu
```

<h4>Linux</h4>
<h5>i386</h5>

```
curl -o tanzu https://storage.googleapis.com/tanzu-cli/artifacts/core/latest/tanzu-core-linux_386 && \
    mv tanzu /usr/local/bin/tanzu && \
    chmod +x /usr/local/bin/tanzu
```

<h6>AMD64</h6>

```
curl -o tanzu https://storage.googleapis.com/tanzu-cli/artifacts/core/latest/tanzu-core-linux_amd64 && \
    mv tanzu /usr/local/bin/tanzu && \
    chmod +x /usr/local/bin/tanzu
```

<h7>Windows</h7>

```
https://storage.googleapis.com/tanzu-cli/artifacts/core/latest/tanzu-core-windows_amd64.exe
```

<h1>2. Builder</h1>
Builder is an admin plugin that needs to be installed.

```
tanzu plugin repo add -n admin -b tanzu-cli-admin-plugins -p artifacts-admin
tanzu plugin install builder
```

<h1>3. Checkout Repo</h1>

```
git clone git@gitlab.eng.vmware.com:nexus/cli.git
cd cli
git checkout nexus-cli-dev
```

<h1>4. Make</h1>
Once you have the tanzu CLI installed, run the following to have the plugins installed.

```
# Commands:
# make version - show information about current version
# make build - builds the plugins and stores the targets in artifacts
# make install - installs the plugins in tanzu framework
```

This will build a local repository under ./artifacts.

```
make build
```

Plugins can be installed from this repository using:

```
make install
```

<h1>5. Usage nexus CLI </h1>

```
# nexus --help
nexus cli to execute tsm operations

Usage:
  nexus [command]

Available Commands:
  app         Sample application installer
  apply       Apply Servicemesh configuration from file
  cluster     Servicemesh cluster features
  completion  Generate the autocompletion script for the specified shell
  config      Servicemesh configuration features
  datamodel   Datamodel installer and uninstaller
  delete      Delete Servicemesh configuration from file
  gns         Servicemesh global namespace features
  help        Help about any command
  login       Login to csp
  runtime     Runtime installer and uninstaller

Flags:
  -h, --help   help for nexus

Use "nexus [command] --help" for more information about a command.

```

<h1>6. Usage tanzu CLI </h1>

You should see servicemesh plugin installed under Tanzu.

```
# tanzu servicemesh --help
service mesh cluster features

Usage:
  tanzu servicemesh [command]

Available Commands:
  cluster     Servicemesh cluster features

Flags:
  -h, --help   help for servicemesh

Use "servicemesh [command] --help" for more information about a command.
```

```
tanzu servicemesh cluster --help
Servicemesh cluster features

Usage:
  tanzu servicemesh cluster [command]

Available Commands:
  create      Create a cluster object
  delete      Delete a cluster
  get         Get a cluster
  list        List clusters

Flags:
  -h, --help   help for cluster

Use "servicemesh cluster [command] --help" for more information about a command.
```

```
 tanzu servicemesh cluster get --help
Get a cluster

Usage:
  tanzu servicemesh cluster get CLUSTER_NAME [flags]

Flags:
  -h, --help            help for get
  -o, --output string   Output formart. Supported formats: json|yaml
```

```
tanzu servicemesh cluster list --help
List clusters

Usage:
  tanzu servicemesh cluster list [flags]

Flags:
  -h, --help            help for list
  -o, --output string   Output formart. Supported formats: json|yaml
 ```


<h1>7. Saas Authentication</h1>

```
tanzu servicemesh login -t <API Token> -s <Saas-DNS-Name>
 
Example: tanzu servicemesh login -t HJGsagdfjhiB683MBXIVXZr9VCEO5zr5jVAlcHAXfXumTo812H9dt9i4HGNk9z1235gjsd -s staging-0.servicemesh.biz
```

If you don't have APIToken, please generate from My Account → API Tokens here: https://console-stg.cloud.vmware.com/csp/gateway/portal/

Note: 
    The  .servicemesh.config and .servicemesh.server files will be stored in the user's $HOME directory

<h1>7. GNS</h1>
<h2>Create GNS</h2>
Example GNS YAML file : https://gitlab.eng.vmware.com/nsx-allspark_users/tsm-cli/-/blob/TSM-3661-decalartive-api/gns-test.yaml

```
metadata:
    project: default
    global-namespace: mygns    <---------- Change as you create new gns
spec:
  name: mygns <---------- Change as you create new gns
  domain: mygns.com.   <---------- Change as you create new gns
  use_shared_gateway: true
  mtls_enforced: false,
  matchingConditions:
  - namespace:
      match: mygns <---------- Change as you create new gns
      type: EXACT
    cluster:
      match: mygns <---------- Change as you create new gns
      type: EXACT
  api_discovery_enabled: true
```

```
tanzu servicemesh apply -f mygns.yaml
```

<h3>Delete GNS</h3>

```
tanzu servicemesh delete -f mygns.yaml
```
<h1>8. GNS Public Service</h1>
<h2>Pre-requisites</h2>
Configure DNS and  Domains using Tanzu Admin Integrations
Retrieve the external_dns_id from API Explorer

```
1.GET v1alpha1/external-domain-name-server-list
2.Make sure the external-id is belonged to you from this GET v1alpha1/external-domain-name-servers/{external-id-from step 1} 
```

<h3>Create GNS Public Service</h3>
<h4>Create Health-Check-Id</h4>

Example Yaml :  https://gitlab.eng.vmware.com/nsx-allspark_users/tsm-cli/-/blob/TSM-3661-decalartive-api/health-checks-test.yaml

```
metadata:
  project: default
spec:
  name: tsmcli.servicemesh.biz
  protocol: HTTP
  domain: tsmcli.servicemesh.biz
  port: 3000
  path: "/"
  healthThreshold: 3
  certificate_id: ''
  external_port: 80
  interval: 10
```

```
tanzu servicemesh apply -f health-checks-test.yaml
```

<h5>Create GNS Public Service</h6>

Example Yaml File : https://gitlab.eng.vmware.com/nsx-allspark_users/tsm-cli/-/blob/TSM-3661-decalartive-api/public-service-test.yaml
```
metadata:
  project: default
  global-namespace: mygns
  public-service: tsm-cli.servicemesh.biz
spec:
  fqdn: tsm-cli.servicemesh.biz
  name: ''
  external_port: 80
  external_protocol: HTTP
  ttl: 300
  public_domain:
    external_dns_id: 1a6747ba-0774-4c8f-a374-108b7446a4c9.  <-------- Update from the step 8.1 pre-requistes
    primary_domain: servicemesh.biz
    sub_domain: tsm-cli
    certificate_id: ''
  ha_policy: ''
  gslb:
    type: ROUND_ROBIN
  wildcard_certificate_id: ''
  healthcheck_ids:
  - efec3c50-fba6-4c7c-aae5-55d1685bc2f0.  <----- Update from the 8.2.1 Response
  ```
```
tanzu servicemesh apply -f public-service-test.yaml
```
<h6>Create Route to GNS Public Service</h6>
Example Yaml file: https://gitlab.eng.vmware.com/nsx-allspark_users/tsm-cli/-/blob/TSM-3661-decalartive-api/route-test.yaml

```
metadata:
  project: default
  global-namespace: mygns
  public-service: tsm-cli.servicemesh.biz
  route: shopping.3000
spec:
  paths:
  - "/"
  target: shopping
  target_port: 3000
```
```
tanzu servicemesh apply -f route-shopping.yaml
```
<h7>Delete GNS Public Service</h7>

```
tanzu servicemesh delete -f public-service-test.yaml
```

Note: This deletes the public-service and route only.  The corresponding health-checks will not be deleted, we have yet to implement the same. There is no side-effect only the unused health-checks will exist.

