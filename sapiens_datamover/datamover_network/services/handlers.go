package services

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_jobs/action"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_jobs/action/scripting"
	"strings"
)

type DataMoverHandlers struct {
	handlers    []*datamover_commons.SettingsNetHandler
	logger      *datamover_commons.Logger
	fnVarEngine *gg_fnvars.FnVarsEngine

	executors map[string]*scripting.ScriptController
}

func NewDataMoverHandlers(handlers []*datamover_commons.SettingsNetHandler, logger *datamover_commons.Logger, fnVarEngine *gg_fnvars.FnVarsEngine) (instance *DataMoverHandlers) {
	instance = new(DataMoverHandlers)
	instance.handlers = handlers
	instance.logger = logger
	instance.fnVarEngine = fnVarEngine
	instance.executors = make(map[string]*scripting.ScriptController)

	instance.init()

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverHandlers) CanHandle(method, endpoint string) (*datamover_commons.SettingsNetHandler, bool) {
	if nil != instance && len(instance.handlers) > 0 {
		for _, handler := range instance.handlers {
			if handler.Enabled && (len(handler.Method) == 0 || handler.Method == method || handler.Method == "all") {
				if strings.HasPrefix(endpoint, handler.Endpoint) {
					return handler, true
				}
			}
		}
	}
	return nil, false
}

func (instance *DataMoverHandlers) Handle(method, endpoint string, params interface{}) (response interface{}, err error) {
	if nil != instance {
		if handler, ok := instance.CanHandle(method, endpoint); ok {
			filename := gg.Paths.WorkspacePath(handler.Handler)
			script, ioErr := gg.IO.ReadTextFromFile(filename)
			if nil != ioErr {
				err = ioErr
				return
			}
			// ready lo run the script
			if len(script) > 0 {
				executor := instance.executor(gg.Paths.Dir(filename))
				if nil != executor {
					return executor.Run(gg.Paths.FileName(filename, false), script, instance.solveParams(params))
				}
			}
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverHandlers) init() {

}

func (instance *DataMoverHandlers) executor(root string) *scripting.ScriptController {
	if nil != instance && nil != instance.executors {
		if item, ok := instance.executors[root]; ok {
			return item
		}
		instance.executors[root] = scripting.NewScriptController(root)
		return instance.executors[root]
	}
	return nil
}

func (instance *DataMoverHandlers) solveParams(params interface{}) (response map[string]interface{}) {
	response = gg.Convert.ToMap(params)

	// context data (existing from native action command)
	data := make([]map[string]interface{}, 0)
	if idata, ok := response["data"]; ok {
		if dd, ok := idata.([]map[string]interface{}); ok {
			data = dd
		}
	}
	if idatabase, ok := response["database"]; ok {
		database := gg.Convert.ToMap(idatabase)
		// fields mapping
		schema := make(map[string]interface{})
		if ischema, ok := database["schema"]; ok {
			schema = ischema.(map[string]interface{})
		}
		if icommand, ok := database["command"]; ok {
			command := gg.Convert.ToString(icommand)
			c, e := instance.fnVarEngine.SolveText(command)
			if nil == e {
				c = action.PreProcessSpecials(c, response)
				database["command"] = c
				commands := make([]string, 0)
				if len(data) > 0 {
					for _, item := range data {
						commands = append(commands, action.ToSQLStatement(c, item, schema))
					}
				} else {
					commands = append(commands, c)
				}
				database["commands"] = commands
			}
		}
	}

	return
}
