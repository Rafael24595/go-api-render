package controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const TETRIS_PROJECT = "https://raw.githubusercontent.com/Rafael24595/js-tetris/refs/heads/main"

const JS_RESOURCE = "js_resource"
const JS_RESOURCE_DESCRIPTION = "The name of the JavaScript, HTML, or CSS file to be retrieved from the remote Tetris project repository."

type ControllerSecret struct {
	router *router.Router
	cache  map[string]string
}

func NewControllerSecret(router *router.Router) ControllerSecret {
	instance := ControllerSecret{
		router: router,
		cache:  make(map[string]string),
	}

	router.
		RouteDocument(http.MethodGet, instance.jsTetris, "secret/js-tetris/play", instance.docPlay()).
		RouteDocument(http.MethodGet, instance.jsTetris, "secret/js-tetris/{%s}", instance.docResource())

	return instance
}

func (c *ControllerSecret) docPlay() docs.DocRoute {
	resources := c.docResource()
	return docs.DocRoute{
		Description: "Serves the main HTML page for the hidden Tetris game.",
		Responses:   resources.Responses,
	}
}

func (c *ControllerSecret) docResource() docs.DocRoute {
	return docs.DocRoute{
		Description: "Retrieves a specific static asset (HTML, JS, or CSS file) from a remote GitHub repository that hosts the Tetris game.",
		Parameters: docs.DocParameters{
			JS_RESOURCE: JS_RESOURCE_DESCRIPTION,
		},
		Responses: docs.DocResponses{
			"200": docs.DocText("The raw content of the requested static file (HTML, JS, CSS)."),
			"500": docs.DocText("Error retrieving or reading the remote file."),
		},
	}
}

func (c *ControllerSecret) jsTetris(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	jsResource := r.PathValue(JS_RESOURCE)
	if jsResource == "" {
		jsResource = "index.html"
	}

	lower := strings.ToLower(jsResource)
	if strings.HasSuffix(lower, "html") {
		w.Header().Add("Content-Type", "text/html")
	}
	if strings.HasSuffix(lower, "css") {
		w.Header().Add("Content-Type", "text/css")
	}
	if strings.HasSuffix(lower, "js") {
		w.Header().Add("Content-Type", "application/javascript")
	}

	if cached, ok := c.cache[jsResource]; ok {
		return result.Ok(cached)
	}

	path := fmt.Sprintf("%s/%s", TETRIS_PROJECT, jsResource)
	response, err := http.Get(path)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	resource := string(bodyBytes)

	c.cache[jsResource] = resource

	return result.Ok(resource)
}
