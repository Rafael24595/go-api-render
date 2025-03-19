package controller

import (
	"encoding/json"
	"net/http"

	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

const ID_REQUEST = "id_request"

type ControllerActions struct {
	router     *router.Router
	repository *repository.RequestManager
}

func NewControllerActions(router *router.Router, repository *repository.RequestManager) ControllerActions {
	instance := ControllerActions{
		router:     router,
		repository: repository,
	}

	router.
		Route(http.MethodPost, instance.action, "/api/v1/action")

	return instance
}

func (c *ControllerActions) action(w http.ResponseWriter, r *http.Request, context router.Context) error {
	actionRequest, err := jsonDeserialize[RequestExecuteAction](r)
	if err != nil {
		return err
	}

	actionResponse, err := core_infrastructure.Client().Fetch(actionRequest.Request)
	if err != nil {
		return err
	}

	response := responseAction{
		Request:  actionRequest.Request,
		Response: *actionResponse,
	}

	json.NewEncoder(w).Encode(response)

	return nil
}
