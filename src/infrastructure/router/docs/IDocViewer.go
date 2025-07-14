package docs

import "net/http"

type ParameterType string

const (
	QUERY ParameterType = "query"
	PATH  ParameterType = "path"
)

type DocViewerHandler struct {
	Method      string
	Route       string
	Handler     func(http.ResponseWriter, *http.Request)
	Name        string
	Description string
}

type DocViewerSources struct {
	Name        string `json:"name"`
	Route       string `json:"route"`
	Description string `json:"description"`
}

type DocGroup struct {
	Headers map[string]string
	Cookies map[string]string
}

type DocPayload struct {
	Parameters map[string]string
	Query      map[string]string
	Files      map[string]string
	Request    any
	Responses  map[string]any
}

type DocRoute struct {
	Method     string
	BasePath   string
	Path       string
	Parameters map[string]string
	Query      map[string]string
	Files      map[string]string
	Request    any
	Responses  map[string]any
}

type IDocViewer interface {
	Handlers() []DocViewerHandler
	RegisterGroup(group string, data DocGroup) IDocViewer
	RegisterRoute(route DocRoute) IDocViewer
}

type noDocViewer struct {
}

func (v *noDocViewer) Handlers() []DocViewerHandler {
	return make([]DocViewerHandler, 0)
}

func (v *noDocViewer) RegisterRoute(route DocRoute) IDocViewer {
	return v
}

func (v *noDocViewer) RegisterGroup(group string, data DocGroup) IDocViewer {
	return v
}

func NoDocViewer() IDocViewer {
	return &noDocViewer{}
}
