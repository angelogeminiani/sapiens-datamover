package datamover_commons

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_sms/gg_sms_engine"
	"bitbucket.org/digi-sense/gg-core/gg_email"
)

type PostmanSettings struct {
	SMS   *PostmanSettingsSMS   `json:"sms"`
	Email *PostmanSettingsEmail `json:"email"`
}

type PostmanSettingsSMS struct {
	gg_sms_engine.SMSConfiguration
}

func (instance *PostmanSettingsSMS) String() string {
	return gg.JSON.Stringify(instance)
}

type PostmanSettingsEmail struct {
	Send *PostmanSettingsSend `json:"send"`
	Read *PostmanSettingsRead `json:"read"`
}

type PostmanSettingsSend struct {
	gg_email.SmtpSettings
	Enabled bool `json:"enabled"`
}

func (instance *PostmanSettingsSend) String() string {
	return gg.JSON.Stringify(instance)
}

type PostmanSettingsRead struct {
	Enabled bool `json:"enabled"`
	// add here mailboxer to read emails
}

func (instance *PostmanSettingsRead) String() string {
	return gg.JSON.Stringify(instance)
}
