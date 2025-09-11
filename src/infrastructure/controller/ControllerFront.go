package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
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
		RouteDocument(http.MethodGet, instance.client, "/", instance.docClient())

	return instance
}

func (c *ControllerFront) docClient() docs.DocRoute {
	return docs.DocRoute{
		Description: "Serves frontend static files. Falls back to index.html for SPA routing.",
		Responses: docs.DocResponses{
			"200": docs.DocText(),
		},
	}
}

func (c *ControllerFront) client(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
	path := "./assets/front" + r.URL.Path
	isPackage := strings.HasPrefix(path, "package.json")
	if _, err := os.Stat(path); err != nil || isPackage {
		path = "./assets/front/index.html"
	}

	http.ServeFile(w, r, path)

	return result.Ok(nil)
}
