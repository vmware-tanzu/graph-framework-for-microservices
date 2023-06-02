package authn

import "errors"

type AuthError struct {
	ErrorMsg             string
	RedirectToAuthServer bool
}

func (ae *AuthError) Error() string {
	return ae.ErrorMsg
}

var (
	// ErrTokenNotFound no token/session found in the request
	ErrTokenNotFound = &AuthError{"authentication token/session not found", true}

	// ErrTokenSignatureInvalid indicates a token signature validation failure
	ErrTokenSignatureInvalid = &AuthError{"token signature verification failure", false}

	// ErrTokenExpiredOrNotValidYet indicates that the iat claim value represents an expired token or a token that is not valid yet
	ErrTokenExpiredOrNotValidYet = &AuthError{"token expired or not valid yet", true}

	// ErrTokenFormatInvalid the session is invalid
	ErrTokenFormatInvalid = &AuthError{"invalid authorization header value", false}

	// ErrTokenClaimsInvalid indicates that at least one of the claims present in the token was invalid
	ErrTokenClaimsInvalid = &AuthError{"one or more claims were invalid", false}

	// ErrJwksNotInitialized indicates that the JWKS key fetcher was not initialized
	ErrJwksNotInitialized = &AuthError{"jwks not initialized", false}
)

var (
	// ErrRedirectFailed indicates that a redirection failed
	ErrRedirectFailed = errors.New("redirection failed")

	// ErrIdTokenNotFound indicates that an ID token couldn't be found in the response from the Auth server during token exchange
	ErrIdTokenNotFound = errors.New("no id_token field in oauth2 token")

	// ErrCallbackEndpointCreation indicates a failure in registering the specified OIDC callback endpoint
	ErrCallbackEndpointCreation = errors.New("failed to register the callback endpoint with the Echo server")

	// ErrAuthenticatorInit indicates a failure in initializing the OIDC AuthenticatorObject
	ErrAuthenticatorInit = errors.New("failed to initialize OIDC AuthenticatorObject")
)
