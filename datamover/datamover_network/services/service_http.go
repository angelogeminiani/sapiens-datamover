package services

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/services/webserver"
)

type ServiceHttp struct {
	name     string
	uid      string
	enabled  bool
	closed   bool
	settings *webserver.WebserverSettings

	callback func(netMessage *message.NetworkMessage) interface{}
	server   *webserver.WebController
}

func NewServiceHttp(name string, configuration map[string]interface{}) (instance *ServiceHttp, err error) {
	instance = new(ServiceHttp)
	instance.name = name

	err = instance.init(configuration)
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ServiceHttp) Name() string {
	if nil != instance {
		return instance.name
	}
	return ""
}

func (instance *ServiceHttp) Port() int {
	if nil != instance && nil != instance.server {
		return 0
	}
	return -1
}

func (instance *ServiceHttp) Uid() string {
	if nil != instance {
		return instance.uid
	}
	return ""
}

func (instance *ServiceHttp) Start() (err error) {
	if nil != instance && instance.enabled {
		instance.closed = false
		if nil != instance.server {
			instance.server.Start()
		}
	}
	return
}

func (instance *ServiceHttp) Stop() (err error) {
	if nil != instance && instance.enabled && !instance.closed {
		instance.closed = true
		if nil != instance.server {
			instance.server.Stop()
		}
	}
	return
}

func (instance *ServiceHttp) OnNewMessage(callback func(netMessage *message.NetworkMessage) interface{}) {
	instance.callback = callback
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ServiceHttp) init(configuration map[string]interface{}) (err error) {
	instance.enabled = true
	instance.uid = gg.Rnd.Uuid()
	instance.settings = new(webserver.WebserverSettings)
	instance.settings.Enabled = true
	instance.settings.Auth = nil
	instance.settings.Http = configuration

	// instance.server = webserver.NewWebController()

	return
}
