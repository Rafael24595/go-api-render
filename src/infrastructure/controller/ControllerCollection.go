package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

const COLLECTION = "collection"
const ACTION = "action"
const TARGET = "target"

type ControllerCollection struct {
	router            *router.Router
	managerCollection *repository.ManagerCollection
	managerActions    *repository.ManagerRequest
	repositoryContext repository.IRepositoryContext
}

func NewControllerCollection(
	router *router.Router,
	managerCollection *repository.ManagerCollection,
	managerActions *repository.ManagerRequest,
	repositoryContext repository.IRepositoryContext) ControllerCollection {
	instance := ControllerCollection{
		router:            router,
		managerCollection: managerCollection,
		managerActions:    managerActions,
		repositoryContext: repositoryContext,
	}

	instance.router.
		Route(http.MethodPost, instance.openapi, "/api/v1/import/openapi").
		Route(http.MethodPost, instance.importCollection, "/api/v1/import/collection").
		Route(http.MethodPost, instance.importToCollection, "/api/v1/import/collection/{%s}", COLLECTION).
		Route(http.MethodGet, instance.findCollection, "/api/v1/collection").
		Route(http.MethodPost, instance.insertCollection, "/api/v1/collection").
		Route(http.MethodPut, instance.pushToCollection, "/api/v1/collection").
		Route(http.MethodDelete, instance.deleteCollection, "/api/v1/collection/{%s}", COLLECTION).
		Route(http.MethodPost, instance.cloneCollection, "/api/v1/collection/{%s}/clone", COLLECTION).
		Route(http.MethodDelete, instance.deleteFromCollection, "/api/v1/collection/{%s}/request/{%s}", COLLECTION, ID_REQUEST).
		Route(http.MethodPut, instance.takeFromCollection, "/api/v1/collection/{%s}/request/{%s}", COLLECTION, ID_REQUEST)

	return instance
}

func (c *ControllerCollection) importCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	dtos, err := jsonDeserialize[[]dto.DtoCollection](r)
	if err != nil {
		return err
	}

	collection, err := c.managerCollection.ImportDtoCollections(user, *dtos)
	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(collection)

	return nil
}

func (c *ControllerCollection) importToCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return nil
	}

	dtos, err := jsonDeserialize[[]dto.DtoRequest](r)
	if err != nil {
		return err
	}

	collection, err := c.managerCollection.ImportDtoRequests(user, idCollection, *dtos)
	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(collection)

	return nil
}

func (c *ControllerCollection) openapi(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	collection, err := c.managerCollection.ImportOpenApi(user, data)
	if err != nil {
		return err
	}

	json.NewEncoder(w).Encode(collection)

	return nil
}

func (c *ControllerCollection) findCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	collections := c.managerCollection.FindByOwner(user)

	dtos := make([]dto.DtoCollection, len(collections))
	for i, v := range collections {
		requests := c.managerActions.FindNodes(v.Nodes)
		context, _ := c.repositoryContext.Find(v.Context)
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

	json.NewEncoder(w).Encode(dtos)

	return nil
}

func (c *ControllerCollection) insertCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	collection, err := jsonDeserialize[domain.Collection](r)
	if err != nil {
		return err
	}

	collection = c.managerCollection.Insert(user, collection)

	json.NewEncoder(w).Encode(collection)

	return nil
}

func (c *ControllerCollection) pushToCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	request, err := jsonDeserialize[RequestPushToCollection](r)
	if err != nil {
		return err
	}

	payload := RequestPushToCollectionToPayload(request)

	collection := c.managerCollection.PushToCollection(user, payload)

	json.NewEncoder(w).Encode(collection)

	return nil
}

func (c *ControllerCollection) deleteCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return nil
	}

	collection := c.managerCollection.Delete(user, idCollection)

	json.NewEncoder(w).Encode(collection)

	return nil
}

func (c *ControllerCollection) cloneCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return nil
	}

	payload, err := jsonDeserialize[RequestCloneCollection](r)
	if err != nil {
		return err
	}

	collection := c.managerCollection.CloneCollection(user, idCollection, payload.CollectionName)

	json.NewEncoder(w).Encode(collection)

	return nil
}

func (c *ControllerCollection) deleteFromCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return nil
	}

	idRequest := r.PathValue(ID_REQUEST)
	if idRequest == "" {
		return nil
	}

	collection, _ := c.managerCollection.RemoveFromCollection(user, idCollection, idRequest)

	json.NewEncoder(w).Encode(collection)

	return nil
}

func (c *ControllerCollection) takeFromCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	idCollection := r.PathValue(COLLECTION)
	if idCollection == "" {
		return nil
	}

	idRequest := r.PathValue(ID_REQUEST)
	if idRequest == "" {
		return nil
	}

	collection, _ := c.managerCollection.TakeFromCollection(user, idCollection, idRequest)

	json.NewEncoder(w).Encode(collection)

	return nil
}
