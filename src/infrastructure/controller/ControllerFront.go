package controller

import (
	"net/http"
	"os"

	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

type ControllerFront struct {
	router *router.Router
}

func NewControllerFront(
	router *router.Router) ControllerFront {
	instance := ControllerFront{
		router: router,
	}

	router.
		Route(http.MethodGet, instance.client, "/")

	return instance
}

func (c *ControllerFront) client(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	path := "./assets/front" + r.URL.Path
	if _, err := os.Stat(path); err != nil {
		path = "./assets/front/index.html"
	}
		
	http.ServeFile(w, r, path)

	return result.Ok(nil)
}
