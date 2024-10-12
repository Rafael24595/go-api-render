package controller

import (
	"fmt"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons"
	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-core/src/domain"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	core_repository "github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
)

const (
	ID_REQUEST = "id_request"
)

type Controller interface {
}

type controller struct {
	router              *router.Router
	manager             templates.TemplateManager
	repositoryHisotric  *repository.RequestManager
	repositoryPersisted *repository.RequestManager
}

func NewController(router *router.Router, repositoryHisotric *repository.RequestManager, repositoryPersisted *repository.RequestManager) Controller {
	builder := templates.NewBuilder().
		AddFunction("SayHello", func(name string) string { return fmt.Sprintf("Hello %s!", name) }).
		AddFunctions(map[string]any{
			"Uuid":                   Uuid,
			"Not":                    Not,
			"Concat":                 Concat,
			"String":                 String,
			"Join":                   Join,
			"FormatMilliseconds":     FormatMilliseconds,
			"FormatMillisecondsDate": FormatMillisecondsDate,
			"FormatBytes":            FormatBytes,
			"ParseCookie":            ParseCookie,
		}).
		AddPath("templates")

	instance := controller{
		router:              router,
		manager:             builder.Make(),
		repositoryHisotric:  repositoryHisotric,
		repositoryPersisted: repositoryPersisted,
	}

	instance.router.ResourcesPath("templates").
		Contextualizer(instance.contextualizer).
		ErrorHandler(instance.error).
		Route(http.MethodGet, instance.home, "/").
		Route(http.MethodGet, instance.client, "/client").
		Route(http.MethodPost, instance.request, "/client").
		Route(http.MethodGet, instance.historic, "/client/{%s}", ID_REQUEST)

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
	requests := c.repositoryHisotric.FindOptions(core_repository.FilterOptions[domain.Request]{
		Sort: func(i, j domain.Request) bool {
			return j.Timestamp > i.Timestamp
		},
		To: 10,
	})

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

	if request.Status == domain.Historic {
		request, response = c.repositoryHisotric.Insert(*request, *response)
	} else {
		request, response = c.repositoryPersisted.Insert(*request, *response)
	}

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

func (c *controller) historic(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)
	request, response, ok := c.repositoryHisotric.Find(idRequest)
	if !ok {
		return commons.ApiErrorFrom(404, "Historic request not found.")
	}

	context.Merge(map[string]any{
		"Request": request,
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