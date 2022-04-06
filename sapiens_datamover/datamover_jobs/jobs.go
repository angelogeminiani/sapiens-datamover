package datamover_jobs

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_fnvars"
	"bitbucket.org/digi-sense/gg-core/gg_utils"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_globals"
	"fmt"
	"path"
)

type DataMoverJobsController struct {
	isDebug     bool
	root        string
	logger      *datamover_commons.Logger
	events      *gg_events.Emitter
	fnVarEngine *gg_fnvars.FnVarsEngine

	closed  bool
	globals *datamover_globals.Globals
	jobs    []*DataMoverJob
}

func NewDataMoverJobsController(mode string, root string, logger *datamover_commons.Logger, events *gg_events.Emitter, fnVarEngine *gg_fnvars.FnVarsEngine) (instance *DataMoverJobsController, err error) {
	instance = new(DataMoverJobsController)
	instance.isDebug = mode == "debug"
	instance.root = gg.Paths.WorkspacePath(root)
	instance.logger = logger
	instance.events = events
	instance.fnVarEngine = fnVarEngine
	instance.closed = true

	err = instance.init(mode)

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverJobsController) Start() (err error) {
	if nil != instance {
		for _, job := range instance.jobs {
			err = job.Start()
			if nil != err {
				return
			}
		}
		instance.closed = false
	}
	return
}

func (instance *DataMoverJobsController) Stop() (err error) {
	if nil != instance {
		instance.closed = true
		for _, job := range instance.jobs {
			_ = job.Stop()
		}
	}
	return
}

func (instance *DataMoverJobsController) Run(name string, contextData []map[string]interface{}, contextVariables map[string]interface{}) (err error) {
	if nil != instance && !instance.closed {
		for _, job := range instance.jobs {
			if job.name == name {
				if job.IsRunning() {
					err = gg.Errors.Prefix(datamover_commons.JobAlreadyRunningError, name)
				} else {
					err = job.Run(contextData, contextVariables)
				}
				break
			}
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverJobsController) init(mode string) (err error) {

	// creates root if any
	err = gg.Paths.Mkdir(instance.root + gg_utils.OS_PATH_SEPARATOR)
	if nil != err {
		return
	}

	instance.globals = datamover_globals.NewGlobals(mode)
	instance.logger.Debug(fmt.Sprintf("Loading JOBS from: %s", instance.root))

	count := 0
	// load jobs
	var dirs []string
	dirs, err = gg.Paths.ListDir(instance.root)
	if nil == err {
		for _, dir := range dirs {
			job, jerr := NewDataMoverJob(instance.isDebug, dir, instance.events, instance.fnVarEngine, instance.globals)
			if nil == jerr {
				count++
				instance.logger.Debug(fmt.Sprintf("  * Loaded JOB '%s' (scheduled=%v).", path.Base(dir), job.IsScheduled()))
				instance.jobs = append(instance.jobs, job)
			} else {
				err = gg.Errors.Prefix(err, fmt.Sprintf("Error creating JOB '%s'", dir))
				return
			}
		}
	}

	instance.logger.Debug(fmt.Sprintf("Loaded '%v' JOBS.", count))

	// subscribe events
	instance.events.On(datamover_commons.EventOnNextJobRun, instance.onNextJobRun)

	return
}

// continue chain of jobs
func (instance *DataMoverJobsController) onNextJobRun(event *gg_events.Event) {
	if nil != instance && !instance.closed {
		if parent, b := event.Argument(1).(DataMoverJob); b {
			name := event.ArgumentAsString(0)
			context := event.Argument(2)
			variables := gg.Convert.ToMap(event.Argument(3))
			var err error
			if a, aok := context.([]map[string]interface{}); aok {
				err = instance.Run(name, a, variables)
			} else {
				err = instance.Run(name, nil, variables)
			}
			if nil != err {
				// log on parent
				parent.logger.Error(fmt.Sprintf("Error running '%s': %s", name, err.Error()))
			}
		}
	}
}
