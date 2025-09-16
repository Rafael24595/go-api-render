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

type ControllerRequest struct {
	router            *router.Router
	managerRequest    *repository.ManagerRequest
	managerCollection *repository.ManagerCollection
}

func NewControllerRequest(
	router *router.Router,
	managerRequest *repository.ManagerRequest,
	managerCollection *repository.ManagerCollection) ControllerRequest {
	instance := ControllerRequest{
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

func (c *ControllerRequest) docImportItems() docs.DocRoute {
	return docs.DocRoute{
		Description: "Imports multiple requests into the user's default collection.",
		Request:     docs.DocJsonPayload[[]dto.DtoRequest](),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]string](),
		},
	}
}

func (c *ControllerRequest) importItems(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	dtos, res := router.InputJson[[]dto.DtoRequest](r)
	if res != nil {
		return *res
	}

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	collection = c.managerCollection.ImportDtoRequests(user, collection, dtos)
	nodes := c.managerCollection.FindLiteRequestNodes(user, collection)

	ids := make([]string, len(nodes))
	for i, v := range nodes {
		ids[i] = v.Request.Id
	}

	return result.JsonOk(ids)
}

func (c *ControllerRequest) docSort() docs.DocRoute {
	return docs.DocRoute{
		Description: "Sorts the requests within the user's default collection based on the provided node structure.",
		Request:     docs.DocJsonPayload[requestSortNodes](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerRequest) sort(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	dto, res := router.InputJson[*requestSortNodes](r)
	if res != nil {
		return *res
	}

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	payload := requestSortCollectionToPayload(dto)

	collection = c.managerCollection.SortCollectionRequest(user, collection, payload)

	return result.Ok(collection.Id)
}

func (c *ControllerRequest) docFindAll() docs.DocRoute {
	return docs.DocRoute{
		Description: "Retrieves all request nodes (lite version) from the user's default collection.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]dto.DtoLiteNodeRequest](),
		},
	}
}

func (c *ControllerRequest) findAll(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	dtos := c.managerCollection.FindLiteRequestNodes(user, collection)

	return result.JsonOk(dtos)
}

func (c *ControllerRequest) docInsert() docs.DocRoute {
	return docs.DocRoute{
		Description: "Inserts a new request and its response into the user's default collection.",
		Request:     docs.DocJsonPayload[requestInsertAction](),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseAction](),
		},
	}
}

func (c *ControllerRequest) insert(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	action, res := router.InputJson[requestInsertAction](r)
	if res != nil {
		return *res
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

	return result.JsonOk(dto)
}

func (c *ControllerRequest) docUpdate() docs.DocRoute {
	return docs.DocRoute{
		Description: "Updates an existing request in the user's collection.",
		Request:     docs.DocJsonPayload[dto.DtoRequest](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_REQUEST_DESCRIPTION),
		},
	}
}

func (c *ControllerRequest) update(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	dtoRequest, res := router.InputJson[*dto.DtoRequest](r)
	if res != nil {
		return *res
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

func (c *ControllerRequest) docFind() docs.DocRoute {
	return docs.DocRoute{
		Description: "Finds a specific request and its response by request ID.",
		Parameters: docs.DocParameters{
			ID_REQUEST: ID_REQUEST_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseAction](),
		},
	}
}

func (c *ControllerRequest) find(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	request, response, ok := c.managerRequest.Find(user, idRequest)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	if request == nil {
		return result.Reject(http.StatusNotFound)
	}

	dto := responseAction{
		Request:  *dto.FromRequest(request),
		Response: *dto.FromResponse(response),
	}

	return result.JsonOk(dto)
}

func (c *ControllerRequest) docDelete() docs.DocRoute {
	return docs.DocRoute{
		Description: "Deletes a specific request from the user's collection by ID.",
		Parameters: docs.DocParameters{
			ID_REQUEST: ID_REQUEST_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseAction](),
		},
	}
}

func (c *ControllerRequest) delete(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	idRequest := r.PathValue(ID_REQUEST)

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, actionRequest, actionResponse := c.managerCollection.DeleteRequestFromCollection(user, collection, idRequest)

	if actionRequest == nil && actionResponse == nil {
		return result.Reject(http.StatusNotFound)
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.JsonOk(response)
}
