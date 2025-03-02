package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
	"github.com/Rafael24595/go-collections/collection"
)

type Controller struct {
	router  *router.Router
	manager templates.TemplateManager
}

func NewController(router *router.Router, repository *repository.RequestManager) Controller {
	instance := Controller{
		router:  router,
		manager: templates.NewBuilder().Make(),
	}

	router.
		Contextualizer(instance.contextualizer).
		ErrorHandler(instance.error)

	NewControllerApiClient(router, repository)

	return instance
}

func (c *Controller) contextualizer(w http.ResponseWriter, r *http.Request) (router.Context, error) {
	return collection.DictionaryEmpty[string, any](), nil
}

func (c *Controller) error(w http.ResponseWriter, r *http.Request, context router.Context, err error) {
	w.WriteHeader(http.StatusInternalServerError)
}
