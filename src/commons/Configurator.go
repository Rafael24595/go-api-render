package commons

import (
	"encoding/json"
	"os"

	core_commons "github.com/Rafael24595/go-api-core/src/commons"
	"github.com/Rafael24595/go-api-core/src/commons/dependency"
	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
)

func Initialize() (*configuration.Configuration, *dependency.DependencyContainer) {
	kargs := core_commons.ReadEnv(".env")
	core, container := core_commons.Initialize(kargs)
	frontPackage := ReadFrontPackage()
	config := configuration.Initialize(core, kargs, frontPackage)

	log.Messagef("Display front: %v", config.Front.Enabled)

	return &config, container
}

func ReadFrontPackage() *configuration.FrontPackage {
	file, err := os.Open("assets/front/package.json")
	if err != nil {
		log.Errorf("Error opening package.json: %v", err)
		return &configuration.FrontPackage{
			Version: "",
			Name: "",
		}
	}

	var pkg configuration.FrontPackage
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&pkg); err != nil {
		log.Panicf("Error decoding JSON: %v", err)
	}

	if err := file.Close(); err !=  nil {
		log.Panicf("Error closing file: %v", err)
	}

	return &pkg
}
