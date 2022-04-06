package webserver

import (
	"bitbucket.org/digi-sense/gg-core"
	"bitbucket.org/digi-sense/gg-core-x/gg_http/httpserver"
	"bitbucket.org/digi-sense/gg-core/gg_utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------
//		t y p e
// ---------------------------------------------------------------------------------------------------------------------

type Webserver struct {
	dirWork     string // workspace
	initialized bool
	enabled     bool

	settings       *WebserverSettings
	httpserver     *httpserver.HttpServer
	httpRoot       string
	httpStaticRoot string
	httpAddr       string
	httpsAddr      string
}

func NewWebserver(httpRoot string, settings *WebserverSettings) *Webserver {
	instance := new(Webserver)
	instance.dirWork = gg.Paths.WorkspacePath("./")
	instance.httpRoot = httpRoot
	instance.enabled = false

	_ = instance.init(settings)

	return instance
}

// ---------------------------------------------------------------------------------------------------------------------
//		p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Webserver) HttpRoot() string {
	if nil != instance {
		return instance.httpRoot
	}
	return ""
}

func (instance *Webserver) HttpPath(path string) string {
	if nil != instance {
		return gg.Paths.Concat(instance.httpRoot, path)
	}
	return ""
}

func (instance *Webserver) HttpStaticRoot() string {
	if nil != instance {
		return instance.httpStaticRoot
	}
	return ""
}

func (instance *Webserver) HttpAddress() string {
	if nil != instance {
		return instance.httpAddr
	}
	return ""
}

func (instance *Webserver) HttpAuth() *Authorization {
	if nil != instance && nil != instance.settings {
		return instance.settings.Auth
	}
	return nil
}

func (instance *Webserver) LocalUrl() string {
	if nil != instance && instance.enabled {
		if len(instance.httpAddr) > 0 {
			return fmt.Sprintf("http://localhost%v/", instance.httpAddr)
		} else if len(instance.httpsAddr) > 0 {
			return fmt.Sprintf("https://localhost%v/", instance.httpsAddr)
		}
	}
	return ""
}

func (instance *Webserver) IsHttps() bool {
	return len(instance.httpsAddr) > 0
}

func (instance *Webserver) Settings() *WebserverSettings {
	if nil != instance {
		return instance.settings
	}
	return nil
}

func (instance *Webserver) Handle(method, endpoint string, handler fiber.Handler) {
	switch strings.ToLower(method) {
	case "get":
		instance.httpserver.Get(endpoint, handler)
	case "post":
		instance.httpserver.Post(endpoint, handler)
	case "put":
		instance.httpserver.Put(endpoint, handler)
	case "delete":
		instance.httpserver.Delete(endpoint, handler)
	case "middleware":
		instance.httpserver.Middleware(endpoint, handler)
	default:
		instance.httpserver.All(endpoint, handler)
	}
}

func (instance *Webserver) IsEnabled() bool {
	if nil != instance {
		return instance.enabled
	}
	return false
}

func (instance *Webserver) Start() bool {
	if nil != instance && instance.enabled {
		instance.start()

		return true
	}
	return false
}

func (instance *Webserver) Stop() {
	if nil != instance && instance.enabled {
		_ = instance.httpserver.Stop()
	}
}

func (instance *Webserver) Exit() {
	defer func() {
		if r := recover(); r != nil {
			// recover from panic if any
		}
	}()
	if nil != instance && instance.enabled {
		go func() {
			_ = instance.httpserver.Stop()
		}()
		// wait a while the server close
		time.Sleep(3 * time.Second)
		if instance.httpserver.IsOpen() {
			// brute force close
			os.Exit(0)
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//		p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *Webserver) init(settings *WebserverSettings) (err error) {
	if !instance.initialized {
		instance.initialized = true

		if nil != settings {
			instance.settings = settings
			instance.enabled = settings.Enabled
			if instance.enabled {

				// create webserver
				err = instance.initHttp()
				if nil != err {
					instance.enabled = false
					return err
				}
			}
		} else {
			instance.enabled = false
		}
	}
	return
}

func (instance *Webserver) initHttp() error {
	// create webserver instance
	httpConfig := instance.settings.Http

	instance.httpserver = httpserver.NewHttpServer(instance.httpRoot, instance.handleHttpError, instance.handleHttpLimit)
	err := instance.httpserver.ConfigureFromMap(httpConfig)
	if nil != err {
		return err
	}
	// check defaults
	configuration := instance.httpserver.Configuration()
	if nil != configuration {
		if nil != configuration.Server {
			if configuration.Server.ReadTimeout < 3*time.Second {
				configuration.Server.ReadTimeout = 3 * time.Second
			}
			if configuration.Server.WriteTimeout < 3*time.Second {
				configuration.Server.WriteTimeout = 3 * time.Second
			}
		}
	}

	// parse settings
	hosts := gg.Reflect.GetArray(httpConfig, "hosts")
	for _, host := range hosts {
		isSecure := gg.Reflect.GetBool(host, "tls")
		if isSecure {
			instance.httpsAddr = gg.Reflect.GetString(host, "addr")
			sslCert := gg.Reflect.GetString(host, "ssl_cert")
			sslKey := gg.Reflect.GetString(host, "ssl_key")
			if len(sslKey) > 0 {
				instance.mkDir(sslKey)
			}
			if len(sslCert) > 0 {
				instance.mkDir(sslCert)
			}
		} else {
			instance.httpAddr = gg.Reflect.GetString(host, "addr")
		}
	}
	static := gg.Reflect.GetArray(httpConfig, "static")
	if len(static) > 0 {
		root := gg.Reflect.GetString(static[0], "root")
		instance.httpStaticRoot = gg.Paths.Concat(instance.httpRoot, root)
		instance.mkDir(instance.httpStaticRoot + gg_utils.OS_PATH_SEPARATOR)
	}

	return nil
}

func (instance *Webserver) start() {
	// wait a while before start to allow runtime is ready
	time.Sleep(2 * time.Second)
	instance.httpserver.Start()
}

func (instance *Webserver) handleHttpError(serverError *httpserver.HttpServerError) {

}

func (instance *Webserver) handleHttpLimit(ctx *fiber.Ctx) error {
	return nil
}

func (instance *Webserver) mkDir(path string) {
	if gg.Paths.IsAbs(path) {
		_ = gg.Paths.Mkdir(path)
	} else {
		_ = gg.Paths.Mkdir(gg.Paths.Concat(instance.httpRoot, path))
	}
}
