package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
)

const (
	
)

type ControllerCollection struct {
	router              *router.Router
	manager             templates.TemplateManager
}

func NewControllerCollection(router *router.Router, builder *templates.BuilderManager) ControllerCollection {
	instance := ControllerCollection{
		router:              router,
		manager:             builder.Make(),
	}

	instance.router.
		Route(http.MethodGet, instance.collection, "/collection")

	return instance
}

func (c *ControllerCollection) collection(w http.ResponseWriter, r *http.Request, context router.Context) error {
	context.Merge(map[string]any{
		"BodyTemplate": "collection.html",
		"VariableTypes": domain.VariableTypes(),
	}).PutIfAbsent("Collection", domain.NewCollectionDefault())

	return c.manager.Render(w, "home.html", context)
}
