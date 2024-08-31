package infrastructure

import (
	"fmt"
	"go-api-render/src/infrastructure/router"
	"go-api-render/src/infrastructure/router/templates"
	"net/http"
)

type Controller interface {
}

type controller struct {
	router *router.Router
	manager templates.TemplateManager
}

func NewController() Controller {
	builder := templates.NewBuilder().
		AddFunction("SayHello", func(name string)string{return fmt.Sprintf("Hello %s!", name)}).
		AddPath("templates")

	instance := controller{
		router: router.NewRouter(),
		manager: builder.Make(),
	}

	instance.router.ResourcesPath("templates").
		Route(http.MethodGet, "/", instance.home)

	return instance
}

func (controller controller) home(w http.ResponseWriter, r *router.Request) {
	data := map[string]interface{}{}
	controller.manager.Render(w, "home.html", data)
}