package echo_server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"

	"api-gw/controllers"
	"api-gw/pkg/model"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/common-library.git/pkg/nexus"
)

type EchoServer struct {
	Echo *echo.Echo
}

func InitEcho(stopCh chan struct{}) {
	fmt.Println("Init Echo")
	e := NewEchoServer()
	e.Start(stopCh)
}

func (s *EchoServer) Start(stopCh chan struct{}) {
	// Start watching URI notification
	go func() {
		fmt.Println("RoutesNotification")
		if err := s.RoutesNotification(stopCh); err != nil {
			s.StopServer()
			InitEcho(stopCh)
		}
	}()

	// Start Server
	go func() {
		fmt.Println("Start Echo Again")
		if err := s.Echo.Start(":5000"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error %v", err)
		}
	}()
}

func (s *EchoServer) RegisterRouter(restURI nexus.RestURIs) {
	urlPattern := model.ConstructEchoPathParamURL(restURI.Uri)
	for m := range restURI.Methods {
		fmt.Printf("Registered Router Path %s Method %s\n", urlPattern, m)
		switch m {
		case http.MethodGet:
			s.Echo.GET(urlPattern, func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, World!\n")
			})
		case http.MethodPost:
			s.Echo.POST(urlPattern, func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, World!\n")
			})
		case http.MethodPut:
			s.Echo.PUT(urlPattern, func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, World!\n")
			})
		case http.MethodPatch:
			s.Echo.PATCH(urlPattern, func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, World!\n")
			})
		case http.MethodDelete:
			s.Echo.DELETE(urlPattern, func(c echo.Context) error {
				return c.String(http.StatusOK, "Hello, World!\n")
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
			fmt.Println("Route notification received...")
			for _, v := range restURIs {
				s.RegisterRouter(v)
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
	return &EchoServer{
		// create a new echo_server instance
		Echo: echo.New(),
	}
}
