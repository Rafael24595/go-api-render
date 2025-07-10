package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain/context"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/docs"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

const ID_REQUEST = "id_request"

type ControllerActions struct {
	router *router.Router
}

func NewControllerActions(router *router.Router) ControllerActions {
	instance := ControllerActions{
		router: router,
	}

	router.
		RouteDocument(http.MethodPost, instance.action, "action", docs.IDocPayload{
			Request: requestExecuteAction{},
			Responses: map[string]any{
				"200": responseAction{},
			},
		})

	return instance
}

func (c *ControllerActions) action(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	actionData, err := jsonDeserialize[requestExecuteAction](r)
	if err != nil {
		return result.Err(http.StatusUnprocessableEntity, err)
	}

	actionContext := dto.ToContext(&actionData.Context)
	actionRequest := dto.ToRequest(&actionData.Request)
	actionRequest = context.ProcessRequest(actionRequest, actionContext)

	actionResponse, apiErr := core_infrastructure.Client().Fetch(*actionRequest)
	if apiErr != nil {
		return result.Err(apiErr.Status, err)
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	return result.Ok(response)
}
