package authentication

import (
	"net/http"

	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

type IDPConfig struct {
	ClientId         string   `json:"clientId"`
	ClientSecret     string   `json:"clientSecret"`
	OAuthIssuerUrl   string   `json:"oAuthIssuerUrl"`
	OAuthRedirectUrl string   `json:"oAuthRedirectUrl"`
	Scopes           []string `json:"scopes"`
}

type ValidationProperties struct {
	// InsecureIssuerURLContext allows discovery to work when the issuer_url reported
	// by upstream is mismatched with the discovery URL. This is meant for integration
	// with off-spec providers such as CSP, Azure.
	InsecureIssuerURLContext bool `json:"insecureIssuerURLContext"`

	// SkipIssuerValidation allows skipping verification of the issuer URL when validating
	// an ID/access token. It's useful for off-spec providers, e.g., CSP
	SkipIssuerValidation bool `json:"skipIssuerValidation"`

	// SkipClientIdValidation allows skipping verification of the client ID when validating
	// an ID/access token. It's useful for off-spec providers, e.g., CSP
	SkipClientIdValidation bool `json:"skipClientIdValidation"`

	// SkipClientAudValidation allows skipping verification of the "aud" (audience) claim when validating
	// an ID/access token. It's useful for off-spec providers, e.g., CSP
	SkipClientAudValidation bool `json:"skipClientAudValidation"`
}

var OIDCRestAPISpec = nexus.RestAPISpec{
	Uris: []nexus.RestURIs{
		{
			Uri:     "/api/{Api.apis}/config/{Config.config}/oidc/{OIDC.authentication}",
			Methods: nexus.DefaultHTTPMethodsResponses,
		},
		{
			Uri: "/oidc",
			Methods: nexus.HTTPMethodsResponses{
				http.MethodGet: nexus.DefaultHTTPGETResponses,
			},
		},
	},
}

// OIDC holds state/config associated with authentication.
//
// Nexus Runtime supports authentication function and the state
// associated with it is rooted on the OIDC node.
// nexus-rest-api-gen:OIDCRestAPISpec
type OIDC struct {
	nexus.Node

	// IDP configuration.
	Config IDPConfig

	ValidationProps ValidationProperties
}
