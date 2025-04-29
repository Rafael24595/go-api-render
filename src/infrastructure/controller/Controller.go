package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	auth "github.com/Rafael24595/go-api-render/src/commons/auth/Jwt.go"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/result"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
)

const USER = "user"

type Controller struct {
	router  *router.Router
	manager templates.TemplateManager
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
		manager: templates.NewBuilder().Make(),
	}

	cors := router.EmptyCors().
		AllowedOrigins("*").
		AllowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS").
		AllowedHeaders("Content-Type", "Authorization").
		AllowCredentials()

	route.
		GroupContextualizer(instance.authSoft,
			"/api/v1/user",
			"/api/v1/user/verify",
		).
		GroupContextualizer(instance.authHard,
			"/api/v1/action",
			"/api/v1/import/",
			"/api/v1/context",
			"/api/v1/historic",
			"/api/v1/request",
			"/api/v1/collection",
		).
		Cors(cors)

	if configuration.Instance().Front() {
		NewControllerFront(route)
	}

	NewControllerLogin(route)
	NewControllerActions(route)
	NewControllerStorage(route, managerRequest, managerCollection)
	NewControllerHistoric(route, managerRequest, managerHisotric)
	NewControllerContext(route, managerContext)
	NewControllerCollection(route, managerCollection, managerGroup)

	return instance
}

func (c *Controller) authSoft(w http.ResponseWriter, r *http.Request, context router.Context) result.Result {
	user := domain.ANONYMOUS_OWNER

	token, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		context.Put(USER, user)
		return result.Ok(context)
	}

	claims, err := auth.ValidateJWT(token.Value)
	if err != nil {
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

	timeLeft := time.Until(claims.ExpiresAt.Time)
	if timeLeft < 10*time.Minute {
		defineSession(w, user)
	}

	return result.Ok(context)
}

func (c *Controller) authHard(w http.ResponseWriter, r *http.Request, context router.Context) result.Result {
	res := c.authSoft(w, r, context)
	if _, ok := res.Err(); ok {
		return res
	}

	userInterface, ok := context.Get(USER)
	if !ok {
		return result.Err(http.StatusNotFound, nil)
	}

	username, ok := (*userInterface).(string)
	if !ok {
		return result.Err(http.StatusNotFound, nil)
	}

	sessions := repository.InstanceManagerSession()

	session, exists := sessions.Find(username)
	if !exists {
		return result.Err(http.StatusNotFound, nil)
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
		panic("//TODO: Manage error.")
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
