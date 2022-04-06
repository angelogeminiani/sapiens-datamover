package sapiens_datamover

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_ticker"
	"bitbucket.org/digi-sense/gg-progr-datamover/sapiens_datamover/datamover_commons"
	"fmt"
	"sync"
	"time"
)

const EventStop = datamover_commons.EventOnDoStop

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e
// ---------------------------------------------------------------------------------------------------------------------

type stopMonitor struct {
	roots      []string // where stop file is stored
	stopCmd    string
	events     *gg_events.Emitter
	fileMux    sync.Mutex
	stopTicker *gg_ticker.Ticker
}

// ---------------------------------------------------------------------------------------------------------------------
//		c o n s t r u c t o r
// ---------------------------------------------------------------------------------------------------------------------

func newStopMonitor(roots []string, stopCmd string, events *gg_events.Emitter) *stopMonitor {
	instance := new(stopMonitor)
	instance.roots = roots
	instance.stopCmd = stopCmd
	instance.events = events

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *stopMonitor) Start() {
	if nil != instance && len(instance.stopCmd) > 0 && nil == instance.stopTicker {
		instance.stopTicker = gg_ticker.NewTicker(1*time.Second, func(t *gg_ticker.Ticker) {
			instance.checkStop()
			// instance.logger.Debug("Checking for stop command....")
		})
		instance.stopTicker.Start()
	}
}

func (instance *stopMonitor) Stop() {
	if nil != instance && nil != instance.stopTicker {
		instance.stopTicker.Stop()
		instance.stopTicker = nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *stopMonitor) checkStop() {
	if nil != instance {
		if len(instance.stopCmd) > 0 {
			instance.fileMux.Lock()
			defer instance.fileMux.Unlock()

			cmd := instance.stopCmd

			// check if file exists
			for _, root := range instance.roots {
				cmdFile := gg.Paths.Concat(root, cmd)
				if b, _ := gg.Paths.Exists(cmdFile); b {
					_ = gg.IO.Remove(cmdFile)
					instance.events.EmitAsync(EventStop)
				}
			}
		}
		instance.tick()
	}
}

func (instance *stopMonitor) tick() {
	root := gg.Paths.WorkspacePath("")
	filename := gg.Paths.Concat(root, "datamover.tick")
	_, _ = gg.IO.WriteTextToFile(fmt.Sprintf("%v", time.Now()), filename)
}
