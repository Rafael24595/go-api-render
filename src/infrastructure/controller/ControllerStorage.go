package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

const USER = "user"

type ControllerStorage struct {
	router             *router.Router
	repositoryActions  *repository.RequestManager
	repositoryHisotric repository.IRepositoryHistoric
}

func NewControllerStorage(
		router *router.Router, 
		repository *repository.RequestManager, 
		repositoryHisotric repository.IRepositoryHistoric) ControllerStorage {
	instance := ControllerStorage{
		router:     router,
		repositoryActions: repository,
		repositoryHisotric: repositoryHisotric,
	}

	//TODO: Extract users from token.
	router.
		Route(http.MethodPost, instance.storage, "/api/v1/storage/{%s}", USER).
		Route(http.MethodGet, instance.findAll, "/api/v1/storage/{%s}", USER).
		Route(http.MethodDelete, instance.delete, "/api/v1/storage/{%s}/{%s}", USER, ID_REQUEST).
		Route(http.MethodGet, instance.find, "/api/v1/storage/{%s}/{%s}", USER, ID_REQUEST).
		Route(http.MethodPost, instance.historic, "/api/v1/historic/{%s}", USER).
		Route(http.MethodGet, instance.findHistoric, "/api/v1/historic/{%s}", USER)

	return instance
}

func (c *ControllerStorage) storage(w http.ResponseWriter, r *http.Request, context router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	action, err := jsonDeserialize[RequestInsertAction](r)
	if err != nil {
		return err
	}

	action.Request.Owner = user
	action.Request.Status = domain.FINAL

	actionRequest, actionResponse := c.repositoryActions.Insert(action.Request, &action.Response)

	response := responseAction{
		Request:  *actionRequest,
		Response: *actionResponse,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) historic(w http.ResponseWriter, r *http.Request, context router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	action, err := jsonDeserialize[RequestInsertAction](r)
	if err != nil {
		return err
	}

	if _, err := domain.StatusFromString(string(action.Request.Status)); err != nil {
		action.Request.Status = domain.DRAFT
	}

	if action.Request.Status != domain.DRAFT {
		return nil
	}

	action.Request.Owner = user

	actionRequest, actionResponse := c.repositoryActions.Insert(action.Request, &action.Response)

	step := domain.NewHistoric(actionRequest.Id, user)
	c.repositoryHisotric.Insert(*step)
	//TODO: Implement delete old steps

	response := responseAction{
		Request:  *actionRequest,
		Response: *actionResponse,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) findAll(w http.ResponseWriter, r *http.Request, context router.Context) error {
	actions := c.repositoryActions.FindAll()

	response := responseActionRequests{
		Requests: actions,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) delete(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse := c.repositoryActions.DeleteById(idRequest)
	
	response := responseAction{
		Request:  *actionRequest,
		Response: *actionResponse,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) find(w http.ResponseWriter, r *http.Request, context router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse, ok := c.repositoryActions.Find(idRequest)
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

func (c *ControllerStorage) findHistoric(w http.ResponseWriter, r *http.Request, context router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	steps := c.repositoryHisotric.FindByOwner(user)
	requests := c.repositoryActions.FindSteps(steps)

	json.NewEncoder(w).Encode(requests)

	return nil
}
