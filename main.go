package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/commons"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/controller"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

func main() {
	config, container := commons.Initialize()
	router := router.NewRouter()
	controller.NewController(router,
		container.ManagerRequest,
		container.ManagerContext,
		container.ManagerCollection,
		container.ManagerHistoric,
		container.ManagerGroup)

	go listen(config, router)

	<-config.Signal.Done()

	time.Sleep(1 * time.Second)

	for i := 3; i > 0; i-- {
		log.Messagef("Shutdown in %d...", i)
		time.Sleep(1 * time.Second)
	}

	log.Message("Exiting app.")

	time.Sleep(1 * time.Second)
}

func listen(config *configuration.Configuration, router *router.Router) {
	port := fmt.Sprintf(":%d", config.Port())
	
	var err error
	if config.EnableTLS() {
		portTLS := fmt.Sprintf(":%d", config.PortTLS())
		if config.OnlyTLS() {
			err = router.ListenTLS(portTLS, config.CertTLS(), config.KeyTLS())
		} else {
			err = router.ListenWithTLS(port, portTLS, config.CertTLS(), config.KeyTLS())
		}
	} else {
		err = router.Listen(port)
	}

	if err == nil {
		return
	}

	log.Errorf("Server exited with error: %v", err)
	time.Sleep(3 * time.Second)
	os.Exit(1)
}
