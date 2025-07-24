package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/docs"
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
		RouteDocument(http.MethodPost, instance.importItems, "import/request", instance.docImportItems()).
		RouteDocument(http.MethodPut, instance.sort, "sort/request", instance.docSort()).
		RouteDocument(http.MethodGet, instance.findAll, "request", instance.docFindAll()).
		RouteDocument(http.MethodPost, instance.insert, "request", instance.docInsert()).
		RouteDocument(http.MethodPut, instance.update, "request", instance.docUpdate()).
		RouteDocument(http.MethodGet, instance.find, "request/{%s}", instance.docFind()).
		RouteDocument(http.MethodDelete, instance.delete, "request/{%s}", instance.docDelete())

	return instance
}

func (c *ControllerStorage) docImportItems() docs.DocPayload {
	return docs.DocPayload{
		Description: "Imports multiple requests into the user's default collection.",
		Request: docs.DocStruct([]dto.DtoRequest{}),
		Responses: docs.DocResponses{
			"200": docs.DocStruct([]string{}),
		},
	}
}

func (c *ControllerStorage) importItems(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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

	ids := make([]string, len(nodes))
	for i, v := range nodes {
		ids[i] = v.Request.Id
	}

	return result.Ok(ids)
}

func (c *ControllerStorage) docSort() docs.DocPayload {
	return docs.DocPayload{
		Description: "Sorts the requests within the user's default collection based on the provided node structure.",
		Request: docs.DocStruct(requestSortNodes{}),
		Responses: docs.DocResponses{
			"200": docs.DocStruct(ID_COLLECTION, ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerStorage) sort(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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

	return result.Ok(collection.Id)
}

func (c *ControllerStorage) docFindAll() docs.DocPayload {
	return docs.DocPayload{
		Description: "Retrieves all request nodes (lite version) from the user's default collection.",
		Responses: docs.DocResponses{
			"200": docs.DocStruct([]dto.DtoLiteNodeRequest{}),
		},
	}
}

func (c *ControllerStorage) findAll(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	dtos := c.managerCollection.FindLiteRequestNodes(user, collection)

	return result.Ok(dtos)
}

func (c *ControllerStorage) docInsert() docs.DocPayload {
	return docs.DocPayload{
		Description: "Inserts a new request and its response into the user's default collection.",
		Request: docs.DocStruct(requestInsertAction{}),
		Responses: docs.DocResponses{
			"200": docs.DocStruct(responseAction{}),
		},
	}
}

func (c *ControllerStorage) insert(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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

func (c *ControllerStorage) docUpdate() docs.DocPayload {
	return docs.DocPayload{
		Description: "Updates an existing request in the user's collection.",
		Request: docs.DocStruct(dto.DtoRequest{}),
		Responses: docs.DocResponses{
			"200": docs.DocStruct(ID_REQUEST, ID_REQUEST_DESCRIPTION),
		},
	}
}

func (c *ControllerStorage) update(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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

	return result.Ok(dto.Id)
}

func (c *ControllerStorage) docFind() docs.DocPayload {
	return docs.DocPayload{
		Description: "Finds a specific request and its response by request ID.",
		Parameters: docs.DocParameters{
			ID_REQUEST: ID_REQUEST_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocStruct(responseAction{}),
		},
	}
}

func (c *ControllerStorage) find(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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

func (c *ControllerStorage) docDelete() docs.DocPayload {
	return docs.DocPayload{
		Description: "Deletes a specific request from the user's collection by ID.",
		Parameters: docs.DocParameters{
			ID_REQUEST: ID_REQUEST_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocStruct(responseAction{}),
		},
	}
}

func (c *ControllerStorage) delete(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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
