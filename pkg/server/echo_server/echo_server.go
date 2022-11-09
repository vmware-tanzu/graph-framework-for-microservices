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
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	log "github.com/sirupsen/logrus"

	"api-gw/pkg/config"
	"api-gw/pkg/model"
	"api-gw/pkg/utils"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

var corsmutex = &sync.Mutex{}

type EchoServer struct {
	Echo   *echo.Echo
	Config *config.Config
}

func InitEcho(stopCh chan struct{}, conf *config.Config) *EchoServer {
	log.Infoln("Init Echo")
	e := NewEchoServer(conf)

	if conf.EnableNexusRuntime {
		e.RegisterNexusRoutes()
	}

	if conf.BackendService != "" {
		e.RegisterDeclarativeRoutes()
		e.RegisterDeclarativeRouter()
	}

	common.Mode = os.Getenv("GATEWAY_MODE")
	log.Infof("Gateway Mode: %s", common.Mode)

	e.Start(stopCh)

	return e
}

func (s *EchoServer) StartHTTPServer() {
	port := "80"
	if s.Config.Server.HttpPort != "" {
		port = s.Config.Server.HttpPort
	}

	if err := s.Echo.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
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
	s.Echo.GET("/:datamodel/openapi.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api.Schemas[c.Param("datamodel")])
	})

	// Swagger-UI
	s.Echo.GET("/:datamodel/docs", SwaggerUI)

	_, err := authn.RegisterCallbackHandler(s.Echo)
	if err != nil {
		log.Errorln("Error registering the OIDC callback path")
		// should we panic?
	}
	authn.RegisterLoginEndpoint(s.Echo)
	authn.RegisterRefreshAccessTokenEndpoint(s.Echo)
	authn.RegisterLogoutEndpoint(s.Echo)
	SetUpCors(middleware.DefaultCORSConfig.AllowHeaders, s.Echo)
}

func (s *EchoServer) RegisterDeclarativeRoutes() {
	s.Echo.GET("/declarative/apis", declarative.ApisHandler)
}

func (s *EchoServer) RegisterRouter(restURI nexus.RestURIs) {
	urlPattern := model.ConstructEchoPathParamURL(restURI.Uri)
	for method, codes := range restURI.Methods {
		log.Infof("Registered Router Path %s Method %s\n", urlPattern, method)

		nexusContext := s.GetNexusContext(restURI, codes)
		switch method {
		// in "admin" mode, the responsibility of authentication is offloaded to the nexus-proxy.
		// so we don't need to add the authn.VerifyAuthenticationMiddleware middleware
		case "LIST":
			if common.IsModeAdmin() {
				s.Echo.GET(urlPattern, listHandler, nexusContext)
			} else {
				s.Echo.GET(urlPattern, listHandler, authn.VerifyAuthenticationMiddleware, nexusContext)
			}
		case http.MethodGet:
			if common.IsModeAdmin() {
				s.Echo.GET(urlPattern, getHandler, nexusContext)
			} else {
				s.Echo.GET(urlPattern, getHandler, authn.VerifyAuthenticationMiddleware, nexusContext)
			}
		case http.MethodPut:
			if common.IsModeAdmin() {
				s.Echo.PUT(urlPattern, putHandler, nexusContext)
			} else {
				s.Echo.PUT(urlPattern, putHandler, authn.VerifyAuthenticationMiddleware, nexusContext)
			}
		case http.MethodDelete:
			if common.IsModeAdmin() {
				s.Echo.DELETE(urlPattern, deleteHandler, nexusContext)
			} else {
				s.Echo.DELETE(urlPattern, deleteHandler, authn.VerifyAuthenticationMiddleware, nexusContext)
			}
		}
	}
}

func (s *EchoServer) GetNexusContext(restURI nexus.RestURIs, codes nexus.HTTPCodesResponse) func(next echo.HandlerFunc) echo.HandlerFunc {
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

func (s *EchoServer) GetNexusCrdContext(crdType, groupName, resource string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			nc := &NexusContext{
				Context:   c,
				CrdType:   crdType,
				GroupName: groupName,
				Resource:  resource,
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
	crdContext := s.GetNexusCrdContext(crdType, groupName, crdParts[0])

	// TODO NPT-313 support authentication for kubectl proxy requests
	s.Echo.GET(resourceNamePattern, KubeGetByNameHandler, crdContext)
	s.Echo.GET(resourcePattern, KubeGetHandler, crdContext)
	s.Echo.POST(resourcePattern, KubePostHandler, crdContext)
	s.Echo.DELETE(resourceNamePattern, KubeDeleteHandler, crdContext)
}

func (s *EchoServer) RegisterDeclarativeRouter() {
	for uri, path := range declarative.Paths {
		if path.Get != nil {
			endpointContext := declarative.SetupContext(uri, http.MethodGet, path.Get)

			if endpointContext.Single {
				s.Echo.GET(endpointContext.Uri, declarative.GetHandler, declarative.Middleware(endpointContext, true))
				if endpointContext.ShortUri != "" {
					s.Echo.GET(endpointContext.ShortUri, declarative.GetHandler, declarative.Middleware(endpointContext, true))
					log.Debugf("Registered declarative short get endpoint: %s for uri: %s", endpointContext.ShortUri, uri)
				}

				declarative.AddApisEndpoint(endpointContext)
				log.Debugf("Registered declarative get endpoint: %s for uri: %s", endpointContext.Uri, uri)
			} else {
				s.Echo.GET(endpointContext.Uri, declarative.ListHandler, declarative.Middleware(endpointContext, false))
				if endpointContext.ShortUri != "" {
					s.Echo.GET(endpointContext.ShortUri, declarative.ListHandler, declarative.Middleware(endpointContext, false))
					log.Debugf("Registered declarative short list endpoint: %s for uri: %s", endpointContext.ShortUri, uri)
				}

				declarative.AddApisEndpoint(endpointContext)
				log.Debugf("Registered declarative list endpoint: %s for uri: %s", endpointContext.Uri, uri)
			}
		}

		if path.Put != nil {
			endpointContext := declarative.SetupContext(uri, http.MethodPut, path.Put)
			s.Echo.PUT(endpointContext.Uri, declarative.PutHandler, declarative.Middleware(endpointContext, false))
			if endpointContext.ShortUri != "" {
				s.Echo.PUT(endpointContext.ShortUri, declarative.PutHandler, declarative.Middleware(endpointContext, false))
				log.Debugf("Registered declarative short put endpoint: %s for uri: %s", endpointContext.ShortUri, uri)
			}

			declarative.AddApisEndpoint(endpointContext)
			log.Debugf("Registered declarative put endpoint: %s for uri: %s", endpointContext.Uri, uri)
		}

		if path.Delete != nil {
			endpointContext := declarative.SetupContext(uri, http.MethodDelete, path.Delete)
			s.Echo.DELETE(endpointContext.Uri, declarative.DeleteHandler, declarative.Middleware(endpointContext, true))
			if endpointContext.ShortUri != "" {
				s.Echo.DELETE(endpointContext.ShortUri, declarative.DeleteHandler, declarative.Middleware(endpointContext, true))
				log.Debugf("Registered declarative short delete endpoint: %s for uri: %s", endpointContext.ShortUri, uri)
			}

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
		case CorsNodeEvent := <-model.CorsChan:
			log.Debug("Cors Event received")
			err := HandleCorsNodeUpdate(&CorsNodeEvent, s.Echo)
			if err != nil {
				log.Errorf("error occured while handling CORS node update notification: %s", err)
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

	if conf.EnableNexusRuntime {
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

func CheckCorsOrigin(origin string) (bool, error) {
	if len(model.CorsConfigOrigins) == 0 {
		return false, nil
	}
	for _, domains := range model.CorsConfigOrigins {
		for _, domain := range domains {
			if origin == domain {
				return true, nil
			}
		}
	}
	return false, nil
}

func HandleCorsNodeUpdate(event *model.CorsNodeEvent, e *echo.Echo) error {
	if event == nil {
		log.Warnln("Nil event received")
		return fmt.Errorf("nil type event received")
	}
	corsmutex.Lock()
	defer corsmutex.Unlock()

	if event.Type == model.Delete {
		// delete predicate is already called to remove the object
	} else {
		model.CorsConfigOrigins[event.Cors.Name] = event.Cors.Spec.Origins
		if len(event.Cors.Spec.Headers) != 0 {
			model.CorsConfigHeaders[event.Cors.Name] = event.Cors.Spec.Headers
		}
	}

	var headers []string
	for _, headerArr := range model.CorsConfigHeaders {
		for _, header := range headerArr {
			headers = append(headers, header)
		}
	}
	SetUpCors(headers, e)

	// Add cors on echo server
	return nil
}

func SetUpCors(headers []string, e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowHeaders:    headers,
			AllowMethods:    middleware.DefaultCORSConfig.AllowMethods,
			AllowOriginFunc: CheckCorsOrigin,
		},
	))
}
