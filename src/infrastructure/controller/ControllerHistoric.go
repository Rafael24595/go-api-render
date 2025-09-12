package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
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
		RouteDocument(http.MethodGet, instance.find, "historic", instance.docFind()).
		RouteDocument(http.MethodPost, instance.insert, "historic", instance.docInsert()).
		RouteDocument(http.MethodDelete, instance.delete, "historic/{%s}", instance.docDelete())

	return instance
}

func (c *ControllerHistoric) docFind() docs.DocRoute {
	return docs.DocRoute{
		Description: "Fetches the list of historic requests for the current user, in lightweight format.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload([]dto.DtoLiteNodeRequest{}),
		},
	}
}

func (c *ControllerHistoric) find(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	collection, resultStatus := findHistoricCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	dtos := c.managerHistoric.FindLite(user, collection)

	return result.JsonOk(dtos)
}

func (c *ControllerHistoric) docInsert() docs.DocRoute {
	return docs.DocRoute{
		Description: "Inserts a new request/response pair into the historic collection. If the request is not a draft, the full response will be returned.",
		Request:     docs.DocJsonPayload(requestInsertAction{}),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload(responseAction{}),
		},
	}
}

func (c *ControllerHistoric) insert(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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

		return result.JsonOk(dto)
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

	return result.JsonOk(resultResponse)
}

func (c *ControllerHistoric) docDelete() docs.DocRoute {
	return docs.DocRoute{
		Description: "Deletes a historic request entry by ID. Returns the removed request and response.",
		Parameters: docs.DocParameters{
			ID_REQUEST: ID_REQUEST_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload(responseAction{}),
		},
	}
}

func (c *ControllerHistoric) delete(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	collection, resultStatus := findHistoricCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, actionRequest, actionResponse := c.managerHistoric.Delete(user, collection, idRequest)

	if actionRequest == nil && actionResponse == nil {
		return result.Reject(http.StatusNotFound)
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.JsonOk(response)
}
