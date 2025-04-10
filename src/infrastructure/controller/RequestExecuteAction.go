package controller

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
)

type RequestExecuteAction struct {
	Request dto.DtoRequest `json:"request"`
	Context dto.DtoContext `json:"context"`
}
