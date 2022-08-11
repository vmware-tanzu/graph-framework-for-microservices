package authentication

import (
	"golang-appnet.eng.vmware.com/nexus-sdk/nexus/nexus"
)

// IDPConfig contains the properties of an OIDC app
type IDPConfig struct {
	// provided by the IDP on creation of an OIDC app
	ClientId string `json:"clientId"`
	// provided by the IDP on creation of an OIDC app
	ClientSecret string `json:"clientSecret"`
	// provided by the IDP on creation of an OIDC app
	OAuthIssuerUrl string `json:"oAuthIssuerUrl"`
	// OAuth 2.0 scopes - determines the scope of the issued access tokens
	Scopes []string `json:"scopes"`
	// the URL to which the auth server must redirect to post-authentication
	OAuthRedirectUrl string `json:"oAuthRedirectUrl"`
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

// OIDC holds state/config associated with authentication.
//
// Nexus Runtime supports authentication function and the state
// associated with it is rooted on the OIDC node.
type OIDC struct {
	nexus.Node

	// IDP configuration.
	Config IDPConfig `json:"config"`

	// Properties to control the claim validation done by the gateway
	ValidationProps ValidationProperties `json:"validationProps,omitempty"`

	// JwtClaimUsername specifies the JWT claim within the JWT payload that
	// holds the username (or a unique identifier for a user)
	JwtClaimUsername string `json:"jwtClaimUsername,omitempty"`
}
