package controllers

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
