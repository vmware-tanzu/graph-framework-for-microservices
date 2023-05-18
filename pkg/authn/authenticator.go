package authn

import (
	"api-gw/internal/tenant/registration"
	"api-gw/pkg/common"
	"api-gw/pkg/envoy"
	"api-gw/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	reg_svc "gitlab.eng.vmware.com/nsx-allspark_users/go-protos/pkg/registration-service/global"
	authnexusv1 "golang-appnet.eng.vmware.com/nexus-sdk/api/build/apis/authentication.nexus.vmware.com/v1"
	nexus_client "golang-appnet.eng.vmware.com/nexus-sdk/api/build/nexus-client"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate users using OIDC
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
	WellKnownIssuer            string
	WellKnownJwksUri           string
	Jwks                       *keyfunc.JWKS
	OAuthIssuerURL             string
	SkipClientAudValidation    bool
	SkipIssuerValidation       bool
	SkipClientIdValidation     bool
	RedirectURLRoot            string
	AccessToken                string
	RefreshToken               string
	IdToken                    string
	OAuthIssuerURLRoot         string
	CallbackEndpoint           string
	RefreshAccessTokenEndpoint string
	IsCSP                      bool
}

var (
	AuthenticatorObject *Authenticator
	mutex               = &sync.Mutex{}
)

func IsOidcEnabled() bool {
	log.Debugf("Authenticator Object:%+v", AuthenticatorObject)
	return AuthenticatorObject != nil
}

func HandlerTenantNodeUpdate(event *model.TenantNodeEvent, e *echo.Echo) error {
	if event == nil {
		log.Warnln("Nil event received")
		return nil
	}
	mutex.Lock()
	defer mutex.Unlock()

	if event.Type == model.Upsert {
		log.Infoln(fmt.Sprintf("Registering tenant %s in registration service", event.Tenant.DisplayName()))
		if err := registration.AddTenantToSystem(event.Tenant, event.RegClient); err != nil {
			log.Error(fmt.Sprintf("Registering tenant %s in registration service failed due to %v", event.Tenant.DisplayName(), err))
			return err
		}
	}
	if event.Type == model.Delete {
		tenanName, ok := common.GetTenantDisplayName(event.Tenant.Name)
		if !ok {
			return fmt.Errorf("Tenant displayname not registered")
		}
		envoy.DeleteTenantConfig(tenanName)
		log.Infoln("unregistering tenant in registration service")
		var registration_retry int = 0
		for registration_retry < common.REGISTRATION_RETRIES {
			registration_retry = registration_retry + 1
			status, ok := common.GetTenantState(tenanName)
			if !ok {
				return fmt.Errorf("Could not get Tenantstate object")
			}
			err := common.UnregisterTenant(event.RegClient, tenanName, reg_svc.TenantRequest_License(common.AvailableSkus[status.SKU]))
			if err != nil {
				if registration_retry == common.REGISTRATION_RETRIES {
					return err
				} else {
					log.Debugf("UnRegisterTenant Failed : continue to retry for %d time", registration_retry)
					time.Sleep(common.REGISTRATION_WAIT_TIME)
					continue
				}
			}
			break
		}
		common.DeleteTenantState(tenanName)
		common.DeleteTenantDisplayName(event.Tenant.Name)
		err := common.DeleteUsersForTenant(tenanName)
		if err != nil {
			return err
		}
	}
	return nil
}
func HandleOidcNodeUpdate(event *model.OidcNodeEvent, e *echo.Echo) error {
	if event == nil {
		log.Warnln("Nil event received")
		return nil
	}
	mutex.Lock()
	defer mutex.Unlock()

	// TODO NPT-306 handle deletion such that in case multiple OidcConfig objects were created
	// Currently, however, we support only 1 OIDC object being present
	if event.Type == model.Delete {
		if AuthenticatorObject != nil {
			if AuthenticatorObject.Jwks != nil {
				AuthenticatorObject.Jwks.EndBackground()
			}
			AuthenticatorObject = nil
			log.Infoln("Disabling OIDC...")
		} else {
			log.Debugln("no AuthenticatorObject present, nothing to do")
		}
		err := envoy.DeleteJwtAuthnConfig()
		if err != nil {
			return fmt.Errorf("error deleting envoy jwt authn config: %s", err)
		}
		return nil
	}

	err := validateOidcSpec(event.Oidc.Spec)
	if err != nil {
		return fmt.Errorf("OIDC Spec validation failed due to error: %s", err)
	}

	AuthenticatorObject, err = newAuthenticator(event.Oidc)
	if err != nil {
		log.Errorf("Error initializing OIDC Authenticator: %s\n", err)
		return ErrAuthenticatorInit
	}

	var callbackPath string
	callbackPath, err = RegisterCallbackHandler(e)
	if err != nil {
		log.Errorf("Could not create OIDC callback endpoint from %s: %v\n", AuthenticatorObject.RedirectURL, err)
		return ErrCallbackEndpointCreation
	}
	log.Infoln("Successfully initialized OIDC Authenticator")

	// Update Envoy state
	err = envoy.AddJwtAuthnConfig(&envoy.JwtAuthnConfig{
		Issuer:               AuthenticatorObject.WellKnownIssuer,
		IdpName:              event.Oidc.Name,
		JwksUri:              AuthenticatorObject.WellKnownJwksUri,
		CallbackEndpoint:     callbackPath,
		RefreshTokenEndpoint: AuthenticatorObject.RefreshAccessTokenEndpoint,
		JwtClaimUsername:     event.Oidc.Spec.JwtClaimUsername,
		AccessToken:          AuthenticatorObject.AccessToken,
		CSP:                  event.Oidc.Spec.Config.IsCSP,
	})
	if err != nil {
		return fmt.Errorf("error adding envoy jwt authn config: %s", err)
	}
	return nil
}

func validateOidcSpec(oidc authnexusv1.OIDCSpec) error {
	if oidc.Config.ClientId == "" {
		return fmt.Errorf("empty client ID")
	}
	if oidc.Config.ClientSecret == "" {
		return fmt.Errorf("empty client secret")
	}
	if err := isValidUrl(oidc.Config.OAuthRedirectUrl); err != nil {
		return fmt.Errorf("invalid OAuthRedirectUrl: %s", err)
	}
	if err := isValidUrl(oidc.Config.OAuthIssuerUrl); err != nil {
		return fmt.Errorf("invalid OAuthIssuerUrl: %s", err)
	}
	if len(oidc.Config.Scopes) == 0 {
		return fmt.Errorf("empty scopes")
	}
	return nil
}

func isValidUrl(input string) error {
	uri, err := url.ParseRequestURI(input)
	if err != nil {
		return err
	}
	switch uri.Scheme {
	case "http":
	case "https":
	default:
		return fmt.Errorf("invalid scheme")
	}
	return nil
}

// RegisterCallbackHandler register the OAuth callback URL, also returns the registered URI path
func RegisterCallbackHandler(e *echo.Echo) (string, error) {
	if AuthenticatorObject == nil {
		log.Debugln("AuthenticatorObject is nil, nothing to do")
		return "", nil
	}
	callbackUrl, err := url.ParseRequestURI(AuthenticatorObject.RedirectURL)
	if err != nil {
		return "", fmt.Errorf("Could not create callback endpoint from %s: %v", AuthenticatorObject.RedirectURL, err)
	}
	e.Any(callbackUrl.Path, CallbackHandler)
	log.Debugf("successfully registered callback handler at %s", callbackUrl.Path)
	return callbackUrl.Path, nil
}

// newAuthenticator instantiates the *Authenticator.
func newAuthenticator(oidcNode authnexusv1.OIDC) (*Authenticator, error) {
	log.Infoln("Initializing OIDC Authenticator...")

	var ctx = context.Background()
	if oidcNode.Spec.ValidationProps.InsecureIssuerURLContext {
		ctx = oidc.InsecureIssuerURLContext(ctx, oidcNode.Spec.Config.OAuthIssuerUrl)
	}

	provider, err := oidc.NewProvider(
		ctx,
		oidcNode.Spec.Config.OAuthIssuerUrl,
	)
	if err != nil {
		return nil, err
	}

	wellknownJson, err := getWellKnownJson(oidcNode.Spec.Config.OAuthIssuerUrl)
	if err != nil {
		log.Errorf("Error getting wellknown json for the issuer: %s\n", oidcNode.Spec.Config.OAuthIssuerUrl)
		return nil, err
	}

	jwksUri, ok := wellknownJson["jwks_uri"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to convert wellknown[jwks_uri] to string")
	}

	var jwks *keyfunc.JWKS = nil
	// Create the JWKS from the resource at the given URL.
	keyfuncOptions := keyfunc.Options{
		RefreshInterval:  1 * time.Hour,
		RefreshRateLimit: 1 * time.Hour,
		RefreshErrorHandler: func(err error) {
			log.Errorf("Error while refreshing JWKS: %s\n", err)
		},
		RefreshUnknownKID: true,
	}
	jwks, err = keyfunc.Get(jwksUri, keyfuncOptions)
	if err != nil {
		log.Errorf("Failed to get the JWKS from the given URL: %s\n", err)
		return nil, err
	}

	var issuer string
	issuer, ok = wellknownJson["issuer"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to convert wellknown[issuer] to string")
	}

	// TODO NPT-312 add a validation webhook to validate the OIDC params
	conf := oauth2.Config{
		ClientID:     oidcNode.Spec.Config.ClientId,
		ClientSecret: oidcNode.Spec.Config.ClientSecret,
		RedirectURL:  oidcNode.Spec.Config.OAuthRedirectUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       oidcNode.Spec.Config.Scopes,
	}
	oauthIssuerParsed, err := url.Parse(oidcNode.Spec.Config.OAuthIssuerUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauthIssuerURL Root")
	}
	oauthIssuerRoot := fmt.Sprintf("%s://%s", oauthIssuerParsed.Scheme, oauthIssuerParsed.Host)

	var accessToken, refreshToken, idToken, callbackEndpoint, refreshTokenEndpoint string
	if oidcNode.Spec.Config.IsCSP {
		accessToken = common.CSPAccessTokenStr
		refreshToken = common.CSPRefreshTokenStr
		idToken = common.CSPIdTokenStr
		callbackEndpoint = common.CSPCallBackEndpoint
		refreshTokenEndpoint = common.CSPRefreshAccessTokenEndpoint
	} else {
		accessToken = common.AccessTokenStr
		refreshToken = common.RefreshTokenStr
		idToken = common.IdTokenStr
		callbackEndpoint = common.CallBackEndpoint
		refreshTokenEndpoint = common.RefreshAccessTokenEndpoint
	}
	redirectUrlParsed, _ := url.Parse(conf.RedirectURL)
	redirectUrl := fmt.Sprintf("%s://%s", redirectUrlParsed.Scheme, redirectUrlParsed.Host)
	return &Authenticator{
		Provider:                   provider,
		Config:                     conf,
		WellKnownIssuer:            issuer,
		WellKnownJwksUri:           jwksUri,
		Jwks:                       jwks,
		OAuthIssuerURL:             oidcNode.Spec.Config.OAuthIssuerUrl,
		OAuthIssuerURLRoot:         oauthIssuerRoot,
		SkipIssuerValidation:       oidcNode.Spec.ValidationProps.SkipIssuerValidation,
		SkipClientIdValidation:     oidcNode.Spec.ValidationProps.SkipClientIdValidation,
		SkipClientAudValidation:    oidcNode.Spec.ValidationProps.SkipClientAudValidation,
		AccessToken:                accessToken,
		RefreshToken:               refreshToken,
		IdToken:                    idToken,
		CallbackEndpoint:           callbackEndpoint,
		RefreshAccessTokenEndpoint: refreshTokenEndpoint,
		RedirectURLRoot:            redirectUrl,
		IsCSP:                      oidcNode.Spec.Config.IsCSP,
	}, nil
}

// VerifyAndGetIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyAndGetIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken := token.Extra(common.IdTokenStr)
	if rawIDToken == nil {
		return nil, fmt.Errorf("id_token not found")
	}

	idToken, ok := rawIDToken.(string)
	if !ok {
		return nil, ErrIdTokenNotFound
	}

	config := &oidc.Config{
		ClientID:          a.ClientID,
		SkipIssuerCheck:   AuthenticatorObject.SkipIssuerValidation,
		SkipClientIDCheck: AuthenticatorObject.SkipClientIdValidation,
	}
	verifier := a.Verifier(config)
	if verifier == nil {
		return nil, fmt.Errorf("Failed to create a verifier with config %v", config)
	}
	return verifier.Verify(ctx, idToken)
}

func VerifyAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if IsOidcEnabled() {
			authErr := isAuthenticated(c)
			if authErr != nil {
				if authErr.RedirectToAuthServer {
					// save the current URI to be able to redirect the user to the same URL post auth
					state := c.Request().RequestURI

					// redirect to the authorization server
					err := c.Redirect(http.StatusTemporaryRedirect, AuthenticatorObject.AuthCodeURL(state))
					if err != nil {
						return ErrRedirectFailed
					}
					return nil
				} else {
					return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
				}
			}
		}
		return next(c)
	}
}

func isAuthenticated(c echo.Context) *AuthError {
	if !IsOidcEnabled() {
		return nil
	}

	accessToken, authErr := getTokenInRequest(c, AuthenticatorObject.AccessToken)
	if authErr != nil {
		log.Warnf("Couldn't find %s in request\n", AuthenticatorObject.AccessToken)
		return ErrTokenNotFound
	}

	if AuthenticatorObject.Jwks == nil {
		log.Errorln("jwks not initialized")
		return ErrJwksNotInitialized
	}

	// Parse the JWT and validate the signature
	token, err := jwt.Parse(accessToken, AuthenticatorObject.Jwks.Keyfunc)
	if err != nil {
		log.Errorf("error parsing token: %s\n", err)
		return ErrTokenSignatureInvalid
	}

	if err = token.Claims.Valid(); err != nil {
		log.Errorf("One or more invalid JWT claims found: %s\n", err)
		return ErrTokenExpiredOrNotValidYet
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Errorln("Failed to cast token claims to jwt.MapClaims")
		return ErrTokenFormatInvalid
	}
	if !validateClaims(mapClaims) {
		log.Errorf("Failed to validate JWT claims\n")
		return ErrTokenClaimsInvalid
	}
	return nil
}

func validateClaims(claims jwt.MapClaims) bool {
	// TODO add validation for audience ("aud") claim
	validIss := AuthenticatorObject.SkipIssuerValidation || claims.VerifyIssuer(AuthenticatorObject.OAuthIssuerURL, true)
	validCid := AuthenticatorObject.SkipClientIdValidation || claims["cid"] == AuthenticatorObject.ClientID

	return validIss && validCid
}

// getWellKnownJson uses the OIDC provider's discovery endpoint to learn fetch the IDP metadata and
// return it as an unmarshalled json object
func getWellKnownJson(issuerURL string) (map[string]interface{}, error) {
	wellKnown := strings.TrimSuffix(issuerURL, "/") + "/.well-known/openid-configuration"
	if err := isValidUrl(wellKnown); err != nil {
		return nil, fmt.Errorf("invalid well-known URL for given issuer URL %s; err=%s", issuerURL, err)
	}

	resp, err := http.Get(wellKnown)
	if err != nil {
		log.Errorf("Failed to get JWKS URI from the discovery URL: %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	var jsonObject map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error while reading body of JWKS URI response: %s\n", err)
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, &jsonObject)
	if err != nil {
		log.Errorf("Failed to unmarshall body of JWKS URI response to JSON: %s\n", err)
		return nil, err
	}
	return jsonObject, nil
}

////////// util functions ///////////////

// getTokenInRequest returns the access token from the http request
// it looks for the token in the 'Authorization' header and 'access_token' Cookie
func getTokenInRequest(c echo.Context, name string) (string, *AuthError) {
	token, err := getTokenInBearer(c)
	if err != nil {
		if err != ErrTokenNotFound {
			return "", err
		}
		if cookie, err := c.Cookie(name); err != nil {
			return "", ErrTokenNotFound
		} else {
			return cookie.Value, nil
		}
	}
	return token, nil
}

// getTokenInBearer retrieves the access token from the 'Authorization' header
func getTokenInBearer(c echo.Context) (string, *AuthError) {
	token := c.Request().Header.Get(common.AuthorizationHeader)
	if token == "" {
		return "", ErrTokenNotFound
	}

	items := strings.Split(token, " ")
	if len(items) != 2 {
		return "", ErrTokenFormatInvalid
	}

	if items[0] != common.AuthorizationTypeBearer {
		return "", ErrTokenNotFound
	}
	return items[1], nil
}

func GetIssuer(jwt *nexus_client.AuthenticationOIDC) (string, error) {
	wellKnownJson, err := getWellKnownJson(jwt.Spec.Config.OAuthIssuerUrl)
	if err != nil {
		return "", err
	}
	var issuer string
	issuer, ok := wellKnownJson["issuer"].(string)
	if !ok {
		return "", fmt.Errorf("GetIssuer: failed to convert wellknown[issuer] to string")
	}
	return issuer, nil
}

func GetJwksUri(jwt *nexus_client.AuthenticationOIDC) (string, error) {
	wellKnownJson, err := getWellKnownJson(jwt.Spec.Config.OAuthIssuerUrl)
	if err != nil {
		return "", err
	}
	var jwksUri string
	jwksUri, ok := wellKnownJson["jwks_uri"].(string)
	if !ok {
		return "", fmt.Errorf("GetJwksUri: failed to convert wellknown[jwksUri] to string")
	}
	return jwksUri, nil
}

func GetCallbackEndpoint(jwt *nexus_client.AuthenticationOIDC) (string, error) {
	callbackUrl, err := url.ParseRequestURI(jwt.Spec.Config.OAuthRedirectUrl)
	if err != nil {
		return "", fmt.Errorf("GetCallbackEndpoint: could not create callback endpoint from %s: %v", AuthenticatorObject.RedirectURL, err)
	}
	return callbackUrl.Path, nil
}
