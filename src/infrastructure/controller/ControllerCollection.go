package controller

import (
	"io"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const ID_COLLECTION = "collection"
const ID_COLLECTION_DESCRIPTION = "Collection ID"
const ACTION = "action"
const TARGET = "target"

type ControllerCollection struct {
	router            *router.Router
	managerCollection *repository.ManagerCollection
	managerGroup      *repository.ManagerGroup
}

func NewControllerCollection(
	router *router.Router,
	managerCollection *repository.ManagerCollection,
	managerGroup *repository.ManagerGroup) ControllerCollection {
	instance := ControllerCollection{
		router:            router,
		managerCollection: managerCollection,
		managerGroup:      managerGroup,
	}

	instance.router.
		RouteDocument(http.MethodPost, instance.openApi, "import/openapi", instance.docOpenApi()).
		RouteDocument(http.MethodPost, instance.importItems, "import/collection", instance.docImportItems()).
		RouteDocument(http.MethodPost, instance.importTo, "import/collection/{%s}", instance.docImportTo()).
		RouteDocument(http.MethodPut, instance.sort, "sort/collection", instance.docSort()).
		RouteDocument(http.MethodPut, instance.sortRequests, "sort/collection/{%s}/request", instance.docSortRequests()).
		RouteDocument(http.MethodGet, instance.findAll, "collection", instance.docFindAll()).
		RouteDocument(http.MethodGet, instance.find, "collection/{%s}", instance.docFind()).
		RouteDocument(http.MethodGet, instance.findLite, "collection/{%s}/lite", instance.docFindLite()).
		RouteDocument(http.MethodPost, instance.insert, "collection", instance.docInsert()).
		RouteDocument(http.MethodDelete, instance.delete, "collection/{%s}", instance.docDelete()).
		RouteDocument(http.MethodPost, instance.clone, "collection/{%s}/clone", instance.docClone()).
		RouteDocument(http.MethodPut, instance.collect, "collection", instance.docCollect()).
		RouteDocument(http.MethodPut, instance.take, "collection/{%s}/request/{%s}", instance.docTake()).
		RouteDocument(http.MethodDelete, instance.deleteFrom, "collection/{%s}/request/{%s}", instance.docDeleteFrom())

	return instance
}

func (c *ControllerCollection) docOpenApi() docs.DocRoute {
	return docs.DocRoute{
		Description: "Imports an OpenAPI specification file and converts it into a collection associated with the current user group.",
		Files: docs.DocParameters{
			"file": "OpenAPI file",
		},
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) openApi(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	err = file.Close()
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collection, err := c.managerGroup.ImportOpenApi(user, group, data)
	if err != nil {
		return result.Err(http.StatusBadRequest, err)
	}

	return result.Ok(collection.Id)
}

func (c *ControllerCollection) docImportItems() docs.DocRoute {
	return docs.DocRoute{
		Description: "Imports multiple collection objects and associates them with the current user's group.",
		Request:     docs.DocJsonPayload[[]dto.DtoCollection](),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]string](ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) importItems(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtos, res := router.InputJson[[]dto.DtoCollection](r)
	if res != nil {
		return *res
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collections, err := c.managerGroup.ImportDtoCollections(user, group, dtos...)
	if err != nil {
		return result.Err(http.StatusBadRequest, err)
	}

	ids := make([]string, len(collections))
	for i, v := range collections {
		ids[i] = v.Id
	}

	return result.JsonOk(ids)
}

func (c *ControllerCollection) docImportTo() docs.DocRoute {
	return docs.DocRoute{
		Description: "Imports a list of requests and adds them to an existing collection identified by ID.",
		Parameters: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
		},
		Request: docs.DocJsonPayload[[]dto.DtoRequest](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) importTo(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(ID_COLLECTION)
	if id == "" {
		return result.Reject(http.StatusNotFound)
	}

	dtos, res := router.InputJson[[]dto.DtoRequest](r)
	if res != nil {
		return *res
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collection := c.managerGroup.ImportDtoRequestsById(user, group, id, dtos)

	return result.Ok(collection.Id)
}

func (c *ControllerCollection) docSort() docs.DocRoute {
	return docs.DocRoute{
		Description: "Sorts the collections in the user group according to the provided node structure.",
		Request:     docs.DocJsonPayload[requestSortNodes](),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]string](ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) sort(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dto, res := router.InputJson[*requestSortNodes](r)
	if res != nil {
		return *res
	}

	group, res := findUserGroup(user)
	if res != nil {
		return *res
	}
	payload := requestSortCollectionToPayload(dto)

	group = c.managerGroup.SortCollections(user, group, payload)

	dtos := c.managerCollection.FindLiteCollectionNodes(user, group.Nodes)

	ids := make([]string, len(dtos))
	for i, v := range dtos {
		ids[i] = v.Collection.Id
	}

	return result.JsonOk(ids)
}

func (c *ControllerCollection) docSortRequests() docs.DocRoute {
	return docs.DocRoute{
		Description: "Sorts the requests inside a specific collection based on the provided node order.",
		Parameters: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
		},
		Request: docs.DocJsonPayload[[]requestSortNodes](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) sortRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(ID_COLLECTION)
	if id == "" {
		return result.Reject(http.StatusNotFound)
	}

	dto, res := router.InputJson[*requestSortNodes](r)
	if res != nil {
		return *res
	}

	payload := requestSortCollectionToPayload(dto)

	collection := c.managerCollection.SortCollectionRequestById(user, id, payload)

	return result.Ok(collection.Id)
}

func (c *ControllerCollection) docFindAll() docs.DocRoute {
	return docs.DocRoute{
		Description: "Retrieves all collection summaries (lite version) for the current user's group.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]dto.DtoLiteNodeCollection](),
		},
	}
}

func (c *ControllerCollection) findAll(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	dtos := c.managerCollection.FindLiteCollectionNodes(user, group.Nodes)

	return result.JsonOk(dtos)
}

func (c *ControllerCollection) docFind() docs.DocRoute {
	return docs.DocRoute{
		Description: "Fetches a full collection by ID, including all associated metadata, context and requests.",
		Parameters: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[dto.DtoCollection](),
		},
	}
}

func (c *ControllerCollection) find(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(ID_COLLECTION)
	if id == "" {
		return result.Reject(http.StatusNotFound)
	}

	dto, _ := c.managerCollection.FindDto(user, id)

	return result.JsonOk(dto)
}

func (c *ControllerCollection) docFindLite() docs.DocRoute {
	return docs.DocRoute{
		Description: "Fetches a lite version of a collection by ID, including basic metadata but excluding full context and request data.",
		Parameters: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[dto.DtoLiteCollection](),
		},
	}
}

func (c *ControllerCollection) findLite(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(ID_COLLECTION)
	if id == "" {
		return result.Reject(http.StatusNotFound)
	}

	dto, _ := c.managerCollection.FindDtoLite(user, id)

	return result.JsonOk(dto)
}

func (c *ControllerCollection) docInsert() docs.DocRoute {
	return docs.DocRoute{
		Description: "Inserts a new collection into the current user's group.",
		Request:     docs.DocJsonPayload[domain.Collection](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) insert(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	collection, res := router.InputJson[*domain.Collection](r)
	if res != nil {
		return *res
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collection = c.managerGroup.ImportCollection(user, group, collection)

	return result.Ok(collection.Id)
}

func (c *ControllerCollection) docCollect() docs.DocRoute {
	return docs.DocRoute{
		Description: "Collects a request from a general context and adds it into the user's active collection.",
		Request:     docs.DocJsonPayload[requestPushToCollection](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) docDelete() docs.DocRoute {
	return docs.DocRoute{
		Description: "Deletes a collection by ID from the current user's group.",
		Parameters: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) delete(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(ID_COLLECTION)
	if id == "" {
		return result.Reject(http.StatusNotFound)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collection := c.managerGroup.DeleteCollection(user, group, id)

	return result.Ok(collection.Id)
}

func (c *ControllerCollection) docClone() docs.DocRoute {
	return docs.DocRoute{
		Description: "Creates a duplicate of a collection, assigning it a new name and preserving all its structure and requests.",
		Parameters: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) clone(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	idCollection := r.PathValue(ID_COLLECTION)
	if idCollection == "" {
		return result.Reject(http.StatusNotFound)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	payload, res := router.InputJson[requestCloneCollection](r)
	if res != nil {
		return *res
	}

	_, collection := c.managerGroup.CloneCollection(user, group, idCollection, payload.CollectionName)

	return result.Ok(collection.Id)
}

func (c *ControllerCollection) collect(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	request, res := router.InputJson[*requestPushToCollection](r)
	if res != nil {
		return *res
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	payload := requestPushToCollectionToPayload(request)

	_, collection, _ := c.managerGroup.CollectRequest(user, group, payload)

	return result.Ok(collection.Id)
}

func (c *ControllerCollection) docTake() docs.DocRoute {
	return docs.DocRoute{
		Description: "Moves a request from one collection to the user's active collection.",
		Parameters: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
			ID_REQUEST:    ID_REQUEST_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) take(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	sourceId := r.PathValue(ID_COLLECTION)
	if sourceId == "" {
		return result.Reject(http.StatusNotFound)
	}

	idRequest := r.PathValue(ID_REQUEST)
	if idRequest == "" {
		return result.Reject(http.StatusNotFound)
	}

	target, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	source, _, _ := c.managerCollection.MoveRequestBetweenCollectionsById(user, sourceId, target.Id, idRequest, repository.MOVE)

	return result.Ok(source.Id)
}

func (c *ControllerCollection) docDeleteFrom() docs.DocRoute {
	return docs.DocRoute{
		Description: "Deletes a specific request from a collection by ID.",
		Parameters: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
			ID_REQUEST:    ID_REQUEST_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_COLLECTION_DESCRIPTION),
		},
	}
}

func (c *ControllerCollection) deleteFrom(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	idCollection := r.PathValue(ID_COLLECTION)
	if idCollection == "" {
		return result.Reject(http.StatusNotFound)
	}

	idRequest := r.PathValue(ID_REQUEST)
	if idRequest == "" {
		return result.Reject(http.StatusNotFound)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	collection, _, _ := c.managerCollection.DeleteRequestFromCollectionById(user, idCollection, idRequest)
	c.managerGroup.ResolveCollectionReferences(user, group)

	return result.Ok(collection.Id)
}
