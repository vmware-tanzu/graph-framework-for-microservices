package envoy

import (
	"api-gw/pkg/common"
	"fmt"

	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/durationpb"
)

func makeRoutes(tenantconfigs []*TenantConfig, jwtAuthnConfig *JwtAuthnConfig, upstreams map[string]*UpstreamConfig, headerUpstreams map[string]*HeaderMatchedUpstream) (*route.RouteConfiguration, error) {
	routes, err := getRoutes(TenantConfigs, jwtAuthnConfig, upstreams, headerUpstreams)
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
func getRoutes(tenantconfigs []*TenantConfig, jwtAuthnConfig *JwtAuthnConfig, upstreams map[string]*UpstreamConfig, headerUpstreams map[string]*HeaderMatchedUpstream) ([]*route.Route, error) {
	var routes []*route.Route

	routes = append(routes, defaultRoute())
	routes = append(routes, defaultRouteQueryParams())
	loginRoute := getLoginRoute(jwtAuthnConfig)
	if loginRoute != nil {
		routes = append(routes, loginRoute)
	}

	routes = append(routes, getRefreshTokensRoute())

	callbackRoute, err := getCallbackRoute(jwtAuthnConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get callback route: %s", err)
	}
	if callbackRoute != nil {
		routes = append(routes, callbackRoute)
	}
	if common.IsModeAdmin() {
		routes = append(routes, makeGlobalRoutes()...)
	}

	var added []string

	if common.IsModeAdmin() {
		//This is a hack to add the config for tenant only once
		//This will be removed once the core issue is resolved
		if TenantConfigs != nil {
			for _, tenantConfig := range TenantConfigs {
				log.Debugf(fmt.Sprintf("Adding route for tenant %s", tenantConfig.Name))
				for _, route := range makeTSMRoutes(tenantConfig.Name) {
					routes = append(routes, route)
				}
				added = append(added, tenantConfig.Name)

			}
		}
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
					Timeout: &durationpb.Duration{
						Seconds: 0,
					},
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: hup.Name,
					},
				},
			},
		})
	}
	return routes, nil
}

func makeGlobalRoutes() []*route.Route {
	var routes []*route.Route
	routes = append(routes, &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: "/v0",
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterNexusAdmin,
				},
			},
		},
	})
	routes = append(routes, &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Path{
				Path: "/",
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterNexusAdmin,
				},
			},
		},
	})

	routes = append(routes, &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Prefix{
				Prefix: "/",
			},
			Headers: []*route.HeaderMatcher{
				{
					Name: "static",
					HeaderMatchSpecifier: &route.HeaderMatcher_StringMatch{
						StringMatch: &matcherv3.StringMatcher{
							MatchPattern: &matcherv3.StringMatcher_Exact{
								Exact: "yes",
							},
						},
					},
				},
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: globalUI,
				},
			},
		},
	})
	return routes
}

func defaultRoute() *route.Route {
	return &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Path{
				Path: "/",
			},
		},
		Action: &route.Route_Redirect{
			Redirect: &route.RedirectAction{
				ResponseCode: route.RedirectAction_MOVED_PERMANENTLY,
				PathRewriteSpecifier: &route.RedirectAction_PathRedirect{
					PathRedirect: "/home",
				},
			},
		},
	}
}

func defaultRouteQueryParams() *route.Route {
	return &route.Route{
		Match: &route.RouteMatch{
			QueryParameters: []*route.QueryParameterMatcher{
				{
					Name: "orgLink",
					QueryParameterMatchSpecifier: &route.QueryParameterMatcher_PresentMatch{
						PresentMatch: true,
					},
				},
			},
			PathSpecifier: &route.RouteMatch_Path{
				Path: "/",
			},
		},
		Action: &route.Route_Redirect{
			Redirect: &route.RedirectAction{
				ResponseCode: route.RedirectAction_MOVED_PERMANENTLY,
				PathRewriteSpecifier: &route.RedirectAction_PathRedirect{
					PathRedirect: "/home",
				},
			},
		},
	}
}

func makeTSMRoutes(TenantName string) []*route.Route {
	var headersMatch []*route.HeaderMatcher
	var routes []*route.Route
	match := &matcherv3.StringMatcher{
		MatchPattern: &matcherv3.StringMatcher_Exact{
			Exact: TenantName,
		},
	}

	headersMatch = append(headersMatch, &route.HeaderMatcher{
		Name: "org-id",
		HeaderMatchSpecifier: &route.HeaderMatcher_StringMatch{
			StringMatch: match,
		},
	})

	for _, svc := range cosmosServices {
		var match *route.RouteMatch
		var cluster string
		fmt.Printf("Adding route for service %s for tenant %s", svc.Name, TenantName)
		for _, prefixRoute := range svc.PrefixSettings {
			var prefixRewrite string
			if prefixRoute.PrefixRewrite != "" {
				prefixRewrite = prefixRoute.PrefixRewrite
			} else {
				prefixRewrite = prefixRoute.Prefix
			}
			var overallheadersMatch []*route.HeaderMatcher
			overallheadersMatch = headersMatch
			if prefixRoute.Header != "" {
				overallheadersMatch = append(headersMatch, &route.HeaderMatcher{
					Name: prefixRoute.Header,
					HeaderMatchSpecifier: &route.HeaderMatcher_StringMatch{
						StringMatch: &matcherv3.StringMatcher{
							MatchPattern: &matcherv3.StringMatcher_Exact{
								Exact: "yes",
							},
						},
					},
				})
			}
			if !svc.Global {
				match = &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: prefixRoute.Prefix,
					},
					Headers: overallheadersMatch,
				}
				cluster = fmt.Sprintf("%s-%s", svc.Name, TenantName)
			} else {
				match = &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_Prefix{
						Prefix: prefixRoute.Prefix,
					},
				}
				cluster = svc.Name
			}
			if (jwt != nil && prefixRoute.Prefix != "/login") || (jwt == nil) {
				routes = append(routes, &route.Route{
					Match: match,
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: cluster,
							},
							PrefixRewrite: prefixRewrite,
						},
					},
				})
			}
		}
	}
	return routes
}

func getLoginRoute(jwt *JwtAuthnConfig) *route.Route {
	var routeAction *route.Route_Route
	if jwt == nil {
		if TenantConfigs != nil {
			routeAction = &route.Route_Route{
				Route: &route.RouteAction{
					Timeout: &durationpb.Duration{
						Seconds: 0,
					},
					ClusterSpecifier: &route.RouteAction_Cluster{
						Cluster: globalUI,
					},
				},
			}
		}
	} else if common.IsModeAdmin() {
		routeAction = &route.Route_Route{
			Route: &route.RouteAction{
				Timeout: &durationpb.Duration{
					Seconds: 0,
				},
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: clusterNexusAdmin,
				},
			},
		}
	}
	if routeAction != nil {
		return &route.Route{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Path{
					Path: "/login",
				},
			},
			Action: routeAction,
		}
	}
	return nil
}

func getRefreshTokensRoute() *route.Route {
	if jwt == nil {
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
	} else {
		return &route.Route{
			Match: &route.RouteMatch{
				PathSpecifier: &route.RouteMatch_Path{
					Path: jwt.RefreshTokenEndpoint,
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
