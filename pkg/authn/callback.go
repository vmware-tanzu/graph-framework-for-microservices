package authn

import (
	"api-gw/pkg/common"
	"fmt"
	"net/http"

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

	token, err := authenticator.Exchange(c.Request().Context(), c.QueryParam(codeQueryParam))
	if err != nil {
		errMsg := fmt.Sprintf("Encountered error while exchanging code for token: %s\n", err)
		log.Error(errMsg)
		return echo.NewHTTPError(http.StatusUnauthorized, errMsg)
	}

	_, err = authenticator.VerifyAndGetIDToken(c.Request().Context(), token)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to verify ID Token due to error: %s\n", err)
		log.Error(errMsg)
		return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
	}

	// TODO NPT-307 consider creating an HTTP session and store the tokens within the session rather than
	// setting the tokens themselves into the cookie
	setCookieFromToken(c, token)
	state := c.QueryParam(stateQueryParam)
	if state == common.LoginEndpoint {
		c.Response().Header().Set(common.AccessTokenStr, token.AccessToken)
		c.Response().Header().Set(common.RefreshTokenStr, token.RefreshToken)
		rawIDToken := token.Extra(common.IdTokenStr)
		if rawIDToken == nil {
			log.Errorln("id_token not found")
			c.String(http.StatusUnauthorized, "failed to fetch id_token")
			return fmt.Errorf("failed to fetch id_token")
		} else {
			idToken, ok := rawIDToken.(string)
			if ok {
				c.Response().Header().Set(common.IdTokenStr, idToken)
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

func setCookieFromToken(c echo.Context, token *oauth2.Token) {
	accessTokenCookie := new(http.Cookie)
	accessTokenCookie.Name = common.AccessTokenStr
	accessTokenCookie.Value = token.AccessToken
	accessTokenCookie.Expires = token.Expiry
	c.SetCookie(accessTokenCookie)

	refreshTokenCookie := new(http.Cookie)
	refreshTokenCookie.Name = common.RefreshTokenStr
	refreshTokenCookie.Value = token.RefreshToken
	refreshTokenCookie.Expires = token.Expiry
	c.SetCookie(refreshTokenCookie)

	idTokenCookie := new(http.Cookie)
	idTokenCookie.Name = common.IdTokenStr
	rawIDToken := token.Extra(common.IdTokenStr)
	if rawIDToken == nil {
		log.Errorln("id_token not found")
	} else {
		idToken, ok := rawIDToken.(string)
		if ok {
			idTokenCookie.Value = idToken
			idTokenCookie.Expires = token.Expiry
			c.SetCookie(idTokenCookie)
		} else {
			log.Errorln("Failed to covert rawIDToken to string. Not setting id token cookie")
		}
	}
}
