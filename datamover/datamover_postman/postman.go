package datamover_postman

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x"
	"bitbucket.org/digi-sense/gg-core/gg_email"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"errors"
	"fmt"
)

type DataMoverPostman struct {
	settings *datamover_commons.PostmanSettings
	logger   *datamover_commons.Logger
	events   *gg_events.Emitter

	closed bool
}

func NewPostman(settings *datamover_commons.PostmanSettings, logger *datamover_commons.Logger, events *gg_events.Emitter) (*DataMoverPostman, error) {
	if nil != settings {
		instance := new(DataMoverPostman)
		instance.settings = settings
		instance.logger = logger
		instance.events = events
		instance.closed = true

		return instance, nil
	}
	return nil, errors.New("missing_settings")
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverPostman) Start() (err error) {
	if nil != instance {
		instance.closed = false
	}
	return
}

func (instance *DataMoverPostman) Stop() (err error) {
	if nil != instance {
		instance.closed = true
	}
	return
}

func (instance *DataMoverPostman) SendEmail(to, subject, message string, attachments []string, callback func(err error)) error {
	if nil != instance && !instance.closed {
		if len(to) > 0 {
			sender, err := instance.mailSender()
			if nil == err && nil != sender {
				sender.SendAsync(subject, message, []string{to}, instance.settings.Email.Send.From, attachments, callback)
			}
			return err
		}
	}
	return nil
}

func (instance *DataMoverPostman) NotifyEmail(to, subject, message string, attachments []string, payload map[string]interface{}, isError bool) (response string) {
	if nil != instance && !instance.closed {
		if len(to) > 0 && nil != instance.settings && nil != instance.settings.Email && instance.settings.Email.Send.Enabled {
			defer func() {
				if r := recover(); r != nil {
					// recovered from panic
					m := gg.Strings.Format("[panic] SendEmail: '%s'", r)
					instance.logger.Error(m)
				}
			}()
			sender, err := instance.mailSender()
			if nil == err && nil != sender {
				mm, _ := gg.Formatter.Merge(message, payload)
				ms, _ := gg.Formatter.Merge(subject, payload)
				if len(mm) > 0 {
					sender.SendAsync(ms, mm, []string{to}, instance.settings.Email.Send.From, attachments, func(err error) {
						if nil != err {
							instance.logger.Error("loadzilla_postman.NotifyEmail()", err)
						} else {
							instance.logger.Info(fmt.Sprintf("Sent email to '%s' with subject '%s'", to, ms))
						}
					})
					response = mm
				} else {
					sender.SendAsync(subject, message, []string{to}, instance.settings.Email.Send.From, attachments, func(err error) {
						if nil != err {
							instance.logger.Error("loadzilla_postman.NotifyEmail()", err)
						} else {
							instance.logger.Info(fmt.Sprintf("Sent email to '%s' with subject '%s'", to, subject))
						}
					})
					response = message
				}
			}
		}

	}
	return
}

func (instance *DataMoverPostman) NotifySMS(to, message string, payload map[string]interface{}, isError bool) (response string) {
	if len(to) > 0 && nil != instance.settings && nil != instance.settings.SMS && instance.settings.SMS.Enabled {
		defer func() {
			if r := recover(); r != nil {
				// recovered from panic
				m := gg.Strings.Format("[panic] SendSMS: '%s'", r)
				instance.logger.Error(m)
			}
		}()
		sender := ggx.SMS.NewEngine(instance.settings.SMS.String())
		if nil != sender {
			from := instance.settings.SMS.Provider().Param("from")
			mm, _ := gg.Formatter.Merge(message, payload)
			var err error
			if len(mm) > 0 {
				_, err = sender.SendMessage("default", mm, to, from)
				response = mm
			} else {
				_, err = sender.SendMessage("default", message, to, from)
				response = message
			}
			if nil != err {
				instance.logger.Error(fmt.Sprintf("SendSMS error: '%s'", err))
			}
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverPostman) mailSender() (*gg_email.SmtpSender, error) {
	if nil != instance {
		return gg.Email.NewSender(instance.settings.Email.Send.String())
	}
	return nil, datamover_commons.PanicSystemError
}
