package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	stdlog "log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/docs"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
	"github.com/Rafael24595/go-collections/collection"
)

type Context = collection.IDictionary[string, any]
type contextHandler = func(http.ResponseWriter, *http.Request) (Context, error)
type requestHandler = func(http.ResponseWriter, *http.Request, Context) result.Result
type errorHandler = func(http.ResponseWriter, *http.Request, Context, result.Result)
type startHandler func(w http.ResponseWriter, req *http.Request) bool

type logWriter struct{}

func newLogWriter() *logWriter {
	return &logWriter{}
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	log.Warningf("%s", bytes.TrimSpace(p))
	return len(p), nil
}

type Router struct {
	contextualizer       collection.IDictionary[string, contextHandler]
	groupContextualizers collection.IDictionary[string, collection.Vector[requestHandler]]
	errors               collection.IDictionary[string, errorHandler]
	routes               collection.IDictionary[string, requestHandler]
	basePath             string
	cors                 *Cors
	docViewer            docs.IDocViewer
}

func NewRouter() *Router {
	return &Router{
		contextualizer:       collection.DictionaryEmpty[string, contextHandler](),
		groupContextualizers: collection.DictionaryEmpty[string, collection.Vector[requestHandler]](),
		errors:               collection.DictionaryEmpty[string, errorHandler](),
		routes:               collection.DictionaryEmpty[string, requestHandler](),
		basePath:             "",
		cors:                 EmptyCors(),
		docViewer:            docs.NoDocViewer(),
	}
}

func (r *Router) BasePath(basePath string) *Router {
	r.basePath = basePath
	return r
}

func (r *Router) DocViewer(viewer docs.IDocViewer) *Router {
	for _, v := range viewer.Handlers() {
		pattern := fmt.Sprintf("%s %s", v.Method, v.Route)
		http.HandleFunc(pattern, v.Handler)
	}
	r.docViewer = viewer
	return r
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

func (r *Router) GroupContextualizerDocument(handler requestHandler, doc docs.DocGroup, group ...string) *Router {
	for _, v := range group {
		result, _ := r.groupContextualizers.
			PutIfAbsent(v, *collection.VectorEmpty[requestHandler]())
		result.Append(handler)
		r.groupContextualizers.Put(v, *result)
		r.docViewer.RegisterGroup(v, doc)
	}
	return r
}

func (r *Router) GroupContextualizer(handler requestHandler, group ...string) *Router {
	for _, v := range group {
		result, _ := r.groupContextualizers.
			PutIfAbsent(v, *collection.VectorEmpty[requestHandler]())
		result.Append(handler)
		r.groupContextualizers.Put(v, *result)
	}
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

func (r *Router) RouteDocument(method string, handler requestHandler, pattern string, doc docs.DocPayload) *Router {
	params := make([]any, 0)
	for p := range doc.Parameters {
		params = append(params, p)
	}

	docRoute := docs.DocRoute{
		Method:     method,
		BasePath:   r.basePath,
		Path:       fmt.Sprintf(pattern, params...),
		Parameters: doc.Parameters,
		Files:      doc.Files,
		Query:      doc.Query,
		Request:    doc.Request,
		Responses:  doc.Responses,
	}
	return r.route(method, handler, pattern, docRoute, params...)
}

func (r *Router) Route(method string, handler requestHandler, pattern string, params ...any) *Router {
	doc := docs.DocRoute{
		Method:   method,
		BasePath: r.basePath,
		Path:     fmt.Sprintf(pattern, params...),
	}
	return r.route(method, handler, pattern, doc, params...)
}

func (r *Router) route(method string, handler requestHandler, pattern string, doc docs.DocRoute, params ...any) *Router {
	route := r.patternKey(method, pattern, params...)
	r.routes.Put(route, handler)
	http.HandleFunc(route, r.handler)
	r.docViewer.RegisterRoute(doc)
	return r
}

func (r *Router) Cors(cors *Cors) *Router {
	r.cors = cors
	return r
}

func (r *Router) startHandle(next http.Handler, middlewares []startHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for _, middleware := range middlewares {
			if exit := middleware(w, req); exit {
				return
			}
		}
		next.ServeHTTP(w, req)
	})
}

func (r *Router) secureHandler(portTLS string) startHandler {
	return func(w http.ResponseWriter, req *http.Request) bool {
		if req.TLS != nil {
			return false
		}

		host := req.Host
		if colon := strings.Index(host, ":"); colon != -1 {
			host = host[:colon]
		}
		url := fmt.Sprintf("https://%s%s%s", host, portTLS, req.RequestURI)
		http.Redirect(w, req, url, http.StatusMovedPermanently)

		return true
	}
}

func (r *Router) corsHandler(w http.ResponseWriter, req *http.Request) bool {
	origin := strings.Join(r.cors.allowedOrigins, ", ")

	if origin == "*" {
		origin = req.Header.Get("Origin")
		w.Header().Set("Vary", "Origin")
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(r.cors.allowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(r.cors.allowedHeaders, ", "))
	w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(r.cors.allowCredentials))

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return true
	}

	return false
}

func (r Router) patternKey(method, pattern string, params ...any) string {
	return fmt.Sprintf("%s %s%s", method, r.basePath, fmt.Sprintf(pattern, params...))
}

func (r *Router) Listen(host string) error {
	handlers := []startHandler{
		r.corsHandler,
	}
	return r.listen(host, handlers)
}

func (r *Router) ListenTLS(hostTLS, certTLS, keyTLS string) error {
	handlers := []startHandler{
		r.corsHandler,
	}
	return r.listenTLS(hostTLS, certTLS, keyTLS, handlers)
}

func (r *Router) ListenWithTLS(host, hostTLS, certTLS, keyTLS string) error {
	handlers := []startHandler{
		r.corsHandler,
		r.secureHandler(hostTLS),
	}
	go func() {
		if err := r.listen(host, handlers); err != nil {
			log.Error(err)
		}
	}()

	return r.ListenTLS(hostTLS, certTLS, keyTLS)
}

func (r *Router) listenTLS(hostTLS, certTLS, keyTLS string, handlers []startHandler) error {
	server := &http.Server{
		Addr:     hostTLS,
		Handler:  r.startHandle(http.DefaultServeMux, handlers),
		ErrorLog: stdlog.New(newLogWriter(), "", 0),
	}

	log.Messagef("The app is listen at: %s with TLS", hostTLS)
	return server.ListenAndServeTLS(certTLS, keyTLS)
}

func (r *Router) listen(host string, handlers []startHandler) error {
	server := &http.Server{
		Addr:     host,
		Handler:  r.startHandle(http.DefaultServeMux, handlers),
		ErrorLog: stdlog.New(newLogWriter(), "", 0),
	}

	log.Messagef("The app is listen at: %s", host)
	return server.ListenAndServe()
}

func (r *Router) handler(wrt http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic: %v", r)
		}
	}()

	handler, ok := r.routes.Get(req.Pattern)
	if !ok {
		log.Panics("Request handler not found")
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
			log.Panic(err)
		}
	}

	group := strings.Split(req.Pattern, " ")[1]

	keys := r.groupContextualizers.KeysVector().Filter(func(key string) bool {
		return strings.HasPrefix(group, key)
	})

	for _, key := range keys.Collect() {
		funcs, ok := r.groupContextualizers.Get(key)
		if !ok {
			return
		}

		for _, f := range funcs.Collect() {
			result := f(wrt, req, context)
			if err, ok := result.Err(); ok {
				if err == nil {
					wrt.WriteHeader(result.Status())
					return
				}

				http.Error(wrt, err.Error(), result.Status())
				return
			}
		}
	}

	result := (*handler)(wrt, req, context)
	if response, ok := result.Ok(); ok {
		switch res := response.(type) {
		case string:
			if _, err := wrt.Write([]byte(res)); err != nil {
				log.Error(err)
			}
		case []byte:
			if _, err := wrt.Write(res); err != nil {
				log.Error(err)
			}
		case error:
			http.Error(wrt, res.Error(), http.StatusInternalServerError)
		case nil:
			return
		default:
			wrt.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(wrt).Encode(res); err != nil {
				http.Error(wrt, err.Error(), http.StatusInternalServerError)
			}
		}
		return
	}

	errorHandler, ok := r.errors.Get(req.Pattern)
	if !ok {
		errorHandler, ok = r.errors.Get("$BASE")
	}

	if ok {
		(*errorHandler)(wrt, req, context, result)
		return
	}

	if err, ok := result.Err(); ok && err != nil {
		http.Error(wrt, err.Error(), result.Status())
		return
	}

	wrt.WriteHeader(result.Status())
}
