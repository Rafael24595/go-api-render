package controller

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
)

type responseAction struct {
	Request  dto.DtoRequest  `json:"request"`
	Response dto.DtoResponse `json:"response"`
}

type responseUserData struct {
	Username    string `json:"username"`
	Timestamp   int64  `json:"timestamp"`
	History     string `json:"history"`
	Collection  string `json:"collection"`
	Context     string `json:"context"`
	IsProtected bool   `json:"is_protected"`
	IsAdmin     bool   `json:"is_admin"`
	FirstTime   bool   `json:"first_time"`
}

type responseClientIdentity struct {
	Ip     string `json:"ip"`
	IsHost bool   `json:"is_host"`
}
