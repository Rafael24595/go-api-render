package repository

import (
	"sync"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/response"
)

type RequestManager struct {
	mu        sync.Mutex
	limit     int
	qRequest  request.IRepositoryQuery
	cRequest  request.IRepositoryCommand
	qResponse response.IRepositoryQuery
	cResponse response.IRepositoryCommand
}

func NewRequestManager(qRequest request.IRepositoryQuery, cRequest request.IRepositoryCommand, qResponse response.IRepositoryQuery, cResponse response.IRepositoryCommand) *RequestManager {
	return NewRequestManagerLimited(0, qRequest, cRequest, qResponse, cResponse)
}

func NewRequestManagerLimited(limit int, qRequest request.IRepositoryQuery, cRequest request.IRepositoryCommand, qResponse response.IRepositoryQuery, cResponse response.IRepositoryCommand) *RequestManager {
	return &RequestManager{
		limit:     limit,
		qRequest:  qRequest,
		cRequest:  cRequest,
		qResponse: qResponse,
		cResponse: cResponse,
	}
}

func (m *RequestManager) FindAll() []domain.Request {
	return m.qRequest.FindAll()
}
	
func (m *RequestManager) Find(key string) (*domain.Request, *domain.Response, bool) {
	request, ok := m.qRequest.Find(key)
	if !ok {
		return nil, nil, ok
	}
	response, _ := m.qResponse.Find(key)
	return request, response, ok
}

func (m *RequestManager) FindOptions(options repository.FilterOptions[domain.Request]) []domain.Request {
	return m.qRequest.FindOptions(options)
}

func (m *RequestManager) Insert(request domain.Request, response domain.Response) (*domain.Request, *domain.Response) {
	m.mu.Lock()
	defer m.mu.Unlock()

	requestResult := m.cRequest.Insert(request)

	response.Id = requestResult.Id
	resultResponse := m.cResponse.Insert(response)

	ids := m.cRequest.DeleteOptions(repository.FilterOptions[domain.Request]{
		Sort: func(i, j domain.Request) bool {
			return j.Timestamp > i.Timestamp
		},
		From: m.limit,
	})

	idCollection := collection.FromList(ids)

	m.cResponse.DeleteOptions(repository.FilterOptions[domain.Response]{
		Predicate: func(r domain.Response) bool {
			_, ok := idCollection.Find(func(s string) bool {
				return r.Id == s
			})
			return !ok
		},
	})

	return requestResult, resultResponse
}
