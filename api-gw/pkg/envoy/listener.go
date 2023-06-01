package envoy

import (
	"api-gw/pkg/common"
	"api-gw/pkg/config"
	"fmt"
	"time"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	jwtauthnv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	luav3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/lua/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/any"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

var URLlDecoderFunc string = `function urldecode(s)
s = s:gsub('+', ' ')
:gsub('%%(%x%x)', function(h)
				return string.char(tonumber(h, 16))
				end)
return s
end
`

var CheckStaticHeader string = `
if request_handle:headers():get("x-admin") == "yes" then
	request_handle:logInfo("Skipping checking org id as it is admin route")
else
	if request_handle:headers():get("static") == "yes" then
		request_handle:logInfo("Skipping checking org id to get static Route")
	else
`

var JWTHeaderFilter string = `
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
`

func ConstructJWTFilter(CSP bool, username string, userIdHeader string) *luav3.Lua {
	var addJwtClaimsToHeaderLuaFilter *luav3.Lua
	if CSP {
		addJwtClaimsToHeaderLuaFilter = &luav3.Lua{
			InlineCode: fmt.Sprintf(`
function envoy_on_request(request_handle)
%s
local jwtMetadata = request_handle:streamInfo():dynamicMetadata():get("envoy.filters.http.jwt_authn")
if jwtMetadata ~= nil then
local jwt = jwtMetadata["jwt_payload"]
if next(jwt) ~= nil then
  %s
  local val = jwt["%s"]
  if val ~= nil and type(val) ~= "table" then
	request_handle:headers():remove("%s")
	request_handle:headers():add("%s", val)
  end
end
end
end
end
end`, CheckStaticHeader, JWTHeaderFilter, username, userIdHeader, userIdHeader),
		}
	} else {
		addJwtClaimsToHeaderLuaFilter = &luav3.Lua{
			InlineCode: fmt.Sprintf(`
function envoy_on_request(request_handle)
local jwtMetadata = request_handle:streamInfo():dynamicMetadata():get("envoy.filters.http.jwt_authn")
if jwtMetadata ~= nil then
local jwt = jwtMetadata["jwt_payload"]
if next(jwt) ~= nil then
  local val = jwt["%s"]
  if val ~= nil and type(val) ~= "table" then
	request_handle:headers():remove("%s")
	request_handle:headers():add("%s", val)
  end
end
end
end
`, username, userIdHeader, userIdHeader),
		}
	}
	return addJwtClaimsToHeaderLuaFilter
}

func ConstructOrgIDHeader(CSP bool) *luav3.Lua {
	// this is to validate the user is authenticated
	// When a Request is recieved in NON-CSP , we verify the user cred present in cookie / queryParam and fetch the corresponding Tenant ID
	// Tenant ID is added as header in the request , further routes will be decided as combination of prefix(/tsm or /apis) + org-id( Tenant ID)
	// However there is a exempted route clusters/onboarding which would directly routed using tenant queryparam instead of validation
	attachment := ""
	if !CSP {
		attachment = `
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
		`
	}
	addOrgIDtoHeader := &luav3.Lua{
		InlineCode: fmt.Sprintf(`
		%s
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
		%s
	end
	end
end
				`, URLlDecoderFunc, attachment)}
	return addOrgIDtoHeader

}

func ConstructStaticRoute() *luav3.Lua {
	var staticCondition string = ""
	for _, prefixRoute := range config.GlobalStaticRouteConfig.Prefix {
		if staticCondition == "" {
			staticCondition = fmt.Sprintf("path:startswith'%s'", prefixRoute)
		} else {
			staticCondition = fmt.Sprintf("%s or path:startswith'%s'", staticCondition, prefixRoute)
		}
	}

	for _, suffixRoute := range config.GlobalStaticRouteConfig.Suffix {
		if staticCondition == "" {
			staticCondition = fmt.Sprintf("path:endswith'%s'", suffixRoute)
		} else {
			staticCondition = fmt.Sprintf("%s or path:endswith'%s'", staticCondition, suffixRoute)
		}
	}

	//Validate if the route is static , which would need route to globalUI
	addIfStaticRoute := &luav3.Lua{
		InlineCode: fmt.Sprintf(`
%s
function string:endswith(suffix)
return self:sub(-#suffix) == suffix
end
function string:startswith(prefix)
return self:find(prefix, 1, true) == 1
end
function envoy_on_request(request_handle)
path = request_handle:headers():get(":path")
if %s
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
end`, URLlDecoderFunc, staticCondition),
	}
	return addIfStaticRoute
}
func makeHttpListener() (*listener.Listener, error) {
	router, err := getRouterTypedConfig()
	if err != nil {
		return nil, err
	}
	listenerName := "listener_1"
	log.Debugf("Creating Http Listener for redirect")
	// Adding a redirect from http to Https
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "redirectHttp",

		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: &route.RouteConfiguration{
				Name: "https_redirect",
				VirtualHosts: []*route.VirtualHost{
					{
						Name: "httptohttps",
						Domains: []string{
							"*",
						},
						Routes: []*route.Route{
							{
								Match: &route.RouteMatch{
									PathSpecifier: &route.RouteMatch_Prefix{
										Prefix: "/",
									},
								},
								Action: &route.Route_Redirect{
									Redirect: &route.RedirectAction{
										SchemeRewriteSpecifier: &route.RedirectAction_HttpsRedirect{
											HttpsRedirect: true,
										},
										PortRedirect: 443,
									},
								},
							},
						},
					},
				},
			},
		},
		HttpFilters: []*hcm.HttpFilter{
			{
				Name: wellknown.Router,
				ConfigType: &hcm.HttpFilter_TypedConfig{
					TypedConfig: router,
				},
			},
		},
	}
	var pbst *anypb.Any
	pbst, err = anypb.New(manager)
	if err != nil {
		log.Errorf("failed to create http connection manager: %s", err)
		return nil, err
	}
	Filterchains := []*listener.FilterChain{{
		Filters: []*listener.Filter{{
			Name: wellknown.HTTPConnectionManager,
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: pbst,
			},
		}},
	}}

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: HttpListenerPort,
					},
				},
			},
		},
		FilterChains: Filterchains,
	}, nil
}

func makeRouteListener(jwtAuthnConfig *JwtAuthnConfig) (*listener.Listener, error) {
	listenerName := "listener_0"
	log.Debugf("jwtAuthnConfig: %+v", jwtAuthnConfig)
	httpFilters, err := getHttpFilters(jwtAuthnConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating http filters: %s", err)
	}

	// HTTP filter configuration
	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    makeConfigSource(),
				RouteConfigName: routeDefault,
			},
		},
		HttpFilters: httpFilters,
	}
	var pbst *anypb.Any
	pbst, err = anypb.New(manager)
	if err != nil {
		log.Errorf("failed to create http connection manager: %s", err)
		return nil, err
	}
	var Filterchains []*listener.FilterChain
	var listenerPort int
	if common.IsHttpsEnabled() {
		log.Info("Https Enabled")
		listenerPort = HttpsListenerPort
		certObject := &tlsv3.DownstreamTlsContext{
			CommonTlsContext: &tlsv3.CommonTlsContext{
				TlsCertificates: []*tlsv3.TlsCertificate{
					{
						CertificateChain: &core.DataSource{
							Specifier: &core.DataSource_Filename{
								Filename: "/ssl/cert/cert.pem",
							},
						},
						PrivateKey: &core.DataSource{
							Specifier: &core.DataSource_Filename{
								Filename: "/ssl/cert/key.pem",
							},
						},
					},
				},
			},
		}
		certO, err := anypb.New(certObject)
		if err != nil {
			log.Errorf("failed to create tls certificate manager: %s", err)
			return nil, err
		}

		Transportsocket := &core.TransportSocket{
			Name: "envoy.transport_sockets.tls",
			ConfigType: &core.TransportSocket_TypedConfig{
				TypedConfig: certO,
			}}
		Filterchains = []*listener.FilterChain{{
			TransportSocket: Transportsocket,
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}}
	} else {
		listenerPort = HttpListenerPort
		Filterchains = []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}}
	}
	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(listenerPort),
					},
				},
			},
		},
		FilterChains: Filterchains,
	}, nil
}

func makeConfigSource() *core.ConfigSource {
	source := &core.ConfigSource{}
	source.ResourceApiVersion = resource.DefaultAPIVersion
	source.ConfigSourceSpecifier = &core.ConfigSource_ApiConfigSource{
		ApiConfigSource: &core.ApiConfigSource{
			TransportApiVersion:       resource.DefaultAPIVersion,
			ApiType:                   core.ApiConfigSource_GRPC,
			SetNodeOnFirstMessageOnly: true,
			GrpcServices: []*core.GrpcService{{
				TargetSpecifier: &core.GrpcService_EnvoyGrpc_{
					EnvoyGrpc: &core.GrpcService_EnvoyGrpc{ClusterName: "xds_cluster"},
				},
			}},
		},
	}
	return source
}

func getHttpFilters(jwtAuthnConfig *JwtAuthnConfig) ([]*hcm.HttpFilter, error) {
	router, err := getRouterTypedConfig()
	if err != nil {
		return nil, err
	}
	// this is hack to proceed with lua script
	dummyLuaFilter := &luav3.Lua{
		InlineCode: `
	function envoy_on_request(request_handle)
	  request_handle:headers():add("dummy-header", "nil")
	end
	`,
	}
	dummyLuaFilterTypedConfig, err := anypb.New(dummyLuaFilter)
	if err != nil {
		return nil, fmt.Errorf("error creating dummy lua filter: %s", err)
	}
	if jwtAuthnConfig == nil {
		if TenantConfigs != nil {

			// this is to validate the user is authenticated
			// When a Request is recieved in NON-CSP , we verify the user cred present in cookie / queryParam and fetch the corresponding Tenant ID
			// Tenant ID is added as header in the request , further routes will be decided as combination of prefix(/tsm or /apis) + org-id( Tenant ID)
			// However there is a exempted route clusters/onboarding which would directly routed using tenant queryparam instead of validation
			addOrgIDtoHeader := ConstructOrgIDHeader(false)
			luaVerifyRoute, err := anypb.New(addOrgIDtoHeader)
			if err != nil {
				fmt.Printf("cannot add orgId header parser envoy")
			}
			addIfStatic := ConstructStaticRoute()
			luaVerifyStatic, err := anypb.New(addIfStatic)
			if err != nil {
				return nil, err
			}
			return []*hcm.HttpFilter{
				{
					Name: wellknown.Lua,
					ConfigType: &hcm.HttpFilter_TypedConfig{
						TypedConfig: dummyLuaFilterTypedConfig,
					},
				},
				{
					Name: wellknown.Lua,
					ConfigType: &hcm.HttpFilter_TypedConfig{
						TypedConfig: luaVerifyStatic,
					},
				},
				{
					Name: wellknown.Lua,
					ConfigType: &hcm.HttpFilter_TypedConfig{
						TypedConfig: luaVerifyRoute,
					},
				},
				{
					Name: wellknown.Router,
					ConfigType: &hcm.HttpFilter_TypedConfig{
						TypedConfig: router,
					},
				},
			}, nil
		} else {
			return []*hcm.HttpFilter{
				{
					Name: wellknown.Lua,
					ConfigType: &hcm.HttpFilter_TypedConfig{
						TypedConfig: dummyLuaFilterTypedConfig,
					},
				},
				{
					Name: wellknown.Router,
					ConfigType: &hcm.HttpFilter_TypedConfig{
						TypedConfig: router,
					},
				},
			}, nil
		}
	}

	if jwtAuthnConfig.Issuer == "" || jwtAuthnConfig.JwksUri == "" || jwtAuthnConfig.IdpName == "" || jwtAuthnConfig.CallbackEndpoint == "" {
		return nil, fmt.Errorf("failed to create JWT authn filter: invalid config")
	}

	// JWT Provider validation happens using Token from FromCookies or FromHeaders
	// RequirementRule without Provider will skip JWT verification
	jwtAuthn := &jwtauthnv3.JwtAuthentication{
		Providers: map[string]*jwtauthnv3.JwtProvider{
			jwtAuthnConfig.IdpName: {
				Issuer:      jwtAuthnConfig.Issuer,
				FromCookies: []string{jwtAuthnConfig.AccessToken},
				FromHeaders: []*jwtauthnv3.JwtHeader{{
					Name:        common.AuthorizationHeader,
					ValuePrefix: fmt.Sprintf("%s ", common.AuthorizationTypeBearer),
				}},
				PayloadInMetadata: jwtPayload,
				JwksSourceSpecifier: &jwtauthnv3.JwtProvider_RemoteJwks{
					RemoteJwks: &jwtauthnv3.RemoteJwks{
						HttpUri: &core.HttpUri{
							Uri: jwtAuthnConfig.JwksUri,
							HttpUpstreamType: &core.HttpUri_Cluster{
								Cluster: fmt.Sprintf("%s_jwks_cluster", jwtAuthnConfig.IdpName),
							},
							Timeout: durationpb.New(5 * time.Second),
						},
						CacheDuration: durationpb.New(5 * time.Minute),
					},
				},
			},
		},
		Rules: []*jwtauthnv3.RequirementRule{
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: common.LoginEndpoint,
					},
				},
			},
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Path{
						Path: jwtAuthnConfig.RefreshTokenEndpoint,
					},
				},
			},
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Path{
						Path: common.LogoutEndpoint,
					},
				},
			},
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_SafeRegex{
						SafeRegex: &matcherv3.RegexMatcher{
							EngineType: &matcherv3.RegexMatcher_GoogleRe2{},
							Regex:      ".*/clusters/onboarding-manifest.*",
						},
					},
				},
			},
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/cluster-registration",
					},
				},
			},
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/release-manifests",
					},
				},
			},
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Path{
						Path: jwtAuthnConfig.CallbackEndpoint,
					},
				},
			},
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: "/",
					},
				},
				RequirementType: &jwtauthnv3.RequirementRule_Requires{
					Requires: &jwtauthnv3.JwtRequirement{
						RequiresType: &jwtauthnv3.JwtRequirement_ProviderName{
							ProviderName: jwtAuthnConfig.IdpName,
						},
					},
				},
			},
		},
	}

	jwtAuthnTypedConfig, err := anypb.New(jwtAuthn)
	if err != nil {
		return nil, fmt.Errorf("error getting jwtAuthnTypedConfig: %s", err)
	}

	authFailureHandlerFilter := &luav3.Lua{
		InlineCode: fmt.Sprintf(`
function envoy_on_response(response_handle)
    if response_handle:headers():get(":status") == "401" then
        response_handle:logInfo("Got status 401, redirect to login....")
        response_handle:headers():replace(":status", "307")
        response_handle:headers():add("location", "%s?state="..response_handle:headers():get("www-authenticate"))
    end
end
`, common.LoginEndpoint),
	}
	authFailureHandlerFilterTypedConfig, err := anypb.New(authFailureHandlerFilter)
	if err != nil {
		return nil, fmt.Errorf("error creating auth failure handler lua filter: %s", err)
	}

	// TODO remove workaround
	// https://github.com/envoyproxy/envoy/issues/19910

	//For CSP we need to verify if the user has permission before routing to tenant
	//For oidc no such check is needed
	var addJwtClaimsToHeaderLuaFilter *luav3.Lua
	addJwtClaimsToHeaderLuaFilter = ConstructJWTFilter(jwtAuthnConfig.CSP, jwtAuthnConfig.JwtClaimUsername, xNexusUserIdHeader)
	addJwtClaimsToHeaderLuaFilterTypedConfig, err := anypb.New(addJwtClaimsToHeaderLuaFilter)
	if err != nil {
		return nil, fmt.Errorf("addJwtClaimsToHeaderLuaFilter: error creating typedconfig: %s", err)
	}
	addOrgIDtoHeader := ConstructOrgIDHeader(true)
	luaVerifyRoute, err := anypb.New(addOrgIDtoHeader)
	if err != nil {
		fmt.Printf("cannot add orgId header parser envoy")
	}
	addIfStaticRoute := ConstructStaticRoute()
	luaVerifyStatic, err := anypb.New(addIfStaticRoute)
	if err != nil {
		fmt.Printf("cannot add StaticURL parser envoy")
	}
	return []*hcm.HttpFilter{
		{
			Name: "envoy.filters.http.jwt_authn",
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: jwtAuthnTypedConfig,
			},
		},
		{
			Name: wellknown.Lua,
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: authFailureHandlerFilterTypedConfig,
			},
		},
		{
			Name: wellknown.Lua,
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: luaVerifyStatic,
			},
		},
		{
			Name: wellknown.Lua,
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: dummyLuaFilterTypedConfig,
			},
		},
		{
			Name: wellknown.Lua,
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: addJwtClaimsToHeaderLuaFilterTypedConfig,
			},
		},
		{
			Name: wellknown.Lua,
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: luaVerifyRoute,
			},
		},
		{
			Name: wellknown.Router,
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: router,
			},
		},
	}, nil
}

func getRouterTypedConfig() (*any.Any, error) {
	routerTypedConfig, err := anypb.New(&router.Router{})
	if err != nil {
		return nil, fmt.Errorf("error creating httpfilter router: %s", err)
	}
	return routerTypedConfig, nil
}
