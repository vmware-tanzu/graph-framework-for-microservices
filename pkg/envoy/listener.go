package envoy

import (
	"api-gw/pkg/common"
	"fmt"
	"time"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	jwtauthnv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/jwt_authn/v3"
	luav3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/lua/v3"
	router "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/any"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

func makeHTTPListener(jwtAuthnConfig *JwtAuthnConfig) (*listener.Listener, error) {
	listenerName := "listener_0"

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

	return &listener.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: ListenerPort,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: pbst,
				},
			}},
		}},
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

	if jwtAuthnConfig == nil {
		return []*hcm.HttpFilter{{
			Name: wellknown.Router,
			ConfigType: &hcm.HttpFilter_TypedConfig{
				TypedConfig: router,
			},
		}}, nil
	}

	if jwtAuthnConfig.Issuer == "" || jwtAuthnConfig.JwksUri == "" || jwtAuthnConfig.IdpName == "" || jwtAuthnConfig.CallbackEndpoint == "" {
		return nil, fmt.Errorf("failed to create JWT authn filter: invalid config")
	}

	jwtAuthn := &jwtauthnv3.JwtAuthentication{
		Providers: map[string]*jwtauthnv3.JwtProvider{
			jwtAuthnConfig.IdpName: {
				Issuer:      jwtAuthnConfig.Issuer,
				FromCookies: []string{common.AccessTokenStr},
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
					PathSpecifier: &route.RouteMatch_Path{
						Path: common.LoginEndpoint,
					},
				},
			},
			{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Path{
						Path: common.RefreshAccessTokenEndpoint,
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
        response_handle:logInfo("Got status 401, redirect to login...")
        response_handle:headers():replace(":status", "307")
        response_handle:headers():add("location", "%s")
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

func getRouterTypedConfig() (*any.Any, error) {
	routerTypedConfig, err := anypb.New(&router.Router{})
	if err != nil {
		return nil, fmt.Errorf("error creating httpfilter router: %s", err)
	}
	return routerTypedConfig, nil
}
