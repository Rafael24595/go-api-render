package main

import (
	"fmt"
	"time"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/commons"
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
	port := fmt.Sprintf(":%d", config.Port())

	go func() {
		if err := router.Listen(port); err != nil {
			log.Errorf("Server exited with error: %v", err)
		}
	}()

	<-config.Signal.Done()
	
	time.Sleep(1 * time.Second)

	for i := 3; i > 0; i-- {
		log.Messagef("Shutdown in %d...", i)
		time.Sleep(1 * time.Second)
	}

	log.Message("Exiting app.")

	time.Sleep(1 * time.Second)
}
