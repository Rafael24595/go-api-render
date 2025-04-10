package commons

import (
	core_commons "github.com/Rafael24595/go-api-core/src/commons"
	"github.com/Rafael24595/go-api-core/src/commons/dependency"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
)

func Initialize() (*configuration.Configuration, *dependency.DependencyContainer) {
	kargs := configuration.ReadEnv(".env")
	_, container := core_commons.Initialize(kargs)
	config := configuration.Initialize(kargs)
	return &config, container
}

