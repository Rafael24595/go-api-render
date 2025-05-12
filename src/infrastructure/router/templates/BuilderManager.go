package templates

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/datalist"
	"github.com/Rafael24595/go-collections/collection"
)

type builderTemplate interface {
	makeTemplate() *template.Template
}

type BuilderManager struct {
	templates []string
	functions template.FuncMap
	lists     *datalist.DataListManager
}

func NewBuilder() *BuilderManager {
	return &BuilderManager{
		templates: []string{},
		lists:     datalist.NewDataListManager(),
		functions: map[string]any{},
	}
}

func (manager *BuilderManager) ListManager() *datalist.DataListManager {
	return manager.lists
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
	defaultTemplates := newDefaultTemplates(builder.lists)

	builder.AddFunction("Template", defaultTemplates.userTemplate)
	builder.AddFunction("DataList", defaultTemplates.renderLists)
	builder.AddFunction("Uuid", uuidString)
	builder.AddFunction("String", itemString)
	builder.AddFunction("Not", not)
	builder.AddFunction("Concat", concat)
	builder.AddFunction("Join", join)
	builder.AddFunction("MilisecondsToTime", millisecondsToTime)
	builder.AddFunction("MillisecondsToDate", millisecondsToDate)

	templates := builder.makeTemplate()

	defaultTemplates.defineUserTemplate(templates)

	return TemplateManager{
		builder:   builder,
		templates: templates,
		lists:     builder.lists,
	}
}

func (builder *BuilderManager) makeTemplate() *template.Template {
	templates := template.New("").Funcs(builder.functions)
	files := builder.files()
	if len(files) == 0 {
		return templates
	}
	
	_, err := templates.ParseFiles(files...)
	if err != nil {
		log.Panic(err)
	}

	return templates
}

func (builder *BuilderManager) files() []string {
	files := collection.DictionaryEmpty[string, bool]()
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
			log.Panic(err)
		}
	}

	return files.Keys()
}
