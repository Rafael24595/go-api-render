package templates

import (
	"html/template"
	"net/http"
	"path"

	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

type TemplateManager interface {
	Render(w http.ResponseWriter, tmpl string, context router.Context)
}

type templateManager struct {
	templates *template.Template
	builder   builderTemplate
}

func (manager *templateManager) load() *template.Template {
	if configuration.Instance().Debug() {
		return manager.builder.makeTemplate()
	}
	return manager.templates
}

func (manager *templateManager) Render(w http.ResponseWriter, tmpl string, context router.Context) {
	_, file := path.Split(tmpl)
	t := manager.load().Lookup(file)

	err := t.Execute(w, context.Collect())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
