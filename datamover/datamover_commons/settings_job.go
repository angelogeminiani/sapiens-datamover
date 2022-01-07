package datamover_commons

import "bitbucket.org/digi-sense/gg-core/gg_scheduler"

type DataMoverSettingsJob struct {
	Schedule *DataMoverScheduleSettings `json:"schedule"`
	NextRun  string                     `json:"next_run"` // name of job to run next
	Action   *DataMoverActionSettings   `json:"action"`
}

type DataMoverScheduleSettings struct {
	gg_scheduler.Schedule
}

type DataMoverActionSettings struct {
	Source *DataMoverDatasourceSettings `json:"source"`
	Target *DataMoverDatasourceSettings `json:"target"`
}

type DataMoverDatasourceSettings struct {
	Connection *DataMoverConnectionSettings `json:"connection"`
	Command    string                       `json:"command"`
	Script     string                       `json:"script"`
}

type DataMoverConnectionSettings struct {
	Driver string `json:"driver"`
	Dsn    string `json:"dsn"`
}
