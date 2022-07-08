package envoy

import (
	"time"

	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	log "github.com/sirupsen/logrus"
)

// JwtAuthnConfig holds information that is used to configure envoy's jwt_authn filter
type JwtAuthnConfig struct {
	IdpName          string
	Issuer           string
	JwksUri          string
	CallbackEndpoint string
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
	ListenerPort      = 10000
	clusterNexusAdmin = "nexus-admin"
	routeDefault      = "default"
	svcNexusApiGw     = "nexus-api-gw"
)

func GenerateNewSnapshot(jwtAuthnConfig *JwtAuthnConfig, upstreams map[string]*UpstreamConfig, headerUpstreams map[string]*HeaderMatchedUpstream) (*cachev3.Snapshot, error) {
	httpListener, err := makeHTTPListener(jwtAuthnConfig)
	if err != nil {
		log.Errorf("failed to create http listener: %s", err)
		return nil, err
	}

	var routes *route.RouteConfiguration
	routes, err = makeRoutes(jwtAuthnConfig, upstreams, headerUpstreams)
	if err != nil {
		log.Errorf("failed to create routes: %s", err)
		return nil, err
	}

	var clusters []types.Resource
	clusters, err = makeClusters(jwtAuthnConfig, upstreams, headerUpstreams)
	if err != nil {
		log.Errorf("failed to create clusters: %s", err)
		return nil, err
	}

	snap, _ := cachev3.NewSnapshot(time.Now().String(),
		map[resource.Type][]types.Resource{
			resource.ListenerType: {httpListener},
			resource.RouteType:    {routes},
			resource.ClusterType:  clusters,
		},
	)
	return snap, nil
}
