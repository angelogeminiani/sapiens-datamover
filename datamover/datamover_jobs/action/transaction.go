package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_log"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"fmt"
)

type DataMoverTransaction struct {
	root     string
	logger   *gg_log.Logger
	events   *gg_events.Emitter
	settings []*datamover_commons.DataMoverDatasourceSettings

	fnVarEngine *gg_fnvars.FnVarsEngine
	transaction []*DataMoverAction
}

func NewDataMoverTransaction(root string, logger *gg_log.Logger, events *gg_events.Emitter, settings []*datamover_commons.DataMoverDatasourceSettings) *DataMoverTransaction {
	instance := new(DataMoverTransaction)
	instance.root = root
	instance.logger = logger
	instance.events = events
	instance.settings = settings

	instance.init()

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverTransaction) Execute(context []map[string]interface{}) (interface{}, error) {
	if nil != instance && nil != instance.transaction {
		var err error
		for _, action := range instance.transaction {
			if action.IsValid() {
				context, err = action.Execute(context)
				if nil != err {
					return nil, gg.Errors.Prefix(err, fmt.Sprintf("Action '%s' got error: ", action.uid))
				}
			} else {
				return nil, gg.Errors.Prefix(datamover_commons.ActionInvalidConfigurationError,
					fmt.Sprintf("'%s': ", action.uid))
			}
		}
		return context, nil
	}
	return nil, nil
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverTransaction) init() {
	instance.fnVarEngine = gg.FnVars.NewEngine()

	if nil != instance.settings {
		for _, setting := range instance.settings {
			instance.transaction = append(instance.transaction,
				NewDataMoverAction(instance.root, instance.fnVarEngine, setting))
		}
	}
}
