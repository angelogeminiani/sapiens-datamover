package scripting

import (
	"bitbucket.org/digi-sense/gg-core"
	ggx "bitbucket.org/digi-sense/gg-core-x"
	"bitbucket.org/digi-sense/gg-core-x/gg_log"
	"bitbucket.org/digi-sense/gg-core-x/gg_scripting"
	"bitbucket.org/digi-sense/gg-core/gg_utils"
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

func (instance *ScriptController) RunWithArray(name, script string, context []map[string]interface{}) []map[string]interface{} {
	if len(script) == 0 {
		return context
	}
	env := &gg_scripting.EnvSettings{
		EngineName:    "js",
		ProgramName:   fmt.Sprintf(instance.dir, name),
		ProgramScript: script,
		Context: map[string]interface{}{
			"$data": context,
		},
	}
	ggx.Scripting.SetLogger(instance.logger)
	response, err := ggx.Scripting.Run(env)
	if nil != err {
		return context
	}
	return instance.toArray(response)
}

func (instance *ScriptController) RunWithRow(name, script string, row map[string]interface{}) map[string]interface{} {
	if len(script) == 0 {
		return row
	}
	env := &gg_scripting.EnvSettings{
		EngineName:    "js",
		ProgramName:   fmt.Sprintf(instance.dir, name),
		ProgramScript: script,
		Context: map[string]interface{}{
			"$data": row,
		},
	}
	ggx.Scripting.SetLogger(instance.logger)
	response, err := ggx.Scripting.Run(env)
	if nil != err {
		return row
	}
	return instance.toMap(response)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *ScriptController) toArray(value string) []map[string]interface{} {
	var result []map[string]interface{}
	_ = gg.JSON.Read(value, &result)
	return result
}

func (instance *ScriptController) toMap(value string) map[string]interface{} {
	var result map[string]interface{}
	_ = gg.JSON.Read(value, &result)
	return result
}
