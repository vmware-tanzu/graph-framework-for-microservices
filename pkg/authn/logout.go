package authn

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func RegisterLogoutEndpoint(e *echo.Echo) error {
	e.POST(logoutEndpoint, LogoutHandler)
	log.Debugf("successfully registered logout endpoint at %s", logoutEndpoint)
	return nil
}

func LogoutHandler(c echo.Context) error {
	accessTokenCookie := new(http.Cookie)
	accessTokenCookie.Name = accessTokenStr
	accessTokenCookie.Value = ""
	accessTokenCookie.Expires = time.Unix(0, 0)
	c.SetCookie(accessTokenCookie)

	refreshTokenCookie := new(http.Cookie)
	refreshTokenCookie.Name = refreshTokenStr
	refreshTokenCookie.Value = ""
	refreshTokenCookie.Expires = time.Unix(0, 0)
	c.SetCookie(refreshTokenCookie)

	idTokenCookie := new(http.Cookie)
	idTokenCookie.Name = idTokenStr
	idTokenCookie.Value = ""
	idTokenCookie.Expires = time.Unix(0, 0)
	c.SetCookie(idTokenCookie)

	c.String(http.StatusOK, "")

	return nil
}
