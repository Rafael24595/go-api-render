package main

import (
	"log"

	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/commons/dependency"
	"github.com/Rafael24595/go-api-render/src/infrastructure/controller"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

func main() {
	configuration.Initialize(configuration.ReadEnv(".env"))
	container := dependency.Initialize()
	router := router.NewRouter()
	controller.NewController(router, container.RepositoryActions)
	log.Fatalln(router.Listen(":8080"))
}