package clients

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_http/httpclient"
	"bitbucket.org/digi-sense/gg-core-x/gg_http/httputils"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
	"fmt"
	"strings"
)

type ClientHttp struct {
	method   string
	endpoint string
	auth     *datamover_commons.SettingsNetAuthorization
	secure   bool

	_client   *httpclient.HttpClient
	connected bool
}

func NewClientHttp(endpoint string, settings *datamover_commons.DataMoverNetworkSettings) (instance *ClientHttp, err error) {
	instance = new(ClientHttp)
	instance.method = "post"
	instance.endpoint = endpoint
	instance.auth = settings.Authentication
	instance.secure = settings.Secure

	// parse endpoint
	tokens := gg.Strings.Split(endpoint, "|,")
	if len(tokens) > 1 {
		instance.method = strings.ToLower(tokens[0])
		instance.endpoint = strings.Join(tokens[1:], "")
	}
	// test client
	_, err = instance.client()

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ClientHttp) Send(payload string) (rawMessage interface{}, err error) {
	client, e := instance.client()
	if nil != e {
		err = e
		return
	}

	// send packet message
	var resp *httputils.ResponseData
	switch instance.method {
	case "post":
		msg := instance.buildMessage(payload)
		resp, err = client.Post(instance.endpoint, msg.String())
	case "get":
		resp, err = client.Get(instance.endpoint)
	}
	if nil != err {
		return
	}

	if len(resp.Body) > 0 {
		str := gg.Convert.ToString(resp.Body)
		m := gg.Convert.ToMap(str)
		if e, ok := m["error"]; ok {
			// response error
			err = gg.Errors.Prefix(datamover_commons.PanicSystemError, fmt.Sprintf("%v", e))
			rawMessage = nil
			return
		} else {
			rawMessage = m["response"]
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ClientHttp) client() (*httpclient.HttpClient, error) {
	var err error
	if nil == instance._client {
		instance._client = httpclient.NewHttpClient()

		if instance.secure {
			instance._client.AddHeader("Authorization", instance.auth.String())
		}

	}
	return instance._client, err
}

func (instance *ClientHttp) buildMessage(payload string) *message.NetworkMessage {
	response := new(message.NetworkMessage)
	response.SetAuthorization(instance.auth.String())
	response.SetHeader("referral", "datamover")
	response.Body = payload
	return response
}
