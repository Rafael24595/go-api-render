package controller

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
)

type RequestPushToCollection struct {
	SourceId    string              `json:"source_id"`
	TargetId    string              `json:"target_id"`
	TargetName  string              `json:"target_name"`
	Request     dto.DtoRequest      `json:"request"`
	RequestName string              `json:"request_name"`
	Movement    repository.Movement `json:"move"`
}

func RequestPushToCollectionToPayload(payload *RequestPushToCollection) repository.PayloadPushToCollection {
	request := dto.ToRequest(&payload.Request)
	return repository.PayloadPushToCollection{
		SourceId:    payload.SourceId,
		TargetId:    payload.TargetId,
		TargetName:  payload.TargetName,
		Request:     *request,
		RequestName: payload.RequestName,
		Movement:    payload.Movement,
	}
}
