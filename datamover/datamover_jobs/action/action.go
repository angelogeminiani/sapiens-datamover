package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_globals"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/schema"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/clients"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_network/message"
	"gorm.io/gorm"
	"strings"
)

type DataMoverAction struct {
	root string
	uid  string

	actionSettings *datamover_commons.DataMoverActionSettings
	datasource     *DataMoverDatasource
	_clientNet     clients.ClientNetwork
	globals        *datamover_globals.Globals
}

func NewDataMoverAction(root string, fnVarEngine *gg_fnvars.FnVarsEngine, datasourceSettings *datamover_commons.DataMoverActionSettings, globals *datamover_globals.Globals) (instance *DataMoverAction, err error) {
	instance = new(DataMoverAction)
	instance.root = root
	instance.actionSettings = datasourceSettings
	instance.globals = globals

	if nil != datasourceSettings {
		instance.uid = datasourceSettings.Uid
		connection := datasourceSettings.Connection
		if len(connection.ConnectionsId) > 0 && instance.globals.HasConnections() {
			// replace with global connection
			conn := instance.globals.GetConnection(connection.ConnectionsId)
			if nil != conn {
				connection = conn
			}
		}
		instance.datasource, err = NewDataMoverDatasource(root, fnVarEngine,
			connection, datasourceSettings.Scripts)
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
		client, e := instance.getClientNet()
		if nil != client {
			payload := new(message.NetworkMessagePayload)
			payload.ActionName = "net-command"
			payload.ActionRoot = instance.root
			payload.ActionRootRelative = strings.ReplaceAll(instance.root, gg.Paths.WorkspacePath("./"), ".")
			payload.ActionConfig = instance.actionSettings
			payload.ActionContextData = contextData
			payload.ActionContextVariables = variables
			payload.ActionGlobals = instance.globals
			payload.ActionDatasets = LoadJsDatasets(instance.root) // load datasets for remote transfer

			// execute
			respData, respErr := client.Send(payload.String())
			if nil != respErr {
				if strings.HasPrefix(respErr.Error(), "dial tcp") {
					instance._clientNet = nil // reset client
				}
				err = respErr
				return
			}
			// deserialize
			var res *message.NetworkMessageResponseBody
			err = gg.JSON.Read(gg.Convert.ToString(respData), &res)
			if nil == err {
				// align datasets
				if nil != res.Datasets && len(res.Datasets) > 0 {
					OverwriteJsDatasets(instance.root, res.Datasets)
				}
				// align variables
				gg.Maps.Merge(true, variables, res.Variables)
				// read the body
				err = gg.JSON.Read(gg.Convert.ToString(res.Body), &result)
			}
		} else {
			err = e
		}
	} else {
		// LOCAL
		command := instance.actionSettings.Command
		if len(command) > 0 {
			mapping := instance.actionSettings.FieldsMapping
			result, err = instance.datasource.GetData(command, mapping, contextData, instance.globals.MergeVariables(variables))
		}
	}

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverAction) init() (err error) {
	_, err = instance.datasource.connection()
	if nil != err {
		return
	}

	_ = instance.initNetwork() // network is lazy initialized

	return
}

func (instance *DataMoverAction) initNetwork() (err error) {
	if instance.IsNetworkAction() {
		_, err = instance.getClientNet() // test connection
	}
	return
}

func (instance *DataMoverAction) getClientNet() (clients.ClientNetwork, error) {
	if instance.IsNetworkAction() && nil == instance._clientNet {
		uri, err := instance.actionSettings.Network.Uri()
		if nil != err {
			instance._clientNet = nil
			return instance._clientNet, err
		}
		c, e := clients.BuildNetworkClient(uri,
			instance.actionSettings.Network)
		if nil != e {
			instance._clientNet = nil
			return instance._clientNet, e
		}
		instance._clientNet = c
	}
	return instance._clientNet, nil
}

func (instance *DataMoverAction) connection() (*gorm.DB, error) {
	if nil != instance && nil != instance.datasource {
		return instance.datasource.connection()
	}
	return nil, datamover_commons.PanicSystemError
}
