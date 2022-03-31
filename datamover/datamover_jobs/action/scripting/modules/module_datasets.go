package modules

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_scripting/commons"
	"bitbucket.org/digi-sense/gg-core-x/gg_scripting/modules/defaults/require"
	"github.com/dop251/goja"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const NAME = "datasets"

type ModuleDatasets struct {
	name     string
	root     string
	filename string
	runtime  *goja.Runtime
	datasets map[string][]interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *ModuleDatasets) LoadData(data map[string][]interface{}) error {
	return instance.LoadText(gg.JSON.Stringify(data))
}

func (instance *ModuleDatasets) LoadText(data string) error {
	return gg.JSON.Read(data, &instance.datasets)
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p o s e d
//----------------------------------------------------------------------------------------------------------------------

// datasets.add(string, array)
// add a dataset to internal cache
func (instance *ModuleDatasets) put(call goja.FunctionCall) goja.Value {
	count := call.Arguments
	if len(count) == 2 {
		name := call.Argument(0).String()
		data := call.Argument(1).Export()
		if len(name) > 0 && nil != data {
			arr := gg.Convert.ToArray(data)
			if nil != arr {
				instance.datasets[name] = arr
				instance.save()
			}
		}
	} else {
		panic(instance.runtime.NewTypeError("Wrong number of parameters. Expected 'name', 'data'"))
	}
	return goja.Undefined()
}

func (instance *ModuleDatasets) get(call goja.FunctionCall) goja.Value {
	count := call.Arguments
	if len(count) == 1 {
		name := call.Argument(0).String()
		if len(name) > 0 {
			if data, ok := instance.datasets[name]; ok {
				return instance.runtime.ToValue(data)
			}
		}
	} else {
		panic(instance.runtime.NewTypeError("Wrong number of parameters. Expected 'name'"))
	}
	return goja.Undefined()
}

func (instance *ModuleDatasets) mapFunc(call goja.FunctionCall) goja.Value {
	response := make([]interface{}, 0)
	count := call.Arguments
	if len(count) == 2 {
		name := call.Argument(0).String()
		callback := commons.GetCallbackIfAny(call)
		if len(name) > 0 && nil != callback {
			if data, ok := instance.datasets[name]; ok {
				for _, item := range data {
					resp, err := callback(call.This, instance.runtime.ToValue(item))
					if nil != err {
						panic(instance.runtime.NewTypeError(err))
					}
					if r := resp.Export(); nil != r {
						response = append(response, r)
					}
				}
			}
		}
	} else {
		panic(instance.runtime.NewTypeError("Wrong number of parameters. Expected 'name', 'callback'"))
	}
	return instance.runtime.ToValue(response)
}

func (instance *ModuleDatasets) forFunc(call goja.FunctionCall) goja.Value {
	count := call.Arguments
	if len(count) == 2 {
		name := call.Argument(0).String()
		callback := commons.GetCallbackIfAny(call)
		if len(name) > 0 && nil != callback {
			if data, ok := instance.datasets[name]; ok {
				for _, item := range data {
					resp, err := callback(call.This, instance.runtime.ToValue(item))
					if nil != err {
						panic(instance.runtime.NewTypeError(err))
					}
					if r := resp.Export(); nil != r {
						// exit loop
						return instance.runtime.ToValue(r)
					}
				}
			}
		}
	} else {
		panic(instance.runtime.NewTypeError("Wrong number of parameters. Expected 'name', 'callback'"))
	}
	return goja.Undefined()
}

func (instance *ModuleDatasets) save() goja.Value {
	_, err := gg.IO.WriteTextToFile(gg.JSON.Stringify(instance.datasets), instance.filename)
	if nil != err {
		panic(instance.runtime.NewTypeError(err))
	}
	return instance.runtime.ToValue(instance.filename)
}

func (instance *ModuleDatasets) reset() goja.Value {
	instance.datasets = make(map[string][]interface{})
	return goja.Undefined()
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *ModuleDatasets) init() {
	if ok, _ := gg.Paths.Exists(instance.filename); !ok {
		instance.datasets = make(map[string][]interface{})
		_, _ = gg.IO.WriteTextToFile(gg.JSON.Stringify(instance.datasets), instance.filename)
	} else {
		// load
		_ = gg.JSON.ReadFromFile(instance.filename, &instance.datasets)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func load(runtime *goja.Runtime, module *goja.Object, args ...interface{}) {
	instance := &ModuleDatasets{
		runtime:  runtime,
		datasets: make(map[string][]interface{}),
	}

	if len(args) > 1 {
		workspace := gg.Paths.WorkspacePath("./")
		name := gg.Convert.ToString(gg.Reflect.ValueOf(args[1]).Interface())
		if len(name) > 0 {
			tokens := strings.Split(name, "#")
			instance.name = tokens[0]
			instance.root = gg.Paths.Concat(workspace, "jobs", instance.name)
			instance.filename = gg.Paths.Concat(instance.root, NAME+".json")
		}
	}

	instance.init()

	o := module.Get("exports").(*goja.Object)
	_ = o.Set("put", instance.put)
	_ = o.Set("get", instance.get)
	_ = o.Set("map", instance.mapFunc)
	_ = o.Set("for", instance.forFunc)
	_ = o.Set("save", instance.save)
	_ = o.Set("reset", instance.reset)
}

func EnableModuleDatasets(ctx *commons.RuntimeContext) {
	// register
	require.RegisterNativeModule(NAME, &commons.ModuleInfo{
		Context: ctx,
		Loader:  load,
	})
}
