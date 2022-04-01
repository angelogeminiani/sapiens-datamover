package message

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_globals"
)

type NetworkMessage struct {
	Header map[string]interface{} `json:"header"`
	Body   string                 `json:"body"`
}

func (instance *NetworkMessage) String() string {
	return gg.JSON.Stringify(instance)
}

func (instance *NetworkMessage) GetHeader() map[string]interface{} {
	if nil == instance.Header {
		instance.Header = make(map[string]interface{})
	}
	return instance.Header
}

func (instance *NetworkMessage) SetAuthorization(token string) {
	instance.GetHeader()["Authorization"] = token
}

func (instance *NetworkMessage) GetAuthorization() string {
	return gg.Convert.ToString(instance.GetHeader()["Authorization"])
}

func (instance *NetworkMessage) SetHeader(key, value string) {
	instance.GetHeader()[key] = value
}

type NetworkMessagePayload struct {
	ActionName             string                                     `json:"name"`
	ActionRoot             string                                     `json:"root"`
	ActionRootRelative     string                                     `json:"root_relative"`
	ActionConfig           *datamover_commons.DataMoverActionSettings `json:"settings"`
	ActionContextData      []map[string]interface{}                   `json:"data"`
	ActionContextVariables map[string]interface{}                     `json:"variables"`
	ActionGlobals          *datamover_globals.Globals                 `json:"globals"`
	ActionDatasets         map[string][]interface{}                   `json:"js_datasets"`
}

func (instance *NetworkMessagePayload) String() string {
	return gg.JSON.Stringify(instance)
}

// NetworkMessageResponseBody wrap the response body
type NetworkMessageResponseBody struct {
	Body      interface{}              `json:"body"`
	Variables map[string]interface{}   `json:"variables"`
	Datasets  map[string][]interface{} `json:"js_datasets"`
}
