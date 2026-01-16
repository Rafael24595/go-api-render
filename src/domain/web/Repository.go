package web

type Repository interface {
	Find(id string) (*WebData, bool)
	FindByOwner(owner string) (*WebData, bool)
	Resolve(owner string, webData *WebData) *WebData
	Delete(token *WebData) *WebData
}
