package webserver

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core/gg_events"
	"bitbucket.org/digi-sense/gg-core/gg_utils"
	"bitbucket.org/digi-sense/gg-progr-datamover/datamover/datamover_commons"
	"github.com/gofiber/fiber/v2"
)

type WebController struct {
	root     string
	logger   *datamover_commons.Logger
	events   *gg_events.Emitter
	settings *WebserverSettings

	webserver *Webserver
	websecure *Websecure
}

func NewWebController(mode string, logger *datamover_commons.Logger, events *gg_events.Emitter, settings *WebserverSettings) (instance *WebController) {
	root := gg.Paths.WorkspacePath("./webserver")
	_ = gg.Paths.Mkdir(root + gg_utils.OS_PATH_SEPARATOR)

	instance = new(WebController)
	instance.logger = logger
	instance.events = events
	instance.root = root
	instance.settings = settings

	instance.init(mode, root)

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *WebController) Start() bool {
	if nil != instance {
		instance.webserver.Start()
		_ = instance.websecure.Start()
		return true
	}
	return false
}

func (instance *WebController) Stop() {
	if nil != instance {
		instance.webserver.Stop()
		instance.websecure.Stop()
	}
}

// Handle expose handle method to add more
func (instance *WebController) Handle(method, endpoint string, handler fiber.Handler) {
	if nil != instance && nil != instance.webserver {
		instance.webserver.Handle(method, endpoint, handler)
	}
}

// RegisterNoAuth register command with no http auth
func (instance *WebController) RegisterNoAuth(method, endpoint, command string) {
	if nil != instance && nil != instance.webserver {
		// uid := strings.ToLower(fmt.Sprintf("%s|%s", method, endpoint))
		// instance.externalCommands[uid] = command
		// instance.webserver.Handle(method, endpoint, instance.onNoAuthHandler)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *WebController) init(mode, root string) {

	instance.webserver = NewWebserver("webserver", mode, root)
	if nil != instance.webserver {
		instance.websecure = NewWebsecure(mode, instance.webserver.HttpAuth())

		// API SYS
		instance.webserver.Handle("get", ApiSysVersion, instance.onSysVersion)

	}
}

/** **/
func (instance *WebController) onSysVersion(ctx *fiber.Ctx) error {
	// no auth
	return WriteResponse(ctx, datamover_commons.AppVersion, nil)
}
