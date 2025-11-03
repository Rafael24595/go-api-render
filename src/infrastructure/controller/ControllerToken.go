package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain/token"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const ID_TOKEN = "token"
const ID_TOKEN_DESCRIPTION = "Token ID"

const RAW_TOKEN_DESCRIPTION = "Raw token"

type ControllerToken struct {
	router       *router.Router
	managerToken *repository.ManagerToken
}

func NewControllerToken(
	router *router.Router,
	managerToken *repository.ManagerToken) ControllerToken {
	instance := ControllerToken{
		router:       router,
		managerToken: managerToken,
	}

	router.
		RouteDocument(http.MethodGet, instance.findTokens, "token", instance.docFindTokens()).
		RouteDocument(http.MethodPost, instance.insertTokens, "token", instance.docInsertToken()).
		RouteDocument(http.MethodDelete, instance.deleteToken, "token/{%s}", instance.docDeleteToken())

	return instance
}

func (c *ControllerToken) docFindTokens() docs.DocRoute {
	return docs.DocRoute{
		Description: "",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]token.LiteToken](),
		},
	}
}

func (c *ControllerToken) findTokens(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	tokens := c.managerToken.FindAll(user)
	return result.JsonOk(tokens)
}

func (c *ControllerToken) docInsertToken() docs.DocRoute {
	return docs.DocRoute{
		Description: "",
		Request:     docs.DocJsonPayload[token.LiteToken](),
		Responses: docs.DocResponses{
			"200": docs.DocText(RAW_TOKEN_DESCRIPTION),
		},
	}
}

func (c *ControllerToken) insertTokens(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	token, res := router.InputJson[token.LiteToken](r)
	if res != nil {
		return *res
	}

	raw, tkn := c.managerToken.Insert(user, &token)
	if tkn == nil {
		return result.TextErr(http.StatusInternalServerError, "cannot generate the token")
	}

	return result.Ok(raw)
}

func (c *ControllerToken) docDeleteToken() docs.DocRoute {
	return docs.DocRoute{
		Description: "",
		Parameters: docs.DocParameters{
			ID_TOKEN: ID_TOKEN_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[token.LiteToken](),
		},
	}
}

func (c *ControllerToken) deleteToken(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
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
