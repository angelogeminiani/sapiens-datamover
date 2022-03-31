package datamover_globals

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"fmt"
	"strings"
)

type Globals struct {
	settings *DataMoverGlobalsSettings
}

func NewGlobals(mode string) (instance *Globals) {
	instance = new(Globals)

	instance.init(mode)

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Globals) LoadJson(json string) {
	instance.settings = new(DataMoverGlobalsSettings)
	_ = gg.JSON.Read(json, instance.settings)
}

func (instance *Globals) Clone() *Globals {
	clone := new(Globals)
	clone.LoadJson(instance.settings.String())
	return clone
}

func (instance *Globals) HasConnections() bool {
	if nil != instance && nil != instance.settings {
		return len(instance.settings.Connections) > 0
	}
	return false
}

func (instance *Globals) GetConnection(id string) *datamover_commons.DataMoverConnectionSettings {
	if nil != instance && nil != instance.settings {
		for _, conn := range instance.settings.Connections {
			if strings.ToLower(id) == strings.ToLower(conn.ConnectionsId) {
				return conn
			}
		}
	}
	return nil
}

func (instance *Globals) MergeVariables(variables map[string]interface{}) map[string]interface{} {
	if nil != instance && nil != instance.settings {
		if len(instance.settings.Constants) > 0 {
			// add constants to variables
			return gg.Maps.Merge(false, variables, instance.settings.Constants)
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
		_ = gg.JSON.ReadFromFile(filename, &instance.settings)
	}
	if nil == instance.settings {
		instance.settings = new(DataMoverGlobalsSettings)
	}
}
