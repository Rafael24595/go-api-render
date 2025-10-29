package controller

import (
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/domain/mock"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const OWNER_NAME = "owner"
const OWNER_NAME_DESCRIPTION = "Owner name"
const END_POINT = "point"
const END_POINT_DESCRIPTION = "End point path"

type ControllerMock struct {
	router          *router.Router
	managerEndPoint *repository.ManagerEndPoint
}

func NewControllerMock(router *router.Router, managerEndPoint *repository.ManagerEndPoint) ControllerMock {
	instance := ControllerMock{
		router:          router,
		managerEndPoint: managerEndPoint,
	}

	router.
		RouteDocument(http.MethodGet, instance.call, "mock/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodHead, instance.call, "mock/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodPost, instance.call, "mock/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodPut, instance.call, "mock/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodPatch, instance.call, "mock/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodDelete, instance.call, "mock/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodConnect, instance.call, "mock/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodOptions, instance.call, "mock/{%s}/{%s...}", instance.docMockCall()).
		RouteDocument(http.MethodTrace, instance.call, "mock/{%s}/{%s...}", instance.docMockCall())

	return instance
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

func (c *ControllerMock) findResponse(r *http.Request, endPoint *mock.EndPoint) (*mock.Response, *result.Result) {
	//TODO: Find one response based on the payload value or return the default.
	response := endPoint.DefaultResponse()
	return &response, nil
}
