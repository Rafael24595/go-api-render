package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

type ControllerStorage struct {
	router             *router.Router
	managerActions     *repository.ManagerRequest
}

func NewControllerStorage(
	router *router.Router,
	managerActions *repository.ManagerRequest) ControllerStorage {
	instance := ControllerStorage{
		router:             router,
		managerActions:     managerActions,
	}
	
	router.
		Route(http.MethodPost, instance.importRequests, "/api/v1/import/request").
		Route(http.MethodGet, instance.findRequests, "/api/v1/request").
		Route(http.MethodPost, instance.insertAction, "/api/v1/request").
		Route(http.MethodPut, instance.updateRequest, "/api/v1/request").
		Route(http.MethodDelete, instance.deleteAction, "/api/v1/request/{%s}", ID_REQUEST).
		Route(http.MethodGet, instance.findAction, "/api/v1/request/{%s}", ID_REQUEST)

	return instance
}

func (c *ControllerStorage) importRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	dtos, err := jsonDeserialize[[]dto.DtoRequest](r)
	if err != nil {
		return err
	}

	requests := c.managerActions.ImportDtoRequests(user, *dtos)

	json.NewEncoder(w).Encode(requests)

	return nil
}

func (c *ControllerStorage) findRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)
	status := domain.FINAL

	actions := c.managerActions.FindOwner(user, &status)

	json.NewEncoder(w).Encode(actions)

	return nil
}

func (c *ControllerStorage) insertAction(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	action, err := jsonDeserialize[requestInsertAction](r)
	if err != nil {
		return err
	}

	actionRequest, actionResponse := c.managerActions.Release(user, dto.ToRequest(&action.Request), dto.ToResponse(&action.Response))

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) updateRequest(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	dtoRequest, err := jsonDeserialize[dto.DtoRequest](r)
	if err != nil {
		return err
	}

	request := c.managerActions.Update(user, dto.ToRequest(dtoRequest))

	json.NewEncoder(w).Encode(dto.FromRequest(request))

	return nil
}

func (c *ControllerStorage) deleteAction(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse := c.managerActions.DeleteById(user, idRequest)

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) findAction(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse, ok := c.managerActions.Find(user, idRequest)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	json.NewEncoder(w).Encode(response)

	return nil
}
