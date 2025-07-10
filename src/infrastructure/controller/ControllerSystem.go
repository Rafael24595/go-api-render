package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/docs"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

type ControllerSystem struct {
	router *router.Router
}

func NewControllerSystem(router *router.Router) ControllerSystem {
	instance := ControllerSystem{
		router: router,
	}

	router.
		Route(http.MethodGet, instance.log, "system/log").
		RouteDocument(http.MethodGet, instance.metadata, "system/metadata", docs.IDocPayload{
			Responses: map[string]any{
				"200": responseSystemMetadata{},
			},
		})

	return instance
}

func (c *ControllerSystem) log(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	session, res := findSession(user)
	if res != nil {
		return *res
	}

	if !session.IsAdmin {
		return result.Err(http.StatusUnauthorized, nil)
	}

	return result.Ok(log.Records())
}

func (c *ControllerSystem) metadata(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	conf := configuration.Instance()
	response := makeResponseSystemMetadata(
		conf.SessionId(),
		conf.Timestamp(),
		conf.Release,
		conf.Mod,
		conf.Project,
		conf.Front)
	return result.Ok(response)
}
