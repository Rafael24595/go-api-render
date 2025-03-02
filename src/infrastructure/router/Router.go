package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Rafael24595/go-collections/collection"
)

type Context = collection.IDictionary[string, any]
type contextHandler = func(http.ResponseWriter, *http.Request) (Context, error)
type requestHandler = func(http.ResponseWriter, *http.Request, Context) error
type errorHandler = func(http.ResponseWriter, *http.Request, Context, error)

type Router struct {
	contextualizer       collection.IDictionary[string, contextHandler]
	groupContextualizers collection.IDictionary[string, collection.Vector[contextHandler]]
	errors               collection.IDictionary[string, errorHandler]
	routes               collection.IDictionary[string, requestHandler]
}

func NewRouter() *Router {
	return &Router{
		contextualizer:       collection.DictionaryEmpty[string, contextHandler](),
		groupContextualizers: collection.DictionaryEmpty[string, collection.Vector[contextHandler]](),
		errors:               collection.DictionaryEmpty[string, errorHandler](),
		routes:               collection.DictionaryEmpty[string, requestHandler](),
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
	result, _ := r.groupContextualizers.
		PutIfAbsent(group, *collection.VectorEmpty[contextHandler]())
	result.Append(handler)
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
	handler, ok := r.routes.Get(req.Pattern)
	if !ok {
		panic("//TODO: handler not found.")
	}
	contextualizer, ok := r.contextualizer.Get(req.Pattern)
	if !ok {
		contextualizer, ok = r.contextualizer.Get("$BASE")
	}

	var context Context
	context = collection.DictionaryEmpty[string, any]()
	if ok {
		var err error
		context, err = (*contextualizer)(wrt, req)
		if err != nil {
			panic("//TODO: contextualizer error.")
		}
	}

	group := strings.Split(req.Pattern, " ")[1]

	keys := r.groupContextualizers.KeysVector().Filter(func(key string) bool {
		return strings.HasPrefix(group, key)
	})

	keys.ForEach(func(i int, key string) {
		funcs, ok := r.groupContextualizers.Get(key)
		if !ok {
			return
		}
		
		funcs.ForEach(func(i int, f contextHandler) {
			result, err := f(wrt, req)
			if err != nil {
				panic("//TODO: contextualizer error.")
			}
			context.Merge(result)
		})
	})

	err := (*handler)(wrt, req, context)
	if err == nil {
		return
	}

	errorHandler, ok := r.errors.Get(req.Pattern)
	if !ok {
		errorHandler, ok = r.errors.Get("$BASE")
	}

	if ok {
		(*errorHandler)(wrt, req, context, err)
		return
	}

	http.Error(wrt, err.Error(), http.StatusInternalServerError)
}
