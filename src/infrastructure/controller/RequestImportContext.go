package controller

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
)

type RequestImportContext struct {
	Target dto.DtoContext `json:"target"`
	Source dto.DtoContext `json:"source"`
}
