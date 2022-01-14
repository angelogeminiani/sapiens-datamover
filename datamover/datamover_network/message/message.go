package message

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
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

type NetworkMessagePayload struct {
	ActionRoot             string                                     `json:"root"`
	ActionConfig           *datamover_commons.DataMoverActionSettings `json:"settings"`
	ActionContextData      []map[string]interface{}                   `json:"data"`
	ActionContextVariables map[string]interface{}                     `json:"variables"`
}

func (instance *NetworkMessagePayload) String() string {
	return gg.JSON.Stringify(instance)
}