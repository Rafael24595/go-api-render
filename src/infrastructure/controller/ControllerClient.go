package controller

import (
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

	router.
		GroupContextualizer("/client", instance.clientContext).
		Route(http.MethodGet, instance.home, "/").
		Route(http.MethodGet, instance.client, "/client").
		Route(http.MethodPost, instance.request, "/client").
		Route(http.MethodGet, instance.show, "/client/{%s}", ID_REQUEST).
		Route(http.MethodPost, instance.request, "/client/{%s}", ID_REQUEST).
		Route(http.MethodDelete, instance.remove, "/client/{%s}", ID_REQUEST).
		Route(http.MethodPut, instance.update, "/client/{%s}", ID_REQUEST)

	return instance
}

func (c *ControllerClient) home(w http.ResponseWriter, r *http.Request, context router.Context) error {
	return c.client(w, r, context)
}

func (c *ControllerClient) client(w http.ResponseWriter, r *http.Request, context router.Context) error {
	source := r.URL.Query().Get(constants.SidebarRequest.Type)

	requestsHistoric := []domain.Request{}
	requestsPersisted := []domain.Request{}

	switch source {
	case constants.SidebarRequest.TagHistoric:
		requestsHistoric = c.repositoryHisotric.FindOptions(core_repository.FilterOptions[domain.Request]{
			Sort: func(i, j domain.Request) bool {
				return i.Timestamp > j.Timestamp
			},
			To: 10,
		})
	case constants.SidebarRequest.TagSaved:
		requestsPersisted = c.repositoryPersisted.FindOptions(core_repository.FilterOptions[domain.Request]{
			Sort: func(i, j domain.Request) bool {
				return i.Timestamp > j.Timestamp
			},
		})
	case constants.SidebarRequest.TagCollection:
		
	}


	context.Merge(map[string]any{
		"BodyTemplate":      "client.html",
		"Methods":           domain.HttpMethods(),
		"RequestsHistoric":  requestsHistoric,
		"RequestsPersisted": requestsPersisted,
	}).
		PutIfAbsent("Request", domain.NewRequestEmpty())

	return c.manager.Render(w, "home.html", context)
}

func (c *ControllerClient) request(w http.ResponseWriter, r *http.Request, context router.Context) error {
	constants := configuration.GetConstants()

	request, err := proccessRequest(r)
	if err != nil {
		return err
	}

	doRequest := r.FormValue(constants.Client.DoRequest) == "true"

	var response *domain.Response
	if doRequest {
		response, err = core_infrastructure.Client().Fetch(*request)
		if err != nil {
			return err
		}
	}

	if request.Status == domain.Historic {
		request, response = c.repositoryHisotric.Insert(*request, response)
	} else {
		okHist, _ := c.repositoryHisotric.Exists(request.Id)
		okPers, _ := c.repositoryPersisted.Exists(request.Id)
		if okHist && !okPers {
			request.Id = ""
		}
		request, response = c.repositoryPersisted.Insert(*request, response)
	}

	context.Merge(map[string]any{
		"Request":  request,
		"Response": response,
	})

	return c.client(w, r, context)
}

func (c *ControllerClient) show(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)

	var request *domain.Request
	var response *domain.Response
	var ok bool

	if c.repositoryPersisted.HasPrefix(idRequest) {
		request, response, ok = c.repositoryPersisted.Find(idRequest)
	} else {
		request, response, ok = c.repositoryHisotric.Find(idRequest)
	}

	if !ok {
		return commons.ApiErrorFrom(404, "Request not found.")
	}

	context.Merge(map[string]any{
		"Request":  request,
		"Response": response,
	})

	return c.client(w, r, context)
}

func (c *ControllerClient) remove(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)
	persisted := c.repositoryPersisted.HasPrefix(idRequest)

	var request *domain.Request
	var ok bool
	if persisted {
		request, _, ok = c.repositoryPersisted.Find(idRequest)
	} else {
		request, _, ok = c.repositoryHisotric.Find(idRequest)
	}

	if !ok {
		return commons.ApiErrorFrom(404, "Request not found.")
	}

	if request == nil {
		return c.client(w, r, context)
	}

	if persisted {
		c.repositoryPersisted.Delete(*request)
	} else {
		c.repositoryHisotric.Delete(*request)
	}

	return c.client(w, r, context)
}

func (c *ControllerClient) update(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)
	name := r.FormValue("name")

	request, _, _ := c.repositoryPersisted.Find(idRequest)

	if request == nil {
		return commons.ApiErrorFrom(404, "Request not found.")
	}

	request.Name = name

	c.repositoryPersisted.Insert(*request, nil)

	context.Merge(map[string]any{
		"Request": request,
	})

	return c.client(w, r, context)
}

func (c *ControllerClient) clientContext(w http.ResponseWriter, r *http.Request) (router.Context, error) {
	requestType := r.URL.Query().Get(constants.SidebarRequest.Type)
	clientType := r.URL.Query().Get(constants.Client.Type)
	bodyType := r.URL.Query().Get(constants.Body.Type)
	authType := r.URL.Query().Get(constants.Auth.Type)
	responseType := r.URL.Query().Get(constants.Response.Type)
	return collection.FromMap(map[string]any{
		"RequestType":  requestType,
		"ResponseType": responseType,
		"ClientType":   clientType,
		"AuthType":     authType,
		"BodyType":     bodyType,
	}), nil
}
