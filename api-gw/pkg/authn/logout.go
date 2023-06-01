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
	if IsOidcEnabled() {

		c.SetCookie(common.CreateCookie(AuthenticatorObject.AccessToken, "", time.Unix(0, 0)))
		c.SetCookie(common.CreateCookie(AuthenticatorObject.RefreshToken, "", time.Unix(0, 0)))
		c.SetCookie(common.CreateCookie(AuthenticatorObject.IdToken, "", time.Unix(0, 0)))

		c.String(http.StatusOK, "")
	} else {
		log.Debugln("OIDC not enabled, nothing to do")
		c.String(http.StatusOK, "")
	}
	return nil
}
