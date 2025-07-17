package swagger

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gopkg.in/yaml.v3"
)

const SWAGGER string = "SWAGGER"

type OpenAPI3Viewer struct {
	build      sync.Once
	data       OpenAPI3
	factory    *FactoryStructToSchema
	headers    map[string]map[string]string
	cookies    map[string]map[string]string
	stringData string
}

func InitializeViewer() *OpenAPI3Viewer {
	data, err := loadYAML("swagger.yaml")
	if err != nil {
		log.Error(err)
		data = &OpenAPI3{}
	}

	conf := configuration.Instance()

	data.Servers = []Server{}

	if !conf.OnlyTLS() {
		httpURL := fmt.Sprintf("http://localhost:%d", conf.Port())
		data.Servers = append(data.Servers, Server{
			URL:         httpURL,
			Description: "HTTP server",
		})
	}

	if conf.EnableTLS() {
		httpsURL := fmt.Sprintf("https://localhost:%d", conf.PortTLS())
		data.Servers = append(data.Servers, Server{
			URL:         httpsURL,
			Description: "HTTPS server",
		})
	}

	data.Info.Version = conf.Project.Version

	log.Custom(SWAGGER, "Swagger displayed on /swagger")

	return &OpenAPI3Viewer{
		data:       *data,
		factory:    NewFactoryStructToSchema(),
		headers:    make(map[string]map[string]string),
		cookies:    make(map[string]map[string]string),
		stringData: "",
	}
}

func loadYAML(filename string) (*OpenAPI3, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var doc OpenAPI3
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (v *OpenAPI3Viewer) RegisterGroup(group string, data docs.DocGroup) docs.IDocViewer {
	v.groupHeaders(group, data.Headers)
	v.groupCookies(group, data.Cookies)
	return v
}

func (v *OpenAPI3Viewer) groupHeaders(group string, headers map[string]string) docs.IDocViewer {
	item, ok := v.headers[group]
	if !ok {
		item = make(map[string]string, 0)
	}

	maps.Copy(item, headers)

	v.headers[group] = item

	return v
}

func (v *OpenAPI3Viewer) groupCookies(group string, cookies map[string]string) docs.IDocViewer {
	item, ok := v.cookies[group]
	if !ok {
		item = make(map[string]string, 0)
	}

	maps.Copy(item, cookies)

	v.cookies[group] = item

	return v
}

func (v *OpenAPI3Viewer) Handlers() []docs.DocViewerHandler {
	return []docs.DocViewerHandler{
		{
			Method:      http.MethodGet,
			Route:       "/swagger/",
			Handler:     httpSwagger.WrapHandler,
			Name:        "OAS3",
			Description: "OpenAPI 3.0 view",
		},
		{
			Method:      http.MethodGet,
			Route:       "/swagger/doc.json",
			Handler:     v.doc,
			Name:        "OAS3 JSON",
			Description: "OpenAPI 3.0 definition",
		},
	}
}

func (v *OpenAPI3Viewer) RegisterRoute(route docs.DocRoute) docs.IDocViewer {
	if v.data.Paths == nil {
		v.data.Paths = make(map[string]PathItem)
	}

	path := fmt.Sprintf("%s%s", route.BasePath, route.Path)

	pathItem, ok := v.data.Paths[path]
	if !ok {
		pathItem = PathItem{}
	}

	operation := &Operation{
		Tags:        makeTags(route),
		Parameters:  v.makeParameters(path, route),
		RequestBody: v.makeRequest(route),
		Responses:   v.makeResponses(route),
	}

	switch route.Method {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "DELETE":
		pathItem.Delete = operation
	case "PATCH":
		pathItem.Patch = operation
	case "HEAD":
		pathItem.Head = operation
	case "OPTIONS":
		pathItem.Options = operation
	default:
		log.Warningf("Unsupported HTTP method: %s", route.Method)
	}

	log.Customf(SWAGGER, "Route registered: [%s] %s", route.Method, path)

	v.data.Paths[path] = pathItem
	return v
}

func (v *OpenAPI3Viewer) doc(w http.ResponseWriter, r *http.Request) {
	v.build.Do(func() {
		v.data.Components = *v.factory.Components()
		data, err := json.Marshal(v.data)
		if err != nil {
			log.Error(err)
		}
		v.stringData = string(data)
	})

	_, err := w.Write([]byte(v.stringData))
	if err != nil {
		log.Error(err)
	}
}

func (v *OpenAPI3Viewer) makeParameters(path string, route docs.DocRoute) []Parameter {
	parameters := make([]Parameter, 0)

	for k, h := range v.headers {
		if strings.HasPrefix(path, k) {
			for n, d := range h {
				parameters = append(parameters, v.makeParameter(n, d, "header"))
			}
		}
	}

	for k, h := range v.cookies {
		if strings.HasPrefix(path, k) {
			for n, d := range h {
				parameters = append(parameters, v.makeParameter(n, d, "cookie"))
			}
		}
	}

	if route.Parameters != nil {
		for n, d := range route.Parameters {
			parameters = append(parameters, v.makeParameter(n, d, "path"))
		}
	}

	if route.Query != nil {
		for n, d := range route.Query {
			parameters = append(parameters, v.makeParameter(n, d, "query"))
		}
	}

	return parameters
}

func (v *OpenAPI3Viewer) makeParameter(name, description, category string) Parameter {
	return Parameter{
		Name:        name,
		In:          category,
		Description: description,
		Required:    true,
	}
}

func (v *OpenAPI3Viewer) makeRequest(route docs.DocRoute) *RequestBody {
	content := make(map[string]MediaType)

	if contentType, media := v.makeMainRequest(route); media != nil {
		content[contentType] = *media
	}

	if contentType, media := v.makeFileRequest(route); media != nil {
		content[contentType] = *media
	}

	return &RequestBody{
		Content: content,
	}
}

func (v *OpenAPI3Viewer) makeMainRequest(route docs.DocRoute) (string, *MediaType) {
	if route.Request == nil {
		return "", nil
	}

	_, main, err := v.factory.MakeSchema(route.Request)
	if err != nil {
		log.Error(err)
		return "", nil
	}

	return "application/json", &MediaType{
		Schema: main,
	}
}

func (v *OpenAPI3Viewer) makeFileRequest(route docs.DocRoute) (string, *MediaType) {
	if len(route.Files) == 0 {
		return "", nil
	}

	properties := make(map[string]*Schema)
	for k, d := range route.Files {
		properties[k] = &Schema{
			Type:        "string",
			Format:      "binary",
			Description: d,
		}
	}

	multipart := &Schema{
		Type:       "object",
		Properties: properties,
	}

	return "multipart/form-data", &MediaType{
		Schema: multipart,
	}
}

func (v *OpenAPI3Viewer) makeResponses(route docs.DocRoute) map[string]Response {
	if len(route.Responses) == 0 {
		return nil
	}

	responses := make(map[string]Response)
	for status, response := range route.Responses {
		_, main, err := v.factory.MakeSchema(response)
		if err != nil {
			log.Error(err)
			return nil
		}
		responses[status] = Response{
			Content: map[string]MediaType{
				"application/json": {
					Schema: main,
				},
			},
		}
	}

	return responses
}

func makeTags(route docs.DocRoute) []string {
	tags := make([]string, 0)
	if route.BasePath != "" {
		tags = append(tags, route.BasePath)
	}

	fragments := strings.Split(route.Path, "/")
	if len(fragments) > 0 && fragments[0] != "" && !strings.HasPrefix(fragments[0], "{") {
		tags = append(tags, fragments[0])
	}

	return tags
}
