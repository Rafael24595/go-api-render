package controller

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/domain/action"
	"github.com/Rafael24595/go-api-core/src/domain/mock"
	"github.com/Rafael24595/go-api-core/src/domain/mock/swr"
	"github.com/Rafael24595/go-api-core/src/domain/token"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-collections/collection"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const OWNER_NAME = "owner"
const OWNER_NAME_DESCRIPTION = "Owner name"

const END_POINT = "point"
const END_POINT_DESCRIPTION = "End point path"

const ID_END_POINT = "endpoint"
const ID_END_POINT_DESCRIPTION = "End point ID"

const SWR_INPUT_DESCRIPTION = "SWR sentence"

type ControllerMock struct {
	router          *router.Router
	managerToken    *repository.ManagerToken
	managerEndPoint *repository.ManagerEndPoint
	managerMetrics  *repository.ManagerMetrics
}

func NewControllerMock(
	router *router.Router,
	managerToken *repository.ManagerToken,
	managerEndPoint *repository.ManagerEndPoint,
	managerMetrics *repository.ManagerMetrics) ControllerMock {
	instance := ControllerMock{
		router:          router,
		managerToken:    managerToken,
		managerEndPoint: managerEndPoint,
		managerMetrics:  managerMetrics,
	}

	router.
		RouteDocument(http.MethodPost, instance.bridgeStpToCnd, "bridge/mock/response/to/step", instance.docBridgeStpToCnd()).
		RouteDocument(http.MethodPost, instance.bridgeCndToStp, "bridge/mock/response/from/step", instance.docBridgeCndToStp()).
		RouteDocument(http.MethodGet, instance.bridgeEndToReq, "bridge/mock/endpoint/{%s}/to/request", instance.docBridgeEndToReq()).
		//
		RouteDocument(http.MethodPut, instance.sortEndPoint, "sort/mock/endpoint", instance.docSortEndPoint()).
		//
		RouteDocument(http.MethodGet, instance.exportAll, "export/mock/endpoint", instance.docExportAll()).
		RouteDocument(http.MethodPost, instance.exportMany, "export/mock/endpoint", instance.docExportMany()).
		RouteDocument(http.MethodPost, instance.importMany, "import/mock/endpoint", instance.docImportMany()).
		//
		RouteDocument(http.MethodGet, instance.findAll, "mock/endpoint", instance.docFindAll()).
		RouteDocument(http.MethodGet, instance.find, "mock/endpoint/{%s}", instance.docFind()).
		RouteDocument(http.MethodPost, instance.insert, "mock/endpoint", instance.docInsert()).
		RouteDocument(http.MethodDelete, instance.remove, "mock/endpoint/{%s}", instance.docRemove()).
		//
		RouteDocument(http.MethodGet, instance.findMetrics, "mock/metrics/endpoint/{%s}", instance.docFindMetrics()).
		RouteDocument(http.MethodDelete, instance.removeMetrics, "mock/metrics/endpoint/{%s}", instance.docRemoveMetrics()).
		//
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

func (c *ControllerMock) docBridgeStpToCnd() docs.DocRoute {
	return docs.DocRoute{
		Description: "Translated the text input of SWR into a Steps vector.",
		Request:     docs.DocText(SWR_INPUT_DESCRIPTION),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]swr.Step](),
			"422": docs.DocJsonPayload[[]string](),
		},
	}
}

func (c *ControllerMock) bridgeStpToCnd(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	input, res := router.InputText(r)
	if res != nil {
		return *res
	}

	steps, errs := swr.UnmarshalWithOptions(input, swr.UnmarshalOpts{
		Evalue: true,
	})
	if len(errs) > 0 {
		return result.Err(http.StatusUnprocessableEntity, errs...)
	}

	return result.JsonOk(steps)
}

func (c *ControllerMock) docBridgeCndToStp() docs.DocRoute {
	return docs.DocRoute{
		Description: "Translated the steps input into a SWR sentence.",
		Request:     docs.DocJsonPayload[[]swr.Step](),
		Responses: docs.DocResponses{
			"200": docs.DocText(SWR_INPUT_DESCRIPTION),
			"422": docs.DocJsonPayload[[]string](),
		},
	}
}

func (c *ControllerMock) bridgeCndToStp(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	steps, res := router.InputJson[[]swr.Step](r)
	if res != nil {
		return *res
	}

	output, errs := swr.MarshalWithOptions(steps, swr.MarshalOpts{
		Evalue: true,
	})
	if len(errs) > 0 {
		return result.Err(http.StatusUnprocessableEntity, errs...)
	}

	return result.TextOk(output)
}

func (c *ControllerMock) docBridgeEndToReq() docs.DocRoute {
	return docs.DocRoute{
		Description: "Translated the steps input into a SWR sentence.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(ID_END_POINT, END_POINT_DESCRIPTION),
		},
		Responses: docs.DocResponses{
			"200": docs.DocText(SWR_INPUT_DESCRIPTION),
			"422": docs.DocJsonPayload[[]string](),
		},
	}
}

func (c *ControllerMock) bridgeEndToReq(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	id := r.PathValue(ID_END_POINT)

	if id == "" {
		return result.Reject(http.StatusNotFound)
	}

	endPoint, ok := c.managerEndPoint.Find(user, id)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	request := endPointToRequest(endPoint)
	dto := dto.FromRequest(request)

	return result.JsonOk(dto)
}

func (c *ControllerMock) docSortEndPoint() docs.DocRoute {
	return docs.DocRoute{
		Description: "Sorts the user's end-point collection based on the provided lite structure.",
		Request:     docs.DocJsonPayload[[]requestNode](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_END_POINT_DESCRIPTION),
		},
	}
}

func (c *ControllerMock) sortEndPoint(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	dto, res := router.InputJson[[]requestNode](r)
	if res != nil {
		return *res
	}

	payload := requestNodeToNodeReference(dto...)

	endPoints := c.managerEndPoint.Sort(user, payload)
	ids := collection.MapToVector(endPoints, func(e mock.EndPoint) string {
		return e.Id
	})

	return result.Ok(ids.Collect())
}

func (c *ControllerMock) docExportAll() docs.DocRoute {
	return docs.DocRoute{
		Description: "Export all mock end-points for the authenticated user.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]mock.EndPoint](),
		},
	}
}

func (c *ControllerMock) exportAll(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	endPoints := c.managerEndPoint.Export(user)
	return result.JsonOk(endPoints)
}

func (c *ControllerMock) docExportMany() docs.DocRoute {
	return docs.DocRoute{
		Description: "Export defined mock end-points for the authenticated user.",
		Request:     docs.DocJsonPayload[[]string](),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]mock.EndPoint](),
		},
	}
}

func (c *ControllerMock) exportMany(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	ids, res := router.InputJson[[]string](r)
	if res != nil {
		return *res
	}

	endPoints := c.managerEndPoint.ExportList(user, ids...)
	return result.JsonOk(endPoints)
}

func (c *ControllerMock) docImportMany() docs.DocRoute {
	return docs.DocRoute{
		Description: "Import all provided mock end-points for the authenticated user.",
		Request:     docs.DocJsonPayload[[]mock.EndPoint](),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]string](),
		},
	}
}

func (c *ControllerMock) importMany(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	endPoints, res := router.InputJson[[]mock.EndPoint](r)
	if res != nil {
		return *res
	}

	ids := c.managerEndPoint.Import(user, endPoints)
	return result.JsonOk(ids)
}

func (c *ControllerMock) docFindAll() docs.DocRoute {
	return docs.DocRoute{
		Description: "Retrieves all user mock end-points.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseSignedPaylaod[[]mock.EndPoint]](),
		},
	}
}

func (c *ControllerMock) findAll(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	endPoints := c.managerEndPoint.FindAll(user)

	sign := signPayload(user, endPoints)
	return result.JsonOk(sign)
}

func (c *ControllerMock) docFind() docs.DocRoute {
	return docs.DocRoute{
		Description: "Finds a specific mock end-point by ID.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(ID_END_POINT, END_POINT_DESCRIPTION),
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[mock.EndPointFull](),
		},
	}
}

func (c *ControllerMock) find(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	id := r.PathValue(ID_END_POINT)

	endPoint, ok := c.managerEndPoint.FindFull(user, id)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	return result.JsonOk(endPoint)
}

func (c *ControllerMock) docInsert() docs.DocRoute {
	return docs.DocRoute{
		Description: "Creates a new mock end-poin for the authenticated user.",
		Request:     docs.DocJsonPayload[mock.EndPoint](),
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_END_POINT_DESCRIPTION),
		},
	}
}

func (c *ControllerMock) insert(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	endPoint, res := router.InputJson[mock.EndPointFull](r)
	if res != nil {
		return *res
	}

	oldEndPoint, _ := c.managerEndPoint.Find(user, endPoint.Id)

	endPointResult, errs := c.managerEndPoint.Insert(user, &endPoint)
	if endPointResult == nil {
		return result.Err(http.StatusUnprocessableEntity, errs...)
	}

	if oldEndPoint == nil {
		oldEndPoint = endPointResult
	}

	go c.managerMetrics.ResolveStatus(user, endPointResult, endPointResult)

	return result.Ok(endPointResult.Id)
}

func (c *ControllerMock) docRemove() docs.DocRoute {
	return docs.DocRoute{
		Description: "Deletes a specific mock end-point by ID.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(ID_END_POINT, END_POINT_DESCRIPTION),
		},
		Responses: docs.DocResponses{
			"200": docs.DocText(ID_END_POINT_DESCRIPTION),
		},
	}
}

func (c *ControllerMock) remove(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	id := r.PathValue(ID_END_POINT)

	endPoint := c.managerEndPoint.Delete(user, id)
	if endPoint == nil {
		return result.Reject(http.StatusNotFound)
	}

	return result.JsonOk(endPoint)
}

func (c *ControllerMock) docFindMetrics() docs.DocRoute {
	return docs.DocRoute{
		Description: "Finds a specific mock end-point metrics by ID.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(ID_END_POINT, END_POINT_DESCRIPTION),
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[mock.Metrics](),
		},
	}
}

func (c *ControllerMock) findMetrics(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	id := r.PathValue(ID_END_POINT)

	endPoint, ok := c.managerEndPoint.Find(user, id)
	if !ok {
		return result.Err(http.StatusNotFound)
	}

	metrics, ok := c.managerMetrics.Find(user, endPoint)
	if metrics == nil && !ok {
		return result.JsonOk(mock.EmptyMetrics(endPoint))
	}

	return result.JsonOk(metrics)
}

func (c *ControllerMock) docRemoveMetrics() docs.DocRoute {
	return docs.DocRoute{
		Description: "Deletes a specific mock end-point metrics by ID.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(ID_END_POINT, END_POINT_DESCRIPTION),
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[mock.Metrics](),
		},
	}
}

func (c *ControllerMock) removeMetrics(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)
	id := r.PathValue(ID_END_POINT)

	endPoint, ok := c.managerEndPoint.Find(user, id)
	if !ok {
		return result.JsonOk(http.StatusNotFound)
	}

	metrics := c.managerMetrics.Delete(user, endPoint)

	return result.JsonOk(metrics)
}

func (c *ControllerMock) docMockCall() docs.DocRoute {
	return docs.DocRoute{
		Description: "Executes a mock HTTP action using a custom context and request configuration. This simulates a request as it would be processed by the client, returning the full request and response objects.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(OWNER_NAME, OWNER_NAME_DESCRIPTION),
			docs.Parameter(END_POINT, END_POINT_DESCRIPTION),
		},
	}
}

func (c *ControllerMock) call(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	method := r.Method
	owner := r.PathValue(OWNER_NAME)
	point := r.PathValue(END_POINT)

	start := time.Now().UnixMilli()
	endPoint, ok := c.managerEndPoint.FindByRequest(owner, domain.HttpMethod(method), point)
	if !ok || endPoint == nil || !endPoint.Status {
		return result.Reject(http.StatusNotFound)
	}

	if res := c.authRequest(r, owner, endPoint); res.Err() {
		return res
	}

	response, res := c.findResponse(r, endPoint)
	if res != nil {
		return *res
	}

	w.WriteHeader(response.Code)

	for _, v := range response.Arguments {
		if !v.Status {
			continue
		}

		w.Header().Set(v.Key, v.Value)
	}

	contentType := response.Body.ContentType.ToHeader()
	w.Header().Set("Content-Type", contentType)

	_, err := w.Write([]byte(response.Body.Payload))
	if err != nil {
		log.Errorf("Error writing response: %s", err.Error())
	}

	end := time.Now().UnixMilli()

	go c.managerMetrics.ResolveRequest(owner, endPoint, response, end-start)

	return result.Continue()
}

func (c *ControllerMock) authRequest(r *http.Request, owner string, endPoint *mock.EndPoint) result.Result {
	if !endPoint.Safe ||
		owner == action.ANONYMOUS_OWNER && owner == endPoint.Owner {
		return result.Next()
	}

	cookie, err := r.Cookie(AUTH_TOKEN)
	if err != nil {
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
	payload, res := router.InputText(r)
	if res != nil {
		return nil, res
	}

	headers := collection.MapToDictionary(r.Header, func(k string, h []string) string {
		return strings.Join(h, ", ")
	})

	queries := collection.MapToDictionary(r.URL.Query(), func(k string, h []string) string {
		return strings.Join(h, ", ")
	})

	response, ok := mock.FindResponse(payload, headers.Merge(queries).Collect(), endPoint)
	if !ok {
		res := result.TextErr(http.StatusNotFound, "the resource does not have a defined default response")
		return nil, &res
	}

	return response, nil
}

func endPointToRequest(endPoint *mock.EndPoint) *action.Request {
	server := mockEndPointPath(endPoint)
	request := mock.ToRequest(server, endPoint)

	if endPoint.Safe {
		request.Cookie.Put(AUTH_TOKEN, AUTH_TOKEN)
	}

	return request
}

func mockEndPointPath(endPoint *mock.EndPoint) string {
	config := configuration.Instance()
	protocol := config.DefaultProtocol()
	port := config.DefaultPort()
	return fmt.Sprintf("%s://localhost:%d/api/v1/mock/call/%s", protocol, port, endPoint.Owner)
}
