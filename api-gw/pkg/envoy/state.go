package envoy

import (
	"context"
	"fmt"
	"net"
	"sync"

	cachev3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	cache           cachev3.SnapshotCache
	jwt             *JwtAuthnConfig
	upstreams       map[string]*UpstreamConfig
	headerUpstreams map[string]*HeaderMatchedUpstream
	TenantConfigs   []*TenantConfig
)

const (
	jwtPayload    = "jwt_payload"
	envoyNodeId   = "envoy-nexus-admin"
	xDSListenPort = 18000
)

var jwtMutex sync.Mutex
var upstreamsMutex sync.Mutex
var headerUpstreamMutex sync.Mutex
var refreshEnvoyMutex sync.Mutex
var XDSServer *grpc.Server
var XDSListener net.Listener

func Init(j *JwtAuthnConfig, u map[string]*UpstreamConfig, hu map[string]*HeaderMatchedUpstream, level log.Level) error {
	log.Infof("initializing xDS server...")
	jwt = j
	if u == nil {
		upstreams = make(map[string]*UpstreamConfig)
	} else {
		upstreams = u
	}
	if hu == nil {
		headerUpstreams = make(map[string]*HeaderMatchedUpstream)
	} else {
		headerUpstreams = hu
	}

	logger := log.New()
	log.SetLevel(level)

	// Create a cache
	cache = cachev3.NewSnapshotCache(false, cachev3.IDHash{}, logger)

	// Create the snapshot that we'll serve to Envoy
	snapshot, err := GenerateNewSnapshot(nil, nil, nil, nil)
	if err != nil {
		log.Errorf("failed to generate a new snapshot: %s", err)
		return err
	}
	if err = snapshot.Consistent(); err != nil {
		log.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		return err
	}
	log.Debugf("will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err = cache.SetSnapshot(context.Background(), envoyNodeId, snapshot); err != nil {
		log.Errorf("snapshot error %q for %+v", err, snapshot)
		return err
	}

	//stopCh := make(chan struct{})
	globalCtx := context.Background()
	XDSListener = CreateXDSListener(globalCtx, xDSListenPort)

	srv := server.NewServer(globalCtx, cache, nil)
	XDSServer = RegisterServer(globalCtx, srv, xDSListenPort)

	go func() {
		// Run the xDS server
		RunServer(globalCtx, XDSServer, xDSListenPort)

	}()

	return nil
}

func RefreshEnvoyConfiguration() error {
	refreshEnvoyMutex.Lock()
	defer refreshEnvoyMutex.Unlock()

	log.Debugf("refreshing envoy configuration...")
	snapshot, err := GenerateNewSnapshot(TenantConfigs, jwt, upstreams, headerUpstreams)
	if err != nil {
		log.Errorf("failed to generate a new snapshot: %s", err)
		return err
	}

	if err = snapshot.Consistent(); err != nil {
		log.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		return err
	}
	log.Debugf("will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err = cache.SetSnapshot(context.Background(), envoyNodeId, snapshot); err != nil {
		log.Errorf("snapshot error %q for %+v", err, snapshot)
		return err
	}
	log.Debugf("successfully refreshed envoy configuration")
	return nil
}

func AddTenantConfig(tenantConfig *TenantConfig) error {
	for _, tenantconfig_ := range TenantConfigs {
		if tenantconfig_.Name == tenantConfig.Name {
			err := RefreshEnvoyConfiguration()
			if err != nil {
				return err
			}
			return nil
		}
	}
	TenantConfigs = append(TenantConfigs, tenantConfig)
	err := RefreshEnvoyConfiguration()
	if err != nil {
		return fmt.Errorf("AddTenantConfig: error while refreshing envoy configuration: %s", err)
	}
	return nil
}

func RemoveIndex(s []*TenantConfig, index int) []*TenantConfig {
	return append(s[:index], s[index+1:]...)
}

func DeleteTenantConfig(tenantname string) error {
	var index int
	for i, object := range TenantConfigs {
		if object.Name == tenantname {
			index = i
		}
	}
	TenantConfigs = RemoveIndex(TenantConfigs, index)
	err := RefreshEnvoyConfiguration()
	if err != nil {
		return fmt.Errorf("DeleteTenantConfig: error while refreshing envoy configuration: %s", err)
	}
	return nil

}

func AddJwtAuthnConfig(jwtAuthnConfig *JwtAuthnConfig) error {
	jwtMutex.Lock()
	defer jwtMutex.Unlock()

	jwt = jwtAuthnConfig
	err := RefreshEnvoyConfiguration()
	if err != nil {
		return fmt.Errorf("AddJwtAuthnConfig: error while refreshing envoy configuration: %s", err)
	}
	return nil
}

func DeleteJwtAuthnConfig() error {
	jwtMutex.Lock()
	defer jwtMutex.Unlock()

	jwt = nil
	err := RefreshEnvoyConfiguration()
	if err != nil {
		return fmt.Errorf("DeleteJwtAuthnConfig: error while refreshing envoy configuration: %s", err)
	}
	return nil
}

func AddUpstream(name string, upstream *UpstreamConfig) error {
	upstreamsMutex.Lock()
	defer upstreamsMutex.Unlock()

	upstreams[name] = upstream
	err := RefreshEnvoyConfiguration()
	if err != nil {
		return fmt.Errorf("AddUpstream: error while refreshing envoy configuration: %s", err)
	}
	return nil
}

func DeleteUpstream(name string) error {
	upstreamsMutex.Lock()
	defer upstreamsMutex.Unlock()

	delete(upstreams, name)
	err := RefreshEnvoyConfiguration()
	if err != nil {
		return fmt.Errorf("DeleteUpstream: error while refreshing envoy configuration: %s", err)
	}
	return nil
}

func AddHeaderUpstream(name string, headerUpstream *HeaderMatchedUpstream) error {
	headerUpstreamMutex.Lock()
	defer headerUpstreamMutex.Unlock()

	headerUpstreams[name] = headerUpstream
	err := RefreshEnvoyConfiguration()
	if err != nil {
		return fmt.Errorf("AddHeaderUpstream: error while refreshing envoy configuration: %s", err)
	}
	return nil
}

func DeleteHeaderUpstream(name string) error {
	upstreamsMutex.Lock()
	defer upstreamsMutex.Unlock()

	delete(headerUpstreams, name)
	err := RefreshEnvoyConfiguration()
	if err != nil {
		return fmt.Errorf("DeleteHeaderUpstream: error while refreshing envoy configuration: %s", err)
	}
	return nil
}
