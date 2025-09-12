package controller

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const QUERY_TIME = "time"
const QUERY_TIME_DESCRIPTION = "Time in milliseconds"

const STATUS_CODE = "status"
const STATUS_CODE_DESCRIPTION = "HTTP status code"

type ControllerDev struct {
	router *router.Router
}

func NewControllerDev(router *router.Router) ControllerDev {
	instance := ControllerDev{
		router: router,
	}

	router.
		RouteDocument(http.MethodGet, instance.playground, "dev/playground", instance.doPlayground()).
		RouteDocument(http.MethodPost, instance.paylaod, "dev/print/payload", instance.doPayload())

	return instance
}

func (c *ControllerDev) doPlayground() docs.DocRoute {
	return docs.DocRoute{
		Description: "Simulates a request that can be programmed to change its behavior.",
		Query: docs.DocParameters{
			QUERY_TIME:  QUERY_TIME_DESCRIPTION,
			STATUS_CODE: STATUS_CODE_DESCRIPTION,
		},
		Tags: docs.DocTags("dev", "debug"),
	}
}

func (c *ControllerDev) playground(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	timeValue := r.URL.Query().Get(QUERY_TIME)
	statusValue := r.URL.Query().Get(STATUS_CODE)

	millis, err := strconv.Atoi(timeValue)
	if err == nil && millis > 0 {
		time.Sleep(time.Duration(millis) * time.Millisecond)
	}

	status, err := strconv.Atoi(statusValue)
	if err != nil {
		status = 200
	}

	return result.Accept(status)
}

func (c *ControllerDev) doPayload() docs.DocRoute {
	return docs.DocRoute{
		Description: "Reads and returns the raw request payload as plain text for debugging.",
		Request:     docs.DocText(),
		Tags:        docs.DocTags("dev", "debug"),
	}
}

func (c *ControllerDev) paylaod(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	return result.Ok(string(bodyBytes))
}
