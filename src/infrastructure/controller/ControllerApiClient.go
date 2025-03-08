package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

const ID_REQUEST = "id_request"

type ControllerApiClient struct {
	router     *router.Router
	repository *repository.RequestManager
}

func NewControllerApiClient(router *router.Router, repository *repository.RequestManager) ControllerApiClient {
	instance := ControllerApiClient{
		router:     router,
		repository: repository,
	}

	router.
		Route(http.MethodPost, instance.doAction, "/api/v1/action").
		Route(http.MethodGet, instance.actions, "/api/v1/action").
		Route(http.MethodGet, instance.action, "/api/v1/action/{%s}", ID_REQUEST)

	return instance
}

func (c *ControllerApiClient) doAction(w http.ResponseWriter, r *http.Request, context router.Context) error {
	actionRequest, err := jsonDeserialize[domain.Request](r)
	if err != nil {
		return err
	}

	actionResponse, err := core_infrastructure.Client().Fetch(*actionRequest)
	if err != nil {
		return err
	}

	//actionRequest, actionResponse = c.repository.Insert(*actionRequest, actionResponse)

	response := responseAction{
		Request:  *actionRequest,
		Response: *actionResponse,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerApiClient) actions(w http.ResponseWriter, r *http.Request, context router.Context) error {
	actions := c.repository.FindAll()

	response := responseActionRequests{
		Requests: actions,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerApiClient) action(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse, ok := c.repository.Find(idRequest)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	response := responseAction{
		Request:  *actionRequest,
		Response: *actionResponse,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func jsonDeserialize[T any](r *http.Request) (*T, error) {
	var item T
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
