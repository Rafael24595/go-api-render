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
	managerRequest    *repository.ManagerRequest
	managerContext    *repository.ManagerContext
}

func NewControllerCollection(
	router *router.Router,
	managerCollection *repository.ManagerCollection,
	managerRequest *repository.ManagerRequest,
	managerContext *repository.ManagerContext) ControllerCollection {
	instance := ControllerCollection{
		router:            router,
		managerCollection: managerCollection,
		managerRequest:    managerRequest,
		managerContext:    managerContext,
	}

	instance.router.
		Route(http.MethodPost, instance.openapi, "/api/v1/import/openapi").
		Route(http.MethodPost, instance.importCollection, "/api/v1/import/collection").
		Route(http.MethodPost, instance.importToCollection, "/api/v1/import/collection/{%s}", COLLECTION).
		Route(http.MethodGet, instance.findCollection, "/api/v1/collection").
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

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("file")
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	collection, err := c.managerCollection.ImportOpenApi(user, data)
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

	collection, err := c.managerCollection.ImportDtoCollections(user, *dtos)
	if err != nil {
		return result.Err(http.StatusBadRequest, err)
	}

	return result.Ok(collection)
}

func (c *ControllerCollection) importToCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return result.Ok(nil)
	}

	dtos, err := jsonDeserialize[[]dto.DtoRequest](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	collection := c.managerCollection.ImportDtoRequestsById(user, idCollection, *dtos)

	return result.Ok(collection)
}

func (c *ControllerCollection) findCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	collections := c.managerCollection.FindFreeByOwner(user)

	dtos := make([]dto.DtoCollection, len(collections))
	for i, v := range collections {
		requests := c.managerRequest.FindNodes(v.Nodes)
		context, _ := c.managerContext.Find(user, v.Context)
		dtoContext := dto.FromContext(context)
		dtos[i] = dto.DtoCollection{
			Id:        v.Id,
			Name:      v.Name,
			Timestamp: v.Timestamp,
			Context:   *dtoContext,
			Nodes:     requests,
			Owner:     v.Owner,
			Modified:  v.Modified,
		}
	}

	return result.Ok(dtos)
}

func (c *ControllerCollection) insertCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	collection, err := jsonDeserialize[domain.Collection](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	collection = c.managerCollection.Insert(user, collection)

	return result.Ok(collection)
}

func (c *ControllerCollection) collectRequest(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	request, err := jsonDeserialize[requestPushToCollection](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	payload := requestPushToCollectionToPayload(request)

	collection, _ := c.managerCollection.CollectRequest(user, payload)

	return result.Ok(collection)
}

func (c *ControllerCollection) deleteCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return result.Ok(nil)
	}

	collection := c.managerCollection.Delete(user, idCollection)

	return result.Ok(collection)
}

func (c *ControllerCollection) cloneCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return result.Ok(nil)
	}

	payload, err := jsonDeserialize[requestCloneCollection](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	collection := c.managerCollection.CloneCollection(user, idCollection, payload.CollectionName)

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

	sessions := repository.InstanceManagerSession()
	target, err := sessions.FindUserCollection(user)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
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

	collection, _, _ := c.managerCollection.RemoveRequestFromCollectionById(user, idCollection, idRequest)

	return result.Ok(collection)
}
