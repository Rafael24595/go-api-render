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

func (c *ControllerContext) docImportItem() docs.DocPayload {
	return docs.DocPayload{
		Description: "Imports and merges a new context object for the authenticated user, combining the target and source contexts.",
		Request:     docs.DocJsonStruct(requestImportContext{}),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_CONTEXT_DESCRIPTION),
		},
	}
}

func (c *ControllerContext) importItem(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtos, err := jsonDeserialize[requestImportContext](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	context := c.managerContext.ImportMerge(user, &dtos.Target, &dtos.Source)

	dto := dto.FromContext(context)
	return result.Ok(dto.Id)
}

func (c *ControllerContext) docFindFromUser() docs.DocPayload {
	return docs.DocPayload{
		Description: "Retrieves the current context associated with the authenticated user, based on their collection metadata.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonStruct(dto.DtoContext{}),
		},
	}
}

func (c *ControllerContext) findFromUser(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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

	return result.Ok(dtoContext)
}

func (c *ControllerContext) docUpdate() docs.DocPayload {
	return docs.DocPayload{
		Description: "Updates an existing context object for the authenticated user using the provided context data.",
		Request:     docs.DocJsonStruct(dto.DtoContext{}),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_CONTEXT_DESCRIPTION),
		},
	}
}

func (c *ControllerContext) update(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	dtoContext, err := jsonDeserialize[dto.DtoContext](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	context := dto.ToContext(dtoContext)
	context, _ = c.managerContext.Update(user, context)

	dtoContext = dto.FromContext(context)

	return result.Ok(dtoContext.Id)
}

func (c *ControllerContext) docFind() docs.DocPayload {
	return docs.DocPayload{
		Description: "Retrieves a specific context by its ID for the authenticated user.",
		Parameters: docs.DocParameters{
			ID_CONTEXT: ID_CONTEXT_DESCRIPTION,
		},
	}
}

func (c *ControllerContext) find(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)
	idContext := r.PathValue(ID_CONTEXT)

	context, ok := c.managerContext.Find(user, idContext)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	dtoContext := dto.FromContext(context)

	return result.Ok(dtoContext)
}
