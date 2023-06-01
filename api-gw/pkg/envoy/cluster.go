package envoy

import (
	"api-gw/pkg/common"
	"fmt"
	"net/url"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

func makeTSMTenantClusters(tenantconfigs []*TenantConfig) (clusters []types.Resource, err error) {
	if tenantconfigs == nil {
		return nil, nil
	}
	for _, tenant := range tenantconfigs {
		for _, svc := range cosmosServices {
			if !svc.Global {
				clusterRoute, err := makeCluster(fmt.Sprintf("%s-%s", svc.Name, tenant.Name), fmt.Sprintf("%s.%s", svc.Name, tenant.Name), svc.Port, nil)
				if err != nil {
					return nil, fmt.Errorf("failed to create %s cluster :%s", fmt.Sprintf("%s-%s", svc.Name, tenant.Name), err)
				}
				clusters = append(clusters, clusterRoute)
			}
		}
	}
	return clusters, nil
}

func makeClusters(tenantConfigs []*TenantConfig, jwtAuthnConfig *JwtAuthnConfig, upstreams map[string]*UpstreamConfig, headerUpstreams map[string]*HeaderMatchedUpstream) ([]types.Resource, error) {
	var clusters []types.Resource

	if common.IsModeAdmin() {
		adminCluster, err := makeAdminCluster()
		if err != nil {
			return nil, fmt.Errorf("failed to create admin cluster :%s", err)
		}
		if adminCluster != nil {
			clusters = append(clusters, adminCluster)
		}
	}

	uiCluster, err := makeCluster(globalUI, globalUIsvcname, globalUIPort, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create globalUI cluster: %s", err)
	}
	clusters = append(clusters, uiCluster)

	if tenantConfigs != nil {
		tenantRoutes, err := makeTSMTenantClusters(tenantConfigs)
		if err != nil {
			return nil, fmt.Errorf("failed to create tenantroutes cluster: %s", err)
		}
		clusters = append(clusters, tenantRoutes...)
	}

	var jwksCluster types.Resource
	jwksCluster, err = makeJwksCluster(jwtAuthnConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create jwks cluster: %s", err)
	}
	if jwksCluster != nil {
		clusters = append(clusters, jwksCluster)
	}

	var upstreamClusters []types.Resource
	upstreamClusters, err = makeUpstreamClusters(upstreams)
	if err != nil {
		return nil, fmt.Errorf("failed to create upstream clusters: %s", err)
	}
	if len(upstreamClusters) > 0 {
		clusters = append(clusters, upstreamClusters...)
	}

	var headerUpstreamClusters []types.Resource
	headerUpstreamClusters, err = makeHeaderUpstreamClusters(headerUpstreams)
	if err != nil {
		return nil, fmt.Errorf("failed to create header upstream clusters: %s", err)
	}
	if len(headerUpstreamClusters) > 0 {
		clusters = append(clusters, headerUpstreamClusters...)
	}
	return clusters, nil
}

func makeAdminCluster() (*cluster.Cluster, error) {
	adminCluster, err := makeCluster(clusterNexusAdmin, svcNexusApiGw, 80, nil)
	if err != nil {
		return nil, err
	}
	return adminCluster, nil
}

func makeUpstreamClusters(upstreams map[string]*UpstreamConfig) ([]types.Resource, error) {
	var upstreamClusters []types.Resource

	for _, upstream := range upstreams {
		cluster, err := makeCluster(upstream.Name, upstream.Host, upstream.Port, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create cluster for upstream %s", upstream.Name)
		}
		upstreamClusters = append(upstreamClusters, cluster)
	}
	return upstreamClusters, nil
}

func makeHeaderUpstreamClusters(headerUpstreams map[string]*HeaderMatchedUpstream) ([]types.Resource, error) {
	var headerUpstreamClusters []types.Resource

	for _, headerUpstream := range headerUpstreams {
		cluster, err := makeCluster(headerUpstream.Name, headerUpstream.Host, headerUpstream.Port, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create cluster for header upstream %s", headerUpstream.Name)
		}
		headerUpstreamClusters = append(headerUpstreamClusters, cluster)
	}
	return headerUpstreamClusters, nil
}

func makeCluster(name, host string, port uint32, tlsTransport *core.TransportSocket) (*cluster.Cluster, error) {
	if name == "" || host == "" || port == 0 {
		return nil, fmt.Errorf("invalid upstream found for cluster %s", name)
	}
	return &cluster.Cluster{
		Name:                 name,
		DnsLookupFamily:      cluster.Cluster_V4_PREFERRED,
		ConnectTimeout:       durationpb.New(30 * time.Second),
		ClusterDiscoveryType: &cluster.Cluster_Type{Type: cluster.Cluster_LOGICAL_DNS},
		LbPolicy:             cluster.Cluster_ROUND_ROBIN,
		LoadAssignment:       makeEndpoint(name, host, port),
		TransportSocket:      tlsTransport,
		HttpProtocolOptions: &core.Http1ProtocolOptions{
			AllowChunkedLength: true,
		},
	}, nil
}

func makeJwksCluster(jwtAuthnConfig *JwtAuthnConfig) (types.Resource, error) {
	if jwtAuthnConfig == nil {
		return nil, nil
	}

	host, err := getHostFromUri(jwtAuthnConfig.JwksUri)
	if err != nil {
		return nil, fmt.Errorf("failed to create a jwks cluster for %s: %s", jwtAuthnConfig.IdpName, err)
	}

	clusterName := fmt.Sprintf("%s_jwks_cluster", jwtAuthnConfig.IdpName)
	tls := &tlsv3.UpstreamTlsContext{
		Sni: host,
	}

	var tlsTypedConfig *anypb.Any
	tlsTypedConfig, err = anypb.New(tls)
	if err != nil {
		return nil, fmt.Errorf("failed to create tlsTypedConfig for %s: %s", jwtAuthnConfig.IdpName, err)
	}

	var cluster types.Resource
	cluster, err = makeCluster(clusterName, host, 443, &core.TransportSocket{
		Name: wellknown.TransportSocketTls,
		ConfigType: &core.TransportSocket_TypedConfig{
			TypedConfig: tlsTypedConfig,
		},
	})
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func getHostFromUri(uri string) (string, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		log.Errorf("failed to parse request uri %s: %s", uri, err)
		return "", err
	}
	return u.Hostname(), nil
}

func makeEndpoint(clusterName string, address string, port uint32) *endpoint.ClusterLoadAssignment {
	return &endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*endpoint.LocalityLbEndpoints{{
			LbEndpoints: []*endpoint.LbEndpoint{{
				HostIdentifier: &endpoint.LbEndpoint_Endpoint{
					Endpoint: &endpoint.Endpoint{
						Address: &core.Address{
							Address: &core.Address_SocketAddress{
								SocketAddress: &core.SocketAddress{
									Protocol: core.SocketAddress_TCP,
									Address:  address,
									PortSpecifier: &core.SocketAddress_PortValue{
										PortValue: port,
									},
								},
							},
						},
					},
				},
			}},
		}},
	}
}
