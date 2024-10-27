package templates

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/Rafael24595/go-api-render/src/infrastructure/router/datalist"
)

type defaultTemplates struct {
	listManager   *datalist.DataListManager
	userTemplates *template.Template
	tmpDataList   *template.Template
}

func newDefaultTemplates(listManager *datalist.DataListManager) *defaultTemplates {
	templates := template.New("")

	tmpDataList, err := parseRenderLists(templates)
	if err != nil {
		panic(err.Error())
	}

	return &defaultTemplates{
		listManager:   listManager,
		tmpDataList:   tmpDataList,
	}
}

func parseRenderLists(templates *template.Template) (*template.Template, error) {
	tmpl, err := templates.Parse(`
		{{if . }}
			<datalist id="{{ .Id }}">
				{{range $Key, $Value := .Options }}
					<option value="{{ $Key }}">
				{{end}}
			</datalist>
		{{end}}`)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (m *defaultTemplates) defineUserTemplate(userTemplates *template.Template) *defaultTemplates {
	m.userTemplates = userTemplates
	return m
}

func (m *defaultTemplates) userTemplate(name string, data interface{}) (template.HTML, error) {
	if m.userTemplates == nil {
		println("Not user templates intance defined.")
		return template.HTML(""), nil
	}

	buf := bytes.NewBuffer([]byte{})
	t := m.userTemplates.Lookup(name)

	err := t.ExecuteTemplate(buf, name, data)
	if err != nil {
		return template.HTML(""), err
	}

	ret := template.HTML(buf.String())
	return ret, nil
}

func (m *defaultTemplates) renderLists(names ...string) template.HTML {
	lists := []string{}
	for _, v := range names {
		if dataList, ok := m.listManager.DataList(v); ok {
			list, err := m.renderList(*dataList)
			if err != nil {
				panic(err.Error())
			}
			lists = append(lists, list)
		}
	}

	return template.HTML(strings.Join(lists, "\n"))
}

func (m *defaultTemplates) renderList(dataList datalist.DataList) (string, error) {
	var renderedTemplate bytes.Buffer
	err := m.tmpDataList.Execute(&renderedTemplate, dataList)
	if err != nil {
		return "", err
	}

	return renderedTemplate.String(), nil
}
