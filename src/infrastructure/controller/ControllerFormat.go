package controller

import (
	"net/http"
	"strings"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/domain/context"
	"github.com/Rafael24595/go-api-core/src/domain/formatter"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const SW_INLINE = "inline"
const SW_INLINE_DESCRIPTION = "Inline flag"

const SW_RAW = "raw"
const SW_RAW_DESCRIPTION = "Raw flag"

type ControllerFormat struct {
	router         *router.Router
	managerRequest *repository.ManagerRequest
	managerContext *repository.ManagerContext
}

func NewControllerFormat(
	router *router.Router,
	managerRequest *repository.ManagerRequest,
	managerContext *repository.ManagerContext) ControllerFormat {
	instance := ControllerFormat{
		router:         router,
		managerRequest: managerRequest,
		managerContext: managerContext,
	}

	router.
		RouteDocument(http.MethodGet, instance.curl, "format/{%s}/curl", instance.docCurl())

	return instance
}

func (c *ControllerFormat) docCurl() docs.DocRoute {
	return docs.DocRoute{
		Description: "Executes an HTTP action using a custom context and request configuration. This simulates a request as it would be processed by the client, returning the full request and response objects.",
		Parameters: docs.DocParameters{
			ID_REQUEST: ID_REQUEST_DESCRIPTION,
		},
		Query: docs.DocParameters{
			ID_CONTEXT: ID_CONTEXT_DESCRIPTION,
			SW_INLINE:  STATUS_CODE_DESCRIPTION,
			SW_RAW:     SW_RAW_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseAction](),
		},
	}
}

func (c *ControllerFormat) curl(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	idRequest := r.PathValue(ID_REQUEST)
	if idRequest == "" {
		return result.Reject(http.StatusNotFound)
	}

	request, _, ok := c.managerRequest.Find(user, idRequest)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	swInline := r.URL.Query().Get(SW_INLINE)
	inline := strings.ToLower(swInline) == "true"

	swRaw := r.URL.Query().Get(SW_RAW)
	if strings.ToLower(swRaw) == "true" {
		return c.toCurl(request, inline)
	}

	context_id := r.URL.Query().Get(ID_CONTEXT)
	if context_id == "" {
		collection, resultStatus := findUserCollection(user)
		if resultStatus != nil {
			return *resultStatus
		}
		context_id = collection.Context
	}

	context, ok := c.managerContext.Find(user, context_id)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	return c.toCurlWithContext(context, request, inline)
}

func (c *ControllerFormat) toCurl(request *domain.Request, inline bool) result.Result {
	curl, err := formatter.ToCurl(request, inline)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}
	return result.Ok(curl)
}

func (c *ControllerFormat) toCurlWithContext(context *context.Context, request *domain.Request, inline bool) result.Result {
	curl, err := formatter.ToCurlWithContext(context, request, inline)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	return result.Ok(curl)
}
