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
	router            *router.Router
	managerRequest    *repository.ManagerRequest
	managerCollection *repository.ManagerCollection
}

func NewControllerStorage(
	router *router.Router,
	managerRequest *repository.ManagerRequest,
	managerCollection *repository.ManagerCollection) ControllerStorage {
	instance := ControllerStorage{
		router:            router,
		managerRequest:    managerRequest,
		managerCollection: managerCollection,
	}

	router.
		Route(http.MethodPost, instance.importRequests, "import/request").
		Route(http.MethodPut, instance.sortRequests, "sort/request").
		Route(http.MethodGet, instance.findRequests, "request").
		Route(http.MethodPost, instance.insertAction, "request").
		Route(http.MethodPut, instance.updateRequest, "request").
		Route(http.MethodGet, instance.findAction, "request/{%s}", ID_REQUEST).
		Route(http.MethodDelete, instance.deleteAction, "request/{%s}", ID_REQUEST)

	return instance
}

func (c *ControllerStorage) importRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtos, err := jsonDeserialize[[]dto.DtoRequest](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	collection = c.managerCollection.ImportDtoRequests(user, collection, *dtos)
	nodes := c.managerCollection.FindLiteRequestNodes(user, collection)

	return result.Ok(nodes)
}

func (c *ControllerStorage) sortRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dto, err := jsonDeserialize[requestSortNodes](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	payload := requestSortCollectionToPayload(dto)

	collection = c.managerCollection.SortCollectionRequest(user, collection, payload)

	return result.Ok(collection)
}

func (c *ControllerStorage) findRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	dtos := c.managerCollection.FindLiteRequestNodes(user, collection)

	return result.Ok(dtos)
}

func (c *ControllerStorage) insertAction(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	action, err := jsonDeserialize[requestInsertAction](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	request, response := c.managerRequest.Release(user, dto.ToRequest(&action.Request), dto.ToResponse(&action.Response))

	if request.Status == domain.FINAL {
		collection, resultStatus := findUserCollection(user)
		if resultStatus != nil {
			return *resultStatus
		}
		c.managerCollection.ResolveRequestReferences(user, collection, *request)
	}

	dto := responseAction{
		Request:  *dto.FromRequest(request),
		Response: *dto.FromResponse(response),
	}

	return result.Ok(dto)
}

func (c *ControllerStorage) updateRequest(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtoRequest, err := jsonDeserialize[dto.DtoRequest](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	request := c.managerRequest.Update(user, dto.ToRequest(dtoRequest))
	if request.Status == domain.FINAL {
		collection, resultStatus := findUserCollection(user)
		if resultStatus != nil {
			return *resultStatus
		}
		c.managerCollection.ResolveRequestReferences(user, collection, *request)
	}

	dto := dto.FromRequest(request)

	return result.Ok(dto)
}

func (c *ControllerStorage) findAction(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	request, response, ok := c.managerRequest.Find(user, idRequest)
	if !ok {
		return result.Err(http.StatusNotFound, nil)
	}

	if request == nil {
		return result.Err(http.StatusNotFound, nil)
	}

	dto := responseAction{
		Request:  *dto.FromRequest(request),
		Response: *dto.FromResponse(response),
	}

	return result.Ok(dto)
}

func (c *ControllerStorage) deleteAction(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, actionRequest, actionResponse := c.managerCollection.DeleteRequestFromCollection(user, collection, idRequest)

	if actionRequest == nil && actionResponse == nil {
		return result.Err(http.StatusNotFound, nil)
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.Ok(response)
}
