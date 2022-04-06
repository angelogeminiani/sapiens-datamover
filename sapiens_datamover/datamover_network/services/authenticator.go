package services

import "bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_commons"

type Authenticator struct {
	settings *datamover_commons.SettingsNetAuthorization
}

func NewAuthenticator(settings *datamover_commons.SettingsNetAuthorization) (instance *Authenticator) {
	instance = new(Authenticator)
	if nil != settings {

	}
	return
}

func (instance *Authenticator) Validate(token string) bool {
	if nil != instance.settings {
		return instance.settings.Value == token
	}
	return true
}
