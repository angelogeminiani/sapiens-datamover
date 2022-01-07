package datamover

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_initializer"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_postman"
	"fmt"
	"path/filepath"
	"time"
)

type DataMover struct {
	mode    string
	root    string
	dirWork string

	settings    *datamover_commons.DataMoverSettings
	logger      *datamover_commons.Logger
	stopChan    chan bool
	stopMonitor *stopMonitor
	events      *gg_events.Emitter

	jobs    *datamover_jobs.DataMoverJobsController
	postman *datamover_postman.DataMoverPostman
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMover) Start() (err error) {
	instance.stopChan = make(chan bool, 1)
	if nil != instance.stopMonitor {
		instance.stopMonitor.Start()
	}

	if nil != instance.postman {
		_ = instance.postman.Start()
	}

	if nil != instance.jobs {
		err = instance.jobs.Start()
	}

	return
}

// Stop Try to close gracefully
func (instance *DataMover) Stop() (err error) {
	if nil != instance.stopMonitor {
		instance.stopMonitor.Stop()
	}

	if nil != instance.postman {
		_ = instance.postman.Stop()
	}

	if nil != instance.jobs {
		_ = instance.jobs.Stop()
	}

	time.Sleep(3 * time.Second)
	if nil != instance.stopChan {
		instance.stopChan <- true
		instance.stopChan = nil
	}
	return
}

// Exit application
func (instance *DataMover) Exit() (err error) {
	if nil != instance.stopMonitor {
		instance.stopMonitor.Stop()
	}
	if nil != instance.stopChan {
		instance.stopChan <- true
		instance.stopChan = nil
	}

	return
}

func (instance *DataMover) Join() {
	if nil != instance.stopChan {
		<-instance.stopChan
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMover) doStop(_ *gg_events.Event) {
	_ = instance.Exit()
}

// ---------------------------------------------------------------------------------------------------------------------
//		S T A T I C
// ---------------------------------------------------------------------------------------------------------------------

func LaunchApplication(mode, cmdStop string, args ...interface{}) (instance *DataMover, err error) {
	instance = new(DataMover)
	instance.mode = mode

	// paths
	instance.dirWork = gg.Paths.GetWorkspace(datamover_commons.WpDirWork).GetPath()
	instance.root = filepath.Dir(instance.dirWork)

	// initialize environment
	err = datamover_initializer.Initialize(mode)
	if nil != err {
		return
	}

	instance.settings, err = datamover_commons.NewSettings(mode)
	if nil == err {
		instance.events = gg.Events.NewEmitter(datamover_commons.AppName)
		instance.stopMonitor = newStopMonitor([]string{instance.root, instance.dirWork}, cmdStop, instance.events)
		instance.events.On(datamover_commons.EventOnDoStop, instance.doStop)

		// logger as first parameter
		l := gg.Arrays.GetAt(args, 0, nil)
		instance.logger = datamover_commons.NewLogger(mode, l)

		// POSTMAN
		instance.postman, err = datamover_postman.NewPostman(instance.settings.Postman,
			instance.logger, instance.events)
		if nil != err {
			goto exit
		}

		// JOBS
		instance.jobs, err = datamover_jobs.NewDataMoverJobsController(instance.settings.PathJobs,
			instance.logger, instance.events)
		if nil != err {
			goto exit
		}
	}

	// final log
exit:
	if nil != err {
		instance.logger.Error(fmt.Sprintf("ERROR starting '%s' v%s: %s",
			datamover_commons.AppName, datamover_commons.AppVersion, err.Error()))
	} else {
		instance.logger.Info(fmt.Sprintf("STARTED '%s' v%s with jobs into dir: '%s'",
			datamover_commons.AppName, datamover_commons.AppVersion, instance.settings.PathJobs))
	}
	return
}
