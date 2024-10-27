package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
)

type Controller struct {
	router  *router.Router
	manager templates.TemplateManager
}

func NewController(router *router.Router, repositoryHisotric *repository.RequestManager, repositoryPersisted *repository.RequestManager) Controller {
	builder := templates.NewBuilder().
		AddFunctions(map[string]any{
			"FormatBytes": FormatBytes,
			"ParseCookie": ParseCookie,
			"BodyString":  BodyString,
			"FormatXml":   FormatXml,
			"FormatHtml":  FormatHtml,
			"FormatJson":  FormatJson,
		}).
		AddPath("templates")

	builder.ListManager().
		PutStatic("headers-list", httpHeaders)

	instance := Controller{
		router:  router,
		manager: builder.Make(),
	}

	router.ResourcesPath("templates").
		Contextualizer(instance.contextualizer).
		ErrorHandler(instance.error)

	NewControllerClient(router, builder, repositoryHisotric, repositoryPersisted)
	NewControllerCollection(router, builder)

	return instance
}

func (c *Controller) contextualizer(w http.ResponseWriter, r *http.Request) (router.Context, error) {
	return collection.FromMap(map[string]any{
		"Constants": configuration.GetConstants(),
	}), nil
}

func (c *Controller) error(w http.ResponseWriter, r *http.Request, context router.Context, err error) {
	context.Merge(map[string]any{
		"BodyTemplate": "error.html",
		"Error":        err,
	})

	c.manager.Render(w, "home.html", context)
}
