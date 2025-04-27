package main

import (
	"fmt"
	"log"

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
		container.ManagerHistoric)
	port := fmt.Sprintf(":%d", config.Port())
	log.Fatalln(router.Listen(port))
}
