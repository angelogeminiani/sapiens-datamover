package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/schema"
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

func NewDataMoverAction(root string, fnVarEngine *gg_fnvars.FnVarsEngine, datasourceSettings *datamover_commons.DataMoverActionSettings) (instance *DataMoverAction, err error) {
	instance = new(DataMoverAction)
	instance.root = root
	instance.actionSettings = datasourceSettings

	if nil != datasourceSettings {
		instance.uid = datasourceSettings.Uid
		instance.datasource, err = NewDataMoverDatasource(root, fnVarEngine,
			datasourceSettings.Connection, datasourceSettings.Scripts)
		// init the action
		if nil == err {
			err = instance.init()
		}
	}
	return instance, err
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverAction) Uid() string {
	if nil != instance {
		return instance.uid
	}
	return ""
}

func (instance *DataMoverAction) IsValid() bool {
	return nil != instance.datasource
}

func (instance *DataMoverAction) IsNetworkAction() bool {
	return nil != instance.actionSettings && nil != instance.actionSettings.Network
}

func (instance *DataMoverAction) GetSchema() *schema.DataMoverDatasourceSchema {
	if nil != instance && nil != instance.datasource {
		return instance.datasource.GetSchema()
	}
	return nil
}

func (instance *DataMoverAction) Execute(contextData []map[string]interface{}, variables map[string]interface{}) (result []map[string]interface{}, err error) {
	result = make([]map[string]interface{}, 0)

	if instance.IsNetworkAction() {
		// REMOTE
		if nil != instance.clientNet {
			payload := new(message.NetworkMessagePayload)
			payload.ActionRoot = instance.root
			payload.ActionConfig = instance.actionSettings
			payload.ActionContextData = contextData
			payload.ActionContextVariables = variables
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
			result, err = instance.datasource.GetData(command, contextData, variables)
		}
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
			instance.actionSettings.Network)
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
