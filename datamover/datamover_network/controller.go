package datamover_network

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_globals"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/services"
	"fmt"
)

type DataMoverNetworkController struct {
	mode        string
	logger      *datamover_commons.Logger
	events      *gg_events.Emitter
	fnVarEngine *gg_fnvars.FnVarsEngine

	enabled       bool
	closed        bool
	settings      *datamover_commons.SettingsNet
	authenticator *services.Authenticator // contains default authentication data
	services      []services.NetworkService
}

func NewDataMoverNetworkController(mode string, logger *datamover_commons.Logger, events *gg_events.Emitter, fnVarEngine *gg_fnvars.FnVarsEngine) (instance *DataMoverNetworkController, err error) {
	instance = new(DataMoverNetworkController)
	instance.mode = mode
	instance.enabled = false
	instance.closed = false
	instance.logger = logger
	instance.events = events
	instance.fnVarEngine = fnVarEngine

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
		handlers := services.NewDataMoverHandlers(instance.settings.Handlers, instance.logger, instance.fnVarEngine)
		for _, service := range instance.settings.Services {
			if service.Enabled {
				s, e := services.BuildNetworkService(service.Name, service.Protocol,
					service.ProtocolConfiguration, handlers, instance.logger)
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
	if nil == networkMessage || len(networkMessage.Body) == 0 {
		return gg.Errors.Prefix(datamover_commons.PanicSystemError, "Empty Message is not allowed!")
	}
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
	return instance.executeAction(&payload)
}

// execute remote action locally
func (instance *DataMoverNetworkController) executeAction(payload *message.NetworkMessagePayload) (response interface{}) {
	root := gg.Paths.WorkspacePath(payload.ActionRootRelative)
	_ = gg.Paths.Mkdir(root)

	context := payload.ActionContextData
	variables := payload.ActionContextVariables
	settings := payload.ActionConfig
	globals := payload.ActionGlobals
	datasets := payload.ActionDatasets

	// align datasets
	if len(datasets) > 0 {
		action.OverwriteJsDatasets(root, datasets)
	}
	// remove network setting to avoid remote execution
	settings.Network = nil

	// check local globals
	// remote globals are used only if locally there are no globals
	locGlobals := datamover_globals.NewGlobals(instance.mode)
	if nil == globals || (nil != locGlobals && locGlobals.HasConnections()) {
		// load globals from local
		globals = locGlobals
	}

	// execute the remote command
	fnvars := gg.FnVars.NewEngine()
	exec, err := action.NewDataMoverAction(root, fnvars, settings, globals)
	if nil != err {
		return err
	}
	respData, respErr := exec.Execute(context, variables)
	if nil != respErr {
		response = respErr
	} else {
		response = new(message.NetworkMessageResponseBody)
		response.(*message.NetworkMessageResponseBody).Body = respData
		response.(*message.NetworkMessageResponseBody).Variables = variables
		response.(*message.NetworkMessageResponseBody).Datasets = action.LoadJsDatasets(root)
	}

	return
}
