package datamover_commons

import "bitbucket.org/digi-sense/gg-core/gg_scheduler"

type DataMoverSettingsJob struct {
	Schedule *DataMoverScheduleSettings `json:"schedule"`
	NextRun  string                     `json:"next_run"` // name of job to run next
}

type DataMoverScheduleSettings struct {
	gg_scheduler.Schedule
}
