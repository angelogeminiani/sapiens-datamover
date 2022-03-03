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
	for _, transaction := range instance.Transaction {
		transaction.ScriptsReset()
	}
	text := gg.JSON.Stringify(instance)
	_, err = gg.IO.WriteTextToFile(text, filename)
	return
}

type DataMoverScheduleSettings struct {
	gg_scheduler.Schedule
}

type DataMoverActionSettings struct {
	Uid           string                         `json:"uid"`
	Description   string                         `json:"description"`
	Network       *DataMoverNetworkSettings      `json:"network"`
	Connection    *DataMoverConnectionSettings   `json:"connection"`
	Command       string                         `json:"command"`
	FieldsMapping map[string]interface{}         `json:"fields_mapping"` // optional: use only if source dataset is different from target
	Scripts       *DataMoverActionScriptSettings `json:"scripts"`
}

func (instance *DataMoverActionSettings) ScriptsLoad(root string) (err error) {
	if nil != instance.Scripts {

		// BEFORE
		if len(instance.Scripts.Before) > 0 && strings.Index(instance.Scripts.Before, ".") == 0 {
			instance.Scripts.BeforeFile = instance.Scripts.Before
			t, e := gg.IO.ReadTextFromFile(gg.Paths.Concat(root, instance.Scripts.BeforeFile))
			if nil != e {
				instance.Scripts.Before = ""
				err = e
			} else {
				instance.Scripts.Before = t
			}
		}
		// CONTEXT
		if len(instance.Scripts.Context) > 0 && strings.Index(instance.Scripts.Context, ".") == 0 {
			instance.Scripts.ContextFile = instance.Scripts.Context
			t, e := gg.IO.ReadTextFromFile(gg.Paths.Concat(root, instance.Scripts.ContextFile))
			if nil != e {
				instance.Scripts.Context = ""
				err = e
			} else {
				instance.Scripts.Context = t
			}
		}
		// AFTER
		if len(instance.Scripts.After) > 0 && strings.Index(instance.Scripts.After, ".") == 0 {
			instance.Scripts.AfterFile = instance.Scripts.After
			t, e := gg.IO.ReadTextFromFile(gg.Paths.Concat(root, instance.Scripts.AfterFile))
			if nil != e {
				instance.Scripts.After = ""
				err = e
			} else {
				instance.Scripts.After = t
			}
		}
	}
	return
}

func (instance *DataMoverActionSettings) ScriptsReset() {
	if nil != instance && nil != instance.Scripts {
		if len(instance.Scripts.BeforeFile) > 0 {
			instance.Scripts.Before = instance.Scripts.BeforeFile
		}
		if len(instance.Scripts.ContextFile) > 0 {
			instance.Scripts.Context = instance.Scripts.ContextFile
		}
		if len(instance.Scripts.AfterFile) > 0 {
			instance.Scripts.After = instance.Scripts.AfterFile
		}
	}
}

type DataMoverActionScriptSettings struct {
	Context     string `json:"context"` // a script to run to change context
	Before      string `json:"before"`  // before query
	After       string `json:"after"`   // after query
	ContextFile string `json:"-"`       // a script to run to change context
	BeforeFile  string `json:"-"`       // before query
	AfterFile   string `json:"-"`       // after query
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
