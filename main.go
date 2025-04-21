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
		container.ManagerActions,
		container.ManagerContext,
		container.ManagerCollection,
		container.RepositoryContext,
		container.RepositoryHistoric)
	port := fmt.Sprintf(":%d", config.Port())
	log.Fatalln(router.Listen(port))
}
