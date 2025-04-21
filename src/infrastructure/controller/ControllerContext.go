package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
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

func (c *ControllerContext) importContext(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtos, err := jsonDeserialize[requestImportContext](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	context := c.managerContext.ImportMerge(user, &dtos.Target, &dtos.Source)

	dto := dto.FromContext(context)
	return result.Ok(dto)
}

func (c *ControllerContext) findUserContext(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	context, ok := c.managerContext.FindByOwner(user)
	if !ok {
		return result.Err(http.StatusNotFound, nil)
	}

	dtoContext := dto.FromContext(context)

	return result.Ok(dtoContext)
}

func (c *ControllerContext) insertContext(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtoContext, err := jsonDeserialize[dto.DtoContext](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	context := dto.ToContext(dtoContext)
	context = c.managerContext.InsertFromOwner(user, context)

	dtoContext = dto.FromContext(context)

	return result.Ok(dtoContext)
}

func (c *ControllerContext) findContext(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idContext := r.PathValue(ID_CONTEXT)

	context, ok := c.managerContext.Find(user, idContext)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return result.Err(http.StatusNotFound, nil)
	}

	dtoContext := dto.FromContext(context)

	return result.Ok(dtoContext)
}
