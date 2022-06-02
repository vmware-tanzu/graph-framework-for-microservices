package authn

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

func RegisterRefreshAccessTokenEndpoint(e *echo.Echo) error {
	if authenticator == nil {
		log.Debugln("authenticator is nil, nothing to do")
		return nil
	}
	e.POST(refreshAccessTokenEndpoint, RefreshTokenHandler)
	log.Debugf("successfully registered refresh access token endpoint at %s", refreshAccessTokenEndpoint)
	return nil
}

func RefreshTokenHandler(c echo.Context) error {
	refreshToken, authError := getTokenInRequest(c, refreshTokenStr)
	if authError != nil {
		return fmt.Errorf("error getting refresh_token cookie: %s", authError)
	}

	// use the refresh token to fetch a new set of tokens
	updatedToken, err := authenticator.Config.TokenSource(c.Request().Context(), &oauth2.Token{
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
}
