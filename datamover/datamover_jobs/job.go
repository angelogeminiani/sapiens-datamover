package datamover_jobs

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_log"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_scheduler"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"path"
)

type DataMoverJob struct {
	root   string
	logger *gg_log.Logger
	events *gg_events.Emitter

	name      string
	settings  *datamover_commons.DataMoverSettingsJob
	scheduler *gg_scheduler.Scheduler
}

func NewDataMoverJob(root string, events *gg_events.Emitter) (instance *DataMoverJob, err error) {
	instance = new(DataMoverJob)
	instance.root = root
	instance.events = events

	err = instance.init()

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverJob) IsScheduled() bool {
	return nil != instance && nil != instance.scheduler
}

func (instance *DataMoverJob) IsRunning() bool {
	return nil != instance && nil != instance.scheduler && !instance.scheduler.IsPaused()
}

func (instance *DataMoverJob) Start() (err error) {
	if nil != instance.scheduler {
		instance.scheduler.Start()
	}
	return
}

// Stop Try to close gracefully
func (instance *DataMoverJob) Stop() (err error) {
	if nil != instance.scheduler {
		instance.scheduler.Stop()
	}
	return
}

func (instance *DataMoverJob) Run() (err error) {
	if nil != instance {
		err = instance.run()
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverJob) init() error {
	instance.name = path.Base(instance.root)

	// job logger
	loggerfile := gg.Paths.Concat(instance.root, "logging.log")
	_ = gg.IO.Remove(loggerfile)
	instance.logger = gg_log.NewLogger()
	instance.logger.SetFileName(loggerfile)

	// lookup settings
	text, err := gg.IO.ReadTextFromFile(gg.Paths.Concat(instance.root, "job.json"))
	if nil != err {
		return err
	}
	err = gg.JSON.Read(text, &instance.settings)
	if nil != err {
		return err
	}

	// scheduler
	_ = instance.initScheduler()

	return nil
}

func (instance *DataMoverJob) initScheduler() bool {
	if nil != instance && nil != instance.settings && nil != instance.settings.Schedule &&
		len(instance.settings.Schedule.Timeline) > 0 {
		// valid schedule
		schedule := instance.settings.Schedule
		instance.scheduler = gg_scheduler.NewScheduler()
		instance.scheduler.SetAsync(true) // sync messages
		instance.scheduler.AddSchedule(&gg_scheduler.Schedule{
			Uid:       schedule.Uid,
			StartAt:   schedule.StartAt,
			Timeline:  schedule.Timeline,
			Payload:   schedule.Payload,
			Arguments: schedule.Arguments,
		})
		instance.scheduler.OnSchedule(func(schedule *gg_scheduler.SchedulerTask) {
			instance.scheduler.Pause()
			defer instance.scheduler.Resume()
			_ = instance.run()
		})
		return true
	}
	return false
}

func (instance *DataMoverJob) run() error {
	if nil != instance {
		// execute current job

		// run next
		if len(instance.settings.NextRun) > 0 {
			instance.events.Emit(datamover_commons.EventOnNextJobRun, instance.settings.NextRun, instance)
		}
	}
	return nil
}
