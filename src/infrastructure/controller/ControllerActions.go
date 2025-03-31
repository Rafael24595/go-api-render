package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain/context"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
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
		Route(http.MethodPost, instance.action, "/api/v1/action")

	return instance
}

func (c *ControllerActions) action(w http.ResponseWriter, r *http.Request, ctx router.Context) error {
	actionData, err := jsonDeserialize[RequestExecuteAction](r)
	if err != nil {
		return err
	}

	actionContext := dto.ToContext(&actionData.Context)
	actionRequest := context.ProcessRequest(&actionData.Request, actionContext)

	actionResponse, err := core_infrastructure.Client().Fetch(*actionRequest)
	if err != nil {
		return err
	}

	response := responseAction{
		Request:  *dto.FromRequest(actionRequest),
		Response: *dto.FromResponse(actionResponse),
	}

	json.NewEncoder(w).Encode(response)

	return nil
}
