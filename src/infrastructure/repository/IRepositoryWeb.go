package repository

import (
	"github.com/Rafael24595/go-api-render/src/domain/web"
)

type IRepositoryWeb interface {
	Find(id string) (*web.WebData, bool)
	FindByOwner(owner string) (*web.WebData, bool)
	Resolve(owner string, webData *web.WebData) *web.WebData
	Delete(token *web.WebData) *web.WebData
}
