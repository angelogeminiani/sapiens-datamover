package services

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_network/message"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_network/services/webserver"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"strings"
)

type ServiceHttp struct {
	name     string
	handlers *DataMoverHandlers
	logger   *datamover_commons.Logger
	uid      string
	enabled  bool
	closed   bool
	settings *webserver.WebserverSettings

	callback func(netMessage *message.NetworkMessage) interface{}
	server   *webserver.WebController
}

func NewServiceHttp(name string, configuration map[string]interface{}, handlers *DataMoverHandlers, logger *datamover_commons.Logger) (instance *ServiceHttp, err error) {
	instance = new(ServiceHttp)
	instance.name = name
	instance.handlers = handlers
	instance.logger = logger

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

	instance.server = webserver.NewWebController(instance.logger, instance.settings)
	instance.server.Handle("all", "*", instance.handle)

	return
}

func (instance *ServiceHttp) handle(ctx *fiber.Ctx) error {
	method := strings.ToLower(ctx.Method())
	url := string(ctx.Request().URI().Path()) // /api/v1

	params := webserver.Params(ctx, true)

	// check if is there any handler
	if _, ok := instance.handlers.CanHandle(method, url); ok {
		// clean params
		args := createScriptParams(ctx.Request().URI(), params)
		response, err := instance.handlers.Handle(method, url, args)
		if nil != err {
			return webserver.WriteResponse(ctx, nil, err)
		}
		return webserver.WriteResponse(ctx, response, nil)
	} else {
		// validate native message
		if len(params) > 0 {
			if nm, err := toNetworkMessage(params); nil == err && nil != instance.callback {
				if nil != nm && len(nm.Body) > 0 {
					// handled
					response := instance.callback(nm)
					if e, b := response.(error); b {
						return webserver.WriteResponse(ctx, nil, e)
					} else {
						return webserver.WriteResponse(ctx, response, nil)
					}
				}
			}
		}
	}

	// unhandled
	return webserver.WriteResponse(ctx, nil, gg.Errors.Prefix(datamover_commons.PanicSystemError,
		fmt.Sprintf("No handlers for '%s'", url)))
}

func toNetworkMessage(data interface{}) (response *message.NetworkMessage, err error) {
	err = gg.JSON.Read(gg.Convert.ToString(data), &response)
	return
}

func toNetworkMessagePayload(data interface{}) (response *message.NetworkMessagePayload, err error) {
	err = gg.JSON.Read(gg.Convert.ToString(data), &response)
	return
}

func createScriptParams(uri *fasthttp.URI, params map[string]interface{}) (response map[string]interface{}) {
	response = make(map[string]interface{})

	// add uri params
	query := uri.QueryArgs()
	if nil != query {
		for k, v := range params {
			if query.Has(k) {
				response[k] = v
			}
		}
	}

	// add message params if any (only for native messages)
	if nm, err := toNetworkMessage(params); nil == err {
		if payload, err := toNetworkMessagePayload(nm.Body); nil == err {
			// variables
			if len(payload.ActionContextVariables) > 0 {
				gg.Maps.Merge(false, response, payload.ActionContextVariables)
			}

			// globals
			globals := payload.ActionGlobals

			// database
			config := payload.ActionConfig
			if nil != config && len(config.Command) > 0 && nil != config.Connection {
				// replace config with globals if existing
				if len(config.Connection.ConnectionsId) > 0 && nil != globals && globals.HasConnections() {
					// replace with global connection
					conn := globals.GetConnection(config.Connection.ConnectionsId)
					if nil != conn {
						config.Connection = conn
					}
				}
				mconn := map[string]interface{}{
					"command": config.Command,
					"schema":  config.FieldsMapping,
					"connection": map[string]interface{}{
						"driver": config.Connection.Driver,
						"dsn":    config.Connection.Dsn,
					},
				}

				response["database"] = mconn
			}

			// data (contextData)
			if len(payload.ActionContextData) > 0 {
				response["data"] = payload.ActionContextData
			}
		}
	}

	return
}
