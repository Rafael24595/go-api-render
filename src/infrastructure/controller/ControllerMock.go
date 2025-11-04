package controller

import (
	"net/http"
	"slices"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/domain/mock"
	"github.com/Rafael24595/go-api-core/src/domain/token"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const OWNER_NAME = "owner"
const OWNER_NAME_DESCRIPTION = "Owner name"

const END_POINT = "point"
const END_POINT_DESCRIPTION = "End point path"

const END_POINT_ID = "endpoint"
const END_POINT_ID_DESCRIPTION = "End point ID"

type ControllerMock struct {
	router          *router.Router
	managerToken    *repository.ManagerToken
	managerEndPoint *repository.ManagerEndPoint
}

func NewControllerMock(
	router *router.Router,
	managerToken *repository.ManagerToken,
	managerEndPoint *repository.ManagerEndPoint) ControllerMock {
	instance := ControllerMock{
		router:          router,
		managerToken:    managerToken,
		managerEndPoint: managerEndPoint,
	}

	router.
		RouteDocument(http.MethodPost, instance.insert, "mock", instance.docInsert()).
		RouteDocument(http.MethodGet, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodHead, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodPost, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodPut, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodPatch, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodDelete, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodConnect, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodOptions, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodTrace, instance.call, "mock/call/{%s}/{%s...}", instance.docMockCall())

	return instance
}

func (c *ControllerMock) docInsert() docs.DocRoute {
	return docs.DocRoute{
		Description: "Creates a new mock end-poin for the authenticated user.",
		Request: docs.DocJsonPayload[mock.EndPoint](),
		Responses: docs.DocResponses{
			"200": docs.DocText(END_POINT_ID_DESCRIPTION),
		},
	}
}

func (c *ControllerMock) insert(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	endPoint, res := router.InputJson[mock.EndPoint](r)
	if res != nil {
		return *res
	}

	endPointResult := c.managerEndPoint.Insert(user, &endPoint)
	if endPointResult == nil {
		return result.TextErr(http.StatusInternalServerError, "cannot generate the token")
	}

	return result.Ok(endPointResult.Id)
}

func (c *ControllerMock) docMockCall() docs.DocRoute {
	return docs.DocRoute{
		Description: "Executes a mock HTTP action using a custom context and request configuration. This simulates a request as it would be processed by the client, returning the full request and response objects.",
		Parameters: docs.DocParameters{
			OWNER_NAME: OWNER_NAME_DESCRIPTION,
			END_POINT:  END_POINT_DESCRIPTION,
		},
	}
}

func (c *ControllerMock) call(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	method := r.Method
	owner := r.PathValue(OWNER_NAME)
	point := r.PathValue(END_POINT)

	endPoint, ok := c.managerEndPoint.FindByRequest(owner, domain.HttpMethod(method), point)
	if !ok || endPoint == nil {
		return result.Reject(http.StatusNotFound)
	}

	if res := c.authRequest(r, owner, endPoint); res.Err() {
		return res
	}

	response, res := c.findResponse(r, endPoint)
	if res != nil {
		return *res
	}

	w.WriteHeader(response.Status)

	for k, v := range response.Headers {
		w.Header().Set(k, v)
	}

	_, err := w.Write([]byte(response.Body))
	if err != nil {
		log.Errorf("Error writing response: %s", err.Error())
	}

	return result.Continue()
}

func (c *ControllerMock) authRequest(r *http.Request, owner string, endPoint *mock.EndPoint) result.Result {
	if !endPoint.Safe {
		return result.Next()
	}

	cookie, err := r.Cookie(AUTH_TOKEN)
	if err == nil {
		return result.Err(http.StatusUnauthorized, err)
	}

	if cookie == nil {
		return result.Reject(http.StatusUnauthorized)
	}

	tkn, ok := c.managerToken.FindByToken(owner, cookie.Value)
	if !ok {
		return result.Err(http.StatusForbidden, err)
	}

	if tkn.IsExipred() {
		return result.TextErr(http.StatusUnauthorized, "the provided token has expired")
	}

	if !slices.Contains(tkn.Scopes, token.ScopeMockAPI) {
		return result.TextErr(http.StatusUnauthorized, "the provided token does not have the necessary permissions")
	}

	return result.Next()
}

func (c *ControllerMock) findResponse(r *http.Request, endPoint *mock.EndPoint) (*mock.Response, *result.Result) {
	//TODO: Find one response based on the payload value or return the default.
	response := endPoint.DefaultResponse()
	return &response, nil
}
