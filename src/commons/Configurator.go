package commons

import (
	"context"
	"encoding/json"
	"os"

	core_commons "github.com/Rafael24595/go-api-core/src/commons"
	
	"github.com/Rafael24595/go-api-core/src/commons/local"
	"github.com/Rafael24595/go-log/log"

	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/commons/dependency"
)

func Initialize(ctx context.Context) (*configuration.Configuration, *dependency.DependencyContainer) {
	kargs := core_commons.ReadAllEnv(".env")
	core_conf, core_cont := core_commons.Initialize(ctx, kargs)
	frontPackage := ReadFrontPackage()

	config := configuration.Initialize(core_conf, kargs, frontPackage)
	container := dependency.Initialize(config, *core_cont)

	log.Messagef("Display front: %v", config.Front.Enabled)

	return &config, container
}

func ReadFrontPackage() *configuration.FrontPackage {
	file, err := os.Open("assets/front/package.json")
	if err != nil {
		log.Errorf("Error opening package.json: %v", err)
		return &configuration.FrontPackage{
			Version: "",
			Name:    "",
		}
	}

	var pkg configuration.FrontPackage
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pkg); err != nil {
		local.Panicf("Error decoding JSON: %v", err)
	}

	if err := file.Close(); err != nil {
		local.Panicf("Error closing file: %v", err)
	}

	return &pkg
}
