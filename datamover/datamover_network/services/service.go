package services

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
	"fmt"
)

type NetworkService interface {
	Start() error
	Stop() error
	Name() string
	Port() int
	Uid() string
	OnNewMessage(callback func(netMessage *message.NetworkMessage) interface{})
}

func BuildNetworkService(name, protocol string, settings map[string]interface{}, handlers *DataMoverHandlers, logger *datamover_commons.Logger) (NetworkService, error) {
	switch protocol {
	case "nio":
		return NewServiceNio(name, settings, handlers, logger)
	case "http":
		return NewServiceHttp(name, settings, handlers, logger)
	}
	return nil, gg.Errors.Prefix(datamover_commons.PanicSystemError,
		fmt.Sprintf("protocol '%s' not supported.", protocol))
}
