package envoy

import (
	"api-gw/pkg/common"
	"fmt"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
)

func makeRoutes(jwtAuthnConfig *JwtAuthnConfig, upstreams map[string]*UpstreamConfig, headerUpstreams map[string]*HeaderMatchedUpstream) (*route.RouteConfiguration, error) {
	routes, err := getRoutes(jwtAuthnConfig, upstreams, headerUpstreams)
	if err != nil {
		return nil, fmt.Errorf("failed to build routes: %s", err)
	} else {
		return &route.RouteConfiguration{
			Name: routeDefault,
			VirtualHosts: []*route.VirtualHost{{
				Name:    "nexus-admin-svc",
				Domains: []string{"*"},
				Routes:  routes,
			}},
		}, nil
	}
}

// getRoutes returns an ordered list of routes that envoy will try to match sequentially
func getRoutes(jwtAuthnConfig *JwtAuthnConfig, upstreams map[string]*UpstreamConfig, headerUpstreams map[string]*HeaderMatchedUpstream) ([]*route.Route, error) {
	var routes []*route.Route

	routes = append(routes, getLoginRoute(), getRefreshTokensRoute(), getLogoutRoute())

	callbackRoute, err := getCallbackRoute(jwtAuthnConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get callback route: %s", err)
	}
	if callbackRoute != nil {
		routes = append(routes, callbackRoute)
	}

	var upstreamRoutes []*route.Route
	upstreamRoutes, err = getUpstreamRoutes(upstreams)
	if err != nil {
		return nil, fmt.Errorf("error getting upstream routes: %s", err)
	}
	if len(upstreamRoutes) > 0 {
		routes = append(routes, upstreamRoutes...)
	}

	var headerUpstreamRoutes []*route.Route
	headerUpstreamRoutes, err = getHeaderUpstreamRoutes(headerUpstreams)
	if err != nil {
		return nil, fmt.Errorf("error getting header upstream routes: %s", err)
	}
	if len(headerUpstreamRoutes) > 0 {
		routes = append(routes, headerUpstreamRoutes...)
	}

	routes = append(routes, defaultRoute())
	return routes, nil
}

func getUpstreamRoutes(upstreams map[string]*UpstreamConfig) ([]*route.Route, error) {
	var routes []*route.Route
	for _, upstream := range upstreams {
		if upstream.JwtClaimKey == "" || upstream.JwtClaimValue == "" {
			return nil, fmt.Errorf("invalid jwt match condition")
		}
		routes = append(routes, &route.Route{
			// route based on JWT content
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
				DynamicMetadata: []*matcherv3.MetadataMatcher{
					{
						Filter: "envoy.filters.http.jwt_authn",
						Path: []*matcherv3.MetadataMatcher_PathSegment{
							{
								Segment: &matcherv3.MetadataMatcher_PathSegment_Key{
									Key: jwtPayload,
								},
							},
							{
								Segment: &matcherv3.MetadataMatcher_PathSegment_Key{
									Key: upstream.JwtClaimKey,
								},
							},
						},
						Value: &matcherv3.ValueMatcher{
							MatchPattern: &matcherv3.ValueMatcher_StringMatch{
								StringMatch: &matcherv3.StringMatcher{
									MatchPattern: &matcherv3.StringMatcher_Exact{
										Exact: upstream.JwtClaimValue,
									},
								},
							},
						},
					},
				},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: upstream.Name,
					},
				},
			},
		})
	}
	return routes, nil
}

func getHeaderUpstreamRoutes(headerUpstreams map[string]*HeaderMatchedUpstream) ([]*route.Route, error) {
	var routes []*route.Route
	for _, hup := range headerUpstreams {
		if hup.HeaderName == "" || hup.HeaderValue == "" {
			return nil, fmt.Errorf("found invalid header match condition")
		}
		routes = append(routes, &route.Route{
			// route based on JWT content
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Prefix{
					Prefix: "/",
				},
				Headers: []*route.HeaderMatcher{{
					Name: hup.HeaderName,
					HeaderMatchSpecifier: &route.HeaderMatcher_StringMatch{
						StringMatch: &matcherv3.StringMatcher{
							MatchPattern: &matcherv3.StringMatcher_Exact{
								Exact: hup.HeaderValue,
							},
						},
					},
				}},
			},
			Action: &route.Route_Route{
				Route: &route.RouteAction{
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: hup.Name,
					},
				},
			},
		})
	}
	return routes, nil
}

func defaultRoute() *route.Route {
	return &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: "/",
			},
		},
		Action: &route.Route_DirectResponse{
			DirectResponse: &route.DirectResponseAction{
				Status: 400,
				Body: &core.DataSource{
					Specifier: &core.DataSource_InlineString{
						InlineString: "No Route Found",
					},
				},
			},
		},
	}
}

func getLoginRoute() *route.Route {
	return &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: common.LoginEndpoint,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterNexusAdmin,
				},
			},
		},
	}
}

func getRefreshTokensRoute() *route.Route {
	return &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Path{
				Path: common.RefreshAccessTokenEndpoint,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterNexusAdmin,
				},
			},
		},
	}
}

func getLogoutRoute() *route.Route {
	return &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Path{
				Path: common.LogoutEndpoint,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterNexusAdmin,
				},
			},
		},
	}
}

func getCallbackRoute(jwtAuthnConfig *JwtAuthnConfig) (*route.Route, error) {
	if jwtAuthnConfig == nil {
		return nil, nil
	} else if jwtAuthnConfig.CallbackEndpoint == "" {
		return nil, fmt.Errorf("empty callback path")
	}

	return &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Path{
				Path: jwtAuthnConfig.CallbackEndpoint,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterNexusAdmin,
				},
			},
		},
	}, nil
}
