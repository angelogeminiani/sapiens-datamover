package datamover_jobs

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_utils"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"fmt"
)

type DataMoverJobsController struct {
	isDebug bool
	root    string
	logger  *datamover_commons.Logger
	events  *gg_events.Emitter

	closed bool
	jobs   []*DataMoverJob
}

func NewDataMoverJobsController(debug bool, root string, logger *datamover_commons.Logger, events *gg_events.Emitter) (instance *DataMoverJobsController, err error) {
	instance = new(DataMoverJobsController)
	instance.isDebug = debug
	instance.root = gg.Paths.WorkspacePath(root)
	instance.logger = logger
	instance.events = events
	instance.closed = true

	err = instance.init()

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

func (instance *DataMoverJobsController) init() (err error) {

	// creates root if any
	err = gg.Paths.Mkdir(instance.root + gg_utils.OS_PATH_SEPARATOR)
	if nil != err {
		return
	}

	// load jobs
	var dirs []string
	dirs, err = gg.Paths.ListDir(instance.root)
	if nil == err {
		for _, dir := range dirs {
			job, jerr := NewDataMoverJob(instance.isDebug, dir, instance.events)
			if nil == jerr {
				instance.jobs = append(instance.jobs, job)
			} else {
				err = gg.Errors.Prefix(err, fmt.Sprintf("Error creating JOB '%s'", dir))
				return
			}
		}
	}

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
