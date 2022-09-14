package authn_test

import (
	"api-gw/pkg/authn"
	"api-gw/pkg/common"
	"api-gw/pkg/config"
	"api-gw/pkg/server/echo_server"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authn tests", func() {
	var e *echo_server.EchoServer

	BeforeSuite(func() {
		config.Cfg = &config.Config{
			Server:             config.ServerConfig{},
			EnableNexusRuntime: true,
			BackendService:     "",
		}
		e = echo_server.NewEchoServer(config.Cfg)

	})

	It("should register endpoints", func() {
		c := e.Echo.NewContext(nil, nil)
		authn.RegisterLoginEndpoint(e.Echo)
		e.Echo.Router().Find(http.MethodPost, common.LoginEndpoint, c)
		Expect(c.Path()).To(Equal(common.LoginEndpoint))

		authn.RegisterRefreshAccessTokenEndpoint(e.Echo)
		e.Echo.Router().Find(http.MethodPost, common.RefreshAccessTokenEndpoint, c)
		Expect(c.Path()).To(Equal(common.RefreshAccessTokenEndpoint))

		authn.RegisterLogoutEndpoint(e.Echo)
		e.Echo.Router().Find(http.MethodPost, common.LogoutEndpoint, c)
		Expect(c.Path()).To(Equal(common.LogoutEndpoint))
	})

	It("should handle login query when oicd is disabled", func() {
		authn.RegisterLoginEndpoint(e.Echo)

		req := httptest.NewRequest(http.MethodPost, common.LoginEndpoint, nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)

		err := authn.LoginHandler(c)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
	})

	// TODO NPT-455 seems that mock server is required...
	//It("should handle login query when oicd is enabled", func() {
	//	authn.RegisterLoginEndpoint(e.Echo)
	//
	//	_, err := net.Listen("tcp", "localhost:60606")
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	oicdEvent := &model.OidcNodeEvent{
	//		Oidc: authnexusv1.OIDC{
	//			Spec: authnexusv1.OIDCSpec{
	//				Config: authnexusv1.IDPConfig{
	//					ClientId:         "my id",
	//					ClientSecret:     "I'm so secret",
	//					OAuthIssuerUrl:   "http://localhost:60606",
	//					Scopes:           []string{"scope 1", "scope 2"},
	//					OAuthRedirectUrl: "http://localhost:60606",
	//				},
	//			},
	//		},
	//		Type: model.Upsert,
	//	}
	//
	//	err = authn.HandleOidcNodeUpdate(oicdEvent, e.Echo)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	req := httptest.NewRequest(http.MethodPost, common.LoginEndpoint, nil)
	//	rec := httptest.NewRecorder()
	//	c := e.Echo.NewContext(req, rec)
	//
	//	err = authn.LoginHandler(c)
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(rec.Code).To(Equal(200))
	//})

	It("should handle logout query when oicd is disabled", func() {
		authn.RegisterLogoutEndpoint(e.Echo)

		req := httptest.NewRequest(http.MethodPost, common.LogoutEndpoint, nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)

		err := authn.LogoutHandler(c)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
	})

	It("should refresh token when oicd is disabled", func() {
		authn.RegisterRefreshAccessTokenEndpoint(e.Echo)

		req := httptest.NewRequest(http.MethodPost, common.RefreshAccessTokenEndpoint, nil)
		rec := httptest.NewRecorder()
		c := e.Echo.NewContext(req, rec)

		err := authn.RefreshTokenHandler(c)
		Expect(err).NotTo(HaveOccurred())
		Expect(rec.Code).To(Equal(200))
	})

	It("should register callback endpoint when authenticator is nil", func() {
		s, err := authn.RegisterCallbackHandler(e.Echo)
		Expect(err).NotTo(HaveOccurred())
		Expect(s).To(Equal(""))
	})
})
