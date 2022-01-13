package clients

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_nio"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
)

type ClientNio struct {
	host   string
	port   int
	auth   *datamover_commons.SettingsNetAuthorization
	secure bool

	_client   *gg_nio.NioClient
	connected bool
}

func NewClientNio(host string, port int, settings *datamover_commons.DataMoverNetworkSettings) (instance *ClientNio, err error) {
	instance = new(ClientNio)
	instance.host = host
	instance.port = port
	instance.auth = settings.Authentication
	instance.secure = settings.Secure

	// test client
	_, err = instance.client()

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ClientNio) Send(payload string) (interface{}, error) {
	client, err := instance.client()
	if nil != err {
		return nil, err
	}
	defer client.Close()

	// send packet message
	client.Secure = instance.secure
	resp, e := client.Send(instance.buildMessage(payload))
	rawMessage := resp.Body

	return rawMessage, e
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ClientNio) client() (*gg_nio.NioClient, error) {
	var err error
	if nil == instance._client {
		instance._client = gg.NIO.NewClient(instance.host, instance.port)
		instance._client.Secure = false // avoid encryption overload that slow down data transmission
		instance._client.OnDisconnect(instance.onDisconnect)
		instance._client.OnConnect(instance.onConnect)

		err = instance._client.Open()
	}
	return instance._client, err
}

func (instance *ClientNio) onDisconnect(e *gg_events.Event) {
	instance.connected = false
	instance._client = nil
}

func (instance *ClientNio) onConnect(e *gg_events.Event) {
	instance.connected = true
}

func (instance *ClientNio) buildMessage(payload string) *message.NetworkMessage {
	message := new(message.NetworkMessage)
	message.SetAuthorization(instance.auth.Value)
	message.Body = payload
	return message
}
