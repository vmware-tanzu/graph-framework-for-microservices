package controllers

var tenantRuntimeExample = `
apiVersion: tenantruntime.nexus.vmware.com/v1
kind: Tenant
metadata:
  annotations:
    tenant-registration-status: '{\"State\":\"STATE_REGISTRATION\",\"Status\":\"STATUS_COMPLETED\",\"Message\":\"Registration
      Completed. 5 out of 5 steps completed.\"}'
  creationTimestamp: "2023-05-02T07:30:38Z"
  labels:
    nexus/display_name: 831a06ab-9781-487c-a9da-6c73973e540a
    nexus/is_name_hashed: "true"
    nexuses.api.nexus.vmware.com: default
    runtimes.runtime.nexus.vmware.com: default
  name: 43e32e8f86ba4a5e90e44049f19ab8f73b77cc7f
spec:
  attributes:
    skus:
    - LICENSE_ADVANCE
  awsKmsKeyId: alias/rootca-01-4pydxsym-zay6xfru-uk3fvv6z
  awsS3Bucket: rootca-01-4pydxsym-zay6xfru-uk3fvv6z
  licenseType: classic
  m7Enabled: "false"
  m7InstallationScheduled: InstallationDisabled
  namespace: 831a06ab-9781-487c-a9da-6c73973e540a
  releaseVersion: cosmos-v2
  saasApiDomainName: jrentcl.tsm-dev-07.com
  saasDomainName: jrentcl-internal.tsm-dev-07.com
  streamName: 01-4pydxsym-zay6xfru-uk3fvv6z
  tenantName: 831a06ab-9781-487c-a9da-6c73973e540a
status:
  appStatus:
    installedApplications:
      nexusApps:
        nexus-tenant-runtime:
          oamApp:
            components:
              831a06ab-9781-487c-a9da-6c73973e540a.nexus-tenant-runtime:
                health: Healthy
                name: 831a06ab-9781-487c-a9da-6c73973e540a.nexus-tenant-runtime
                sync: Synced
          state: Running
          stateReason: All components installed
        tsm-tenant:
          oamApp:
            components:
              831a06ab-9781-487c-a9da-6c73973e540a.tsm-tenant-runtime:
                health: Healthy
                name: 831a06ab-9781-487c-a9da-6c73973e540a.tsm-tenant-runtime
                sync: Synced
          state: Running
          stateReason: All components installed
    releaseVersionStatus: Installed version cosmos-v2 successfully
  nexus:
    remoteGeneration: 0
    sourceGeneration: 0
`
var tenantConfigExample = `
apiVersion: tenantconfig.nexus.vmware.com/v1
kind: Tenant
metadata:
  name: tenant1
  labels:
    nexuses.api.nexus.vmware.com: default
    configs.config.nexus.vmware.com: default
    apigateways.apigateway.nexus.vmware.com: default
spec:
  name: tenant1
  skus:
    - advance`

var oidcCrdObjectExample = `
apiVersion: authentication.nexus.vmware.com/v1
kind: OIDC
metadata:
  name: okta
  labels:
    nexuses.api.nexus.vmware.com: default
    configs.config.nexus.vmware.com: default
    apigateways.apigateway.nexus.vmware.com: default
spec:
  config:
    clientId: "XXX"
    clientSecret: "XXX"
    oAuthIssuerUrl: "https://dev-XXX.okta.com/oauth2/default"
    oAuthRedirectUrl: "http://<API-GW-DNS/IP>:<PORT>/<CALLBACK_PATH>"
    scopes:
      - openid
      - profile
      - offline_access
  validationProps:
    insecureIssuerURLContext: false
    skipIssuerValidation: false
    skipClientIdValidation: false
    skipClientAudValidation: false
`

var corsConfigExample = `
kind: CORSConfig
apiVersion: domain.nexus.vmware.com/v1
metadata:
  name: default
spec:
  origins:
    - http://localhost
    - http://test
    - http://testing
  headers:
    - test
`

var proxyRuleHeaderExample = `
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
`

var proxyRuleJwtExample = `
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
`

var nexusExample = `
apiVersion: api.nexus.vmware.com/v1
kind: Nexus
metadata:
  name: default
`

var configExample = `
apiVersion: config.nexus.vmware.com/v1
kind: Config
metadata:
  name: default
`

var routeExample = `
apiVersion: route.nexus.vmware.com/v1
kind: Route
metadata:
  name: custom
spec:
  service:
    name: testappserver
    port: 80
    scheme: Http
  resource:
    name: custom
  uri: "/*"
`
