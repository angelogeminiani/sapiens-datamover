package services

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_nio"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
)

type ServiceNio struct {
	name          string
	handlers      *DataMoverHandlers
	logger        *datamover_commons.Logger
	uid           string
	enabled       bool
	closed        bool
	configuration *datamover_commons.SettingsNetProtocolNio

	callback func(netMessage *message.NetworkMessage) interface{}
	server   *gg_nio.NioServer
}

func NewServiceNio(name string, configuration map[string]interface{}, handlers *DataMoverHandlers, logger *datamover_commons.Logger) (instance *ServiceNio, err error) {
	instance = new(ServiceNio)
	instance.name = name
	instance.handlers = handlers
	instance.logger = logger

	err = instance.init(configuration)
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ServiceNio) Name() string {
	if nil != instance {
		return instance.name
	}
	return ""
}

func (instance *ServiceNio) Port() int {
	if nil != instance && nil != instance.configuration {
		return instance.configuration.Port
	}
	return -1
}

func (instance *ServiceNio) Uid() string {
	if nil != instance {
		return instance.uid
	}
	return ""
}

func (instance *ServiceNio) Start() (err error) {
	if nil != instance && instance.enabled {
		instance.closed = false
		if nil != instance.server {
			err = instance.server.Open()
		}
	}
	return
}

func (instance *ServiceNio) Stop() (err error) {
	if nil != instance && instance.enabled && !instance.closed {
		instance.closed = true
		if nil != instance.server {
			err = instance.server.Close()
		}
	}
	return
}

func (instance *ServiceNio) OnNewMessage(callback func(netMessage *message.NetworkMessage) interface{}) {
	instance.callback = callback
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ServiceNio) init(configuration map[string]interface{}) (err error) {
	err = gg.JSON.Read(gg.JSON.Stringify(configuration), &instance.configuration)
	if nil == err {
		instance.enabled = true
		instance.uid = gg.Rnd.Uuid()

		instance.server = gg.NIO.NewServer(instance.configuration.Port)
		instance.server.OnMessage(instance.onMessage)
	}
	return
}

func (instance *ServiceNio) onMessage(nioMessage *gg_nio.NioMessage) interface{} {
	body := gg.Convert.ToString(nioMessage.Body)
	if len(body) > 0 {
		// parse body into message
		var msg message.NetworkMessage
		jsonErr := gg.JSON.Read(body, &msg)
		if nil != jsonErr {
			return jsonErr
		}
		// handle on central controller
		if nil != instance.callback {
			return instance.callback(&msg)
		}
	}
	return gg.Errors.Prefix(datamover_commons.PanicSystemError, "NIO - UNHANDLED MESSAGE: ")
}
