package authn

import (
	"api-gw/pkg/common"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func RegisterLoginEndpoint(e *echo.Echo) {
	e.Any(common.LoginEndpoint, LoginHandler)
	log.Debugf("successfully registered login endpoint at %s", common.LoginEndpoint)
}

func LoginHandler(c echo.Context) error {
	if isOidcEnabled() {
		// TODO accept a URL to redirect to post login
		err := c.Redirect(http.StatusTemporaryRedirect, authenticator.AuthCodeURL("/"))
		if err != nil {
			return ErrRedirectFailed
		}
	} else {
		log.Debugln("OIDC not enabled, nothing to do")
		c.String(http.StatusOK, "")
	}
	return nil
}
