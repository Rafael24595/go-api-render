package configuration

import (
	core_configuration "github.com/Rafael24595/go-api-core/src/commons/configuration"
	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/commons/utils"
)

const defaultPort = 8080

var instance *Configuration

type Configuration struct {
	core_configuration.Configuration
	Front FrontPackage
	debug bool
	port  int
}

func Initialize(core *core_configuration.Configuration, kargs map[string]utils.Any, frontPackage *FrontPackage) Configuration {
	if instance != nil {
		log.Panics("The configuration is alredy initialized")
	}

	debug, ok := kargs["GO_API_DEBUG"].Bool()
	if !ok {
		debug = false
	}

	port, ok := kargs["GO_API_SERVER_PORT"].Int()
	if !ok {
		log.Messagef("Custom port flag is not defined; using default port %d", defaultPort)
		port = defaultPort
	}

	front, ok := kargs["GO_API_SERVER_FRONT"].Bool()
	if !ok {
		log.Message("Front flag is not defined; the frontend application will not be displayed")
		front = false
	}

	frontPackage.Enabled = front
	if !front {
		frontPackage.Name = ""
		frontPackage.Version = ""
	}

	instance = &Configuration{
		Configuration: *core,
		Front:         *frontPackage,
		debug:         debug,
		port:          port,
	}

	return *instance
}

func Instance() Configuration {
	if instance == nil {
		log.Panics("The configuration is not initialized yet")
	}
	return *instance
}

func (c Configuration) Debug() bool {
	return c.debug
}

func (c Configuration) Port() int {
	return c.port
}
