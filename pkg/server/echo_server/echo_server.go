package echo_server

import (
	"api-gw/pkg/openapi"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"api-gw/controllers"
	"api-gw/pkg/model"
	log "github.com/sirupsen/logrus"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type EchoServer struct {
	Echo *echo.Echo
}

func InitEcho(stopCh chan struct{}) {
	fmt.Println("Init Echo")
	openapi.New()
	e := NewEchoServer()
	e.RegisterRoutes()
	e.Start(stopCh)
}

func (s *EchoServer) Start(stopCh chan struct{}) {
	// Start watching URI notification
	go func() {
		log.Info("RoutesNotification")
		if err := s.RoutesNotification(stopCh); err != nil {
			s.StopServer()
			InitEcho(stopCh)
		}
	}()

	// Start Server
	go func() {
		log.Info("Start Echo Again")
		if err := s.Echo.Start(":5000"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error %v", err)
		}
	}()
}

type NexusContext struct {
	echo.Context
	NexusURI string
	Codes    nexus.HTTPCodesResponse
}

func (s *EchoServer) RegisterRoutes() {
	// OpenAPI route
	s.Echo.GET("/openapi.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, openapi.Schema)
	})
}

func (s *EchoServer) RegisterRouter(restURI nexus.RestURIs) {
	urlPattern := model.ConstructEchoPathParamURL(restURI.Uri)
	for method, codes := range restURI.Methods {
		log.Infof("Registered Router Path %s Method %s\n", urlPattern, method)
		switch method {
		case http.MethodGet:
			s.Echo.GET(urlPattern, getHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					nc := &NexusContext{
						Context:  c,
						NexusURI: restURI.Uri,
						Codes:    codes,
					}
					return next(nc)
				}
			})
		case http.MethodPut:
			s.Echo.PUT(urlPattern, putHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					nc := &NexusContext{
						Context:  c,
						NexusURI: restURI.Uri,
						Codes:    codes,
					}
					return next(nc)
				}
			})
		case http.MethodDelete:
			s.Echo.DELETE(urlPattern, deleteHandler, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					nc := &NexusContext{
						Context:  c,
						NexusURI: restURI.Uri,
						Codes:    codes,
					}
					return next(nc)
				}
			})
		}
	}
}

func (s *EchoServer) RoutesNotification(stopCh chan struct{}) error {
	for {
		select {
		case <-stopCh:
			return fmt.Errorf("stop signal received")
		case restURIs := <-controllers.GlobalRestURIChan:
			log.Println("Route notification received...")
			for _, v := range restURIs {
				s.RegisterRouter(v)
				openapi.AddPath(v)
			}
		}
	}
}

func (s *EchoServer) StopServer() {
	if err := s.Echo.Shutdown(context.Background()); err != nil {
		log.Fatalf("Shutdown signal received")
	} else {
		log.Println("Server exiting")
	}
}

func NewEchoServer() *EchoServer {
	e := echo.New()
	e.Use(middleware.CORS())

	return &EchoServer{
		// create a new echo_server instance
		Echo: e,
	}
}
