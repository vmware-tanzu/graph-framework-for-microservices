package authn

import (
	"api-gw/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	authnexusv1 "gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/api.git/build/apis/authentication.nexus.org/v1"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Authenticator is used to authenticate users using OIDC
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
	Jwks                    *keyfunc.JWKS
	OAuthIssuerURL          string
	SkipClientAudValidation bool
	SkipIssuerValidation    bool
	SkipClientIdValidation  bool
}

var (
	authenticator *Authenticator
	mutex         = &sync.Mutex{}
)

func isOidcEnabled() bool {
	return authenticator != nil
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
		if authenticator != nil {
			authenticator.Jwks.EndBackground()
			authenticator = nil
			log.Infoln("Disabling OIDC...")
			return nil
		} else {
			log.Debugln("no authenticator present, nothing to do")
			return nil
		}
	}

	err := validateOidcSpec(event.Oidc.Spec)
	if err != nil {
		return fmt.Errorf("OIDC Spec validation failed due to error: %s", err)
	}

	authenticator, err = newAuthenticator(event.Oidc)
	if err != nil {
		log.Errorf("Error initializing OIDC Authenticator: %s\n", err)
		return ErrAuthenticatorInit
	}

	err = RegisterCallbackHandler(e)
	if err != nil {
		log.Errorf("Could not create OIDC callback endpoint from %s: %v\n", authenticator.RedirectURL, err)
		return ErrCallbackEndpointCreation
	}

	log.Infoln("Successfully initialized OIDC Authenticator")
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

// RegisterCallbackHandler register the OAuth callback URL
func RegisterCallbackHandler(e *echo.Echo) error {
	if authenticator == nil {
		log.Debugln("authenticator is nil, nothing to do")
		return nil
	}
	callbackUrl, err := url.ParseRequestURI(authenticator.RedirectURL)
	if err != nil {
		return fmt.Errorf("Could not create callback endpoint from %s: %v\n", authenticator.RedirectURL, err)
	}
	e.Any(callbackUrl.Path, CallbackHandler)
	log.Debugf("successfully registered callback handler at %s", callbackUrl.Path)
	return nil
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

	jwksUri, err := getJwksUri(oidcNode.Spec.Config.OAuthIssuerUrl)
	if err != nil {
		log.Errorf("Error getting JWKS URI for the issuer: %s\n", oidcNode.Spec.Config.OAuthIssuerUrl)
		return nil, err
	}

	// Create the JWKS from the resource at the given URL.
	keyfuncOptions := keyfunc.Options{
		RefreshInterval:  1 * time.Hour,
		RefreshRateLimit: 1 * time.Hour,
		RefreshErrorHandler: func(err error) {
			log.Errorf("Error while refreshing JWKS: %s\n", err)
		},
		RefreshUnknownKID: true,
	}
	jwks, err := keyfunc.Get(jwksUri, keyfuncOptions)
	if err != nil {
		log.Errorf("Failed to get the JWKS from the given URL: %s\n", err)
		return nil, err
	}

	// TODO NPT-312 add a validation webhook to validate the OIDC params
	conf := oauth2.Config{
		ClientID:     oidcNode.Spec.Config.ClientId,
		ClientSecret: oidcNode.Spec.Config.ClientSecret,
		RedirectURL:  oidcNode.Spec.Config.OAuthRedirectUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       oidcNode.Spec.Config.Scopes,
	}

	return &Authenticator{
		Provider:                provider,
		Config:                  conf,
		Jwks:                    jwks,
		OAuthIssuerURL:          oidcNode.Spec.Config.OAuthIssuerUrl,
		SkipIssuerValidation:    oidcNode.Spec.ValidationProps.SkipIssuerValidation,
		SkipClientIdValidation:  oidcNode.Spec.ValidationProps.SkipClientIdValidation,
		SkipClientAudValidation: oidcNode.Spec.ValidationProps.SkipClientAudValidation,
	}, nil
}

// VerifyAndGetIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyAndGetIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken := token.Extra(idTokenStr)
	if rawIDToken == nil {
		return nil, fmt.Errorf("id_token not found")
	}

	idToken, ok := rawIDToken.(string)
	if !ok {
		return nil, ErrIdTokenNotFound
	}

	config := &oidc.Config{
		ClientID:          a.ClientID,
		SkipIssuerCheck:   authenticator.SkipIssuerValidation,
		SkipClientIDCheck: authenticator.SkipClientIdValidation,
	}
	verifier := a.Verifier(config)
	if verifier == nil {
		return nil, fmt.Errorf("Failed to create a verifier with config %v\n", config)
	}
	return verifier.Verify(ctx, idToken)
}

func VerifyAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if isOidcEnabled() {
			authErr := isAuthenticated(c)
			if authErr != nil {
				if authErr.RedirectToAuthServer {
					// save the current URI to be able to redirect the user to the same URL post auth
					state := c.Request().RequestURI

					// redirect to the authorization server
					err := c.Redirect(http.StatusTemporaryRedirect, authenticator.AuthCodeURL(state))
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
	accessToken, authErr := getTokenInRequest(c, accessTokenStr)
	if authErr != nil {
		log.Warnf("Couldn't find %s in request\n", accessTokenStr)
		return ErrTokenNotFound
	}

	// Parse the JWT and validate the signature
	token, err := jwt.Parse(accessToken, authenticator.Jwks.Keyfunc)
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
	validIss := authenticator.SkipIssuerValidation || claims.VerifyIssuer(authenticator.OAuthIssuerURL, true)
	validCid := authenticator.SkipClientIdValidation || claims["cid"] == authenticator.ClientID

	return validIss && validCid
}

// getJwksUri uses the OIDC provider's discovery endpoint to learn the JWKS URI (where the signing keys are published)
func getJwksUri(issuerURL string) (string, error) {
	wellKnown := strings.TrimSuffix(issuerURL, "/") + "/.well-known/openid-configuration"
	if err := isValidUrl(wellKnown); err != nil {
		return "", fmt.Errorf("invalid well-known URL for given issuer URL %s; err=%s", issuerURL, err)
	}

	resp, err := http.Get(wellKnown)
	if err != nil {
		log.Errorf("Failed to get JWKS URI from the discovery URL: %s\n", err)
		return "", err
	}
	defer resp.Body.Close()

	var jsonObject map[string]interface{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error while reading body of JWKS URI response: %s\n", err)
		return "", err
	}
	err = json.Unmarshal(bodyBytes, &jsonObject)
	if err != nil {
		log.Errorf("Failed to unmarshall body of JWKS URI response to JSON: %s\n", err)
		return "", err
	}
	jwksUri := jsonObject["jwks_uri"]
	if jwksUri != nil {
		uri, ok := jsonObject["jwks_uri"].(string)
		if !ok {
			return "", fmt.Errorf("jwks_uri field not found")
		}
		if uri == "" {
			return "", fmt.Errorf("jwks_uri empty. aborting")
		}
		return uri, nil
	}
	return "", fmt.Errorf("jwks_uri field not found")
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
	token := c.Request().Header.Get(authorizationHeader)
	if token == "" {
		return "", ErrTokenNotFound
	}

	items := strings.Split(token, " ")
	if len(items) != 2 {
		return "", ErrTokenFormatInvalid
	}

	if items[0] != authorizationTypeBearer {
		return "", ErrTokenNotFound
	}
	return items[1], nil
}
