package templates

import (
	"html/template"
	"net/http"
	"path"

	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
)

type TemplateManager struct {
	templates *template.Template
	builder   builderTemplate
}

func (manager *TemplateManager) load() *template.Template {
	if configuration.Instance().Debug() {
		return manager.builder.makeTemplate()
	}
	return manager.templates
}

func (manager *TemplateManager) Render(w http.ResponseWriter, tmpl string, context router.Context) error {
	_, file := path.Split(tmpl)
	t := manager.load().Lookup(file)

	err := t.Execute(w, context.Collect())
	if err != nil {
		//TODO: Replace with log.
		println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return nil
}
