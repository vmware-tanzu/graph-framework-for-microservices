package common

const (
	GlobalUISvcName               = "http://allspark-ui"
	CSPPermissionName             = "CSP_PERMISSION_NAME"
	CSPServiceID                  = "CSP_SERVICE_ID"
	AuthorizationTypeBearer       = "Bearer"
	AuthorizationHeader           = "Authorization"
	AccessTokenStr                = "access_token"
	RefreshTokenStr               = "refresh_token"
	RefreshAccessTokenEndpoint    = "/refreshTokens"
	CSPAccessTokenStr             = "csp-auth-token"
	IdTokenStr                    = "id_token"
	CSPIdTokenStr                 = "csp-id-token"
	CSPRefreshTokenStr            = "csp-refresh-token"
	CSPRefreshAccessTokenEndpoint = "/v0/cspauth/refresh-token"
	CSPCallBackEndpoint           = "/v0/cspauth/callback"
	CallBackEndpoint              = "/callback"
	LogoutEndpoint                = "/logout"
	LoginEndpoint                 = "/login"
	ClusterManifestEndpoint       = "/clusters/onboarding-manifest"
)
