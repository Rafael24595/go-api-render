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
		RouteDocument(http.MethodGet, instance.log, "system/log", instance.docLog()).
		RouteDocument(http.MethodGet, instance.metadata, "system/metadata", instance.docMetadata())

	return instance
}

func (c *ControllerSystem) docLog() docs.DocPayload {
	return docs.DocPayload{
		Description: "Returns all server-side application logs. Only accessible by admin users.",
		Responses: map[string]docs.DocItemStruct{
			"200": docs.DocStruct([]log.Record{}),
		},
	}
}

func (c *ControllerSystem) log(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	user := findUser(ctx)

	session, res := findSession(user)
	if res != nil {
		return *res
	}

	if !session.IsAdmin {
		return result.Err(http.StatusForbidden, nil)
	}

	return result.Ok(log.Records())
}

func (c *ControllerSystem) docMetadata() docs.DocPayload {
	return docs.DocPayload{
		Description: "Returns runtime system metadata including session ID, timestamp, release version, and frontend status.",
		Responses: map[string]docs.DocItemStruct{
			"200": docs.DocStruct(responseSystemMetadata{}),
		},
	}
}

func (c *ControllerSystem) metadata(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	conf := configuration.Instance()
	response := makeResponseSystemMetadata(
		conf.SessionId(),
		conf.Timestamp(),
		conf.Release,
		conf.Mod,
		conf.Project,
		conf.Front,
		c.router.ViewerSources())
	return result.Ok(response)
}
