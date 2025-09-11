package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs/swagger"

	web_log "github.com/Rafael24595/go-api-render/src/commons/log"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/commons"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/controller"
)

func main() {
	config, container := commons.Initialize()

	webLog := web_log.NewWebLog()

	route := router.NewRouter()
	route.Logger(webLog)

	if config.Dev() {
		route = addOAPIViewer(config, route, webLog)
	}

	controller.NewController(route,
		container.ManagerRequest,
		container.ManagerContext,
		container.ManagerCollection,
		container.ManagerHistoric,
		container.ManagerGroup)

	go listen(config, route)

	<-config.Signal.Done()

	time.Sleep(1 * time.Second)

	for i := 3; i > 0; i-- {
		log.Messagef("Shutdown in %d...", i)
		time.Sleep(1 * time.Second)
	}

	log.Message("Exiting app.")

	time.Sleep(1 * time.Second)
}

func addOAPIViewer(config *configuration.Configuration, route *router.Router, webLog web_log.WebLog) *router.Router {
	options := swagger.OpenAPI3ViewerOptions{
		Version:   config.Project.Version,
		EnableTLS: config.EnableTLS(),
		OnlyTLS:   config.OnlyTLS(),
		Port:      config.Port(),
		PortTLS:   config.PortTLS(),
		FileYML:   "swagger.yaml",
	}

	viewer := swagger.NewViewer()
	viewer.Logger(webLog)
	viewer.Load(options)

	return route.DocViewer(viewer)
}

func listen(config *configuration.Configuration, route *router.Router) {
	err := serve(config, route)
	if err == nil {
		return
	}

	log.Errorf("Server exited with error: %v", err)
	time.Sleep(3 * time.Second)
	os.Exit(1)
}

func serve(config *configuration.Configuration, route *router.Router) error {
	port := fmt.Sprintf(":%d", config.Port())
	if !config.EnableTLS() {
		return route.Listen(port)
	}

	portTLS := fmt.Sprintf(":%d", config.PortTLS())
	if config.OnlyTLS() {
		return route.ListenTLS(portTLS, config.CertTLS(), config.KeyTLS())
	}

	return route.ListenWithTLS(port, portTLS, config.CertTLS(), config.KeyTLS())
}
