package repository

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/response"
)

type RequestManager struct {
	mu        sync.Mutex
	prefix    string
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
		prefix:    "",
		limit:     limit,
		qRequest:  qRequest,
		cRequest:  cRequest,
		qResponse: qResponse,
		cResponse: cResponse,
	}
}

func (m *RequestManager) HasPrefix(id string) bool {
	return strings.HasPrefix(id, fmt.Sprintf("%s-", m.prefix))
}

func (m *RequestManager) SetPrefix(prefix string) *RequestManager {
	m.prefix = prefix
	m.qRequest.SetPrefix(prefix)
	m.qResponse.SetPrefix(prefix)
	return m
}

func (m *RequestManager) Exists(key string) (bool, bool) {
	_, okReq := m.qRequest.Find(key)
	_, okRes := m.qResponse.Find(key)
	return okReq, okRes
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

func (m *RequestManager) Insert(request domain.Request, response *domain.Response) (*domain.Request, *domain.Response) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ids, requestResult := m.insertRequest(request)
	resultResponse := m.insertResponse(ids, requestResult, response)

	return requestResult, resultResponse
}

func (m *RequestManager) insertRequest(request domain.Request) ([]string, *domain.Request) {
	result := m.cRequest.Insert(request)

	ids := m.cRequest.DeleteOptions(repository.FilterOptions[domain.Request]{
		Sort: func(i, j domain.Request) bool {
			return j.Timestamp > i.Timestamp
		},
		From: m.limit,
	})

	return ids, result
}

func (m *RequestManager) insertResponse(ids []string, request *domain.Request, response *domain.Response) *domain.Response {
	if response == nil {
		return nil
	}

	response.Id = request.Id

	result := m.cResponse.Insert(*response)

	idCollection := collection.FromList(ids)

	m.cResponse.DeleteOptions(repository.FilterOptions[domain.Response]{
		Predicate: func(r domain.Response) bool {
			_, ok := idCollection.Find(func(s string) bool {
				return r.Id == s
			})
			return !ok
		},
	})

	return result
}

func (m *RequestManager) Delete(request domain.Request) (*domain.Request, *domain.Response) {
	requestResult := m.cRequest.Delete(request)
	response, ok := m.qResponse.Find(request.Id)
	if ok {
		m.cResponse.Delete(*response)
	}
	return requestResult, response
}
