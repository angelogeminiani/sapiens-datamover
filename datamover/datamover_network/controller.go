package datamover_network

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/services"
	"fmt"
)

type DataMoverNetworkController struct {
	logger *datamover_commons.Logger
	events *gg_events.Emitter

	enabled       bool
	closed        bool
	settings      *datamover_commons.SettingsNet
	authenticator *services.Authenticator // contains default authentication data
	services      []services.NetworkService
}

func NewDataMoverNetworkController(mode string, logger *datamover_commons.Logger, events *gg_events.Emitter) (instance *DataMoverNetworkController, err error) {
	instance = new(DataMoverNetworkController)
	instance.enabled = false
	instance.closed = false
	instance.logger = logger
	instance.events = events

	err = instance.init(mode)
	if nil == err {
		err = instance.initServices()
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverNetworkController) IsEnabled() bool {
	if nil != instance {
		return instance.enabled
	}
	return false
}

func (instance *DataMoverNetworkController) Count() int {
	if nil != instance && nil != instance.services {
		return len(instance.services)
	}
	return -1
}

func (instance *DataMoverNetworkController) Start() (err error) {
	if nil != instance && instance.enabled {
		instance.closed = false
		instance.logger.Info("STARTING ENABLED SERVICES:")
		for _, s := range instance.services {
			err = s.Start()
			if nil != err {
				instance.logger.Error(fmt.Sprintf("  * Service '%s' error: %s", s.Name(), err.Error()))
				break
			} else {
				instance.logger.Info(fmt.Sprintf("  * Service '%s' started on port '%v'.", s.Name(), s.Port()))
			}
		}

		// log success or errors
		if nil == err {
			instance.logger.Info("ALL SERVICES STARTED WITH NO ERRORS")
		} else {
			instance.logger.Error("SOME SERVICES DID NOT START PROPERLY!")
		}
	}
	return
}

func (instance *DataMoverNetworkController) Stop() (err error) {
	if nil != instance && instance.enabled && !instance.closed {
		instance.closed = true
		for _, s := range instance.services {
			err = s.Stop()
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverNetworkController) init(mode string) error {
	filename := gg.Paths.WorkspacePath(fmt.Sprintf("services.%s.json", mode))
	if b, _ := gg.Paths.Exists(filename); b {

		err := gg.JSON.ReadFromFile(filename, &instance.settings)
		if nil != err {
			return err
		}
		instance.enabled = instance.settings.Enabled

		instance.authenticator = services.NewAuthenticator(instance.settings.Authorization)
	} else {
		instance.enabled = false
	}
	return nil
}

func (instance *DataMoverNetworkController) initServices() error {
	if nil != instance && instance.enabled {
		for _, service := range instance.settings.Services {
			if service.Enabled {
				s, e := services.BuildNetworkService(service.Name, service.Protocol, service.ProtocolConfiguration)
				if nil != e {
					return e
				}
				instance.services = append(instance.services, s)
				s.OnNewMessage(instance.handleMessage) // central callback error
			}
		}
	}
	return nil
}

func (instance *DataMoverNetworkController) handleMessage(networkMessage *message.NetworkMessage) interface{} {
	authData := networkMessage.GetAuthorization()
	if len(authData) > 0 && nil != instance.authenticator {
		if !instance.authenticator.Validate(authData) {
			return gg.Errors.Prefix(datamover_commons.PanicSystemError, fmt.Sprintf("Token '%s' not authorized: ", authData))
		}
	}

	var payload message.NetworkMessagePayload
	err := gg.JSON.Read(networkMessage.Body, &payload)
	if nil != err {
		return err
	}

	// ready to handle message
	return executeAction(&payload)
}

func executeAction(payload *message.NetworkMessagePayload) interface{} {
	root := payload.ActionRoot
	fnvars := gg.FnVars.NewEngine()
	context := payload.ActionContextData
	variables := payload.ActionContextVariables
	settings := payload.ActionConfig
	settings.Network = nil // remove network setting to avoid remote execution
	exec, err := action.NewDataMoverAction(root, fnvars, settings)
	if nil != err {
		return err
	}
	respData, respErr := exec.Execute(context, variables)
	if nil != respErr {
		return respErr
	}

	return respData
}
