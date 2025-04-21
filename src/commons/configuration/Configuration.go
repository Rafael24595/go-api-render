package configuration

import (
	core_configuration "github.com/Rafael24595/go-api-core/src/commons/configuration"
	"github.com/Rafael24595/go-api-core/src/commons/utils"
)

var instance *Configuration

type Configuration struct {
	core_configuration.Configuration
	debug bool
	port  int
	front bool
}

func Initialize(core *core_configuration.Configuration, kargs map[string]utils.Any) Configuration {
	if instance != nil {
		panic("")
	}

	debug, ok := kargs["GO_API_DEBUG"].Bool()
	if !ok {
		debug = false
	}

	port, ok := kargs["GO_API_SERVER_PORT"].Int()
	if !ok {
		//TODO: log
		port = 8080
	}

	front, ok := kargs["GO_API_SERVER_FRONT"].Bool()
	if !ok {
		//TODO: log
		front = false
	}

	instance = &Configuration{
		Configuration: *core,
		debug:         debug,
		port:          port,
		front:         front,
	}

	return *instance
}

func Instance() Configuration {
	if instance == nil {
		panic("")
	}
	return *instance
}

func (c Configuration) Debug() bool {
	return c.debug
}

func (c Configuration) Port() int {
	return c.port
}

func (c Configuration) Front() bool {
	return c.front
}
