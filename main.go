package main

import (
	"context"
	"os"
	"time"

	"github.com/Rafael24595/go-api-render/src/commons"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/controller"
	"github.com/Rafael24595/go-log/log"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs/swagger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	config, container := commons.Initialize(ctx)

	route := router.NewRouter()
	route.Logger(log.Default())

	if config.Dev() {
		route = addOAPIViewer(config, route)
	}

	controller.NewController(route,
		container.ManagerRequest,
		container.ManagerContext,
		container.ManagerCollection,
		container.ManagerHistoric,
		container.ManagerGroup,
		container.ManagerEndPoint,
		container.ManagerMetrics,
		container.ManagerToken,
		container.ManagerSessionData,
		container.ManagerWeb)

	go listen(config, route)

	<-config.Signal.Done()

	time.Sleep(1 * time.Second)

	for i := 3; i > 0; i-- {
		log.Messagef("Shutdown in %d...", i)
		time.Sleep(1 * time.Second)
	}

	log.Message("Exiting app.")

	cancel()

	time.Sleep(1 * time.Second)
}

func addOAPIViewer(config *configuration.Configuration, route *router.Router) *router.Router {
	options := swagger.OpenAPI3ViewerOptions{
		Version:   config.Project.Version,
		EnableTLS: config.EnableTLS(),
		OnlyTLS:   config.OnlyTLS(),
		Port:      config.Port(),
		PortTLS:   config.PortTLS(),
		FileYML:   "swagger.yaml",
	}

	viewer := swagger.NewViewer()
	viewer.Logger(log.Default())
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
	if !config.EnableTLS() {
		return route.Listen(config.Port())
	}

	if config.OnlyTLS() {
		return route.ListenTLS(config.PortTLS(), config.CertTLS(), config.KeyTLS())
	}

	return route.ListenWithTLS(config.Port(), config.PortTLS(), config.CertTLS(), config.KeyTLS())
}
