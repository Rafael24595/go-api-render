package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons"
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

type ControllerClient struct {
	router              *router.Router
	manager             templates.TemplateManager
	repositoryHisotric  *repository.RequestManager
	repositoryPersisted *repository.RequestManager
}

func NewControllerClient(router *router.Router, builder *templates.BuilderManager, repositoryHisotric *repository.RequestManager, repositoryPersisted *repository.RequestManager) ControllerClient {
	instance := ControllerClient{
		router:              router,
		manager:             builder.Make(),
		repositoryHisotric:  repositoryHisotric,
		repositoryPersisted: repositoryPersisted,
	}

	instance.router.
		Route(http.MethodGet, instance.home, "/").
		Route(http.MethodGet, instance.client, "/client").
		Route(http.MethodPost, instance.request, "/client").
		Route(http.MethodGet, instance.historic, "/client/{%s}", ID_REQUEST).
		Route(http.MethodDelete, instance.remove, "/client/{%s}", ID_REQUEST)

	return instance
}

func (c *ControllerClient) home(w http.ResponseWriter, r *http.Request, context router.Context) error {
	return c.client(w ,r, context)
}

func (c *ControllerClient) client(w http.ResponseWriter, r *http.Request, context router.Context) error {
	requests := c.repositoryHisotric.FindOptions(core_repository.FilterOptions[domain.Request]{
		Sort: func(i, j domain.Request) bool {
			return i.Timestamp > j.Timestamp
		},
		To: 10,
	})

	context.Merge(map[string]any{
		"BodyTemplate": "client.html",
		"Methods":  domain.HttpMethods(),
		"Requests": requests,
	}).
		PutIfAbsent("Request", domain.NewRequestEmpty())

	return c.manager.Render(w, "home.html", context)
}

func (c *ControllerClient) request(w http.ResponseWriter, r *http.Request, context router.Context) error {
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
		"AuthType":   authType,
		"BodyType":   bodyType,
	})

	return c.client(w, r, context)
}

func (c *ControllerClient) historic(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)
	request, response, ok := c.repositoryHisotric.Find(idRequest)
	if !ok {
		return commons.ApiErrorFrom(404, "Historic request not found.")
	}

	context.Merge(map[string]any{
		"Request":  request,
		"Response": response,
	})

	return c.client(w, r, context)
}

func (c *ControllerClient) remove(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)
	request, _, ok := c.repositoryHisotric.Find(idRequest)
	if !ok {
		return commons.ApiErrorFrom(404, "Historic request not found.")
	}

	if request != nil {
		c.repositoryHisotric.Delete(*request)
	}

	return c.client(w, r, context)
}
