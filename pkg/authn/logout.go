package authn

import (
	"api-gw/pkg/common"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func RegisterLogoutEndpoint(e *echo.Echo) {
	e.POST(common.LogoutEndpoint, LogoutHandler)
	log.Debugf("successfully registered logout endpoint at %s", common.LogoutEndpoint)
}

func LogoutHandler(c echo.Context) error {
	if isOidcEnabled() {
		accessTokenCookie := new(http.Cookie)
		accessTokenCookie.Name = common.AccessTokenStr
		accessTokenCookie.Value = ""
		accessTokenCookie.Expires = time.Unix(0, 0)
		c.SetCookie(accessTokenCookie)

		refreshTokenCookie := new(http.Cookie)
		refreshTokenCookie.Name = common.RefreshTokenStr
		refreshTokenCookie.Value = ""
		refreshTokenCookie.Expires = time.Unix(0, 0)
		c.SetCookie(refreshTokenCookie)

		idTokenCookie := new(http.Cookie)
		idTokenCookie.Name = common.IdTokenStr
		idTokenCookie.Value = ""
		idTokenCookie.Expires = time.Unix(0, 0)
		c.SetCookie(idTokenCookie)

		c.String(http.StatusOK, "")
	} else {
		log.Debugln("OIDC not enabled, nothing to do")
		c.String(http.StatusOK, "")
	}
	return nil
}
