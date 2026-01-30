package controller

import (
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/Rafael24595/go-api-core/src/application/manager"
	"github.com/Rafael24595/go-api-core/src/application/session"
	"github.com/Rafael24595/go-api-core/src/domain/action"
	"github.com/Rafael24595/go-api-core/src/domain/collection"
	"github.com/Rafael24595/go-api-core/src/domain/group"
	domain_session "github.com/Rafael24595/go-api-core/src/domain/session"
	"github.com/Rafael24595/go-api-core/src/domain/token"
	render_manager "github.com/Rafael24595/go-api-render/src/application/manager"
	auth "github.com/Rafael24595/go-api-render/src/commons/auth/Jwt.go"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-web/router"
	"github.com/Rafael24595/go-web/router/docs"
	"github.com/Rafael24595/go-web/router/result"
)

const USER = "user"

const (
	AUTH_401 = "Invalid or expired authentication token"
	AUTH_404 = "User does not exist or session is invalid"
	AUTH_406 = "Password update required"
)

const BASE_PATH = "/api/v1/"

type Controller struct {
	router       *router.Router
	managerToken *manager.ManagerToken
}

func NewController(
	route *router.Router,
	managerRequest *manager.ManagerRequest,
	managerContext *manager.ManagerContext,
	managerCollection *manager.ManagerCollection,
	managerHisotric *manager.ManagerHistoric,
	managerGroup *manager.ManagerGroup,
	managerEndPoint *manager.ManagerEndPoint,
	managerMetrics *manager.ManagerMetrics,
	managerToken *manager.ManagerToken,
	managerSessionData *session.ManagerSessionData,
	managerWeb *render_manager.ManagerWeb,
) Controller {
	conf := configuration.Instance()

	instance := Controller{
		router:       route,
		managerToken: managerToken,
	}

	if conf.Front.Enabled {
		NewControllerFront(route)
	}

	laxAuth := instance.laxAuth
	if conf.EnableUserToken() {
		laxAuth = router.FallbackHandlers(instance.authToken, laxAuth)
	}

	strictAuth := router.ValidateHandlers(laxAuth, instance.authStrict)

	route.
		BasePath(BASE_PATH).
		GroupContextualizerDocument(laxAuth, docAuthLax,
			"user",
			"user/verify",
		).
		GroupContextualizerDocument(strictAuth, docAuthStrict,
			"system/log",
			"system/cmd",
			"action",
			"import",
			"sort",
			"context",
			"historic",
			"request",
			"collection",
			"curl",
			"token",
			"mock/endpoint",
			"mock/metrics/",
			"bridge/mock/endpoint",
		).
		Cors(router.PermissiveCors())

	if configuration.Instance().Dev() {
		NewControllerDev(route)
	}

	if configuration.Instance().EnableSecrets() {
		NewControllerSecret(route)
	}

	NewControllerSystem(route)
	NewControllerLogin(route, managerWeb)
	NewControllerActions(route)
	NewControllerRequest(route, managerRequest, managerCollection, managerSessionData)
	NewControllerHistoric(route, managerRequest, managerHisotric, managerSessionData)
	NewControllerContext(route, managerContext, managerSessionData)
	NewControllerCollection(route, managerCollection, managerGroup, managerSessionData)
	NewControllerCurl(route, managerRequest, managerCollection, managerGroup,
		managerContext, managerEndPoint, managerSessionData)
	NewControllerMock(route, managerToken, managerEndPoint, managerMetrics)
	NewControllerToken(route, managerToken)

	return instance
}

var docAuthLax = docs.DocGroup{
	Cookies: docs.DocParameters{
		AUTH_COOKIE: AUTH_COOKIE_DESCRIPTION,
		AUTH_TOKEN:  AUTH_TOKEN_DESCRIPTION,
	},
	Responses: docs.DocResponses{
		"401": docs.DocText(AUTH_401),
		"404": docs.DocText(AUTH_404),
	},
}

func (c *Controller) authToken(w http.ResponseWriter, r *http.Request, context *router.Context) result.Result {
	cookie, err := r.Cookie(AUTH_TOKEN)
	if err != nil {
		return result.Err(http.StatusUnauthorized, err)
	}

	if cookie == nil {
		return result.Reject(http.StatusUnauthorized)
	}

	tkn, ok := c.managerToken.FindGlobal(cookie.Value)
	if !ok {
		return result.Err(http.StatusForbidden, err)
	}

	if tkn.IsExipred() {
		return result.TextErr(http.StatusUnauthorized, "the provided token has expired")
	}

	if !slices.Contains(tkn.Scopes, token.ScopeAPIToken) {
		return result.TextErr(http.StatusUnauthorized, "the provided token does not have the necessary permissions")
	}

	context.Put(USER, tkn.Owner)

	return result.Ok(context)
}

func (c *Controller) laxAuth(w http.ResponseWriter, r *http.Request, context *router.Context) result.Result {
	user := action.ANONYMOUS_OWNER

	token, err := r.Cookie(AUTH_COOKIE)
	if err != nil {
		context.Put(USER, user)
		return result.Ok(context)
	}

	claims, err := auth.ValidateJWT(token.Value)
	if err != nil {
		closeSession(w)
		if claims != nil && 0 >= time.Until(claims.ExpiresAt.Time) {
			return result.Err(498, errors.New("token expired"))
		}
		return result.Err(http.StatusUnauthorized, err)
	}

	user = claims.Username

	sessions := session.InstanceManagerSession()
	_, exists := sessions.Find(user)
	if !exists {
		err = errors.New("user not exists")
		return result.Err(http.StatusNotFound, err)
	}

	context.Put(USER, user)

	return result.Ok(context)
}

var docAuthStrict = docs.DocGroup{
	Cookies: docs.DocParameters{
		AUTH_COOKIE: AUTH_COOKIE_DESCRIPTION,
	},
	Responses: docs.DocResponses{
		"401": docs.DocText(AUTH_401),
		"404": docs.DocText(AUTH_404),
		"406": docs.DocText(AUTH_406),
	},
}

func (c *Controller) authStrict(w http.ResponseWriter, r *http.Request, context *router.Context) result.Result {
	userAny, ok := context.Get(USER)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	username, ok := userAny.String()
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	sessions := session.InstanceManagerSession()

	session, exists := sessions.Find(username)
	if !exists {
		return result.Reject(http.StatusNotFound)
	}

	if session.IsNotVerified() {
		return result.Err(http.StatusNotAcceptable, errors.New("password update required"))
	}

	return result.Ok(context)
}

func findUser(ctx *router.Context) string {
	return ctx.Getz(USER).
		Stringd(action.ANONYMOUS_OWNER)
}

func findSession(user string) (*domain_session.Session, *result.Result) {
	sessions := session.InstanceManagerSession()
	session, ok := sessions.Find(user)
	if !ok {
		result := result.Reject(http.StatusUnauthorized)
		return nil, &result
	}
	return session, nil
}

func findTransientCollection(user string, client *session.ManagerSessionData) (*collection.Collection, *result.Result) {
	collection, err := client.FindTransient(user)
	if err != nil {
		result := result.Err(http.StatusInternalServerError, err)
		return nil, &result
	}
	return collection, nil
}

func findPersistentCollection(user string, client *session.ManagerSessionData) (*collection.Collection, *result.Result) {
	collection, err := client.FindPersistent(user)
	if err != nil {
		result := result.Err(http.StatusInternalServerError, err)
		return nil, &result
	}

	return collection, nil
}

func findUserCollections(user string, client *session.ManagerSessionData) (*group.Group, *result.Result) {
	group, err := client.FindCollections(user)
	if err != nil {
		result := result.Err(http.StatusInternalServerError, err)
		return nil, &result
	}

	return group, nil
}
