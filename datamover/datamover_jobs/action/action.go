package action

import (
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"gorm.io/gorm"
)

type DataMoverAction struct {
	root string
	uid  string

	datasource *DataMoverDatasource
}

func NewDataMoverAction(root string, datasourceSettings *datamover_commons.DataMoverDatasourceSettings) *DataMoverAction {
	instance := new(DataMoverAction)
	instance.root = root
	if nil != datasourceSettings {
		instance.uid = datasourceSettings.Uid

		instance.datasource = NewDataMoverDatasource(root, datasourceSettings)

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

func (instance *DataMoverAction) Execute(context []map[string]interface{}) (result []map[string]interface{}, err error) {
	result = make([]map[string]interface{}, 0)

	command := instance.datasource.settings.Command
	if len(command) > 0 {
		result, err = instance.datasource.GetData(context, command)
		if nil != err {
			return
		}
	}

	script := instance.datasource.settings.Script
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
	return nil
}

func (instance *DataMoverAction) connection() (*gorm.DB, error) {
	if nil != instance && nil != instance.datasource {
		return instance.datasource.connection()
	}
	return nil, datamover_commons.PanicSystemError
}
