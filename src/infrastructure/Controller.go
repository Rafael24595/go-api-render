package infrastructure

import (
	"fmt"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-core/src/domain"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
)

type Controller interface {
}

type controller struct {
	router                   *router.Router
	manager                  templates.TemplateManager
	queryRepositoryHistoric  request.RepositoryQuery
	queryRepositoryPersisted request.RepositoryQuery
	commandRepository        *request.MemoryCommandManager
}

func NewController(router *router.Router, queryRepositoryHistoric request.RepositoryQuery, queryRepositoryPersisted request.RepositoryQuery, commandRepository *request.MemoryCommandManager) Controller {
	builder := templates.NewBuilder().
		AddFunction("SayHello", func(name string) string { return fmt.Sprintf("Hello %s!", name) }).
		AddFunctions(map[string]any{
			"ToString": ToString,
			"Uuid":     Uuid,
			"Not":      Not,
			"Concat":   Concat,
			"String":   String,
		}).
		AddPath("templates").
		AddPath("templates/**").
		AddPath("templates/**/**")

	instance := controller{
		router:                   router,
		manager:                  builder.Make(),
		queryRepositoryHistoric:  queryRepositoryHistoric,
		queryRepositoryPersisted: queryRepositoryPersisted,
		commandRepository:        commandRepository,
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
	return collection.FromMap(map[string]any{
		"Constants": configuration.GetConstants(),
	}), nil
}

func (c *controller) home(w http.ResponseWriter, r *http.Request, context router.Context) error {
	return c.manager.Render(w, "home.html", context)
}

func (c *controller) client(w http.ResponseWriter, r *http.Request, context router.Context) error {
	requests := c.queryRepositoryHistoric.FindAll()

	context.Merge(map[string]any{
		"Methods":  domain.HttpMethods(),
		"Requests": requests,
	}).
	PutIfAbsent("Request", domain.NewRequestEmpty())

	return c.manager.Render(w, "client/client.html", context)
}

func (c *controller) request(w http.ResponseWriter, r *http.Request, context router.Context) error {
	constants := configuration.GetConstants()

	request, err := proccessRequestAnonymous(r)
	if err != nil {
		return err
	}

	response, err := core_infrastructure.Client().Fetch(*request)
	if err != nil {
		return err
	}

	clientType := r.Form.Get(constants.Client.Type)
	bodyType := r.Form.Get(constants.Body.Type)
	authStatus := r.Form.Get(constants.Auth.Enabled) == "on"
	authType := r.Form.Get(constants.Auth.Type)

	request = c.commandRepository.Insert(*request)

	context.Merge(map[string]any{
		"Request":    request,
		"Response":   response,
		"ClientType": clientType,
		"AuthStatus": authStatus,
		"AuthType":   authType,
		"BodyType":   bodyType,
	})

	return c.client(w, r, context)
}

func (c *controller) error(w http.ResponseWriter, r *http.Request, context router.Context, err error) {
	context.Merge(map[string]any{
		"Error": err,
	})
	c.manager.Render(w, "error.html", context)
}
