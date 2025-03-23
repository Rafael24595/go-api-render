package controller

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
)

type responseAction struct {
	Request  dto.DtoRequest  `json:"request"`
	Response dto.DtoResponse `json:"response"`
}
