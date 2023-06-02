package envoy

import (
	"api-gw/pkg/common"
	"time"

	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	log "github.com/sirupsen/logrus"
)

type CustomHeader struct {
	FieldName string
	Header    string
}

// JwtAuthnConfig holds information that is used to configure envoy's jwt_authn filter
type JwtAuthnConfig struct {
	IdpName              string
	Issuer               string
	JwksUri              string
	CallbackEndpoint     string
	RefreshTokenEndpoint string
	JwtClaimUsername     string
	CSP                  bool
	AccessToken          string
}

type TenantConfig struct {
	Name   string
	Status bool
}

type PrefixSettings struct {
	Prefix        string
	PrefixRewrite string
	Header        string
}
type CosmosService struct {
	Name            string
	Port            uint32
	Svc             string
	PrefixSettings  []PrefixSettings
	ExactPath       []PrefixSettings
	Global          bool
	AdditionalMatch map[string]string
	HeaderMatches   map[string]string
}

// Routes should be added in this order to avoid route ambuiguity issues
// Example the / default route is added at bottom to avoid confusion between /tsm to go to tenant api-gateway instead of api-gateway
var cosmosServices []CosmosService = []CosmosService{
	{
		Name: "nexus-api-gw",
		Port: 80,
		Svc:  "nexus-api-gw",
		PrefixSettings: []PrefixSettings{
			{
				Prefix: "/declarative/",
			},
			{
				Prefix: "/apis",
			},
			{
				Prefix:        "/tsm/explorer/",
				PrefixRewrite: "/explorer/",
			},
		},
	},
	{
		Name:   "allspark-ui",
		Port:   80,
		Svc:    "allspark-ui",
		Global: true,
		PrefixSettings: []PrefixSettings{
			{
				Prefix: "/home",
			},
			{
				Prefix: "/login",
			},
		},
		HeaderMatches: map[string]string{
			"static": "yes",
		},
	},
	{
		Name: "tenant-api-gw",
		Port: 3000,
		Svc:  "tenant-api-gw",
		PrefixSettings: []PrefixSettings{
			{
				Prefix:        "/tsm/",
				PrefixRewrite: "/",
			},
		},
	},
	{
		Name: "local-api-gateway",
		Port: 3000,
		Svc:  "local-api-gateway",
		PrefixSettings: []PrefixSettings{
			{
				Prefix:        "/local/",
				PrefixRewrite: "/",
			},
		},
	},
	{
		Name: "nexus-api-gw",
		Port: 80,
		Svc:  "nexus-api-gw",
		PrefixSettings: []PrefixSettings{
			{
				Prefix: "/",
			},
		},
	},
}

// UpstreamConfig defines a condition (jwt_payload[JwtClaimKey] == JwtClaimValue), if matched,
// causes the request to be proxied to the upstream Host:Port
type UpstreamConfig struct {
	Name          string
	JwtClaimKey   string
	JwtClaimValue string
	Host          string
	Port          uint32
}

// HeaderMatchedUpstream defines a condition (headers[HeaderName] == HeaderValue), if matched,
// causes the request to be proxied to the upstream Host:Port
type HeaderMatchedUpstream struct {
	Name        string
	HeaderName  string
	HeaderValue string
	Host        string
	Port        uint32
}

const (
	//Keeping this as 10000 and 10001 as k8s version above 1.24 requires to use non admin ports
	HttpListenerPort   = 10000
	HttpsListenerPort  = 10001
	clusterNexusAdmin  = "nexus-admin"
	routeDefault       = "default"
	svcNexusApiGw      = "nexus-api-gw"
	xNexusUserIdHeader = "x-user-id"
	globalUI           = "allspark-ui"
	globalUIPort       = 80
	globalUIsvcname    = "allspark-ui"
)

func GenerateNewSnapshot(tenantConfigs []*TenantConfig, jwtAuthnConfig *JwtAuthnConfig, upstreams map[string]*UpstreamConfig, headerUpstreams map[string]*HeaderMatchedUpstream) (*cachev3.Snapshot, error) {
	routeListener, err := makeRouteListener(jwtAuthnConfig)
	if err != nil {
		log.Errorf("failed to create route listener: %s", err)
		return nil, err
	}

	//This is for adding redirection
	var httpListener *listenerv3.Listener
	if common.IsHttpsEnabled() {
		httpListener, err = makeHttpListener()
		if err != nil {
			log.Errorf("failed to create http listener: %s", err)
			return nil, err
		}
	}

	var routes *route.RouteConfiguration
	routes, err = makeRoutes(tenantConfigs, jwtAuthnConfig, upstreams, headerUpstreams)
	if err != nil {
		log.Errorf("failed to create routes: %s", err)
		return nil, err
	}

	var clusters []types.Resource
	clusters, err = makeClusters(tenantConfigs, jwtAuthnConfig, upstreams, headerUpstreams)
	if err != nil {
		log.Errorf("failed to create clusters: %s", err)
		return nil, err
	}
	var snap *cachev3.Snapshot

	if common.IsHttpsEnabled() {
		snap, _ = cachev3.NewSnapshot(time.Now().String(),
			map[resource.Type][]types.Resource{
				resource.ListenerType: {routeListener, httpListener},
				resource.RouteType:    {routes},
				resource.ClusterType:  clusters,
			},
		)
	} else {
		snap, _ = cachev3.NewSnapshot(time.Now().String(),
			map[resource.Type][]types.Resource{
				resource.ListenerType: {routeListener},
				resource.RouteType:    {routes},
				resource.ClusterType:  clusters,
			},
		)
	}
	return snap, nil
}
