package authn

import (
	"api-gw/pkg/common"
	"net/http"
	"regexp"
	"strings"

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
		var state, parsed_url string
		state = c.QueryParam("state")
		// Example: state: Bearer%20realm=%22http://localhost:10000/api/v1/namespaces%22
		// split the string by '=' as seperator to get the URL ("http://localhost:10000/api/v1/namespaces")
		// check if length is 2 to get the 2nd phrase( URL) or get URL directly ( beacuse user can pass state directly)
		// trim '"' in URL and get URL only to pass it to authenticator
		if state != "" {
			full_url := strings.Split(state, "=")
			parsed_url = full_url[0]
			if len(full_url) > 1 {
				parsed_url = full_url[1]
			}
			state = parsed_url
			trimmed_version := regexp.MustCompile(`^"(.*)"$`).ReplaceAllString(parsed_url, `$1`)
			if trimmed_version != "" {
				state = trimmed_version
			}
		} else {
			state = "/"
		}
		err := c.Redirect(http.StatusTemporaryRedirect, authenticator.AuthCodeURL(state))
		if err != nil {
			return ErrRedirectFailed
		}
	} else {
		log.Debugln("OIDC not enabled, nothing to do")
		c.String(http.StatusOK, "")
	}
	return nil
}
