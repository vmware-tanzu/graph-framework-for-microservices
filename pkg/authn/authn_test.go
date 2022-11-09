package authn_test

import (
	"api-gw/pkg/authn"
	"api-gw/pkg/client"
	"api-gw/pkg/common"
	"api-gw/pkg/config"
	"api-gw/pkg/envoy"
	"api-gw/pkg/model"
	"api-gw/pkg/server/echo_server"
	"api-gw/pkg/utils"
	"context"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	adminnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/admin.nexus.org/v1"
	apinexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/api.nexus.org/v1"
	apigatewaynexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/apigateway.nexus.org/v1"
	confignexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/config.nexus.org/v1"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"

	authnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.org/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Authn tests", func() {
	var e *echo_server.EchoServer
	var server testServer
	var serverURL string

	BeforeSuite(func() {
		config.Cfg = &config.Config{
			Server:             config.ServerConfig{},
			EnableNexusRuntime: true,
			BackendService:     "",
		}
		e = echo_server.NewEchoServer(config.Cfg)

		serverURL = server.run()
	})

	AfterSuite(func() {
		server.close()

		err := authn.HandleOidcNodeUpdate(&model.OidcNodeEvent{
			Oidc: authnexusv1.OIDC{},
			Type: model.Delete,
		}, e.Echo)
		Expect(err).NotTo(HaveOccurred())
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

	Context("oidc disabled", func() {
		It("should handle login query when oidc is disabled", func() {
			authn.RegisterLoginEndpoint(e.Echo)

			req := httptest.NewRequest(http.MethodPost, common.LoginEndpoint, nil)
			rec := httptest.NewRecorder()
			c := e.Echo.NewContext(req, rec)

			err := authn.LoginHandler(c)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(200))
		})

		It("should handle logout query when oidc is disabled", func() {
			authn.RegisterLogoutEndpoint(e.Echo)

			req := httptest.NewRequest(http.MethodPost, common.LogoutEndpoint, nil)
			rec := httptest.NewRecorder()
			c := e.Echo.NewContext(req, rec)

			err := authn.LogoutHandler(c)
			Expect(err).NotTo(HaveOccurred())
			Expect(rec.Code).To(Equal(200))
		})

		It("should refresh token when oidc is disabled", func() {
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

	Context("oidc enabled", func() {
		Context("blank envoy config", func() {

			err := envoy.Init(nil, nil, nil, 1)
			Expect(err).NotTo(HaveOccurred())

			It("should return 307 for login handler when when oidc is enabled but "+
				"envoy config is not initialized yet", func() {
				authn.RegisterLoginEndpoint(e.Echo)
				oidcEvent := &model.OidcNodeEvent{
					Oidc: authnexusv1.OIDC{
						ObjectMeta: metav1.ObjectMeta{
							Name: "my_name_is_luka",
						},
						Spec: authnexusv1.OIDCSpec{
							Config: authnexusv1.IDPConfig{
								ClientId:         "my id",
								ClientSecret:     "I'm so secret",
								OAuthIssuerUrl:   serverURL,
								Scopes:           []string{"scope 1", "scope 2"},
								OAuthRedirectUrl: serverURL + "/callback",
							},
						},
					},
					Type: model.Upsert,
				}

				err = authn.HandleOidcNodeUpdate(oidcEvent, e.Echo)
				Expect(err).NotTo(HaveOccurred())

				req := httptest.NewRequest(http.MethodPost, common.LoginEndpoint, nil)
				rec := httptest.NewRecorder()
				c := e.Echo.NewContext(req, rec)

				err = authn.LoginHandler(c)
				Expect(err).NotTo(HaveOccurred())
				Expect(rec.Code).To(Equal(307))
			})

			It("should return 200 for logout handler when when oidc is enabled but "+
				"envoy config is not initialized yet", func() {
				authn.RegisterLogoutEndpoint(e.Echo)
				oidcEvent := &model.OidcNodeEvent{
					Oidc: authnexusv1.OIDC{
						ObjectMeta: metav1.ObjectMeta{
							Name: "my_name_is_luka",
						},
						Spec: authnexusv1.OIDCSpec{
							Config: authnexusv1.IDPConfig{
								ClientId:         "my id",
								ClientSecret:     "I'm so secret",
								OAuthIssuerUrl:   serverURL,
								Scopes:           []string{"scope1", "scope2"},
								OAuthRedirectUrl: serverURL + "/callback",
							},
							ValidationProps: authnexusv1.ValidationProperties{
								InsecureIssuerURLContext: false,
								SkipIssuerValidation:     true,
								SkipClientIdValidation:   true,
								SkipClientAudValidation:  false,
							},
						},
					},
					Type: model.Upsert,
				}

				err := authn.HandleOidcNodeUpdate(oidcEvent, e.Echo)
				Expect(err).NotTo(HaveOccurred())

				req := httptest.NewRequest(http.MethodPost, common.LogoutEndpoint, nil)
				rec := httptest.NewRecorder()
				c := e.Echo.NewContext(req, rec)

				err = authn.LogoutHandler(c)
				Expect(err).NotTo(HaveOccurred())
				Expect(rec.Code).To(Equal(200))
			})
		})

		Context("envoy config", func() {
			log.SetLevel(log.DebugLevel)

			It("should setup envoy params", func() {
				client.NexusClient = nexus_client.NewFakeClient()

				_, err := client.NexusClient.Api().CreateNexusByName(context.TODO(), &apinexusv1.Nexus{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = client.NexusClient.Config().CreateConfigByName(context.TODO(), &confignexusv1.Config{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = client.NexusClient.Apigateway().CreateApiGatewayByName(context.TODO(), &apigatewaynexusv1.ApiGateway{
					ObjectMeta: metav1.ObjectMeta{
						Name: "default",
					},
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = client.NexusClient.Authentication().CreateOIDCByName(context.TODO(), &authnexusv1.OIDC{
					ObjectMeta: metav1.ObjectMeta{
						Name: "oidc-1",
					},
					Spec: authnexusv1.OIDCSpec{
						Config: authnexusv1.IDPConfig{
							ClientId:         "my id 2",
							ClientSecret:     "I'm so secret",
							OAuthIssuerUrl:   serverURL,
							Scopes:           []string{"scope 1", "scope 2"},
							OAuthRedirectUrl: serverURL + "/callback",
						},
						ValidationProps: authnexusv1.ValidationProperties{
							InsecureIssuerURLContext: false,
							SkipIssuerValidation:     true,
							SkipClientIdValidation:   true,
							SkipClientAudValidation:  false,
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = client.NexusClient.Admin().CreateProxyRuleByName(context.TODO(), &adminnexusv1.ProxyRule{
					ObjectMeta: metav1.ObjectMeta{
						Name: "proxy-rule-1",
					},
					Spec: adminnexusv1.ProxyRuleSpec{
						MatchCondition: adminnexusv1.MatchCondition{
							Type:  "header",
							Key:   "x-tenant",
							Value: "t-1",
						},
						Upstream: adminnexusv1.Upstream{
							Scheme: "http",
							Host:   "127.0.0.1",
							Port:   80,
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = client.NexusClient.Admin().CreateProxyRuleByName(context.TODO(), &adminnexusv1.ProxyRule{
					ObjectMeta: metav1.ObjectMeta{
						Name: "proxy-rule-2",
					},
					Spec: adminnexusv1.ProxyRuleSpec{
						MatchCondition: adminnexusv1.MatchCondition{
							Type:  "jwt",
							Key:   "foo",
							Value: "bar",
						},
						Upstream: adminnexusv1.Upstream{
							Scheme: "http",
							Host:   "127.0.0.1",
							Port:   80,
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				jwtConfig, upstreamConfigs, headerConfigs, err := utils.GetEnvoyInitParams()
				Expect(err).NotTo(HaveOccurred())

				Expect(jwtConfig).ToNot(BeNil())
				Expect(upstreamConfigs).ToNot(BeNil())
				Expect(headerConfigs).ToNot(BeEmpty())

			})

			It("should return 307 when access token is not in the request", func() {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				e.Echo.Use(authn.VerifyAuthenticationMiddleware)
				e.Echo.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(307))
			})

			It("should return 200 when valid jwt token is supplied to request", func() {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMiwiY2lkIjoidGVzdCIsImlzcyI6ImlzcyJ9.EP7B_RLsde1VYA5qJfBsEqzZH6X7nOoN2cAFmwzeUMEQvj7aQNs1V_K0x1R_rjXfrYs0yr_Ft7Eyeo77Go2MVo8fLbWgFb9pUpL8dSHcOEf4E7cxd-opQc9t7_h5gVXjpR--U2Wd-Tx2zmg9U-OjTlhtUvmotK7b5gi80lqTqx5xRXfIWjF2i2iVVl-0q_ZE9wvkLmZFTuZTgg_8Ve0V54CnKz25X8AG0OWyVaO1TxqK9B69ll6lS71LhNrfgilNBrCNCox9cnFADZ6iwS56dKZFgWlQUuKCrmI_sHyn-6jQD-hZSH0TpRV90uCVwXQ1a2fS0nhmQX--FngS5PLtPg")
				rec := httptest.NewRecorder()
				e.Echo.Use(authn.VerifyAuthenticationMiddleware)
				e.Echo.GET("/", func(c echo.Context) error {
					return c.NoContent(200)
				})
				e.Echo.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(200))
			})

			It("should return 401 when valid jwt token is invalid", func() {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJjaWQiOiJ0ZXN0IiwiaXNzIjoiaXNzIn0.YNGcMFCQ0xWucA2HpmeddQplK2i1KeKw4oQkBxF4eqQ")
				rec := httptest.NewRecorder()
				e.Echo.Use(authn.VerifyAuthenticationMiddleware)
				e.Echo.ServeHTTP(rec, req)

				Expect(rec.Code).To(Equal(401))
			})
		})

		Context("callback handler", func() {
			It("should return 307 when valid token is sent", func() {
				req := httptest.NewRequest(http.MethodGet, "/?code=123&state=/", nil)
				rec := httptest.NewRecorder()
				c := e.Echo.NewContext(req, rec)
				c.SetPath("?code=123&state=/")

				err := authn.CallbackHandler(c)
				Expect(err).NotTo(HaveOccurred())

				Expect(rec.Code).To(Equal(307))
				Expect(rec.Header().Get("Location")).To(Equal("/"))
			})

			It("should return error when error param is sent", func() {
				req := httptest.NewRequest(http.MethodGet, "/?error=test", nil)
				rec := httptest.NewRecorder()
				c := e.Echo.NewContext(req, rec)
				c.SetPath("/?error=test")

				err := authn.CallbackHandler(c)
				Expect(err).To(HaveOccurred())
			})

			It("should return error when code param is not provided", func() {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rec := httptest.NewRecorder()
				c := e.Echo.NewContext(req, rec)
				c.SetPath("/")

				err := authn.CallbackHandler(c)
				Expect(err).To(HaveOccurred())
			})

			It("should return 200 when state param is set to login endpoint", func() {
				req := httptest.NewRequest(http.MethodGet, "/?code=123&state="+common.LoginEndpoint, nil)
				rec := httptest.NewRecorder()
				c := e.Echo.NewContext(req, rec)
				c.SetPath("?code=123&state=" + common.LoginEndpoint)

				err := authn.CallbackHandler(c)
				Expect(err).NotTo(HaveOccurred())
				Expect(rec.Code).To(Equal(200))
				Expect(rec.Header().Get("Access_token")).To(Equal("test"))
			})
		})

		It("should return 200 when token is successfully refreshed", func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InRlc3QifQ.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJjaWQiOiJ0ZXN0IiwiaXNzIjoiaXNzIn0.YNGcMFCQ0xWucA2HpmeddQplK2i1KeKw4oQkBxF4eqQ")
			rec := httptest.NewRecorder()
			c := e.Echo.NewContext(req, rec)

			err := authn.RefreshTokenHandler(c)
			Expect(err).NotTo(HaveOccurred())

			Expect(rec.Code).To(Equal(200))
		})
	})
})
