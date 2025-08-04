package controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
)

const TETRIS_PROJECT = "https://raw.githubusercontent.com/Rafael24595/js-tetris/refs/heads/main"
const JS_RESOURCE = "js_resource"

type ControllerSecret struct {
	router *router.Router
	cache map[string]string
}

func NewControllerSecret(router *router.Router) ControllerSecret {
	instance := ControllerSecret{
		router: router,
		cache: make(map[string]string),
	}

	router.
		Route(http.MethodGet, instance.jsTetris, "secret/js-tetris/play").
		Route(http.MethodGet, instance.jsTetris, "secret/js-tetris/{%s}", JS_RESOURCE)

	return instance
}

func (c *ControllerSecret) jsTetris(w http.ResponseWriter, r *http.Request, ctx router.Context) result.Result {
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
		return result.Err(500, err)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return result.Err(500, err)
	}

	resource := string(bodyBytes)

	c.cache[jsResource] = resource
	
	return result.Ok(resource)
}
