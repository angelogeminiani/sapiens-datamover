package datamover_commons

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_scheduler"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_jobs/action/schema"
	"net/url"
	"strings"
)

type DataMoverSettingsJob struct {
	Schedule    *DataMoverScheduleSettings `json:"schedule"`
	NextRun     string                     `json:"next_run"` // name of job to run next
	Transaction []*DataMoverActionSettings `json:"transaction"`
	Variables   map[string]interface{}     `json:"variables"`
}

func (instance *DataMoverSettingsJob) SaveToFile(filename string) (err error) {
	text := gg.JSON.Stringify(instance)
	_, err = gg.IO.WriteTextToFile(text, filename)
	return
}

type DataMoverScheduleSettings struct {
	gg_scheduler.Schedule
}

type DataMoverActionSettings struct {
	Uid         string                         `json:"uid"`
	Description string                         `json:"description"`
	Network     *DataMoverNetworkSettings      `json:"network"`
	Connection  *DataMoverConnectionSettings   `json:"connection"`
	Command     string                         `json:"command"`
	Scripts     *DataMoverActionScriptSettings `json:"scripts"`
}

func (instance *DataMoverActionSettings) NormalizeScripts(root string) (err error) {
	if nil != instance.Scripts {

		// CONTEXT
		if len(instance.Scripts.Context) > 0 && strings.Index(instance.Scripts.Context, ".") == 0 {
			t, e := gg.IO.ReadTextFromFile(gg.Paths.Concat(root, instance.Scripts.Context))
			if nil != e {
				instance.Scripts.Context = ""
				err = e
			} else {
				instance.Scripts.Context = t
			}
		}
	}
	return
}

type DataMoverActionScriptSettings struct {
	Context string `json:"context"` // a script to run to change context
}

type DataMoverNetworkSettings struct {
	Host           string                    `json:"host"`
	Authentication *SettingsNetAuthorization `json:"authorization"`
	Secure         bool                      `json:"secure"`
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
