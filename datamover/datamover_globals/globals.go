package datamover_globals

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"fmt"
	"strings"
)

type Globals struct {
	Settings *DataMoverGlobalsSettings `json:"settings"`
}

func NewGlobals(mode string) (instance *Globals) {
	instance = new(Globals)

	instance.init(mode)

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Globals) String() string {
	return gg.JSON.Stringify(&instance)
}

func (instance *Globals) LoadJson(json string) {
	instance.Settings = new(DataMoverGlobalsSettings)
	_ = gg.JSON.Read(json, instance.Settings)
}

func (instance *Globals) Clone() *Globals {
	clone := new(Globals)
	clone.LoadJson(instance.Settings.String())
	return clone
}

func (instance *Globals) HasConnections() bool {
	if nil != instance && nil != instance.Settings {
		return len(instance.Settings.Connections) > 0
	}
	return false
}

func (instance *Globals) GetConnection(id string) *datamover_commons.DataMoverConnectionSettings {
	if nil != instance && nil != instance.Settings {
		for _, conn := range instance.Settings.Connections {
			if strings.ToLower(id) == strings.ToLower(conn.ConnectionsId) {
				return conn
			}
		}
	}
	return nil
}

func (instance *Globals) MergeVariables(variables map[string]interface{}) map[string]interface{} {
	if nil != instance && nil != instance.Settings {
		if len(instance.Settings.Constants) > 0 {
			// add constants to variables
			return gg.Maps.Merge(false, variables, instance.Settings.Constants)
		}
	}
	return variables
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Globals) init(mode string) {
	filename := gg.Paths.WorkspacePath(fmt.Sprintf("globals.%s.json", mode))
	if ok, _ := gg.Paths.Exists(filename); ok {
		_ = gg.JSON.ReadFromFile(filename, &instance.Settings)
	}
	if nil == instance.Settings {
		instance.Settings = new(DataMoverGlobalsSettings)
	}
}
