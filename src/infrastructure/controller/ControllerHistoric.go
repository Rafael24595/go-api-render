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
	router          *router.Router
	managerRequest  *repository.ManagerRequest
	managerHistoric *repository.ManagerHistoric
}

func NewControllerHistoric(
	router *router.Router,
	managerRequest *repository.ManagerRequest,
	managerHistoric *repository.ManagerHistoric) ControllerHistoric {
	instance := ControllerHistoric{
		router:          router,
		managerRequest:  managerRequest,
		managerHistoric: managerHistoric,
	}

	router.
		Route(http.MethodGet, instance.findHistoric, "historic").
		Route(http.MethodPost, instance.insertHistoric, "historic").
		Route(http.MethodDelete, instance.deleteHistoric, "historic/{%s}", ID_REQUEST)

	return instance
}

func (c *ControllerHistoric) insertHistoric(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	action, err := jsonDeserialize[requestInsertAction](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	if action.Request.Status != domain.DRAFT {
		response := c.managerRequest.InsertResponse(user, dto.ToResponse(&action.Response))

		var dtoResponse *dto.DtoResponse
		if response != nil {
			dtoResponse = dto.FromResponse(response)
		}

		dto := responseAction{
			Request:  action.Request,
			Response: *dtoResponse,
		}

		return result.Ok(dto)
	}

	collection, resultStatus := findHistoricCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, request, response := c.managerHistoric.Insert(user, collection, dto.ToRequest(&action.Request), dto.ToResponse(&action.Response))

	resultResponse := responseAction{
		Request:  *dto.FromRequest(request),
		Response: *dto.FromResponse(response),
	}

	return result.Ok(resultResponse)
}

func (c *ControllerHistoric) findHistoric(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	collection, resultStatus := findHistoricCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	dtos := c.managerHistoric.FindLite(user, collection)

	return result.Ok(dtos)
}

func (c *ControllerHistoric) deleteHistoric(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	collection, resultStatus := findHistoricCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, actionRequest, actionResponse := c.managerHistoric.Delete(user, collection, idRequest)

	if actionRequest == nil && actionResponse == nil {
		return result.Err(http.StatusNotFound, nil)
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.Ok(response)
}
