package echo_server

import (
	"api-gw/pkg/authn"
	"api-gw/pkg/common"
	"api-gw/pkg/openapi/api"
	"api-gw/pkg/openapi/combined"
	"api-gw/pkg/openapi/declarative"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	log "github.com/sirupsen/logrus"

	"api-gw/pkg/config"
	"api-gw/pkg/model"
	"api-gw/pkg/utils"

	userv1 "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/apis/user.nexus.vmware.com/v1"
	nexus_client "github.com/vmware-tanzu/graph-framework-for-microservices/api/build/nexus-client"

	"github.com/vmware-tanzu/graph-framework-for-microservices/nexus/nexus"
)

type TenantData struct {
	TenantName string `json:"tenantName" form:"tenantName"`
	Sku        string `json:"sku,omitempty" form:"sku,omitempty"`
}

type UserLogin struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

var corsmutex = &sync.Mutex{}
var TotalHttpServerRestartCounter = 0
var HttpServerRestartFromOpenApiSpecUpdateCounter = 0

type EchoServer struct {
	Echo        *echo.Echo
	Config      *config.Config
	Client      *kubernetes.Clientset
	NexusClient *nexus_client.Clientset
	k8sProxy    *httputil.ReverseProxy
}

func InitEcho(stopCh chan struct{}, conf *config.Config, client *kubernetes.Clientset, nexusClient *nexus_client.Clientset) *EchoServer {
	log.Infoln("Init Echo")
	e := NewEchoServer(conf, client, nexusClient)

	if conf.EnableNexusRuntime {
		e.RegisterNexusRoutes()
	}

	if conf.BackendService != "" {
		declarative.Setup(declarative.OpenApiSpecFile)
		e.RegisterDeclarativeRoutes()
		e.RegisterDeclarativeRouter()
	}
	e.RegisterDebug()
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
				InitEcho(stopCh, s.Config, s.Client, s.NexusClient)
				TotalHttpServerRestartCounter++
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

func (s *EchoServer) RegisterDebug() {
	s.Echo.GET("/debug/all", DebugAllHandler)
}

func (s *EchoServer) RegisterCosmosAdminRoutes() {

	s.Echo.Any("/v0/temp/authmode", s.GetAuthmode)
	// Code removed for brevity
	// Added REST API endpoint for tenant creation and added validaion with SKU type to fail tenant creation request with invalid SKU
	// Currently creating the TenantConfig CR using CreateTenantIfNotExists method from common and this fails request if tenant already exists
	s.Echo.GET("/v0/tenants/status", s.GetTenantStatusHandler)
	s.Echo.PUT("/v0/tenants/instance", s.TenantCreateHandler)
	// Added REST API Endpoint to delete tenant using tenantid ex: /v0/tenants/instance/tenantname
	// Checking if tenant exists before deletion , if it does not exists , we will suppress failure and send 200 response
	s.Echo.DELETE("/v0/tenants/instance/:tenantid", s.DeleteTenantHander)
	//Added REST API Endpoint to create user and store hashed password as part of SPEC and add it to cache
	//Adding to Cache is for facilitating login call in Private SAAS environments
	//Added validations for empty fields (tenant, user and password)
	s.Echo.POST("/v0/users", s.CreateUserHandler)
	//Added VERSION call endpoint to get version of Admin components
	//Added Version Call struct with GRPC and Http module support , for GRPC it calls version Proto endpoint and provide JSON
	s.Echo.Any("/v0/version", s.VersionHandler)

	// Routing to UI for unhandled routes in admin gateway mode.
	uiProxyUrl, err := url.Parse(common.GlobalUISvcName)
	if err != nil {
		log.Warnf("Could not parse proxy URL: %v", err)
	}
	uiProxy := httputil.NewSingleHostReverseProxy(uiProxyUrl)
	s.Echo.Any("/*", echo.WrapHandler(uiProxy))

	s.Echo.Any("/v0/users/login", s.UserLoginHandler)
	//Added REST API Endpoint for querying user details
	s.Echo.GET("/v0/users/:userid", s.GetUserHandler)
	s.Echo.DELETE("/v0/users/:userid", s.DeleteUserHandler)

	s.Echo.Any("/v0/users/validate", s.ValidateUserHandler)
	s.Echo.GET("/", echo.WrapHandler(uiProxy))
	s.Echo.Any("/v0/user/preferences", s.GetUserPreferencesHandler)

	s.Echo.Any("/v0/cspauth/token", s.CSPTokenHander)
	s.Echo.Any(common.CSP_ORG_REDIRECT_URL, s.DiscoveryHandler)

}

func (s *EchoServer) RegisterNexusRoutes() {
	// OpenAPI route
	s.Echo.GET("/:datamodel/openapi.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, api.Schemas[c.Param("datamodel")])
	})

	s.Echo.GET("/explorer/openapi.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, combined.CombinedSpecs())
	})

	// Swagger-UI
	s.Echo.GET("/:datamodel/docs", SwaggerUI)
	if common.IsModeAdmin() {
		s.RegisterCosmosAdminRoutes()
	}

	_, err := authn.RegisterCallbackHandler(s.Echo)
	if err != nil {
		log.Errorln("Error registering the OIDC callback path")
		// should we panic?
	}
	authn.RegisterLoginEndpoint(s.Echo)

	//Endpoints to add for CSP

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
		case http.MethodPatch:
			if common.IsModeAdmin() {
				s.Echo.PATCH(urlPattern, patchHandler, nexusContext)
			} else {
				s.Echo.PATCH(urlPattern, patchHandler, authn.VerifyAuthenticationMiddleware, nexusContext)
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
				if httpCodesResponse, ok := v.Methods[http.MethodPut]; ok {
					v.Methods[http.MethodPatch] = httpCodesResponse
				}
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
		case tenantEvent := <-model.TenantEvent:
			log.Debugln("TenantConfig notification received")
			err := authn.HandlerTenantNodeUpdate(&tenantEvent, s.Echo)
			if err != nil {
				log.Errorf("error occurred while handling Tenant node update notification: %s", err)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Echo.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown signal received")
	} else {
		log.Debugln("Server exiting")
	}

	address := ":80"
	if s.Config.Server.HttpPort != "" {
		address = ":" + s.Config.Server.HttpPort
	}

	if utils.IsServerConfigValid(s.Config) && utils.IsFileExists(s.Config.Server.CertPath) && utils.IsFileExists(s.Config.Server.KeyPath) {
		address = s.Config.Server.Address
	}

	ok := false
	timeout := time.Now().Add(30 * time.Second)
	for time.Now().Before(timeout) {
		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err != nil {
			//informative log. When port is free then error will occur
			log.Debugf("StopServer: DialTimeout err: %v\n", err)
		}

		if conn == nil {
			ok = true
			break
		} else {
			conn.Close()
			time.Sleep(100 * time.Millisecond)
		}
	}
	if !ok {
		log.Fatalf("Error occured while stopping echo server. TCP port is busy")
	}
}

func NewEchoServer(conf *config.Config, client *kubernetes.Clientset, nexusClient *nexus_client.Clientset) *EchoServer {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	var k8sProxy *httputil.ReverseProxy
	if conf.EnableNexusRuntime {
		// Setup proxy to api server
		k8sProxy = kubeSetupProxy(e)
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "ACCESS[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}\n",
	}))

	return &EchoServer{
		// create a new echo_server instance
		Echo:        e,
		Config:      conf,
		Client:      client,
		NexusClient: nexusClient,
		k8sProxy:    k8sProxy,
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

func extractUnverifiedClaims(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		log.Debugf("unverified token not able to parse")
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("could not parse JWT Token")
}

type AssignedInstance struct {
	Url string `json:"url"`
}
type Services struct {
	AssignedInstance []AssignedInstance `json:"allOrgInstances"`
}
type Results struct {
	Services []Services `json:"services"`
}

func (s *EchoServer) CreateUserHandler(c echo.Context) error {
	var creds map[string]string
	err := json.NewDecoder(c.Request().Body).Decode(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid request body")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if creds["username"] == "" || creds["password"] == "" {
		c.JSON(http.StatusBadRequest, "Please provide username and password")
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide username and password")
	}
	if creds["tenantId"] == "" {
		c.JSON(http.StatusBadRequest, "Please provide tenantID")
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide tenantID")
	}
	passSha := md5.Sum([]byte(creds["password"]))
	passwordHashed := hex.EncodeToString(passSha[:])
	userObj := userv1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: creds["username"],
		},
		Spec: userv1.UserSpec{
			Username:  creds["username"],
			Mail:      creds["email"],
			FirstName: creds["firstName"],
			LastName:  creds["lastName"],
			Password:  passwordHashed,
			TenantId:  creds["tenantId"],
		},
	}
	err = common.CreateUser(s.NexusClient, creds["tenantId"], userObj)
	if err != nil {
		c.JSON(http.StatusBadGateway, map[string]interface{}{"error": fmt.Sprintf("Could not create user due to %s", err)})
		return echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("Could not create user due to %s", err))
	}

	return c.JSON(http.StatusCreated, "User created")
}

func (s *EchoServer) UserLoginHandler(c echo.Context) error {
	u := new(UserLogin)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request due to %s", err))
	}
	var token string
	userObjSpec, ok := common.GetUser(u.Username)
	if !ok {
		return echo.NewHTTPError(http.StatusBadGateway, "could not fetch user")
	}
	if userObjSpec.Username == u.Username {
		passSha := md5.Sum([]byte(u.Password))
		passwordHashed := hex.EncodeToString(passSha[:])
		if userObjSpec.Password == passwordHashed {
			credString := fmt.Sprintf("%s:%s", userObjSpec.Username, userObjSpec.Password)
			rawEncodedCred := base64.StdEncoding.EncodeToString([]byte(credString))
			token = rawEncodedCred
		}
	}

	if token == "" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": "Please check if username/password is correct"})
	}

	accessTokenCookie := common.CreateCookie("token", token, time.Time{})
	c.SetCookie(accessTokenCookie)

	// This is the expected token format
	data := map[string]interface{}{
		"id":     token,
		"userId": u.Username,
	}

	return c.JSON(http.StatusOK, data)
}

func (s *EchoServer) TenantCreateHandler(c echo.Context) error {
	tenant := new(TenantData)
	if err := c.Bind(tenant); err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("bad request due to %s", err))
	}
	SKU := tenant.Sku
	if SKU == "" {
		SKU = os.Getenv("DEFAULT_SKU")
	}
	err := common.CreateTenantIfNotExists(s.NexusClient, tenant.TenantName, SKU)
	if err != nil {
		return c.JSON(http.StatusBadGateway, fmt.Sprintf("could not create Tenant %s due to %s", tenant.TenantName, err))
	}

	respdata := map[string]string{
		"tenantID": tenant.TenantName,
	}
	return c.JSON(http.StatusCreated, respdata)
}

func (s *EchoServer) DeleteTenantHander(c echo.Context) error {
	found, err := common.CheckTenantIfExists(s.NexusClient, c.Param("tenantid"))
	if err != nil {
		return c.JSON(http.StatusBadGateway, fmt.Sprintf("could not get tenats due to %s", err))
	}
	if found {
		configObj, err := common.GetConfigNode(s.NexusClient, "default")
		if err != nil {
			return c.JSON(http.StatusBadGateway, fmt.Sprintf("could not get tenants due to %s", err))
		}
		err = configObj.DeleteTenant(context.Background(), c.Param("tenantid"))
		if err != nil {
			return c.JSON(http.StatusBadGateway, fmt.Sprintf("could not delete tenant %s due to %s", c.Param("tenantid"), err))
		}

	}
	return c.JSON(http.StatusOK, "")

}

func (s *EchoServer) VersionHandler(c echo.Context) error {
	var results []interface{}
	for _, v := range utils.VersionCalls {
		result, err := v.GetVersion()
		if err != nil {
			log.Errorf("could not get version for %s due to %s", v.Service, err)
		}
		if result != nil {
			results = append(results, result)
		}
	}
	return c.JSON(http.StatusOK, results)
}

func (s *EchoServer) DeleteUserHandler(c echo.Context) error {
	_, ok := common.GetUser(c.Param("userid"))
	if !ok {
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "User is deleted already"})
	}
	err := common.DeleteUserObject(s.NexusClient, c.Param("userid"))
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]interface{}{"error": fmt.Sprintf("could not delete used due to %s", err)})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "user deleted"})

}

func (s *EchoServer) GetUserHandler(c echo.Context) error {
	//USE credentials secret object to get the username and pass the username here
	//Flash creates user with all details at runtime
	spec, present := common.GetUser(c.Param("userid"))
	if !present {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": "user not found"})
	}
	data := map[string]interface{}{
		"username":  spec.Username,
		"firstName": spec.FirstName,
		"lastName":  spec.LastName,
		"email":     spec.Mail,
		"id":        spec.Username,
		"tenantId":  spec.TenantId,
		"name":      spec.Username,
		"realm":     "admin",
	}

	return c.JSON(http.StatusOK, data)

}

func (s *EchoServer) ValidateUserHandler(c echo.Context) error {

	// Get the value of the "token" query parameter.
	var token string
	token = c.QueryParam("token")
	if token == "" {
		// if token queryparam is empty Getting token from cookie
		tokenObj, err := c.Request().Cookie("token")
		if err != nil {
			return c.JSON(http.StatusForbidden, map[string]interface{}{"error": fmt.Sprintf("could not get token from cookie due to %s", err)})
		}
		token = tokenObj.Value
	}
	rawDecodedCred, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": fmt.Sprintf("could not decode token due to %s", err)})
	}

	username := common.GetUserNameFromToken(string(rawDecodedCred))

	// Get the user spec from cache.
	spec, present := common.GetUser(username)
	if !present {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": fmt.Sprintf("user not found due to %s", err)})
	}

	return c.JSON(http.StatusOK, spec.TenantId)

}

func (s *EchoServer) GetTenantStatusHandler(c echo.Context) error {
	var tenantId string
	tenantId = c.Request().URL.Query().Get("tenantid")
	if tenantId == "" {
		tenantId = c.Request().Header.Get("org-id")
	}
	httpStatus, resultJson := common.GetServableTenantStatus(tenantId)
	if httpStatus == 503 {
		return c.JSON(httpStatus, resultJson)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"lifecycle": map[string]string{"state": "LIVE"}})

}

func (s *EchoServer) GetUserPreferencesHandler(c echo.Context) error {
	if authn.IsOidcEnabled() {
		//verifyuserhasadminororgmemberaccess
		accessToken, err := c.Request().Cookie(authn.AuthenticatorObject.AccessToken)
		token := accessToken.Value
		if err != nil {
			return c.JSON(http.StatusForbidden, "invalid token")
		}
		if authn.AuthenticatorObject.Jwks == nil {
			log.Error("Authentication provider not configured")
			return c.JSON(http.StatusBadGateway, map[string]string{"error": "could not validate token"})
		} else {
			claimsObj, err := jwt.Parse(token, authn.AuthenticatorObject.Jwks.Keyfunc)
			if err != nil {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "invalid token"})
			}
			access := common.VerifyPermissions(token, claimsObj.Claims, common.Permissions)
			if !access {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "invalid token"})
			} else {
				_, access = authn.GetAssignedInstance(token, claimsObj.Claims.(jwt.MapClaims)["context_name"].(string))
				if access {
					httpStatus, resultJson := common.GetServableTenantStatus(claimsObj.Claims.(jwt.MapClaims)["context_name"].(string))
					if httpStatus == 503 {
						return c.JSON(httpStatus, resultJson)
					}
					return c.JSON(http.StatusOK, claimsObj.Claims.(jwt.MapClaims))
				} else {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "permission to tenant not found"})
				}
			}
		}
	} else {
		var i interface{}
		return c.JSON(http.StatusOK, i)
	}

}

func (s *EchoServer) GetAuthmode(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{
		"csp": authn.IsOidcEnabled(),
	})

}

func (s *EchoServer) DiscoveryHandler(c echo.Context) error {
	if authn.IsOidcEnabled() {
		var discoveryURL, state string

		queryParams := c.Request().URL.Query()
		state = queryParams.Get("targetUri")
		parsed, err := url.Parse(state)
		if err == nil {
			state = fmt.Sprintf("%s://%s/", parsed.Scheme, parsed.Host)
		}
		cspGwURL := authn.AuthenticatorObject.AuthCodeURL(queryParams.Get("state"))
		clientID := authn.AuthenticatorObject.ClientID
		redirectURI := fmt.Sprintf("%s/%s", state, common.CallBackEndpoint)
		orgLink := queryParams.Get("orgLink")
		if orgLink != "" {
			discoveryURL = fmt.Sprintf("%s&orgLink=%s&client_id=%s&state=%s&redirect_uri=%s", cspGwURL, orgLink, clientID, state, redirectURI)
		}
		return c.Redirect(http.StatusTemporaryRedirect, discoveryURL)
	}

	return nil

}

func (s *EchoServer) CSPTokenHander(c echo.Context) error {
	if authn.IsOidcEnabled() {
		referrer := c.Request().Referer()
		authTokenValue, _ := c.Request().Cookie(authn.AuthenticatorObject.AccessToken)
		authToken := authTokenValue.Value
		if authn.AuthenticatorObject.Jwks == nil {
			log.Error("Authentication provider not configured")
			return c.JSON(http.StatusBadGateway, "authenticator not available")
		} else {
			claimsObj, _ := jwt.Parse(authToken, authn.AuthenticatorObject.Jwks.Keyfunc)
			if referrer != "" {
				urlReferrer, err := url.ParseQuery(referrer)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, fmt.Errorf("Could not parse url in referrer"))
				}
				for key, value := range urlReferrer {
					if key == "orgLink" {
						obj := strings.Split(value[0], "/")
						if claimsObj.Claims.(jwt.MapClaims)["context_name"] == obj[len(obj)-1] {
							return c.JSON(http.StatusPermanentRedirect, fmt.Sprintf("%s?redirect_uri=%s&orgLink=%s", common.CSP_ORG_REDIRECT_URL, authn.AuthenticatorObject.RedirectURL, value))
						}
					}

				}
			}
			return c.JSON(http.StatusOK, map[string]interface{}{"claims": claimsObj.Claims.(jwt.MapClaims), "value": authToken})
		}
	}
	return nil

}

func WatchForOpenApiSpecChanges(stopCh chan struct{}, openApiSpecDir string, openApiSpecFile string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorln("NewWatcher failed: ", err)
		return
	}

	go func() {
		for {
			_, err := os.Stat(openApiSpecFile)

			if err != nil { //openApiSpec file does not exist
				er := watcher.Add(openApiSpecDir)
				if er != nil {
					log.Panicln("Unable to add watcher for ", openApiSpecDir, ": ", er.Error())
				}
				log.Debugln("Watching: ", openApiSpecDir)
				fileDoesNotExist := true
				for fileDoesNotExist {
					select {
					case event := <-watcher.Events:
						if event.Op == fsnotify.Create && event.Name == openApiSpecFile {
							log.Debugln("Restarting echo server because openApi spec file is created")
							stopCh <- struct{}{}
							HttpServerRestartFromOpenApiSpecUpdateCounter++
							fileDoesNotExist = false
							watcher.Remove(openApiSpecDir)
							break
						} else {
							log.Traceln("Received Event on dir watch: " + event.Op.String() + " on file " + event.Name)
						}
					case e := <-watcher.Errors:
						if e != nil {
							log.Errorln("Error:", e)
							fileDoesNotExist = false
							watcher.Remove(openApiSpecDir)
							break
						}
					}
				}
			} else { //openApiSpec file exists
				er := watcher.Add(openApiSpecFile)
				if er != nil {
					log.Panicln("Unable to add watcher for ", openApiSpecFile, ": ", er.Error())
				}
				log.Debugln("Watching:", openApiSpecFile)
				fileExist := true
				for fileExist {
					select {
					case event := <-watcher.Events:
						if event.Op == fsnotify.Write && event.Name == openApiSpecFile {
							log.Debugln("Restarting echo server because openApi spec file is updated")
							stopCh <- struct{}{}
							HttpServerRestartFromOpenApiSpecUpdateCounter++
						}
						if event.Op == fsnotify.Remove && event.Name == openApiSpecFile {
							log.Debugln("Restarting echo server because openApi spec file is removed")
							stopCh <- struct{}{}
							HttpServerRestartFromOpenApiSpecUpdateCounter++
							fileExist = false
							watcher.Remove(openApiSpecFile)
							break
						} else {
							log.Traceln("Received Event on file watch: " + event.Op.String() + " on file " + event.Name)
						}

					case e := <-watcher.Errors:
						if e != nil {
							log.Errorln("Error:", e)
							fileExist = false
							watcher.Remove(openApiSpecFile)
							break
						}
					}
				}
			}

		}
	}()
}
