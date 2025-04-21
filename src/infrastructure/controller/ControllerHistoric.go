package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

type ControllerHistoric struct {
	router             *router.Router
	managerActions     *repository.ManagerRequest
	repositoryHisotric repository.IRepositoryHistoric
}

func NewControllerHistoric(
	router *router.Router,
	managerActions *repository.ManagerRequest,
	repositoryHisotric repository.IRepositoryHistoric) ControllerHistoric {
	instance := ControllerHistoric{
		router:             router,
		managerActions:     managerActions,
		repositoryHisotric: repositoryHisotric,
	}

	router.
		Route(http.MethodPost, instance.insertHistoric, "/api/v1/historic").
		Route(http.MethodGet, instance.findHistoric, "/api/v1/historic")

	return instance
}

func (c *ControllerHistoric) insertHistoric(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	action, err := jsonDeserialize[requestInsertAction](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	if action.Request.Status != domain.DRAFT {
		response := c.managerActions.InsertResponse(user, dto.ToResponse(&action.Response))

		dto := responseAction{
			Request:  action.Request,
			Response: *dto.FromResponse(response),
		}

		return result.Ok(dto)
	}

	actionRequest, actionResponse := c.managerActions.Insert(user, dto.ToRequest(&action.Request), dto.ToResponse(&action.Response))

	step := domain.NewHistoric(actionRequest.Id, user)
	c.repositoryHisotric.Insert(*step)
	//TODO: Implement delete old steps

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.Ok(response)
}

func (c *ControllerHistoric) findHistoric(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	steps := c.repositoryHisotric.FindByOwner(user)
	requests := c.managerActions.FindSteps(steps)

	dtos := make([]dto.DtoRequest, len(requests))
	for i, v := range requests {
		dtos[i] = *dto.FromRequest(&v)
	}

	return result.Ok(dtos)
}
