package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
)

type Context = *collection.CollectionMap[string, any]
type contextHandler = func(http.ResponseWriter, *http.Request) (Context, error)
type requestHandler = func(http.ResponseWriter, *http.Request, Context) error
type errorHandler = func(http.ResponseWriter, *http.Request, Context, error)

type Router struct {
	contextualizer       collection.CollectionMap[string, contextHandler]
	groupContextualizers collection.CollectionMap[string, collection.CollectionList[contextHandler]]
	errors               collection.CollectionMap[string, errorHandler]
	routes               collection.CollectionMap[string, requestHandler]
}

func NewRouter() *Router {
	return &Router{
		contextualizer:       *collection.EmptyMap[string, contextHandler](),
		groupContextualizers: *collection.EmptyMap[string, collection.CollectionList[contextHandler]](),
		errors:               *collection.EmptyMap[string, errorHandler](),
		routes:               *collection.EmptyMap[string, requestHandler](),
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

func (r *Router) GroupContextualizer(group string, handler contextHandler) *Router {
	result := r.groupContextualizers.
		ComputeIfAbsent(group, *collection.EmptyList[contextHandler]()).
		Append(handler)
	r.groupContextualizers.Put(group, *result)
	return r
}

func (r *Router) ErrorHandler(handler errorHandler) *Router {
	r.errors.Put("$BASE", handler)
	return r
}

func (r *Router) RouteOptions(method string, handler requestHandler, contextualizer *contextHandler, error *errorHandler, pattern string, params ...any) *Router {
	route := r.patternKey(method, pattern, params...)
	if contextualizer != nil {
		r.contextualizer.Put(route, *contextualizer)
	}
	if error != nil {
		r.errors.Put(route, *error)
	}
	return r.Route(method, handler, pattern, params...)
}

func (r *Router) Route(method string, handler requestHandler, pattern string, params ...any) *Router {
	route := r.patternKey(method, pattern, params...)
	r.routes.Put(route, handler)
	http.HandleFunc(route, r.handler)
	return r
}

func (r Router) patternKey(method, pattern string, params ...any) string {
	return fmt.Sprintf("%s %s", method, fmt.Sprintf(pattern, params...))
}

func (r *Router) Listen(host string) error {
	println(fmt.Sprintf("Listen at: %s", host))
	return http.ListenAndServe(host, nil)
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

	group := strings.Split(req.Pattern, " ")[1]

	keys := r.groupContextualizers.KeysCollection().Filter(func(key string) bool {
		return strings.HasPrefix(group, key)
	})

	keys.ForEach(func(i int, key string) {
		funcs, ok := r.groupContextualizers.Find(key)
		if !ok {
			return
		}
		
		funcs.ForEach(func(i int, f contextHandler) {
			result, err := f(wrt, req)
			if err != nil {
				panic("//TODO: contextualizer error.")
			}
			context.Merge(result.Collect())
		})
	})

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
		return
	}

	http.Error(wrt, err.Error(), http.StatusInternalServerError)
}
