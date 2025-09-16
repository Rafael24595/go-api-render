package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const ID_CONTEXT = "id_context"
const ID_CONTEXT_DESCRIPTION = "Context ID"

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
		RouteDocument(http.MethodPost, instance.importItem, "import/context", instance.docImportItem()).
		RouteDocument(http.MethodGet, instance.findFromUser, "context", instance.docFindFromUser()).
		RouteDocument(http.MethodPut, instance.update, "context", instance.docUpdate()).
		RouteDocument(http.MethodGet, instance.find, "context/{%s}", instance.docFind())

	return instance
}

func (c *ControllerContext) docImportItem() docs.DocRoute {
	return docs.DocRoute{
		Description: "Imports and merges a new context object for the authenticated user, combining the target and source contexts.",
		Request:     docs.DocJsonPayload[requestImportContext](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_CONTEXT_DESCRIPTION),
		},
	}
}

func (c *ControllerContext) importItem(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	dtos, res := router.InputJson[requestImportContext](r)
	if res != nil {
		return *res
	}

	context := c.managerContext.ImportMerge(user, &dtos.Target, &dtos.Source)

	dto := dto.FromContext(context)
	return result.Ok(dto.Id)
}

func (c *ControllerContext) docFindFromUser() docs.DocRoute {
	return docs.DocRoute{
		Description: "Retrieves the current context associated with the authenticated user, based on their collection metadata.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[dto.DtoContext](),
		},
	}
}

func (c *ControllerContext) findFromUser(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	collection, resultStatus := findUserCollection(user)
	if resultStatus != nil {
		return *resultStatus
	}

	context, ok := c.managerContext.Find(user, collection.Context)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	dtoContext := dto.FromContext(context)

	return result.JsonOk(dtoContext)
}

func (c *ControllerContext) docUpdate() docs.DocRoute {
	return docs.DocRoute{
		Description: "Updates an existing context object for the authenticated user using the provided context data.",
		Request:     docs.DocJsonPayload[dto.DtoContext](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_CONTEXT_DESCRIPTION),
		},
	}
}

func (c *ControllerContext) update(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	dtoContext, res := router.InputJson[*dto.DtoContext](r)
	if res != nil {
		return *res
	}

	context := dto.ToContext(dtoContext)
	context, _ = c.managerContext.Update(user, context)

	dtoContext = dto.FromContext(context)

	return result.Ok(dtoContext.Id)
}

func (c *ControllerContext) docFind() docs.DocRoute {
	return docs.DocRoute{
		Description: "Retrieves a specific context by its ID for the authenticated user.",
		Parameters: docs.DocParameters{
			ID_CONTEXT: ID_CONTEXT_DESCRIPTION,
		},
	}
}

func (c *ControllerContext) find(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	idContext := r.PathValue(ID_CONTEXT)

	context, ok := c.managerContext.Find(user, idContext)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	dtoContext := dto.FromContext(context)

	return result.JsonOk(dtoContext)
}
