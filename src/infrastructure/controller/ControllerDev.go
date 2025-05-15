package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

type ControllerDev struct {
	router *router.Router
}

func NewControllerDev(router *router.Router) ControllerDev {
	instance := ControllerDev{
		router: router,
	}

	router.
		Route(http.MethodGet, instance.wait, "/api/v1/dev/wait")

	return instance
}

func (c *ControllerDev) wait(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	timeValue := r.URL.Query().Get("time")

	millis, err := strconv.Atoi(timeValue)
	if err == nil && millis > 0 {
		time.Sleep(time.Duration(millis) * time.Millisecond)
	}

	return result.Ok(nil)
}
