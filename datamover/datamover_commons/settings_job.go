package datamover_commons

import (
	"bitbucket.org/digi-sense/gg-core/gg_scheduler"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/schema"
)

type DataMoverSettingsJob struct {
	Schedule    *DataMoverScheduleSettings     `json:"schedule"`
	NextRun     string                         `json:"next_run"` // name of job to run next
	Transaction []*DataMoverDatasourceSettings `json:"transaction"`
}

type DataMoverScheduleSettings struct {
	gg_scheduler.Schedule
}

type DataMoverDatasourceSettings struct {
	Uid         string                       `json:"uid"`
	Description string                       `json:"description"`
	Connection  *DataMoverConnectionSettings `json:"connection"`
	Command     string                       `json:"command"`
	Script      string                       `json:"script"`
}

type DataMoverConnectionSettings struct {
	Driver string                            `json:"driver"`
	Dsn    string                            `json:"dsn"`
	Schema *schema.DataMoverDatasourceSchema `json:"schema"`
}
