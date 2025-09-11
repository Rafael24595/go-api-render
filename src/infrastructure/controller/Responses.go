package controller

import (
	core_configuration "github.com/Rafael24595/go-api-core/src/commons/configuration"
	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-web/router/docs"
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

type responseSystemMetadata struct {
	SessionId     string                  `json:"session_id"`
	SessionTime   int64                   `json:"session_time"`
	CoreName      string                  `json:"core_name"`
	CoreVersion   string                  `json:"core_version"`
	CoreReplace   bool                    `json:"core_replace"`
	RenderRelease string                  `json:"render_release"`
	RenderName    string                  `json:"render_name"`
	RenderVersion string                  `json:"render_version"`
	FrontName     string                  `json:"front_name"`
	FrontVersion  string                  `json:"front_version"`
	ViewerSources []docs.DocViewerSources `json:"viewer_sources"`
	EnableSecrets bool                    `json:"enable_secrets"`
}

func makeResponseSystemMetadata(sessionId string, timestamp int64,
	release *core_configuration.Release,
	mod core_configuration.Mod,
	project core_configuration.Project,
	front configuration.FrontPackage,
	viewer []docs.DocViewerSources,
	enableSecrets bool) responseSystemMetadata {
	core, ok := mod.Dependencies["github.com/Rafael24595/go-api-core"]
	if !ok {
		log.Panics("Core dependency is not defined")
	}

	renderRelease := project.Version
	if release != nil && release.TagName != "" {
		renderRelease = release.TagName
	}

	return responseSystemMetadata{
		SessionId:     sessionId,
		SessionTime:   timestamp,
		CoreName:      core.Module,
		CoreVersion:   core.Version,
		CoreReplace:   core.Replace != "",
		RenderRelease: renderRelease,
		RenderName:    mod.Module,
		RenderVersion: project.Version,
		FrontName:     front.Name,
		FrontVersion:  front.Version,
		ViewerSources: viewer,
		EnableSecrets: enableSecrets,
	}
}
