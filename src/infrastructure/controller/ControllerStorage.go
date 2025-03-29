package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

const USER = "user"
const COLLECTION = "collection"

type ControllerStorage struct {
	router               *router.Router
	repositoryContext    repository.IRepositoryContext
	repositoryActions    *repository.RequestManager
	repositoryHisotric   repository.IRepositoryHistoric
	repositoryCollection repository.IRepositoryCollection
}

func NewControllerStorage(
	router *router.Router,
	repositoryContext repository.IRepositoryContext,
	repository *repository.RequestManager,
	repositoryHisotric repository.IRepositoryHistoric,
	repositoryCollection repository.IRepositoryCollection) ControllerStorage {
	instance := ControllerStorage{
		router:             router,
		repositoryContext:  repositoryContext,
		repositoryActions:  repository,
		repositoryHisotric: repositoryHisotric,
		repositoryCollection: repositoryCollection,
	}

	//TODO: Extract users from token.
	router.
		Route(http.MethodGet, instance.findContext, "/api/v1/context/{%s}", USER).
		Route(http.MethodPost, instance.insertContext, "/api/v1/context/{%s}", USER).
		Route(http.MethodPost, instance.storage, "/api/v1/storage/{%s}", USER).
		Route(http.MethodGet, instance.findAll, "/api/v1/storage/{%s}", USER).
		Route(http.MethodDelete, instance.delete, "/api/v1/storage/{%s}/{%s}", USER, ID_REQUEST).
		Route(http.MethodGet, instance.find, "/api/v1/storage/{%s}/{%s}", USER, ID_REQUEST).
		Route(http.MethodPost, instance.historic, "/api/v1/historic/{%s}", USER).
		Route(http.MethodGet, instance.findHistoric, "/api/v1/historic/{%s}", USER).
		Route(http.MethodGet, instance.findCollection, "/api/v1/collection/{%s}", USER).
		Route(http.MethodPost, instance.insertCollection, "/api/v1/collection/{%s}", USER)

	return instance
}

func (c *ControllerStorage) findContext(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	context, ok := c.repositoryContext.FindByOwner(user)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	dtoContext := dto.FromContext(context)

	json.NewEncoder(w).Encode(dtoContext)

	return nil
}

func (c *ControllerStorage) insertContext(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	dtoContext, err := jsonDeserialize[dto.DtoContext](r)
	if err != nil {
		return err
	}

	context := dto.ToContext(dtoContext)
	context = c.repositoryContext.Insert(user, context)

	dtoContext = dto.FromContext(context)

	json.NewEncoder(w).Encode(dtoContext)

	return nil
}

func (c *ControllerStorage) storage(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	action, err := jsonDeserialize[RequestInsertAction](r)
	if err != nil {
		return err
	}

	actionRequest, actionResponse := c.repositoryActions.Release(user, dto.ToRequest(&action.Request), dto.ToResponse(&action.Response))

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) historic(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	action, err := jsonDeserialize[RequestInsertAction](r)
	if err != nil {
		return err
	}

	if action.Request.Status != domain.DRAFT {
		return nil
	}

	actionRequest, actionResponse := c.repositoryActions.Insert(user, dto.ToRequest(&action.Request), dto.ToResponse(&action.Response))

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

func (c *ControllerStorage) findAll(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	actions := c.repositoryActions.FindOptions(repository.FilterOptions[domain.Request]{
		Predicate: func(r domain.Request) bool {
			return r.Status == domain.FINAL
		},
	})

	json.NewEncoder(w).Encode(actions)

	return nil
}

func (c *ControllerStorage) delete(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse := c.repositoryActions.DeleteById(idRequest)

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) find(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	idRequest := r.PathValue(ID_REQUEST)

	actionRequest, actionResponse, ok := c.repositoryActions.Find(idRequest)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *ControllerStorage) findHistoric(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	steps := c.repositoryHisotric.FindByOwner(user)
	requests := c.repositoryActions.FindSteps(steps)

	json.NewEncoder(w).Encode(requests)

	return nil
}

func (c *ControllerStorage) findCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	collections := c.repositoryCollection.FindByOwner(user)

	dtos := make([]dto.DtoCollection, len(collections))
	for i, v := range collections {
		requests := c.repositoryActions.FindNodes(v.Nodes)
		context, _ := c.repositoryContext.FindByCollection(v.Owner, v.Name)
		dtoContext := dto.FromContext(context)
		dtos[i] = dto.DtoCollection{
			Id: v.Id,
			Name: v.Name,
			Timestamp: v.Timestamp,
			Context: *dtoContext,
			Nodes: requests,
			Owner: v.Owner,
			Modified: v.Modified,
		}
	}


	json.NewEncoder(w).Encode(dtos)

	return nil
}

func (c *ControllerStorage) insertCollection(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := r.PathValue(USER)
	if user == "" {
		user = domain.ANONYMOUS_OWNER
	}

	collection, err := jsonDeserialize[domain.Collection](r)
	if err != nil {
		return err
	}

	collection = c.repositoryCollection.Insert(user, collection)

	json.NewEncoder(w).Encode(collection)

	return nil
}
