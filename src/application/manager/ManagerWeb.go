package manager

import (
	"github.com/Rafael24595/go-api-render/src/domain/web"
)

type ManagerWeb struct {
	web web.Repository
}

func NewManagerWeb(web web.Repository) *ManagerWeb {
	return &ManagerWeb{
		web: web,
	}
}

func (m *ManagerWeb) FindByOwner(owner string) (*web.WebData, bool) {
	if result, ok := m.web.FindByOwner(owner); ok && result != nil {
		return result, true
	}
	return m.Resolve(owner, web.EmptyWebData(owner)), true
}

func (m *ManagerWeb) Resolve(owner string, webData *web.WebData) *web.WebData {
	if owner != "" && webData.Owner != owner {
		webData, _ = m.FindByOwner(owner)
		return webData
	}

	if result, ok := m.web.FindByOwner(owner); ok && result != nil {
		result.Data = webData.Data
		webData = result
	}
	return m.web.Resolve(owner, webData)
}

func (m *ManagerWeb) Delete(owner string) (*web.WebData, bool) {
	webData, ok := m.FindByOwner(owner)
	if !ok {
		return nil, false
	}
	return m.web.Delete(webData), true
}
