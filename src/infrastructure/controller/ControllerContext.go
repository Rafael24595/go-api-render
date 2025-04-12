package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

const ID_CONTEXT = "id_context"

type ControllerContext struct {
	router         *router.Router
	managerContext *repository.ManagerContext
}

func NewControllerContext(
	router *router.Router,
	managerContext *repository.ManagerContext) ControllerContext {
	instance := ControllerContext{
		router:         router,
		managerContext: managerContext,
	}

	instance.router.
		Route(http.MethodPost, instance.importContext, "/api/v1/import/context").
		Route(http.MethodGet, instance.findUserContext, "/api/v1/context").
		Route(http.MethodPost, instance.insertContext, "/api/v1/context").
		Route(http.MethodGet, instance.findContext, "/api/v1/context/{%s}", ID_CONTEXT)

	return instance
}

func (c *ControllerContext) importContext(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	dtos, err := jsonDeserialize[requestImportContext](r)
	if err != nil {
		return err
	}

	context := c.managerContext.ImportMerge(user, &dtos.Target, &dtos.Source)

	json.NewEncoder(w).Encode(dto.FromContext(context))

	return nil
}

func (c *ControllerContext) findUserContext(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)

	context, ok := c.managerContext.FindByOwner(user)
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
	context = c.managerContext.InsertFromOwner(user, context)

	dtoContext = dto.FromContext(context)

	json.NewEncoder(w).Encode(dtoContext)

	return nil
}

func (c *ControllerContext) findContext(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	user := findUser(ctx)
	idContext := r.PathValue(ID_CONTEXT)

	context, ok := c.managerContext.Find(user, idContext)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return nil
	}

	dtoContext := dto.FromContext(context)

	json.NewEncoder(w).Encode(dtoContext)

	return nil
}

