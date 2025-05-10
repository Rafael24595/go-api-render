package controller

import (
	"github.com/Rafael24595/go-api-core/src/commons/configuration"
	"github.com/Rafael24595/go-api-core/src/commons/log"
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

type responseSystemMetadata struct {
	SessionId     string `json:"session_id"`
	SessionTime   int64  `json:"session_time"`
	CoreName      string `json:"core_name"`
	CoreVersion   string `json:"core_version"`
	CoreReplace   bool   `json:"core_replace"`
	RenderName    string `json:"render_name"`
	RenderVersion string `json:"render_version"`
}

func makeResponseSystemMetadata(sessionId string, timestamp int64, mod configuration.Mod, project configuration.Project) responseSystemMetadata {
	core, ok := mod.Dependencies["github.com/Rafael24595/go-api-core"]
	if !ok {
		log.Panics("Core dependency is not defined")
	}

	return responseSystemMetadata{
		SessionId:     sessionId,
		SessionTime:   timestamp,
		CoreName:      core.Module,
		CoreVersion:   core.Version,
		CoreReplace:   core.Replace != "",
		RenderName:    mod.Module,
		RenderVersion: project.Version,
	}
}
