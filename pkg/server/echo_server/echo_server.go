package echo_server

import (
	"api-gw/pkg/authn"
	"api-gw/pkg/common"
	"api-gw/pkg/openapi/api"
	"api-gw/pkg/openapi/declarative"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	log "github.com/sirupsen/logrus"

	"api-gw/pkg/config"
	"api-gw/pkg/model"
	"api-gw/pkg/utils"

	openMiddleware "github.com/go-openapi/runtime/middleware"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type EchoServer struct {
	Echo   *echo.Echo
	Config *config.Config
}

func InitEcho(stopCh chan struct{}, conf *config.Config) {
	log.Infoln("Init Echo")
	e := NewEchoServer(conf)

	if config.Cfg.EnableNexusRuntime {
		e.RegisterNexusRoutes()
	}

	if config.Cfg.BackendService != "" {
		e.RegisterDeclarativeRoutes()
		e.RegisterDeclarativeRouter()
	}

	common.Mode = os.Getenv("GATEWAY_MODE")
	log.Infof("Gateway Mode: %s", common.Mode)

	e.Start(stopCh)
}

func (s *EchoServer) StartHTTPServer() {
	if err := s.Echo.Start(":80"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error %v", err)
	}
}

func (s *EchoServer) Start(stopCh chan struct{}) {
	if config.Cfg.EnableNexusRuntime {
		// Start watching URI notification
		go func() {
			log.Debug("NodeUpdateNotifications")
			if err := s.NodeUpdateNotifications(stopCh); err != nil {
				s.StopServer()
				InitEcho(stopCh, s.Config)
			}
		}()
	}

	// Start Server
	go func() {
		log.Info("Start Echo Server")
		if utils.IsServerConfigValid(s.Config) && utils.IsFileExists(s.Config.Server.CertPath) && utils.IsFileExists(s.Config.Server.KeyPath) {
			log.Infof("Server Config %v", s.Config.Server)
			log.Info("Start TLS Server")
			if err := s.Echo.StartTLS(s.Config.Server.Address, s.Config.Server.CertPath, s.Config.Server.KeyPath); err != nil && err != http.ErrServerClosed {
				log.Fatalf("TLS Server error %v", err)
			}
		} else {
			log.Info("Certificates or TLS port not configured correctly, hence starting the HTTP Server")
			s.StartHTTPServer()
		}
	}()
}

type NexusContext struct {
	echo.Context
	NexusURI string
	Codes    nexus.HTTPCodesResponse

	// Kube
	CrdType   string
	GroupName string
	Resource  string
}

func (s *EchoServer) RegisterNexusRoutes() {
	// OpenAPI route
	s.Echo.GET("/openapi.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api.Schema)
	})

	// Swagger-UI
	opts := openMiddleware.SwaggerUIOpts{
		SpecURL: "/openapi.json",
		Title:   "API Gateway Documentation",
	}
	s.Echo.GET("/docs", echo.WrapHandler(openMiddleware.SwaggerUI(opts, nil)))

	_, err := authn.RegisterCallbackHandler(s.Echo)
	if err != nil {
		log.Errorln("Error registering the OIDC callback path")
		// should we panic?
	}
	authn.RegisterLoginEndpoint(s.Echo)
	authn.RegisterRefreshAccessTokenEndpoint(s.Echo)
	authn.RegisterLogoutEndpoint(s.Echo)
}

func (s *EchoServer) RegisterDeclarativeRoutes() {
	s.Echo.GET("/declarative/apis", func(c echo.Context) error {
		return c.JSON(200, declarative.ApisList)
	})
}

func (s *EchoServer) RegisterRouter(restURI nexus.RestURIs) {
	urlPattern := model.ConstructEchoPathParamURL(restURI.Uri)
	for method, codes := range restURI.Methods {
		log.Infof("Registered Router Path %s Method %s\n", urlPattern, method)
		switch method {
		// in "admin" mode, the responsibility of authentication is offloaded to the nexus-proxy.
		// so we don't need to add the authn.VerifyAuthenticationMiddleware middleware
		case http.MethodGet:
			if common.IsModeAdmin() {
				s.Echo.GET(urlPattern, getHandler, getNexusContext(restURI, codes))
			} else {
				s.Echo.GET(urlPattern, getHandler, authn.VerifyAuthenticationMiddleware, getNexusContext(restURI, codes))
			}
		case http.MethodPut:
			if common.IsModeAdmin() {
				s.Echo.PUT(urlPattern, putHandler, getNexusContext(restURI, codes))
			} else {
				s.Echo.PUT(urlPattern, putHandler, authn.VerifyAuthenticationMiddleware, getNexusContext(restURI, codes))
			}
		case http.MethodDelete:
			if common.IsModeAdmin() {
				s.Echo.DELETE(urlPattern, deleteHandler, getNexusContext(restURI, codes))
			} else {
				s.Echo.DELETE(urlPattern, deleteHandler, authn.VerifyAuthenticationMiddleware, getNexusContext(restURI, codes))
			}
		}
	}
}

func getNexusContext(restURI nexus.RestURIs, codes nexus.HTTPCodesResponse) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:  c,
				NexusURI: restURI.Uri,
				Codes:    codes,
			}
			return next(nc)
		}
	}
}

func (s *EchoServer) RegisterCrdRouter(crdType string) {
	crdParts := strings.Split(crdType, ".")
	groupName := strings.Join(crdParts[1:], ".")
	resourcePattern := fmt.Sprintf("/apis/%s/v1/%s", groupName, crdParts[0])
	resourceNamePattern := resourcePattern + "/:name"

	// TODO NPT-313 support authentication for kubectl proxy requests
	s.Echo.GET(resourceNamePattern, kubeGetByNameHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  crdParts[0],
			}
			return next(nc)
		}
	})
	s.Echo.GET(resourcePattern, kubeGetHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  crdParts[0],
			}
			return next(nc)
		}
	})
	s.Echo.POST(resourcePattern, kubePostHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  crdParts[0],
			}
			return next(nc)
		}
	})
	s.Echo.DELETE(resourceNamePattern, kubeDeleteHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  crdParts[0],
			}
			return next(nc)
		}
	})
}

func (s *EchoServer) RegisterDeclarativeRouter() {
	for uri, path := range declarative.Paths {
		if path.Get != nil {
			endpointContext := declarative.SetupContext(uri, http.MethodGet, path.Get)

			if endpointContext.Single {
				s.Echo.GET(endpointContext.Uri, declarative.GetHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
					return func(c echo.Context) error {
						endpointContext.Context = c
						endpointContext.Single = true
						return next(endpointContext)
					}
				})

				declarative.AddApisEndpoint(endpointContext)
				log.Debugf("Registered declarative get endpoint: %s for uri: %s", endpointContext.Uri, uri)
			} else {
				s.Echo.GET(endpointContext.Uri, declarative.ListHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
					return func(c echo.Context) error {
						endpointContext.Context = c
						endpointContext.Single = false
						return next(endpointContext)
					}
				})

				declarative.AddApisEndpoint(endpointContext)
				log.Debugf("Registered declarative list endpoint: %s for uri: %s", endpointContext.Uri, uri)
			}
		}

		if path.Put != nil {
			endpointContext := declarative.SetupContext(uri, http.MethodPut, path.Put)

			s.Echo.PUT(endpointContext.Uri, declarative.PutHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					endpointContext.Context = c
					return next(endpointContext)
				}
			})

			declarative.AddApisEndpoint(endpointContext)
			log.Debugf("Registered declarative put endpoint: %s for uri: %s", endpointContext.Uri, uri)
		}

		if path.Delete != nil {
			endpointContext := declarative.SetupContext(uri, http.MethodDelete, path.Delete)

			s.Echo.DELETE(endpointContext.Uri, declarative.DeleteHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					endpointContext.Context = c
					endpointContext.Single = true
					return next(endpointContext)
				}
			})

			declarative.AddApisEndpoint(endpointContext)
			log.Debugf("Registered declarative delete endpoint: %s for uri: %s", endpointContext.Uri, uri)
		}
	}
}

func (s *EchoServer) NodeUpdateNotifications(stopCh chan struct{}) error {
	for {
		select {
		case <-stopCh:
			return fmt.Errorf("stop signal received")
		case restURIs := <-model.RestURIChan:
			log.Debugln("Rest route notification received")
			for _, v := range restURIs {
				s.RegisterRouter(v)
			}
		case crdType := <-model.CrdTypeChan:
			log.Debugln("CRD route notification received")
			s.RegisterCrdRouter(crdType)
		case oidcNodeEvent := <-model.OidcChan:
			log.Debugln("OIDC notification received")
			err := authn.HandleOidcNodeUpdate(&oidcNodeEvent, s.Echo)
			if err != nil {
				log.Errorf("error occurred while handling OIDC node update notification: %s", err)
			}
		}
	}
}

func (s *EchoServer) StopServer() {
	if err := s.Echo.Shutdown(context.Background()); err != nil {
		log.Fatalf("Shutdown signal received")
	} else {
		log.Debugln("Server exiting")
	}
}

func NewEchoServer(conf *config.Config) *EchoServer {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	if config.Cfg.EnableNexusRuntime {
		// Setup proxy to api server
		kubeSetupProxy(e)

	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "ACCESS[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}\n",
	}))

	return &EchoServer{
		// create a new echo_server instance
		Echo:   e,
		Config: conf,
	}
}
