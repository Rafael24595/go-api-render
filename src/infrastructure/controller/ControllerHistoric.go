package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
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

func (c *ControllerHistoric) insertHistoric(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	action, err := jsonDeserialize[requestInsertAction](r)
	if err != nil {
		return err
	}

	if action.Request.Status != domain.DRAFT {
		c.managerActions.InsertResponse(user, dto.ToResponse(&action.Response))
		return nil
	}

	actionRequest, actionResponse := c.managerActions.Insert(user, dto.ToRequest(&action.Request), dto.ToResponse(&action.Response))

	step := domain.NewHistoric(actionRequest.Id, user)
	c.repositoryHisotric.Insert(*step)
	//TODO: Implement delete old steps

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerHistoric) findHistoric(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	steps := c.repositoryHisotric.FindByOwner(user)
	requests := c.managerActions.FindSteps(steps)

	json.NewEncoder(w).Encode(requests)

	return nil
}
