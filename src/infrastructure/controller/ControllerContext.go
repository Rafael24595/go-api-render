package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/docs"
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
		Route(http.MethodPost, instance.importContext, "import/context").
		Route(http.MethodGet, instance.findUserContext, "context").
		Route(http.MethodPut, instance.updateContext, "context").
		RouteDocument(http.MethodGet, instance.findContext, "context/{%s}", docs.DocPayload{
			Parameters: map[string]string{
				ID_CONTEXT: "",
			},
		})

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

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	context, ok := c.managerContext.Find(user, collection.Context)
	if !ok {
		return result.Err(http.StatusNotFound, nil)
	}

	dtoContext := dto.FromContext(context)

	return result.Ok(dtoContext)
}

func (c *ControllerContext) updateContext(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtoContext, err := jsonDeserialize[dto.DtoContext](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	context := dto.ToContext(dtoContext)
	context, _ = c.managerContext.Update(user, context)

	dtoContext = dto.FromContext(context)

	return result.Ok(dtoContext)
}

func (c *ControllerContext) findContext(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idContext := r.PathValue(ID_CONTEXT)

	context, ok := c.managerContext.Find(user, idContext)
	if !ok {
		return result.Err(http.StatusNotFound, nil)
	}

	dtoContext := dto.FromContext(context)

	return result.Ok(dtoContext)
}
