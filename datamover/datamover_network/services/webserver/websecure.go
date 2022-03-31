package webserver

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_auth0"
	"github.com/gofiber/fiber/v2"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e
// ---------------------------------------------------------------------------------------------------------------------

type Websecure struct {
	apiAuthorization *Authorization
	enabled          bool
	settings         *gg_auth0.Auth0Config
	auth0            *gg_auth0.Auth0
}

func NewWebsecure(mode string, apiAuthorization *Authorization) *Websecure {
	instance := new(Websecure)
	instance.apiAuthorization = apiAuthorization
	instance.init(mode)

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Websecure) IsEnabled() bool {
	if nil != instance {
		return instance.enabled
	}
	return false
}

func (instance *Websecure) Start() error {
	if nil != instance && instance.enabled {
		return instance.auth0.Open()
	}
	return nil
}

func (instance *Websecure) Stop() {
	if nil != instance && instance.enabled {
		_ = instance.auth0.Close()
	}
}

func (instance *Websecure) AuthenticateRequest(ctx *fiber.Ctx, assert bool) bool {
	if nil != ctx && instance.IsEnabled() {
		auth0 := instance.auth0
		authSettings := instance.apiAuthorization
		if nil != authSettings && len(authSettings.Type) > 0 {
			// get check token
			requiredAuthToken, err := getAuthToken(authSettings)
			if nil != err {
				_ = WriteResponse(ctx, nil, gg.Errors.Prefix(HttpUnauthorizedError, err.Error()))
				return false
			}
			requiredAuthMode := authSettings.Type
			if len(requiredAuthMode) > 0 && requiredAuthMode != "none" {
				// get access token
				token := getAuthenticationToken(ctx)
				if len(token) == 0 {
					token = getApplicationToken(ctx)
				}
				if len(token) == 0 && assert {
					_ = WriteResponse(ctx, nil, gg.Errors.Prefix(HttpUnauthorizedError, "Missing Authorization Token:"))
					return false
				}
				if len(requiredAuthToken) > 0 {
					// direct check
					if token != requiredAuthToken {
						_ = WriteResponse(ctx, nil, HttpUnauthorizedError)
						return false
					}
				} else {
					// access token validation
					if nil != auth0 {
						if b, _ := auth0.TokenValidate(token); !b {
							// token expired or invalid
							_ = WriteResponse(ctx, nil, AccessTokenExpiredError)
							return false
						}
					}
				}
			}
		}
	}
	return true
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Websecure) init(mode string) {
	settings, err := loadSecureSettings(mode)
	if nil == err {
		instance.enabled = true
		instance.settings = settings
		instance.auth0 = gg_auth0.NewAuth0(settings)
	} else {
		instance.enabled = false
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//		S T A T I C
// ---------------------------------------------------------------------------------------------------------------------

func loadSecureSettings(mode string) (*gg_auth0.Auth0Config, error) {
	path := gg.Paths.WorkspacePath("websecure." + mode + ".json")
	settings := new(gg_auth0.Auth0Config)
	text, err := gg.IO.ReadTextFromFile(path)
	if nil != err {
		return settings, err
	}
	err = gg.JSON.Read(text, &settings)
	return settings, err
}
