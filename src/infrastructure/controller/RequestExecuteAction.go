package controller

import (
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
)

type RequestExecuteAction struct {
	Request domain.Request `json:"request"`
	Context dto.DtoContext `json:"context"`
}
