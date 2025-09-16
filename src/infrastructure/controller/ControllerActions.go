package controller

import (
	"net/http"

	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const ID_REQUEST = "id_request"
const ID_REQUEST_DESCRIPTION = "Request ID"

type ControllerActions struct {
	router *router.Router
}

func NewControllerActions(router *router.Router) ControllerActions {
	instance := ControllerActions{
		router: router,
	}

	router.
		RouteDocument(http.MethodPost, instance.action, "action", instance.docAction())

	return instance
}

func (c *ControllerActions) docAction() docs.DocRoute {
	return docs.DocRoute{
		Description: "Executes an HTTP action using a custom context and request configuration. This simulates a request as it would be processed by the client, returning the full request and response objects.",
		Request:     docs.DocJsonPayload[requestExecuteAction](),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseAction](),
		},
	}
}

func (c *ControllerActions) action(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	actionData, res := router.InputJson[requestExecuteAction](r)
	if res != nil {
		return *res
	}

	actionContext := dto.ToContext(&actionData.Context)
	actionRequest := dto.ToRequest(&actionData.Request)

	actionResponse, err := core_infrastructure.Client().
		FetchWithContext(actionContext, actionRequest)
	if err != nil {
		return result.Err(err.Status, err)
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.JsonOk(response)
}
