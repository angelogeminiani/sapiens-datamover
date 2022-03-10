package action

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_log"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_globals"
	"fmt"
)

type DataMoverTransaction struct {
	root      string
	logger    *gg_log.Logger
	events    *gg_events.Emitter
	settings  []*datamover_commons.DataMoverActionSettings
	variables map[string]interface{}
	globals   *datamover_globals.Globals

	fnVarEngine *gg_fnvars.FnVarsEngine
	transaction []*DataMoverAction
}

func NewDataMoverTransaction(root string, logger *gg_log.Logger, events *gg_events.Emitter, settings []*datamover_commons.DataMoverActionSettings, variables map[string]interface{}, globals *datamover_globals.Globals) (instance *DataMoverTransaction, err error) {
	instance = new(DataMoverTransaction)
	instance.root = root
	instance.logger = logger
	instance.events = events
	instance.settings = settings
	instance.variables = variables
	instance.globals = globals

	err = instance.init()

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverTransaction) Execute(contextData []map[string]interface{}, contextVariables map[string]interface{}) (responseData []map[string]interface{}, responseVariables map[string]interface{}, err error) {
	if nil != instance && nil != instance.transaction {
		if nil == contextVariables {
			contextVariables = make(map[string]interface{})
		}
		if len(contextVariables) == 0 {
			_ = gg.Maps.Merge(true, contextVariables, instance.variables)
		}
		for _, action := range instance.transaction {
			if action.IsValid() {
				contextData, err = action.Execute(contextData, contextVariables)
				if nil != err {
					err = gg.Errors.Prefix(err, fmt.Sprintf("Action '%s' got error: ", action.uid))
					return
				}
			} else {
				err = gg.Errors.Prefix(datamover_commons.ActionInvalidConfigurationError,
					fmt.Sprintf("'%s': ", action.uid))
				return
			}
		}
		responseData = contextData
		responseVariables = contextVariables
		return
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverTransaction) init() error {
	instance.fnVarEngine = gg.FnVars.NewEngine()

	if nil != instance.settings {
		for _, setting := range instance.settings {
			err := setting.ScriptsLoad(instance.root)
			if nil != err {
				instance.logger.Warn(err)
			}
			action, actionErr := NewDataMoverAction(instance.root, instance.fnVarEngine, setting, instance.globals)
			if nil != actionErr {
				return actionErr
			}
			instance.transaction = append(instance.transaction, action)
			// debug info about schema
			schema := action.GetSchema()
			if nil != schema && instance.logger.GetLevel() == gg_log.LEVEL_DEBUG {
				filename := gg.Paths.Concat(instance.root, fmt.Sprintf("schema.%s.json", action.uid))
				_, _ = gg.IO.WriteTextToFile(schema.String(), filename)
			}
		}
	}

	return nil
}
