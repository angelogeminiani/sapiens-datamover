package datamover_globals

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_commons"
)

type DataMoverGlobalsSettings struct {
	Connections []*datamover_commons.DataMoverConnectionSettings `json:"connections"`
	Constants   map[string]interface{}                           `json:"constants"`
}

func (instance *DataMoverGlobalsSettings) String() string {
	return gg.JSON.Stringify(instance)
}
