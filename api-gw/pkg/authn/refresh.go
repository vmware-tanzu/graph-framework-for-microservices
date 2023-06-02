package authn

import (
	"api-gw/pkg/common"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func RegisterRefreshAccessTokenEndpoint(e *echo.Echo) {
	if IsOidcEnabled() {
		e.Any(AuthenticatorObject.RefreshAccessTokenEndpoint, RefreshTokenHandler)
	} else {
		e.Any(common.RefreshAccessTokenEndpoint, RefreshTokenHandler)
	}
	log.Debugf("successfully registered refresh access token endpoint at %s", common.RefreshAccessTokenEndpoint)
}

func RefreshTokenHandler(c echo.Context) error {
	if IsOidcEnabled() {
		refreshToken, authError := getTokenInRequest(c, AuthenticatorObject.RefreshToken)
		if authError != nil {
			return fmt.Errorf("error getting refresh_token cookie: %s", authError)
		}

		// use the refresh token to fetch a new set of tokens
		updatedToken, err := AuthenticatorObject.Config.TokenSource(c.Request().Context(), &oauth2.Token{
			RefreshToken: refreshToken,
		}).Token()
		if err != nil {
			errMsg := fmt.Sprintf("Encountered error while refreshing tokens: %s\n", err)
			log.Error(errMsg)
			return echo.NewHTTPError(http.StatusInternalServerError, errMsg)
		}
		setCookieFromToken(c, updatedToken)
		log.Debugln("Successfully refreshed tokens and updated cookies")
		return nil
	} else {
		log.Debugln("OIDC not enabled, nothing to do")
		c.String(http.StatusOK, "")
		return nil
	}
}
