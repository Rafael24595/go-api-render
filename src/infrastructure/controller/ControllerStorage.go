package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

type ControllerStorage struct {
	router         *router.Router
	managerActions *repository.ManagerRequest
}

func NewControllerStorage(
	router *router.Router,
	managerActions *repository.ManagerRequest) ControllerStorage {
	instance := ControllerStorage{
		router:         router,
		managerActions: managerActions,
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

func (c *ControllerStorage) importRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtos, err := jsonDeserialize[[]dto.DtoRequest](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	requests := c.managerActions.ImportDtoRequests(user, *dtos)

	return result.Ok(requests)
}

func (c *ControllerStorage) findRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	status := domain.FINAL

	requests := c.managerActions.FindOwner(user, &status)

	dtos := make([]dto.DtoRequest, len(requests))
	for i, v := range requests {
		dtos[i] = *dto.FromRequest(&v)
	}

	return result.Ok(dtos)
}

func (c *ControllerStorage) insertAction(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	action, err := jsonDeserialize[requestInsertAction](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	actionRequest, actionResponse := c.managerActions.Release(user, dto.ToRequest(&action.Request), dto.ToResponse(&action.Response))

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.Ok(response)
}

func (c *ControllerStorage) updateRequest(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtoRequest, err := jsonDeserialize[dto.DtoRequest](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	request := c.managerActions.Update(user, dto.ToRequest(dtoRequest))

	dto := dto.FromRequest(request)

	return result.Ok(dto)
}

func (c *ControllerStorage) deleteAction(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse := c.managerActions.DeleteById(user, idRequest)

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.Ok(response)
}

func (c *ControllerStorage) findAction(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse, ok := c.managerActions.Find(user, idRequest)
	if !ok {
		return result.Err(http.StatusNotFound, nil)
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.Ok(response)
}
