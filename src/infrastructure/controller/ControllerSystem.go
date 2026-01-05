package controller

import (
	"net/http"
	"strconv"

	"github.com/Rafael24595/go-api-core/src/commons/command"
	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/commons/session"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const CMD_SUCCESS_RESPONSE = "Output message"
const CMD_EXCEPTION_RESPONSE = "Error message"

const CMD_DESCRIPTION = "Command sentence"

const CMD_QUERY_POSITION = "step"
const CMD_QUERY_POSITION_DESCRIPTION = "Step value"

type ControllerSystem struct {
	router *router.Router
}

func NewControllerSystem(router *router.Router) ControllerSystem {
	instance := ControllerSystem{
		router: router,
	}

	router.
		RouteDocument(http.MethodGet, instance.log, "system/log", instance.docLog()).
		RouteDocument(http.MethodPost, instance.cmdExec, "system/cmd/exec", instance.docCmdExec()).
		RouteDocument(http.MethodPost, instance.cmdComp, "system/cmd/comp", instance.docCmdComp()).
		RouteDocument(http.MethodGet, instance.metadata, "system/metadata", instance.docMetadata())

	return instance
}

func (c *ControllerSystem) docLog() docs.DocRoute {
	return docs.DocRoute{
		Description: "Returns all server-side application logs. Only accessible by admin users.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[[]log.Record](),
		},
	}
}

func (c *ControllerSystem) log(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	sess, res := findSession(user)
	if res != nil {
		return *res
	}

	if !sess.HasRole(session.ROLE_ADMIN) {
		return result.Reject(http.StatusForbidden)
	}

	return result.JsonOk(log.Records())
}

func (c *ControllerSystem) docCmdExec() docs.DocRoute {
	return docs.DocRoute{
		Description: "Executes a system command; requires administrator privileges.",
		Request:     docs.DocText(CMD_DESCRIPTION),
		Responses: docs.DocResponses{
			"200": docs.DocText(CMD_SUCCESS_RESPONSE),
			"500": docs.DocText(CMD_EXCEPTION_RESPONSE),
		},
	}
}

func (c *ControllerSystem) cmdExec(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	res := c.hasCmdPrivileges(user)
	if res != nil {
		return *res
	}

	cmd, res := router.InputText(r)
	if res != nil {
		return *res
	}

	message, err := command.Exec(user, cmd)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	return result.TextOk(message)
}

func (c *ControllerSystem) docCmdComp() docs.DocRoute {
	return docs.DocRoute{
		Description: "Executes a system command; requires administrator privileges.",
		Query: docs.DocParameters{
			CMD_QUERY_POSITION: CMD_QUERY_POSITION_DESCRIPTION,
		},
		Request: docs.DocText(CMD_DESCRIPTION),
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[cmdCompHelp](),
		},
	}
}

func (c *ControllerSystem) cmdComp(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	user := findUser(ctx)

	res := c.hasCmdPrivileges(user)
	if res != nil {
		return *res
	}

	step := -1
	if raw := r.URL.Query().Get(CMD_QUERY_POSITION); raw != "" {
		if result, err := strconv.ParseInt(raw, 10, 0); err == nil {
			step = int(result)
		}
	}

	cmd, res := router.InputText(r)
	if res != nil {
		return *res
	}

	data, err := command.Comp(user, cmd, step)
	if err != nil {
		return result.Err(http.StatusInternalServerError, err)
	}

	response := cmdCompHelp{
		Message:     data.Message,
		Application: data.Application,
		Position:    data.Position,
		Lenght:      data.Lenght,
	}

	return result.JsonOk(response)
}

func (c *ControllerSystem) docMetadata() docs.DocRoute {
	return docs.DocRoute{
		Description: "Returns runtime system metadata including session ID, timestamp, release version, and frontend status.",
		Responses: docs.DocResponses{
			"200": docs.DocJsonPayload[responseSystemMetadata](),
		},
	}
}

func (c *ControllerSystem) metadata(w http.ResponseWriter, r *http.Request, ctx *router.Context) result.Result {
	conf := configuration.Instance()
	response := makeResponseSystemMetadata(
		conf.SessionId(),
		conf.Timestamp(),
		conf.Release,
		conf.Mod,
		conf.Project,
		conf.Front,
		c.router.ViewerSources(),
		conf.EnableSecrets(),
	)

	return result.JsonOk(response)
}

func (c *ControllerSystem) hasCmdPrivileges(user string) *result.Result {
	sess, res := findSession(user)
	if res != nil {
		return res
	}

	if !sess.HasRole(session.ROLE_ADMIN) {
		res := result.TextErr(http.StatusForbidden, "the user does not have privileges to execute cmd actions")
		return &res
	}

	return nil
}
