package docs

import "net/http"

type ParameterType string

const (
	QUERY ParameterType = "query"
	PATH  ParameterType = "path"
)

type IDocViewerHandler struct {
	Method  string
	Route   string
	Handler func(http.ResponseWriter, *http.Request)
}

type IDocGroup struct {
	Headers map[string]string
	Cookies map[string]string
}

type IDocPayload struct {
	Parameters map[string]string
	Query      map[string]string
	Request    any
	Responses  map[string]any
}

type IDocRoute struct {
	Method     string
	BasePath   string
	Path       string
	Parameters map[string]string
	Query      map[string]string
	Request    any
	Responses  map[string]any
}

type IDocViewer interface {
	Handlers() []IDocViewerHandler
	RegisterGroup(group string, data IDocGroup) IDocViewer
	RegisterRoute(route IDocRoute) IDocViewer
}

type noDocViewer struct {
}

func (v *noDocViewer) Handlers() []IDocViewerHandler {
	return make([]IDocViewerHandler, 0)
}

func (v *noDocViewer) RegisterRoute(route IDocRoute) IDocViewer {
	return v
}

func (v *noDocViewer) RegisterGroup(group string, data IDocGroup) IDocViewer {
	return v
}

func NoDocViewer() IDocViewer {
	return &noDocViewer{}
}
