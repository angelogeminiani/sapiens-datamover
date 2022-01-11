package datamover_commons

import (
	"bitbucket.org/digi-sense/gg-core/gg_scheduler"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/schema"
	"net/url"
)

type DataMoverSettingsJob struct {
	Schedule    *DataMoverScheduleSettings `json:"schedule"`
	NextRun     string                     `json:"next_run"` // name of job to run next
	Transaction []*DataMoverActionSettings `json:"transaction"`
}

type DataMoverScheduleSettings struct {
	gg_scheduler.Schedule
}

type DataMoverActionSettings struct {
	Uid         string                       `json:"uid"`
	Description string                       `json:"description"`
	Network     *DataMoverNetworkSettings    `json:"network"`
	Connection  *DataMoverConnectionSettings `json:"connection"`
	Command     string                       `json:"command"`
	Script      string                       `json:"script"`
}

type DataMoverNetworkSettings struct {
	Host           string                    `json:"host"`
	Authentication *SettingsNetAuthorization `json:"authorization"`
}

type DataMoverConnectionSettings struct {
	Driver string                            `json:"driver"`
	Dsn    string                            `json:"dsn"`
	Schema *schema.DataMoverDatasourceSchema `json:"schema"`
}

func (instance *DataMoverNetworkSettings) Uri() (uri *url.URL, err error) {
	text := instance.Host // nio://127.0.0.1:10001
	uri, err = url.Parse(text)
	return
}
