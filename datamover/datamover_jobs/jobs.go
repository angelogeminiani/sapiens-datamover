package datamover_jobs

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_utils"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"fmt"
)

type DataMoverJobsController struct {
	root   string
	logger *datamover_commons.Logger
	events *gg_events.Emitter

	closed bool
	jobs   []*DataMoverJob
}

func NewDataMoverJobsController(root string, logger *datamover_commons.Logger, events *gg_events.Emitter) (instance *DataMoverJobsController, err error) {
	instance = new(DataMoverJobsController)
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

func (instance *DataMoverJobsController) Run(name string) (err error) {
	if nil != instance && !instance.closed {
		for _, job := range instance.jobs {
			if job.name == name {
				if job.IsRunning() {
					err = gg.Errors.Prefix(datamover_commons.JobAlreadyRunningError, name)
				} else {
					err = job.Run()
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
			job, jerr := NewDataMoverJob(dir, instance.events)
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

func (instance *DataMoverJobsController) onNextJobRun(event *gg_events.Event) {
	if nil != instance && !instance.closed {
		if parent, b := event.Argument(1).(DataMoverJob); b {
			name := event.ArgumentAsString(0)
			err := instance.Run(name)
			if nil != err {
				// log on parent
				parent.logger.Error(fmt.Sprintf("Error running '%s': %s", name, err.Error()))
			}
		}
	}
}
