package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/clients"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
	"gorm.io/gorm"
)

type DataMoverAction struct {
	root string
	uid  string

	actionSettings *datamover_commons.DataMoverActionSettings
	datasource     *DataMoverDatasource
	clientNet      clients.ClientNetwork
}

func NewDataMoverAction(root string, fnVarEngine *gg_fnvars.FnVarsEngine, datasourceSettings *datamover_commons.DataMoverActionSettings) *DataMoverAction {
	instance := new(DataMoverAction)
	instance.root = root
	instance.actionSettings = datasourceSettings

	if nil != datasourceSettings {
		instance.uid = datasourceSettings.Uid
		instance.datasource = NewDataMoverDatasource(root, fnVarEngine,
			datasourceSettings.Connection, datasourceSettings.Script)

		_ = instance.init()
	}
	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverAction) IsValid() bool {
	return nil != instance.datasource
}

func (instance *DataMoverAction) IsNetworkAction() bool {
	return nil != instance.actionSettings && nil != instance.actionSettings.Network
}

func (instance *DataMoverAction) Execute(context []map[string]interface{}) (result []map[string]interface{}, err error) {
	result = make([]map[string]interface{}, 0)

	if instance.IsNetworkAction() {
		// REMOTE
		if nil != instance.clientNet {
			payload := new(message.NetworkMessagePayload)
			payload.ActionRoot = instance.root
			payload.ActionConfig = instance.actionSettings
			payload.ActionContext = context
			respData, respErr := instance.clientNet.Send(payload.String())
			if nil != respErr {
				err = respErr
				return
			}
			// deserialize
			err = gg.JSON.Read(gg.Convert.ToString(respData), &result)
		}
	} else {
		// LOCAL
		command := instance.actionSettings.Command
		if len(command) > 0 {
			result, err = instance.datasource.GetData(context, command)
		}
	}

	if nil != err {
		return
	}

	script := instance.actionSettings.Script
	if len(script) > 0 {
		
	}

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverAction) init() error {
	_, err := instance.datasource.connection()
	if nil != err {
		return err
	}

	err = instance.initNetwork()
	if nil != err {
		return err
	}

	return nil
}

func (instance *DataMoverAction) initNetwork() error {
	if instance.IsNetworkAction() {
		uri, err := instance.actionSettings.Network.Uri()
		if nil != err {
			return err
		}
		c, e := clients.BuildNetworkClient(uri,
			instance.actionSettings.Network.Authentication)
		if nil != e {
			return e
		}
		instance.clientNet = c
	}
	return nil
}

func (instance *DataMoverAction) connection() (*gorm.DB, error) {
	if nil != instance && nil != instance.datasource {
		return instance.datasource.connection()
	}
	return nil, datamover_commons.PanicSystemError
}
