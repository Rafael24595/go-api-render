package repository

import (
	web_domain "github.com/Rafael24595/go-api-render/src/domain/web"
)

type ManagerWeb struct {
	web IRepositoryWeb
}

func NewManagerWeb(web IRepositoryWeb) *ManagerWeb {
	return &ManagerWeb{
		web: web,
	}
}

func (m *ManagerWeb) FindByOwner(owner string) (*web_domain.WebData, bool) {
	if result, ok := m.web.FindByOwner(owner); ok && result != nil {
		return result, true
	}
	return m.Resolve(owner, web_domain.EmptyWebData(owner)), true
}

func (m *ManagerWeb) Resolve(owner string, webData *web_domain.WebData) *web_domain.WebData {
	return m.web.Resolve(owner, webData)
}

func (m *ManagerWeb) Delete(owner string) (*web_domain.WebData, bool) {
	webData, ok := m.FindByOwner(owner)
	if !ok {
		return nil, false
	}
	return m.web.Delete(webData), true
}
