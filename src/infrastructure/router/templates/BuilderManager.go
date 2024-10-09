package templates

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
)

type builderTemplate interface {
	makeTemplate() *template.Template 
}

type BuilderManager struct {
	templates []string
	functions template.FuncMap
}

func NewBuilder() *BuilderManager {
	return &BuilderManager{
		templates: []string{},
		functions: map[string]any{},
	}
}

func (builder *BuilderManager) AddPath(path string) *BuilderManager {
	builder.templates = append(builder.templates, path)
	return builder
}

func (builder *BuilderManager) AddFunctions(funcs map[string]any) *BuilderManager {
	for k, v := range funcs {
		builder.functions[k] = v
	}
	return builder
}

func (builder *BuilderManager) AddFunction(key string, value any) *BuilderManager {
	builder.functions[key] = value
	return builder
}

func (builder *BuilderManager) Make() TemplateManager {
	return TemplateManager{
		templates: builder.makeTemplate(),
		builder: builder,
	}
}

func (builder *BuilderManager) makeTemplate() *template.Template {
	templates := template.New("").Funcs(builder.functions)
	templates.ParseFiles(builder.files()...)
	return templates
}

func (builder *BuilderManager) files() []string {
	files := collection.EmptyMap[string, bool]()
	for _, t := range builder.templates {
		err := filepath.WalkDir(t, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			
			if d.IsDir() || !strings.HasSuffix(d.Name(), ".html") {
				return nil
			}
			
			files.Put(path, true)
	
			return nil
		})

		if err != nil {
			panic(err)
		}
	}

	return files.Keys()
}