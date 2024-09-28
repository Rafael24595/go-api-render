package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-core/src/domain"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
)

type Controller interface {
}

type controller struct {
	router            *router.Router
	manager           templates.TemplateManager
	queryRepository   request.QueryRepository
	commandRepository request.CommandRepository
}

func NewController(router *router.Router, queryRepository request.QueryRepository, commandRepository request.CommandRepository) Controller {
	builder := templates.NewBuilder().
		AddFunction("SayHello", func(name string) string { return fmt.Sprintf("Hello %s!", name) }).
		AddPath("templates").
		AddPath("templates/**").
		AddPath("templates/**/**")

	instance := controller{
		router:            router,
		manager:           builder.Make(),
		queryRepository:   queryRepository,
		commandRepository: commandRepository,
	}

	instance.router.ResourcesPath("templates").
		Contextualizer(instance.contextualizer).
		ErrorHandler(instance.error).
		Route(http.MethodGet, "/", instance.home).
		Route(http.MethodGet, "/client", instance.client).
		Route(http.MethodPost, "/client", instance.request)

	return instance
}

func (c *controller) contextualizer(w http.ResponseWriter, r *http.Request) (router.Context, error) {
	return collection.EmptyMap[string, any](), nil
}

func (c *controller) home(w http.ResponseWriter, r *http.Request, context router.Context) error {
	return c.manager.Render(w, "home.html", context)
}

func (c *controller) client(w http.ResponseWriter, r *http.Request, context router.Context) error {
	requests := c.queryRepository.FindAll()

	context.Merge(map[string]any{
		"Methods":      domain.HttpMethods(),
		"Requests":     requests,
	})

	return c.manager.Render(w, "client/client.html", context)
}

func (c *controller) request(w http.ResponseWriter, r *http.Request, context router.Context) error {
	request, err := proccessRequest(r)
	if err != nil {
		return err
	}

	response, err := core_infrastructure.Client().Fetch(*request)
	if err != nil {
		return err
	}

	context.Merge(map[string]any{
		"Response": response,
	})

	return c.client(w, r, context)
}

func (c *controller) error(w http.ResponseWriter, r *http.Request, context router.Context, err error) {
	context.Merge(map[string]any{
		"Error": err,
	})
	c.manager.Render(w, "error.html", context)
}
