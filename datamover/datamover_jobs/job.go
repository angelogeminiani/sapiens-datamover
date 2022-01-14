package datamover_jobs

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_log"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_scheduler"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action"
	"fmt"
	"path"
)

type DataMoverJob struct {
	isDebug bool
	root    string
	logger  *gg_log.Logger
	events  *gg_events.Emitter

	name         string
	settings     *datamover_commons.DataMoverSettingsJob
	scheduler    *gg_scheduler.Scheduler
	_transaction *action.DataMoverTransaction
}

func NewDataMoverJob(debug bool, root string, events *gg_events.Emitter) (instance *DataMoverJob, err error) {
	instance = new(DataMoverJob)
	instance.isDebug = debug
	instance.root = root
	instance.events = events

	err = instance.init()
	if nil == err {
		instance.logger.Info(fmt.Sprintf("'%s' IS READY!", instance.name))
	} else {
		instance.logger.Info(fmt.Sprintf("'%s' EXIT WITH ERROR: ", err))
	}
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

func (instance *DataMoverJob) Run(contextData []map[string]interface{}, contextVariables map[string]interface{}) (err error) {
	if nil != instance {
		err = instance.run(contextData, contextVariables)
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
	if instance.isDebug {
		instance.logger.SetLevel(gg_log.LEVEL_DEBUG)
	}
	instance.logger.Info(fmt.Sprintf("INITIALIZING '%s'", instance.name))

	// lookup settings
	text, err := gg.IO.ReadTextFromFile(gg.Paths.Concat(instance.root, "job.json"))
	if nil != err {
		return err
	}
	err = gg.JSON.Read(text, &instance.settings)
	if nil != err {
		return err
	}

	if nil != instance.settings {
		if instance.isDebug {
			instance.logger.Info("* settings loaded.")
		}

		// internal variables
		if nil == instance.settings.Variables {
			instance.settings.Variables = make(map[string]interface{})
		}
		/**
		if _, b := instance.settings.Variables[datamover_commons.VarLIMIT]; !b {
			instance.settings.Variables[datamover_commons.VarLIMIT] = 10
		}
		if _, b := instance.settings.Variables[datamover_commons.VarOFFSET]; !b {
			instance.settings.Variables[datamover_commons.VarOFFSET] = 0
		}**/
		instance.logger.Debug(fmt.Sprintf("* Variables: %s", gg.JSON.Stringify(instance.settings.Variables)))

		// scheduler
		if instance.initScheduler() {
			if instance.isDebug {
				instance.logger.Info("* scheduler enabled.")
			}
		} else {
			instance.logger.Warn("* scheduler not enabled.")
		}
	} else {
		instance.logger.Warn("* SETTINGS ARE NOT VALID!")
	}

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
			err := instance.run(nil, nil)
			if nil != err {
				instance.logger.Error(err)
			}
		})
		return true
	}
	return false
}

func (instance *DataMoverJob) transaction() (*action.DataMoverTransaction, error) {
	var err error
	if nil == instance._transaction {
		if nil != instance.settings {
			instance._transaction, err = action.NewDataMoverTransaction(instance.root, instance.logger,
				instance.events, instance.settings.Transaction, instance.settings.Variables)
		} else {
			return nil, gg.Errors.Prefix(datamover_commons.PanicSystemError,
				fmt.Sprintf("Misconfiguration in JOB '%s' settings", instance.name))
		}
	}
	return instance._transaction, err
}

func (instance *DataMoverJob) run(contextData []map[string]interface{}, contextVariables map[string]interface{}) error {
	if nil != instance {
		transaction, err := instance.transaction()
		if nil != err {
			return err
		}

		if nil != transaction {
			// execute current job
			response, variables, err := transaction.Execute(contextData, contextVariables)
			if nil != err {
				return err
			}

			// update settings on file
			instance.incrementVariables(variables)

			// run next
			if len(instance.settings.NextRun) > 0 {
				instance.events.Emit(datamover_commons.EventOnNextJobRun,
					instance.settings.NextRun, instance, response, variables)
			}
		}
	}
	return nil
}

func (instance *DataMoverJob) incrementVariables(variables map[string]interface{}) {
	if len(variables) > 0 {
		var limit, offset interface{}
		var ok bool
		if limit, ok = variables[datamover_commons.VarLIMIT]; ok {
			instance.settings.Variables[datamover_commons.VarLIMIT] = limit
		}
		if offset, ok = variables[datamover_commons.VarOFFSET]; ok {
			instance.settings.Variables[datamover_commons.VarOFFSET] = gg.Convert.ToInt(offset) + gg.Convert.ToInt(limit)
		}

		// replace configuration file
		_ = instance.settings.SaveToFile(gg.Paths.Concat(instance.root, "job.json"))
	}
}
