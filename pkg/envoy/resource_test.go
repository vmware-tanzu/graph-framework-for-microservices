package envoy_test

import (
	"api-gw/pkg/envoy"

	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("nexus proxy tests", func() {
	It("initialize nexus-proxy state with nils", func() {
		logLevel, _ := log.ParseLevel("debug")
		envoy.Init(nil, nil, nil, logLevel)
		snap, err := envoy.GenerateNewSnapshot(nil, nil, nil)
		Expect(snap).NotTo(BeNil())
		Expect(err).To(BeNil())
	})

	It("initialize nexus-proxy state with Header-based routing rules", func() {
		headerUpstreams := map[string]*envoy.HeaderMatchedUpstream{"test1": &envoy.HeaderMatchedUpstream{
			Name:        "test1",
			HeaderName:  "x-tenant",
			HeaderValue: "1",
			Host:        "example.com",
			Port:        80,
		}}
		snap, err := envoy.GenerateNewSnapshot(nil, nil, headerUpstreams)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l := snap.GetResources(resource.ListenerType)
		r := snap.GetResources(resource.RouteType)
		c := snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(1))

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok := r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts := routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(5))

		Expect(len(c)).To(Equal(2))

		// add another upstream and check assertions
		// we expect a route and a cluster to have gotten added
		headerUpstreams["test2"] = &envoy.HeaderMatchedUpstream{
			Name:        "test2",
			HeaderName:  "x-tenant",
			HeaderValue: "2",
			Host:        "google.com",
			Port:        443,
		}
		snap, err = envoy.GenerateNewSnapshot(nil, nil, headerUpstreams)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l = snap.GetResources(resource.ListenerType)
		r = snap.GetResources(resource.RouteType)
		c = snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(1))

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok = r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts = routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(6))

		Expect(len(c)).To(Equal(3))
	})

	It("initialize nexus-proxy state with JWT authn and JWT-claim based routing rules", func() {
		jwt := &envoy.JwtAuthnConfig{
			IdpName:          "csp",
			Issuer:           "https://csp.url",
			JwksUri:          "https://csp.url/jwks",
			CallbackEndpoint: "http://nexus.proxy.url/callback",
			JwtClaimUsername: "username",
		}
		upstreams := map[string]*envoy.UpstreamConfig{"test1": &envoy.UpstreamConfig{
			Name:          "test1",
			JwtClaimKey:   "username",
			JwtClaimValue: "foo@example.com",
			Host:          "example.com",
			Port:          80,
		}}

		snap, err := envoy.GenerateNewSnapshot(jwt, upstreams, nil)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l := snap.GetResources(resource.ListenerType)
		r := snap.GetResources(resource.RouteType)
		c := snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(1))
		// TODO assert that the jwt_authn filter is configured

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok := r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts := routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(6))

		Expect(len(c)).To(Equal(3))

		// add another upstream and check assertions
		// we expect a route and a cluster to have gotten added
		upstreams["test2"] = &envoy.UpstreamConfig{
			Name:          "test2",
			JwtClaimKey:   "username",
			JwtClaimValue: "bar@example.com",
			Host:          "google.com",
			Port:          443,
		}
		snap, err = envoy.GenerateNewSnapshot(jwt, upstreams, nil)
		Expect(err).To(BeNil())
		Expect(snap).ToNot(BeNil())
		l = snap.GetResources(resource.ListenerType)
		r = snap.GetResources(resource.RouteType)
		c = snap.GetResources(resource.ClusterType)
		Expect(len(l)).To(Equal(1))
		// TODO assert that the jwt_authn filter is configured

		Expect(len(r)).To(Equal(1))
		Expect(r["default"]).ToNot(BeNil())
		routes, ok = r["default"].(*routev3.RouteConfiguration)
		Expect(ok).To(Equal(true))
		Expect(len(routes.GetVirtualHosts())).To(Equal(1))
		virtualHosts = routes.GetVirtualHosts()[0]
		Expect(len(virtualHosts.GetRoutes())).To(Equal(7))

		Expect(len(c)).To(Equal(4))

	})
})
