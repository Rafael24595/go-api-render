package controller

import "github.com/Rafael24595/go-api-core/src/domain"

type responseActionRequests struct {
	Requests []domain.Request `json:"requests"`
}

type responseAction struct {
	Request  domain.Request  `json:"request"`
	Response domain.Response `json:"response"`
}
