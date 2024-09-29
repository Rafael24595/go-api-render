package templates

import (
	"fmt"
	"html/template"
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
	for _, t := range builder.templates {
		templates = template.Must(templates.ParseGlob(fmt.Sprintf("%s/*.html", t)))
	}
	return templates
}