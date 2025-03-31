package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

const ()

type ControllerContext struct {
	router            *router.Router
	repositoryContext repository.IRepositoryContext
}

func NewControllerContext(
	router *router.Router,
	repositoryContext repository.IRepositoryContext) ControllerContext {
	instance := ControllerContext{
		router:            router,
		repositoryContext: repositoryContext,
	}

	instance.router.
		Route(http.MethodGet, instance.findContext, "/api/v1/context").
		Route(http.MethodPost, instance.insertContext, "/api/v1/context")

	return instance
}

func (c *ControllerContext) findContext(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	context, ok := c.repositoryContext.FindByOwner(user)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	dtoContext := dto.FromContext(context)

	json.NewEncoder(w).Encode(dtoContext)

	return nil
}

func (c *ControllerContext) insertContext(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	dtoContext, err := jsonDeserialize[dto.DtoContext](r)
	if err != nil {
		return err
	}

	context := dto.ToContext(dtoContext)
	context = c.repositoryContext.InsertFromOwner(user, context)

	dtoContext = dto.FromContext(context)

	json.NewEncoder(w).Encode(dtoContext)

	return nil
}
