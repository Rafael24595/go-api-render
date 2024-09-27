package router

import (
	"fmt"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
)

type Context = *collection.CollectionMap[string, any]
type contextHandler = func(http.ResponseWriter, *http.Request) (Context, error)
type requestHandler = func(http.ResponseWriter, *http.Request, Context) error
type errorHandler = func(http.ResponseWriter, *http.Request, Context, error)

type Router struct {
	contextualizer collection.CollectionMap[string, contextHandler]
	errors collection.CollectionMap[string, errorHandler]
	routes collection.CollectionMap[string, requestHandler]
}

func NewRouter() *Router {
	return &Router {
		contextualizer: *collection.EmptyMap[string, contextHandler](),
		errors: *collection.EmptyMap[string, errorHandler](),
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

func (r *Router) ErrorHandler(handler errorHandler) *Router {
	r.errors.Put("$BASE", handler)
	return r
}

func (r *Router) RouteOptions(method string, pattern string, handler requestHandler, contextualizer *contextHandler, error *errorHandler) *Router {
	route := r.patternKey(method, pattern)
	if contextualizer != nil {
		r.contextualizer.Put(route, *contextualizer)
	}
	if error != nil {
		r.errors.Put(route, *error)
	}
	return r.Route(method, pattern, handler)
}

func (r *Router) Route(method string, pattern string, handler requestHandler) *Router {
	route := r.patternKey(method, pattern)
	r.routes.Put(route, handler)
	http.HandleFunc(route, r.handler)
	return r
}

func (r Router) patternKey(method, pattern string) string {
	return fmt.Sprintf("%s %s", method, pattern)
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

	err := (*handler)(wrt, req, context)
	if err == nil {
		return
	}

	errorHandler, ok := r.errors.Find(req.Pattern)
	if !ok {
		errorHandler, ok = r.errors.Find("$BASE")
	}

	if ok {
		(*errorHandler)(wrt, req, context, err)
	}

	http.Error(wrt, err.Error(), http.StatusInternalServerError)
}
