package controller

import (
	"net/http"
	"strings"

	"github.com/Rafael24595/go-api-core/src/domain/action"
	"github.com/Rafael24595/go-api-core/src/domain/context"
	"github.com/Rafael24595/go-api-core/src/domain/formatter/curl"
	"github.com/Rafael24595/go-api-core/src/domain/mock"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const SW_INLINE = "inline"
const SW_INLINE_DESCRIPTION = "Inline flag"

const SW_RAW = "raw"
const SW_RAW_DESCRIPTION = "Raw flag"

const CURL_COMMAND_DESCRIPTION = "CURL commands"

type ControllerCurl struct {
	router            *router.Router
	managerRequest    *repository.ManagerRequest
	managerCollection *repository.ManagerCollection
	managerGroup      *repository.ManagerGroup
	managerContext    *repository.ManagerContext
	managerEndPoint   *repository.ManagerEndPoint
	managerClientData *repository.ManagerClientData
}

func NewControllerCurl(
	router *router.Router,
	managerRequest *repository.ManagerRequest,
	managerCollection *repository.ManagerCollection,
	managerGroup *repository.ManagerGroup,
	managerContext *repository.ManagerContext,
	managerEndPoint *repository.ManagerEndPoint,
	managerClientData *repository.ManagerClientData) ControllerCurl {
	instance := ControllerCurl{
		router:            router,
		managerRequest:    managerRequest,
		managerGroup:      managerGroup,
		managerCollection: managerCollection,
		managerContext:    managerContext,
		managerEndPoint:   managerEndPoint,
		managerClientData: managerClientData,
	}

	router.
		RouteDocument(http.MethodGet, instance.encodeCurl, "curl/request/{%s}", instance.docEncodeCurl()).
		RouteDocument(http.MethodPost, instance.decodeCurl, "curl/request", instance.docDecodeCurl()).
		RouteDocument(http.MethodGet, instance.encodeEndPoint, "curl/endpoint/{%s}", instance.docEncodeEndPoint())

	return instance
}

func (c *ControllerCurl) docEncodeCurl() docs.DocRoute {
	return docs.DocRoute{
		Description: "Generates a cURL command representing a previously saved HTTP request. Optionally applies a specific context for variable resolution and environment configuration. Supports raw and inline modes for flexible output formatting.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(ID_REQUEST, ID_REQUEST_DESCRIPTION),
		},
		Query: docs.DocParameters{
			ID_CONTEXT: ID_CONTEXT_DESCRIPTION,
			SW_INLINE:  SW_INLINE_DESCRIPTION,
			SW_RAW:     SW_RAW_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseAction](),
		},
	}
}

func (c *ControllerCurl) encodeCurl(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	idRequest := r.PathValue(ID_REQUEST)
	if idRequest == "" {
		return result.Reject(http.StatusNotFound)
	}

	request, _, ok := c.managerRequest.Find(user, idRequest)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	swInline := r.URL.Query().Get(SW_INLINE)
	inline := strings.ToLower(swInline) == "true"

	swRaw := r.URL.Query().Get(SW_RAW)
	if strings.ToLower(swRaw) == "true" {
		return c.toCurl(request, inline)
	}

	context_id := r.URL.Query().Get(ID_CONTEXT)
	if context_id == "" {
		collection, resultStatus := findPersistentCollection(user, c.managerClientData)
		if resultStatus != nil {
			return *resultStatus
		}
		context_id = collection.Context
	}

	context, ok := c.managerContext.Find(user, context_id)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	return c.toCurlWithContext(context, request, inline)
}

func (c *ControllerCurl) docDecodeCurl() docs.DocRoute {
	return docs.DocRoute{
		Description: "Parses and imports an HTTP request from a cURL command. The provided cURL string is decoded into a structured request object and stored in the user's general collection.",
		Query: docs.DocParameters{
			ID_COLLECTION: ID_COLLECTION_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]string](CURL_COMMAND_DESCRIPTION),
		},
	}
}

func (c *ControllerCurl) decodeCurl(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	curls, res := router.InputJson[[]string](r)
	if res != nil {
		return *res
	}

	reqs := make([]action.Request, 0)
	for _, v := range curls {
		req, err := curl.Unmarshal([]byte(v))
		if err != nil {
			return result.Err(http.StatusUnprocessableEntity, err)
		}
		reqs = append(reqs, *req)
	}

	if coll := r.URL.Query().Get(ID_COLLECTION); coll != "" {
		return c.decodeCurlToCollection(user, coll, reqs)
	}

	return c.decodeCurlToGlobal(user, reqs)
}

func (c *ControllerCurl) docEncodeEndPoint() docs.DocRoute {
	return docs.DocRoute{
		Description: "Generates a cURL command representing a previously saved mock end-point.",
		Parameters: docs.DocOrderParameters{
			docs.Parameter(ID_END_POINT, ID_END_POINT_DESCRIPTION),
		},
		Query: docs.DocParameters{
			SW_INLINE: SW_INLINE_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseAction](),
		},
	}
}

func (c *ControllerCurl) encodeEndPoint(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	idEndPoint := r.PathValue(ID_END_POINT)
	if idEndPoint == "" {
		return result.Reject(http.StatusNotFound)
	}

	endPoint, ok := c.managerEndPoint.Find(user, idEndPoint)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	swInline := r.URL.Query().Get(SW_INLINE)
	inline := strings.ToLower(swInline) == "true"

	return c.endPointToCurl(endPoint, inline)
}

func (c *ControllerCurl) toCurl(request *action.Request, inline bool) result.Result {
	curl, err := curl.Marshal(request, inline)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}
	return result.Ok(curl)
}

func (c *ControllerCurl) toCurlWithContext(context *context.Context, request *action.Request, inline bool) result.Result {
	curl, err := curl.MarshalContext(context, request, inline)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	return result.Ok(curl)
}

func (c *ControllerCurl) endPointToCurl(endPoint *mock.EndPoint, inline bool) result.Result {
	request := endPointToRequest(endPoint)
	
	curl, err := curl.Marshal(request, inline)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	return result.Ok(curl)
}

func (c *ControllerCurl) decodeCurlToCollection(user, coll string, reqs []action.Request) result.Result {
	group, resultStatus := findUserCollections(user, c.managerClientData)
	if resultStatus != nil {
		return *resultStatus
	}

	c.managerGroup.ImportRequestsById(user, group, coll, reqs...)

	return result.Ok(reqs)
}

func (c *ControllerCurl) decodeCurlToGlobal(user string, reqs []action.Request) result.Result {
	collection, resultStatus := findPersistentCollection(user, c.managerClientData)
	if resultStatus != nil {
		return *resultStatus
	}

	c.managerCollection.ImportRequests(user, collection, reqs...)

	return result.Ok(reqs)
}
