package templates

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"strings"
)

type TemplateManager interface {
	Render(w http.ResponseWriter, tmpl string, data interface{})
}

type templateManager struct {
	templates *template.Template
	builder   builderTemplate
}

func (manager *templateManager) load() *template.Template {
	debug := strings.ToLower(os.Getenv("DEBUG")) == "true"
	if debug {
		return manager.builder.makeTemplate()
	}
	return manager.templates
}

func (manager *templateManager) Render(w http.ResponseWriter, tmpl string, data interface{}) {
	_, file := path.Split(tmpl)
	t := manager.load().Lookup(file)

	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
