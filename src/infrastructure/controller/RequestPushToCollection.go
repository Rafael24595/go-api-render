package controller

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
)

type RequestPushToCollection struct {
	Request        dto.DtoRequest `json:"request"`
	RequestName    string         `json:"request_name"`
	CollectionId   string         `json:"collection_id"`
	CollectionName string         `json:"collection_name"`
}
