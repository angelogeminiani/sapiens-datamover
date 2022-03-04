package scripting

import (
	"bitbucket.org/digi-sense/gg-core"
	ggx "bitbucket.org/digi-sense/gg-core-x"
	"bitbucket.org/digi-sense/gg-core-x/gg_log"
	"bitbucket.org/digi-sense/gg-core-x/gg_scripting"
	"bitbucket.org/digi-sense/gg-core-x/gg_scripting/commons"
	"bitbucket.org/digi-sense/gg-core/gg_utils"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/scripting/modules"
	"fmt"
	"path"
)

type ScriptController struct {
	root   string
	dir    string
	logger *gg_log.Logger
}

func NewScriptController(root string) (instance *ScriptController) {
	instance = new(ScriptController)
	instance.root = root
	instance.dir = path.Base(root)

	_ = gg.Paths.Mkdir(instance.root + gg_utils.OS_PATH_SEPARATOR)
	logfile := gg.Paths.Concat(instance.root, "logging-script.log")
	_ = gg.IO.Remove(logfile)
	instance.logger = gg_log.NewLogger()
	instance.logger.SetFileName(logfile)

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ScriptController) RunWithArray(name, script string, contextData []map[string]interface{}, contextVariables map[string]interface{}) (data []map[string]interface{}, variables map[string]interface{}) {
	data = contextData
	variables = contextVariables
	if len(script) == 0 {
		return
	}

	env := environment(fmt.Sprintf("%s#%s", instance.dir, name), script, data, variables)

	ggx.Scripting.SetLogger(instance.logger)
	response, err := ggx.Scripting.Run(env)

	if nil == err {
		mResponse := instance.toMap(response)
		if d, b := mResponse["data"]; b {
			dd := instance.toArrayOfMap(d)
			if nil != dd {
				data = dd
			}
		}
		if v, b := mResponse["variables"]; b {
			if mm, ok := v.(map[string]interface{}); ok {
				variables = mm
			}
		}
	} else {
		instance.logger.Error(fmt.Sprintf("Error executing script '%s': %s", gg.Paths.FileName(script, true), err.Error()))
	}

	return
}

// not used
func (instance *ScriptController) _RunWithRow(name, script string, contextData map[string]interface{}, contextVariables map[string]interface{}) (data, variables map[string]interface{}) {
	data = contextData
	variables = contextVariables
	if len(script) == 0 {
		return
	}

	env := environment(fmt.Sprintf(instance.dir, name), script, data, variables)

	ggx.Scripting.SetLogger(instance.logger)
	response, err := ggx.Scripting.Run(env)

	if nil == err {
		mResponse := instance.toMap(response)
		if d, b := mResponse["data"]; b {
			if dd, ok := d.(map[string]interface{}); ok {
				data = dd
			}
		}
		if v, b := mResponse["variables"]; b {
			if mm, ok := v.(map[string]interface{}); ok {
				variables = mm
			}
		}
	} else {
		instance.logger.Error(fmt.Sprintf("Error executing script '%s': %s", gg.Paths.FileName(script, true), err.Error()))
	}

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ScriptController) toArray(value string) []map[string]interface{} {
	var result []map[string]interface{}
	_ = gg.JSON.Read(value, &result)
	return result
}

func (instance *ScriptController) toArrayOfMap(value interface{}) (response []map[string]interface{}) {
	if nil != value {
		a := gg.Convert.ToArray(value)
		for _, item := range a {
			response = append(response, gg.Convert.ToMap(item))
		}
	}
	return
}

func (instance *ScriptController) toMap(value string) map[string]interface{} {
	var result map[string]interface{}
	_ = gg.JSON.Read(value, &result)
	return result
}

// ---------------------------------------------------------------------------------------------------------------------
//	S T A T I C
// ---------------------------------------------------------------------------------------------------------------------

func environment(programName, script string, data interface{}, variables map[string]interface{}) *gg_scripting.EnvSettings {
	return &gg_scripting.EnvSettings{
		EngineName:    "js",
		ProgramName:   programName,
		ProgramScript: script,
		OnReady: func(rtContext *commons.RuntimeContext) {
			modules.EnableModuleDatasets(rtContext)
		},
		Context: map[string]interface{}{
			"$data":      data,
			"$variables": variables,
		},
	}
}