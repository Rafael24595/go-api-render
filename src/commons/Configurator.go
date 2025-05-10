package commons

import (
	core_commons "github.com/Rafael24595/go-api-core/src/commons"
	"github.com/Rafael24595/go-api-core/src/commons/dependency"
	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
)

func Initialize() (*configuration.Configuration, *dependency.DependencyContainer) {
	kargs := core_commons.ReadEnv(".env")
	core, container := core_commons.Initialize(kargs)
	config := configuration.Initialize(core, kargs)

	log.Messagef("Display front: %v", config.Front())

	return &config, container
}

