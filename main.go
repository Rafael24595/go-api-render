package main

import (
	"log"

	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/commons/dependency"
	"github.com/Rafael24595/go-api-render/src/infrastructure"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

func main() {
	configuration.Initialize(configuration.ReadEnv(".env"))
	container := dependency.Initialize()
	router := router.NewRouter()
	infrastructure.NewController(router, container.RequestQueryHistoric, container.RequestQueryPersisted, container.RequestCommandManager)
	log.Fatalln(router.Listen(":8080"))
}