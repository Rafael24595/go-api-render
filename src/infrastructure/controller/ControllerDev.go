package controller

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/docs"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

const QUERY_TIME = "time"
const QUERY_TIME_DESCRIPTION = "Time in milliseconds"

type ControllerDev struct {
	router *router.Router
}

func NewControllerDev(router *router.Router) ControllerDev {
	instance := ControllerDev{
		router: router,
	}

	router.
		RouteDocument(http.MethodGet, instance.wait, "dev/wait", instance.doWait()).
		RouteDocument(http.MethodPost, instance.paylaod, "dev/print/payload", instance.doPayload())

	return instance
}

func (c *ControllerDev) doWait() docs.DocPayload {
	return docs.DocPayload{
		Description: "Simulates a delayed response by waiting for a specified number of milliseconds.",
		Query: docs.DocParameters{
			QUERY_TIME: QUERY_TIME_DESCRIPTION,
		},
		Tags: docs.DocTags("dev", "debug"),
	}
}

func (c *ControllerDev) wait(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	timeValue := r.URL.Query().Get(QUERY_TIME)

	millis, err := strconv.Atoi(timeValue)
	if err == nil && millis > 0 {
		time.Sleep(time.Duration(millis) * time.Millisecond)
	}

	return result.Ok(nil)
}

func (c *ControllerDev) doPayload() docs.DocPayload {
	return docs.DocPayload{
		Description: "Reads and returns the raw request payload as plain text for debugging.",
		Request:     docs.DocText(),
		Tags:        docs.DocTags("dev", "debug"),
	}
}

func (c *ControllerDev) paylaod(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return result.Err(500, err)
	}

	return result.Ok(string(bodyBytes))
}
