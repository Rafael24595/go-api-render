package docs

import (
	"net/http"
	"strings"
)

type ParameterType string
type DocResponses map[string]DocItemStruct
type DocParameters map[string]string

const (
	QUERY ParameterType = "query"
	PATH  ParameterType = "path"
)

type MediaType string

const (
	JSON MediaType = "application/json"
	XML  MediaType = "application/xml"
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
	Headers   DocParameters
	Cookies   DocParameters
	Responses DocResponses
}

type DocPayload struct {
	Description string
	Parameters  DocParameters
	Query       DocParameters
	Files       DocParameters
	Cookies     DocParameters
	Request     DocItemStruct
	Responses   DocResponses
	Tags        *[]string
}

type DocRoute struct {
	Description string
	Method      string
	BasePath    string
	Path        string
	Parameters  DocParameters
	Query       DocParameters
	Files       DocParameters
	Cookies     DocParameters
	Request     DocItemStruct
	Responses   DocResponses
	Tags        *[]string
}

type DocItemStruct struct {
	Item        any
	MediaType   MediaType
	Description string
}

func DocXmlStruct(item any, description ...string) DocItemStruct {
	return docStruct(item, XML, description...)
}

func DocJsonStruct(item any, description ...string) DocItemStruct {
	return docStruct(item, JSON, description...)
}

func docStruct(item any, media MediaType, description ...string) DocItemStruct {
	return DocItemStruct{
		Item:        item,
		MediaType:   media,
		Description: strings.Join(description, ""),
	}
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

func DocTags(tags ...string) *[]string {
	return &tags
}
