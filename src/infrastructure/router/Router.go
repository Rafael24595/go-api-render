package router

import (
	"fmt"
	"net/http"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router {
	}
}

func (router *Router) ResourcesPath(path string) *Router {
	fs := http.FileServer(http.Dir(path))
	route := fmt.Sprintf("/%s/", path)
    http.Handle(fmt.Sprintf("GET %s", route), http.StripPrefix(route, fs))
	return router
}

func (router *Router) Route(method string, pattern  string, handler func(http.ResponseWriter, *http.Request)) *Router {
	http.HandleFunc(fmt.Sprintf("%s %s", method, pattern), handler)
	return router
}

func (router *Router) Listen(host string) error {
	return http.ListenAndServe(":8080", nil)
}