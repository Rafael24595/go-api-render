package router

import (
	"fmt"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
)

type Context = *collection.CollectionMap[string, any]
type contextHandler = func(http.ResponseWriter, *http.Request) (Context, error)
type requestHandler = func(http.ResponseWriter, *http.Request, Context)

type Router struct {
	contextualizer collection.CollectionMap[string, contextHandler]
	routes collection.CollectionMap[string, requestHandler]
}

func NewRouter() *Router {
	return &Router {
		contextualizer: *collection.EmptyMap[string, contextHandler](),
		routes: *collection.EmptyMap[string, requestHandler](),
	}
}

func (r *Router) ResourcesPath(path string) *Router {
	fs := http.FileServer(http.Dir(path))
	route := fmt.Sprintf("/%s/", path)
    http.Handle(fmt.Sprintf("GET %s", route), http.StripPrefix(route, fs))
	return r
}

func (r *Router) Contextualizer(handler contextHandler) *Router {
	r.contextualizer.Put("$BASE", handler)
	return r
}

func (r *Router) Route(method string, pattern  string, handler requestHandler, contextualizer *contextHandler) *Router {
	route := fmt.Sprintf("%s %s", method, pattern)
	r.routes.Put(route, handler)
	http.HandleFunc(route, r.handler)
	return r
}

func (r *Router) Listen(host string) error {
	return http.ListenAndServe(":8080", nil)
}

func (r *Router) handler(wrt http.ResponseWriter, req *http.Request) {
	handler, ok := r.routes.Find(req.Pattern)
	if !ok {
		panic("//TODO: handler not found.")
	}
	contextualizer, ok := r.contextualizer.Find(req.Pattern)
	if !ok {
		contextualizer, ok = r.contextualizer.Find("$BASE")
	}

	context := collection.EmptyMap[string, any]()
	if ok {
		var err error
		context, err = (*contextualizer)(wrt, req)
		if err != nil {
			panic("//TODO: contextualizer error.")
		}
	}

	(*handler)(wrt, req, context)
}