package action

import (
	"bitbucket.org/digi-sense/gg-core-x/gg_log"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
)

type DataMoverAction struct {
	root     string
	logger   *gg_log.Logger
	events   *gg_events.Emitter
	settings *datamover_commons.DataMoverActionSettings

	source *DataMoverDatasource
	target *DataMoverDatasource
}

func NewDataMoverAction(root string, logger *gg_log.Logger, events *gg_events.Emitter, settings *datamover_commons.DataMoverActionSettings) *DataMoverAction {
	instance := new(DataMoverAction)
	instance.root = root
	instance.logger = logger
	instance.events = events
	instance.settings = settings

	instance.init()

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverAction) Run() (interface{}, error) {
	if nil != instance {

	}
	return nil, nil
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DataMoverAction) init() {
	instance.source = NewDataMoverDatasource("source", instance.settings.Source)
	instance.target = NewDataMoverDatasource("target", instance.settings.Target)
}
