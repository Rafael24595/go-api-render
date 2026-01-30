package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/application/manager"
	"github.com/Rafael24595/go-api-core/src/domain/token"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const AUTH_TOKEN = "go_user_token"
const AUTH_TOKEN_DESCRIPTION = "User cookie token"

const ID_TOKEN = "token"
const ID_TOKEN_DESCRIPTION = "Token ID"

const RAW_TOKEN_DESCRIPTION = "Raw token"

type ControllerToken struct {
	router       *router.Router
	managerToken *manager.ManagerToken
}

func NewControllerToken(
	router *router.Router,
	managerToken *manager.ManagerToken,
) ControllerToken {
	instance := ControllerToken{
		router:       router,
		managerToken: managerToken,
	}

	router.
		RouteDocument(http.MethodGet, instance.scopes, "scopes", instance.docScopes()).
		//
		RouteDocument(http.MethodGet, instance.find, "token", instance.docFind()).
		RouteDocument(http.MethodPost, instance.insert, "token", instance.docInsert()).
		RouteDocument(http.MethodDelete, instance.delete, "token/{%s}", instance.docDelete())

	return instance
}

func (c *ControllerToken) docScopes() docs.DocRoute {
	return docs.DocRoute{
		Description: "Returns all token scopes.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]token.Scope](),
		},
	}
}

func (c *ControllerToken) scopes(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	return result.JsonOk(token.ListScopes())
}

func (c *ControllerToken) docFind() docs.DocRoute {
	return docs.DocRoute{
		Description: "Returns all tokens associated with the authenticated user.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]token.LiteToken](),
		},
	}
}

func (c *ControllerToken) find(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	tokens := c.managerToken.FindAll(user)
	return result.JsonOk(tokens)
}

func (c *ControllerToken) docInsert() docs.DocRoute {
	return docs.DocRoute{
		Description: "Creates a new token for the authenticated user.",
		Request:     docs.DocJsonPayload[token.LiteToken](),
		Responses: docs.DocResponses{
			"200": docs.DocText(RAW_TOKEN_DESCRIPTION),
		},
	}
}

func (c *ControllerToken) insert(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	token, res := router.InputJson[token.LiteToken](r)
	if res != nil {
		return *res
	}

	if token.Name == "" {
		return result.TextErr(http.StatusUnprocessableEntity, "the token name is not specified")
	}

	raw, tkn := c.managerToken.Insert(user, &token)
	if tkn == nil {
		return result.TextErr(http.StatusInternalServerError, "cannot generate the token")
	}

	return result.Ok(raw)
}

func (c *ControllerToken) docDelete() docs.DocRoute {
	return docs.DocRoute{
		Description: "Deletes the token specified by its ID.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(ID_TOKEN, ID_TOKEN_DESCRIPTION),
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[token.LiteToken](),
		},
	}
}

func (c *ControllerToken) delete(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	id := r.PathValue(ID_TOKEN)
	if id == "" {
		return result.Reject(http.StatusNotFound)
	}

	token, ok := c.managerToken.DeleteById(user, id)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	return result.Ok(token)
}
