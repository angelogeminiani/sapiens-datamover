package webserver

import (
	"errors"
)

type Authorization struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type WebserverSettings struct {
	Enabled bool                   `json:"enabled"`
	Http    map[string]interface{} `json:"http"`
	Auth    *Authorization         `json:"auth"`
}

var (
	HttpUnauthorizedError       = errors.New("unauthorized")          // 401
	HttpInvalidCredentialsError = errors.New("invalid_credentials")   // 401
	AccessTokenExpiredError     = errors.New("access_token_expired")  // 403
	RefreshTokenExpiredError    = errors.New("refresh_token_expired") // 401
	AccessTokenInvalidError     = errors.New("access_token_invalid")  // 401
	HttpUnsupportedApiError     = errors.New("unsupported_api")
)

// ---------------------------------------------------------------------------------------------------------------------
//	a p i
// ---------------------------------------------------------------------------------------------------------------------

// API v1
var (
	ApiSysVersion = "/api/v1/sys_version"

	// USERS_AUTH_WALLET    = "/api/v1/" + strings.Replace(ggxcommons.CmdUsersAuthWallet, "_", "/", 1)

)
