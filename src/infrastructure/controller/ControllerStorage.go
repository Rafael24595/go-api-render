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

	router.
		Route(http.MethodPost, instance.storage, "/api/v1/storage/{%s}", USER).
		Route(http.MethodPost, instance.historic, "/api/v1/historic/{%s}", USER)

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

	if action.Request.Status != domain.DRAFT {
		return nil
	}

	action.Request.Owner = user

	actionRequest, actionResponse := c.repositoryActions.Insert(action.Request, &action.Response)

	step := domain.NewHistoric(actionRequest.Id)
	c.repositoryHisotric.Insert(*step)
	//TODO: Implement delete old steps

	response := responseAction{
		Request:  *actionRequest,
		Response: *actionResponse,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}
