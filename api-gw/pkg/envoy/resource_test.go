package envoy_test

import (
	"api-gw/pkg/common"
	"api-gw/pkg/config"
	"api-gw/pkg/envoy"
	"strings"

	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

const JWTFilterNonCSPExpected = `
function envoy_on_request(request_handle)
local jwtMetadata = request_handle:streamInfo():dynamicMetadata():get("envoy.filters.http.jwt_authn")
if jwtMetadata ~= nil then
local jwt = jwtMetadata["jwt_payload"]
if next(jwt) ~= nil then
  local val = jwt["x-user-name"]
  if val ~= nil and type(val) ~= "table" then
	request_handle:headers():remove("x-user-id")
	request_handle:headers():add("x-user-id", val)
  end
end
end
end`

const JWTFilterCSPExpected = `

function envoy_on_request(request_handle)

if request_handle:headers():get("x-admin") == "yes" then
	request_handle:logInfo("Skipping checking org id as it is admin route")
else
	if request_handle:headers():get("static") == "yes" then
		request_handle:logInfo("Skipping checking org id to get static Route")
	else

local jwtMetadata = request_handle:streamInfo():dynamicMetadata():get("envoy.filters.http.jwt_authn")
if jwtMetadata ~= nil then
local jwt = jwtMetadata["jwt_payload"]
if next(jwt) ~= nil then
  
			local context_name = jwt["context_name"]
			path = "/v0/user/preferences"
			local headers, body = request_handle:httpCall(
				"nexus-admin",
				{
				[":authority"] = "nexus-api-gw",
				[":method"] = "GET",
				[":path"] = path,
				["cookie"] = request_handle:headers():get("cookie")
				},
				"",
				1000,
				false)
				if headers[":status"] == "200" then
					request_handle:headers():add("org-id", context_name)
				end
				if headers[":status"] == "503" then
					request_handle:headers():add("org-id", context_name)
					request_handle:headers():get("x-admin", "yes")
					request_handle:headers():remove(":path")
					request_handle:headers():remove(":method")
					request_handle:headers():add(":method", "GET")
					request_handle:headers():add(":path", "/v0/tenants/status")
				end

  local val = jwt["x-user-name"]
  if val ~= nil and type(val) ~= "table" then
	request_handle:headers():remove("x-user-id")
	request_handle:headers():add("x-user-id", val)
  end
end
end
end
end
end
`

const OrgIDRouteExpected = `
function urldecode(s)
s = s:gsub('+', ' ')
:gsub('%%(%x%x)', function(h)
				return string.char(tonumber(h, 16))
				end)
return s
end

function parseurl(s)
local ans = {}
for k,v in s:gmatch('([^&=?]-)=([^&?]+)' ) do
	ans[ k ] = urldecode(v)
end
return ans
end
function string:startswith(prefix)
return self:find(prefix, 1, true) == 1
end
function envoy_on_request(request_handle)
path = request_handle:headers():get(":path")
local query_params = parseurl(path)
local token = query_params["token"]
if token  == nil or token == '' then
	token = query_params["access_token"]
end
if request_handle:headers():get("x-admin") == "yes" then
  request_handle:logInfo("Skipping checking org id as it is admin route")
else
	if request_handle:headers():get("static") == "yes" then
	request_handle:logInfo("Skipping checking org id to get static Route")
	else
		if string.match(path,'clusters/onboarding--manifest')
		then
			local tenant = query_params["tenant"]
			request_handle:headers():add("org-id",tenant)
		end
		
		if token == nil or token == '' then
		request_handle:logInfo("without token")
		path = "/v0/users/validate"
		else
		path = "/v0/users/validate?token="..token
		end
			request_handle:logCritical(path)
			local headers, body = request_handle:httpCall(
				"nexus-admin",
				{
				[":authority"] = "nexus-api-gw",
				[":method"] = "GET",
				[":path"] = path,
				["cookie"] = request_handle:headers():get("cookie")
				},
				"",
				1000,
				false)
				if headers[":status"] == "200" then
					local tenant = body:gsub('"', ""):gsub("%s+", "")
					request_handle:headers():add("org-id",tenant)
				end
			request_handle:headers():add("user-id","test")
		
	end
	end
end
`
const OrgIDCSPRouteExpected = `
function urldecode(s)
s = s:gsub('+', ' ')
:gsub('%%(%x%x)', function(h)
				return string.char(tonumber(h, 16))
				end)
return s
end

function parseurl(s)
local ans = {}
for k,v in s:gmatch('([^&=?]-)=([^&?]+)' ) do
	ans[ k ] = urldecode(v)
end
return ans
end
function string:startswith(prefix)
return self:find(prefix, 1, true) == 1
end
function envoy_on_request(request_handle)
path = request_handle:headers():get(":path")
local query_params = parseurl(path)
local token = query_params["token"]
if token  == nil or token == '' then
	token = query_params["access_token"]
end
if request_handle:headers():get("x-admin") == "yes" then
  request_handle:logInfo("Skipping checking org id as it is admin route")
else
	if request_handle:headers():get("static") == "yes" then
	request_handle:logInfo("Skipping checking org id to get static Route")
	else
		if string.match(path,'clusters/onboarding--manifest')
		then
			local tenant = query_params["tenant"]
			request_handle:headers():add("org-id",tenant)
		end
		
	end
	end
end
`

const StaticRouteConfigExpected = `
function urldecode(s)
s = s:gsub('+', ' ')
:gsub('%%(%x%x)', function(h)
				return string.char(tonumber(h, 16))
				end)
return s
end

function string:endswith(suffix)
return self:sub(-#suffix) == suffix
end
function string:startswith(prefix)
return self:find(prefix, 1, true) == 1
end
function envoy_on_request(request_handle)
path = request_handle:headers():get(":path")
if path:startswith'/home' or path:startswith'/allspark-static' or path:endswith'js' or path:endswith'css' or path:endswith'png'
then
if path:startswith"/apis" == false and path:startswith"/declarative" == false
then
	request_handle:headers():add("static","yes")
end
end
if path:startswith"/v0" then
		request_handle:headers():add("x-admin", "yes")
end
if path == "/" then
       request_handle:headers():add("x-admin", "yes")
end
end
`

var _ = Describe("nexus proxy tests", func() {
	AfterSuite(func() {
		envoy.XDSListener.Close()
	})

	It("initialize nexus-proxy state with nils", func() {
		config.GlobalStaticRouteConfig = &config.GlobalStaticRoutes{
			Suffix: []string{"js", "css", "png"},
			Prefix: []string{"/home", "/allspark-static"},
		}
		logLevel, _ := log.ParseLevel("debug")
		envoy.Init(nil, nil, nil, logLevel)
		snap, err := envoy.GenerateNewSnapshot(nil, nil, nil, nil)
		Expect(snap).NotTo(BeNil())
		Expect(err).To(BeNil())

		staticRouteLua := envoy.ConstructStaticRoute()
		Expect(strings.TrimSpace(staticRouteLua.InlineCode)).To(Equal(strings.TrimSpace(StaticRouteConfigExpected)))

	})

	It("initialize nexus-proxy with http enabled", func() {
		tenantConfigs := []*envoy.TenantConfig{
			{
				Name:   "test",
				Status: false,
			},
		}
		common.Mode = "admin"
		snap, err := envoy.GenerateNewSnapshot(tenantConfigs, nil, nil, nil)
		Expect(err).To(BeNil())
		Expect(snap).NotTo(BeNil())
		l := snap.GetResources(resource.ListenerType)
		Expect(len(l)).To(Equal(1))

	})

	It("initialize nexus-proxy state with tenantconfig routing rules", func() {
		common.SSLEnabled = "true"
		tenantConfigs := []*envoy.TenantConfig{
			{
				Name:   "test",
				Status: false,
			},
		}
		common.Mode = "admin"
		snap, err := envoy.GenerateNewSnapshot(tenantConfigs, nil, nil, nil)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l := snap.GetResources(resource.ListenerType)
		r := snap.GetResources(resource.RouteType)
		c := snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(2))
		listener, ok := l["listener_0"].(*listener.Listener)
		Expect(ok).To(Equal(true))
		Expect(listener.FilterChains[0].TransportSocket.Name).To(Equal("envoy.transport_sockets.tls"))

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok := r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts := routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(6))

		Expect(len(c)).To(Equal(5))
	})

	It("initialize nexus-proxy state with Header-based routing rules", func() {
		headerUpstreams := map[string]*envoy.HeaderMatchedUpstream{"test1": {
			Name:        "test1",
			HeaderName:  "x-tenant",
			HeaderValue: "1",
			Host:        "example.com",
			Port:        80,
		}}
		snap, err := envoy.GenerateNewSnapshot(nil, nil, nil, headerUpstreams)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l := snap.GetResources(resource.ListenerType)
		r := snap.GetResources(resource.RouteType)
		c := snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(2))

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok := r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts := routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(7))

		Expect(len(c)).To(Equal(3))

		// add another upstream and check assertions
		// we expect a route and a cluster to have gotten added
		headerUpstreams["test2"] = &envoy.HeaderMatchedUpstream{
			Name:        "test2",
			HeaderName:  "x-tenant",
			HeaderValue: "2",
			Host:        "google.com",
			Port:        443,
		}
		snap, err = envoy.GenerateNewSnapshot(nil, nil, nil, headerUpstreams)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l = snap.GetResources(resource.ListenerType)
		r = snap.GetResources(resource.RouteType)
		c = snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(2))

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok = r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts = routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(8))

		Expect(len(c)).To(Equal(4))
	})

	It("initialize nexus-proxy state with JWT authn and JWT-claim based routing rules", func() {
		jwt := &envoy.JwtAuthnConfig{
			IdpName:              "csp",
			Issuer:               "https://csp.url",
			JwksUri:              "https://csp.url/jwks",
			CallbackEndpoint:     "http://nexus.proxy.url/callback",
			JwtClaimUsername:     "username",
			RefreshTokenEndpoint: common.RefreshAccessTokenEndpoint,
			AccessToken:          common.AccessTokenStr,
		}
		upstreams := map[string]*envoy.UpstreamConfig{"test1": {
			Name:          "test1",
			JwtClaimKey:   "username",
			JwtClaimValue: "foo@example.com",
			Host:          "example.com",
			Port:          80,
		}}

		snap, err := envoy.GenerateNewSnapshot([]*envoy.TenantConfig{},
			jwt, upstreams, nil)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l := snap.GetResources(resource.ListenerType)
		r := snap.GetResources(resource.RouteType)
		c := snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(2))
		// TODO assert that the jwt_authn filter is configured

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok := r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts := routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(9))

		Expect(len(c)).To(Equal(4))

		// add another upstream and check assertions
		// we expect a route and a cluster to have gotten added
		upstreams["test2"] = &envoy.UpstreamConfig{
			Name:          "test2",
			JwtClaimKey:   "username",
			JwtClaimValue: "bar@example.com",
			Host:          "google.com",
			Port:          443,
		}
		snap, err = envoy.GenerateNewSnapshot(nil, jwt, upstreams, nil)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l = snap.GetResources(resource.ListenerType)
		r = snap.GetResources(resource.RouteType)
		c = snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(2))
		// TODO assert that the jwt_authn filter is configured

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok = r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts = routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(10))

		Expect(len(c)).To(Equal(5))

		err = envoy.AddTenantConfig(&envoy.TenantConfig{
			Name: "test2",
		})
		Expect(err).To(BeNil())

		err = envoy.DeleteTenantConfig("test2")
		Expect(err).To(BeNil())

		orgID := envoy.ConstructOrgIDHeader(false)
		Expect(strings.TrimSpace(orgID.InlineCode)).To(BeEquivalentTo(strings.TrimSpace(OrgIDRouteExpected)))

		orgIDCSP := envoy.ConstructOrgIDHeader(true)
		Expect(strings.TrimSpace(orgIDCSP.InlineCode)).To(BeEquivalentTo(strings.TrimSpace(OrgIDCSPRouteExpected)))

		JWTFilter := envoy.ConstructJWTFilter(true, "x-user-name", "x-user-id")
		Expect(strings.TrimSpace(JWTFilter.InlineCode)).To(BeEquivalentTo(strings.TrimSpace(JWTFilterCSPExpected)))

		JWTFilterNonCSP := envoy.ConstructJWTFilter(false, "x-user-name", "x-user-id")
		Expect(strings.TrimSpace(JWTFilterNonCSP.InlineCode)).To(BeEquivalentTo(strings.TrimSpace(JWTFilterNonCSPExpected)))
	})
})
