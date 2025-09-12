package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/commons/session"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
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
	router  *router.Router
}

func NewController(
	route *router.Router,
	managerRequest *repository.ManagerRequest,
	managerContext *repository.ManagerContext,
	managerCollection *repository.ManagerCollection,
	managerHisotric *repository.ManagerHistoric,
	managerGroup *repository.ManagerGroup) Controller {
	instance := Controller{
		router:  route,
	}

	if configuration.Instance().Front.Enabled {
		NewControllerFront(route)
	}

	route.
		BasePath(BASE_PATH).
		GroupContextualizerDocument(instance.authSoft, docAuthSoft,
			"user",
			"user/verify",
		).
		GroupContextualizerDocument(instance.authHard, docAuthHard,
			"system/log",
			"action",
			"import",
			"sort",
			"context",
			"historic",
			"request",
			"collection",
			"format",
		).
		Cors(router.PermissiveCors())

	if configuration.Instance().Dev() {
		NewControllerDev(route)
	}

	if configuration.Instance().EnableSecrets() {
		NewControllerSecret(route)
	}

	NewControllerSystem(route)
	NewControllerLogin(route)
	NewControllerActions(route)
	NewControllerRequest(route, managerRequest, managerCollection)
	NewControllerHistoric(route, managerRequest, managerHisotric)
	NewControllerContext(route, managerContext)
	NewControllerCollection(route, managerCollection, managerGroup)
	NewControllerFormat(route, managerRequest, managerContext)

	return instance
}

var docAuthSoft = docs.DocGroup{
	Cookies: docs.DocParameters{
		AUTH_COOKIE: AUTH_COOKIE_DESCRIPTION,
	},
	Responses: docs.DocResponses{
		"401": docs.DocText(AUTH_401),
		"404": docs.DocText(AUTH_404),
	},
}

func (c *Controller) authSoft(w http.ResponseWriter, r *http.Request, context router.Context) result.Result {
	user := domain.ANONYMOUS_OWNER

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

	sessions := repository.InstanceManagerSession()
	_, exists := sessions.Find(user)
	if !exists {
		err = errors.New("user not exists")
		return result.Err(http.StatusNotFound, err)
	}

	context.Put(USER, user)

	return result.Ok(context)
}

var docAuthHard = docs.DocGroup{
	Cookies: docs.DocParameters{
		AUTH_COOKIE: AUTH_COOKIE_DESCRIPTION,
	},
	Responses: docs.DocResponses{
		"401": docs.DocText(AUTH_401),
		"404": docs.DocText(AUTH_404),
		"406": docs.DocText(AUTH_406),
	},
}

func (c *Controller) authHard(w http.ResponseWriter, r *http.Request, context router.Context) result.Result {
	res := c.authSoft(w, r, context)
	if res.Err() {
		return res
	}

	userInterface, ok := context.Get(USER)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	username, ok := (*userInterface).(string)
	if !ok {
		return result.Reject(http.StatusNotFound)
	}

	sessions := repository.InstanceManagerSession()

	session, exists := sessions.Find(username)
	if !exists {
		return result.Reject(http.StatusNotFound)
	}

	if session.IsNotVerified() {
		return result.Err(http.StatusNotAcceptable, errors.New("password update required"))
	}

	sessions.Visited(session)

	return result.Ok(context)
}

func findUser(ctx router.Context) string {
	userInterface, exists := ctx.Get(USER)
	if !exists {
		return domain.ANONYMOUS_OWNER
	}

	user, ok := (*userInterface).(string)
	if !ok {
		log.Panics("The user type must be a string")
	}

	return user
}

func jsonDeserialize[T any](r *http.Request) (*T, error) {
	var item T
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func findSession(user string) (*session.Session, *result.Result) {
	sessions := repository.InstanceManagerSession()
	session, ok := sessions.Find(user)
	if !ok {
		result := result.Reject(http.StatusUnauthorized)
		return nil, &result
	}
	return session, nil
}

func findHistoricCollection(user string) (*domain.Collection, *result.Result) {
	sessions := repository.InstanceManagerSession()
	collection, err := sessions.FindUserHistoric(user)
	if err != nil {
		result := result.Err(http.StatusInternalServerError, err)
		return nil, &result
	}
	return collection, nil
}

func findUserCollection(user string) (*domain.Collection, *result.Result) {
	sessions := repository.InstanceManagerSession()
	collection, err := sessions.FindUserCollection(user)
	if err != nil {
		result := result.Err(http.StatusInternalServerError, err)
		return nil, &result
	}

	return collection, nil
}

func findUserGroup(user string) (*domain.Group, *result.Result) {
	sessions := repository.InstanceManagerSession()
	group, err := sessions.FindUserGroup(user)
	if err != nil {
		result := result.Err(http.StatusInternalServerError, err)
		return nil, &result
	}

	return group, nil
}
