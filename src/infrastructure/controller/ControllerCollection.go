package controller

import (
	"io"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

const COLLECTION = "collection"
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
		Route(http.MethodPost, instance.openapi, "/api/v1/import/openapi").
		Route(http.MethodPost, instance.importCollection, "/api/v1/import/collection").
		Route(http.MethodPost, instance.importToCollection, "/api/v1/import/collection/{%s}", COLLECTION).
		Route(http.MethodPut, instance.sortCollections, "/api/v1/sort/collection").
		Route(http.MethodPut, instance.sortRequests, "/api/v1/sort/collection/{%s}/request", COLLECTION).
		Route(http.MethodGet, instance.findCollections, "/api/v1/collection").
		Route(http.MethodGet, instance.findCollection, "/api/v1/collection/{%s}", COLLECTION).
		Route(http.MethodPost, instance.insertCollection, "/api/v1/collection").
		Route(http.MethodPut, instance.collectRequest, "/api/v1/collection").
		Route(http.MethodDelete, instance.deleteCollection, "/api/v1/collection/{%s}", COLLECTION).
		Route(http.MethodPost, instance.cloneCollection, "/api/v1/collection/{%s}/clone", COLLECTION).
		Route(http.MethodPut, instance.takeFromCollection, "/api/v1/collection/{%s}/request/{%s}", COLLECTION, ID_REQUEST).
		Route(http.MethodDelete, instance.deleteFromCollection, "/api/v1/collection/{%s}/request/{%s}", COLLECTION, ID_REQUEST)

	return instance
}

func (c *ControllerCollection) openapi(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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

	return result.Ok(collection)
}

func (c *ControllerCollection) importCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtos, err := jsonDeserialize[[]dto.DtoCollection](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collections, err := c.managerGroup.ImportDtoCollections(user, group, *dtos...)
	if err != nil {
		return result.Err(http.StatusBadRequest, err)
	}

	return result.Ok(collections)
}

func (c *ControllerCollection) importToCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(COLLECTION)
	if id == "" {
		return result.Ok(nil)
	}

	dtos, err := jsonDeserialize[[]dto.DtoRequest](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collection := c.managerGroup.ImportDtoRequestsById(user, group, id, *dtos)

	return result.Ok(collection)
}

func (c *ControllerCollection) sortCollections(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dto, err := jsonDeserialize[requestSortNodes](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}
	payload := requestSortCollectionToPayload(dto)

	group = c.managerGroup.SortCollections(user, group, payload)
	
	dtos := c.managerCollection.FindLiteCollectionNodes(user, group.Nodes)

	return result.Ok(dtos)
}

func (c *ControllerCollection) sortRequests(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(COLLECTION)
	if id == "" {
		return result.Ok(nil)
	}

	dto, err := jsonDeserialize[requestSortNodes](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	payload := requestSortCollectionToPayload(dto)

	collection := c.managerCollection.SortCollectionRequestById(user, id, payload)

	return result.Ok(collection)
}

func (c *ControllerCollection) findCollections(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	dtos := c.managerCollection.FindLiteCollectionNodes(user, group.Nodes)

	return result.Ok(dtos)
}

func (c *ControllerCollection) findCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(COLLECTION)
	if id == "" {
		return result.Ok(nil)
	}

	dto, _ := c.managerCollection.FindDto(user, id)

	return result.Ok(dto)
}

func (c *ControllerCollection) insertCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	collection, err := jsonDeserialize[domain.Collection](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collection = c.managerGroup.ImportCollection(user, group, collection)

	return result.Ok(collection)
}

func (c *ControllerCollection) collectRequest(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	request, err := jsonDeserialize[requestPushToCollection](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	payload := requestPushToCollectionToPayload(request)

	_, collection, _ := c.managerGroup.CollectRequest(user, group, payload)

	return result.Ok(collection)
}

func (c *ControllerCollection) deleteCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(COLLECTION)
	if id == "" {
		return result.Ok(nil)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	_, collection := c.managerGroup.DeleteCollection(user, group, id)

	return result.Ok(collection)
}

func (c *ControllerCollection) cloneCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return result.Ok(nil)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	payload, err := jsonDeserialize[requestCloneCollection](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	_, collection := c.managerGroup.CloneCollection(user, group, idCollection, payload.CollectionName)

	return result.Ok(collection)
}

func (c *ControllerCollection) takeFromCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	sourceId := r.PathValue(COLLECTION)
	if sourceId == "" {
		return result.Ok(nil)
	}

	idRequest := r.PathValue(ID_REQUEST)
	if idRequest == "" {
		return result.Ok(nil)
	}

	target, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	source, _, _ := c.managerCollection.MoveRequestBetweenCollectionsById(user, sourceId, target.Id, idRequest, repository.MOVE)

	return result.Ok(source)
}

func (c *ControllerCollection) deleteFromCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return result.Ok(nil)
	}

	idRequest := r.PathValue(ID_REQUEST)
	if idRequest == "" {
		return result.Ok(nil)
	}

	group, resultStatus := findUserGroup(user)
	if resultStatus != nil {
		return *resultStatus
	}

	collection, _, _ := c.managerCollection.DeleteRequestFromCollectionById(user, idCollection, idRequest)
	c.managerGroup.ResolveCollectionReferences(user, group)

	return result.Ok(collection)
}
