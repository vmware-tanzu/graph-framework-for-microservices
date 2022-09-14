package controllers

var oidcCrdObjectExample = `
apiVersion: authentication.nexus.org/v1
kind: OIDC
metadata:
  name: okta
  labels:
    nexuses.api.nexus.org: default
    configs.config.nexus.org: default
    apigateways.apigateway.nexus.org: default
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

var proxyRuleHeaderExample = `
apiVersion: admin.nexus.org/v1
kind: ProxyRule
metadata:
  name: header-based
  labels:
    nexuses.api.nexus.org: default
    configs.config.nexus.org: default
    apigateways.apigateway.nexus.org: default
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
apiVersion: admin.nexus.org/v1
kind: ProxyRule
metadata:
  name: csp
  labels:
    nexuses.api.nexus.org: default
    configs.config.nexus.org: default
    apigateways.apigateway.nexus.org: default
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
apiVersion: apis.nexus.org/v1
kind: Nexus
metadata:
  name: default
`

var configExample = `
apiVersion: config.nexus.org/v1
kind: Config
metadata:
  name: default
`

var routeExample = `
apiVersion: Routes.nexus.org/v1
kind: Route
metadata:
  name: custom
spec:
  service:
    Name: testappserver
    Port: 80
    Scheme: Http
  resource:
    Name: custom
  uri: "/*"
`
