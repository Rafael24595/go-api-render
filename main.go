package main

import (
	"log"
	"github.com/Rafael24595/go-api-render/src/infrastructure"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

func main() {
	router := router.NewRouter()
	infrastructure.NewController(router)
	log.Fatalln(router.Listen(":8080"))
}