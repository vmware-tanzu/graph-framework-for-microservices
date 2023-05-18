package authn

import (
	csptenant "api-gw/internal/tenant/csp"
	"api-gw/pkg/client"
	"api-gw/pkg/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// CallbackHandler is the handler for the OAuth callback. We expect to receive an authorization "code"
// and "state" as query params. We exchange the auth code for id/access tokens and set id/access/refresh token cookies
func CallbackHandler(c echo.Context) error {
	const errorQueryParam = "error"
	const codeQueryParam = "code"
	const stateQueryParam = "state"

	log.Debugln("In callback handler...")

	if c.QueryParam(errorQueryParam) != "" {
		errMsg := fmt.Sprintf("authorization server returned an error: %s", c.QueryParam(errorQueryParam))
		log.Errorf(errMsg)
		return echo.NewHTTPError(http.StatusUnauthorized, errMsg)
	}

	// Make sure the 'code' was provided
	if c.QueryParam(codeQueryParam) == "" {
		errMsg := "Received empty authorization code"
		log.Error(errMsg)
		return echo.NewHTTPError(http.StatusUnauthorized, errMsg)
	}

	token, err := AuthenticatorObject.Exchange(c.Request().Context(), c.QueryParam(codeQueryParam))
	if err != nil {
		errMsg := fmt.Sprintf("Encountered error while exchanging code for token: %s\n", err)
		log.Error(errMsg)
		return echo.NewHTTPError(http.StatusUnauthorized, errMsg)
	}

	_, err = AuthenticatorObject.VerifyAndGetIDToken(c.Request().Context(), token)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to verify ID Token due to error: %s\n", err)
		log.Error(errMsg)
		return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
	}

	if AuthenticatorObject.IsCSP {
		claimsObj, _ := jwt.Parse(token.AccessToken, AuthenticatorObject.Jwks.Keyfunc)

		orgId := claimsObj.Claims.(jwt.MapClaims)["context_name"]
		access := common.VerifyPermissions(token.AccessToken, claimsObj.Claims, common.Permissions)
		if access {
			_, hasAccess := GetAssignedInstance(token.AccessToken, orgId.(string))
			if hasAccess {
				// Check if tenant Exists
				if found, err := common.CheckTenantIfExists(client.NexusClient, orgId.(string)); !found {
					if err != nil {
						c.String(http.StatusBadGateway, fmt.Sprintf("could not check if tenant exists %s due to err: %s", orgId.(string), err))
					}
					var sku, productId string
					//Getting serviceownerToken everytime before fallingback to empty string
					serviceOwnerToken := common.GetCSPServiceOwnerToken()
					if serviceOwnerToken == "" {
						sku = os.Getenv("DEFAULT_SKU")
					} else {
						tenant := csptenant.InitCSPTenant(
							serviceOwnerToken,
							common.CSP_SERVICE_ID,
							AuthenticatorObject.OAuthIssuerURLRoot,
						)
						productId, err = tenant.ProductID(orgId.(string))
						if err != nil {
							productId = os.Getenv("DEFAULT_SKU")
						} else {
							license := common.ConvertProductIDtoLicense(productId)
							if license == "" {
								sku = os.Getenv("DEFAULT_SKU")
							}
							sku = license
						}
					}
					common.AddTenantState(orgId.(string), common.TenantState{
						Status:        common.CREATING,
						Message:       "Tenant  creation in progress",
						CreationStart: time.Now().Format(time.RFC3339Nano),
						SKU:           sku,
					})
					fmt.Println(fmt.Sprintf("Trying to create SKU using %s for tenantId: %s", sku, orgId.(string)))
					if err := common.CreateTenantIfNotExists(client.NexusClient, orgId.(string), sku); err != nil {
						c.String(http.StatusBadGateway, "Could not create tenant")
					}
				}

			}
		}
	}

	// TODO NPT-307 consider creating an HTTP session and store the tokens within the session rather than
	// setting the tokens themselves into the cookie
	setCookieFromToken(c, token)
	state := c.QueryParam(stateQueryParam)
	if len(strings.Split(state, "?")) > 1 {
		state = strings.Split(state, "?")[0]
		state = strings.Trim(state, "\"")
		state = fmt.Sprintf("%shome", state)
	}
	if state == common.LoginEndpoint {
		c.Response().Header().Set(AuthenticatorObject.AccessToken, token.AccessToken)
		c.Response().Header().Set(AuthenticatorObject.RefreshToken, token.RefreshToken)
		rawIDToken := token.Extra(common.IdTokenStr)
		if rawIDToken == nil {
			log.Errorln("id_token not found")
			c.String(http.StatusUnauthorized, "failed to fetch id_token")
			return fmt.Errorf("failed to fetch id_token")
		} else {
			idToken, ok := rawIDToken.(string)
			if ok {
				c.Response().Header().Set(AuthenticatorObject.IdToken, idToken)
			} else {
				c.String(http.StatusUnauthorized, "invalid id_token")
				return fmt.Errorf("invalid id_token")
			}
		}
		c.String(http.StatusOK, "Login successful")
		return nil
	}

	c.Redirect(http.StatusTemporaryRedirect, state)
	return nil
}

// add csp-auth-token as another cookie
func setCookieFromToken(c echo.Context, token *oauth2.Token) {
	accessTokenCookie := common.CreateCookie(AuthenticatorObject.AccessToken, token.AccessToken, token.Expiry)
	c.SetCookie(accessTokenCookie)

	refreshTokenCookie := common.CreateCookie(AuthenticatorObject.RefreshToken, token.RefreshToken, token.Expiry)
	c.SetCookie(refreshTokenCookie)

	rawIDToken := token.Extra(common.IdTokenStr)
	if rawIDToken == nil {
		log.Errorln("id_token not found")
	} else {
		idToken, ok := rawIDToken.(string)
		if ok {
			idTokenCookie := common.CreateCookie(AuthenticatorObject.IdToken, idToken, token.Expiry)
			c.SetCookie(idTokenCookie)
		} else {
			log.Errorln("Failed to covert rawIDToken to string. Not setting id token cookie")
		}
	}
}

func GetAssignedInstance(token string, tenantId string) (url string, hasAccess bool) {
	queryParams := fmt.Sprintf("includeSubOrgServices=false&serviceDefinitionId=%s", common.CSP_SERVICE_ID)
	getInstanceURL := fmt.Sprintf("%s?%s", common.GenerateServiceDefinitionURL(AuthenticatorObject.OAuthIssuerURLRoot, tenantId), queryParams)
	log.Debugf("Reaching %s to get the assignedInstance config", getInstanceURL)
	req, err := http.NewRequest("GET", getInstanceURL, http.NoBody)
	if err != nil {
		log.Errorf("could not make http GET request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("csp-auth-token", token)
	client := http.Client{}
	res, err := client.Do(req)
	if res.StatusCode != 200 {
		log.Errorf("Could not get assignedInstance config for tenant %s statusCode: %s", tenantId, res.Status)
		return "", false
	} else if err != nil {
		log.Errorf("Could not get assignedInstance config for tenant %s due to %v", tenantId, err)
		return "", false
	} else {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Errorf("Error while reading body of ServiceDefinition URI response: %s\n", err)
			return "", false
		}
		var data map[string]interface{}
		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			log.Errorf("Error while reading body of ServiceDefinition URI response: %s\n", err)
			return "", false
		}

		if data["results"] == nil {
			log.Debugf("Proceeding as the User for tenant %s has permission satisfied and found legacy service instance..", tenantId)
			return "", true
		}

		results := data["results"].([]interface{})
		if len(results) == 0 {
			log.Debugf("Proceeding as the User for tenant %s has permission satisfied and no additional instances found..", tenantId)
			return "", true
		}
		var assignedInstance string
		services, ok := results[0].(map[string]interface{})["services"]
		if ok {
			orgI, ok := services.([]interface{})[0].(map[string]interface{})["allOrgInstances"]
			if ok {
				assignedInstance, ok = orgI.([]interface{})[0].(map[string]interface{})["url"].(string)
				if ok {
					if assignedInstance != AuthenticatorObject.RedirectURLRoot {
						log.Debugf("Redirecting User for tenant %s to the correct instance %s", tenantId, assignedInstance)
						return assignedInstance, false
					}
				}
			}

		}

	}
	return "", true
}
