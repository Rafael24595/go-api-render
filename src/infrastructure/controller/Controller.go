package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	auth "github.com/Rafael24595/go-api-render/src/commons/auth/Jwt.go"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
)

const USER = "user"

type Controller struct {
	router  *router.Router
	manager templates.TemplateManager
}

func NewController(
	route *router.Router,
	managerActions *repository.ManagerRequest,
	managerContext *repository.ManagerContext,
	managerCollection *repository.ManagerCollection,
	repositoryContext repository.IRepositoryContext,
	repositoryHisotric repository.IRepositoryHistoric) Controller {
	instance := Controller{
		router:  route,
		manager: templates.NewBuilder().Make(),
	}

	cors := router.EmptyCors().
		AllowedOrigins("*").
		AllowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS").
		AllowedHeaders("Content-Type", "Authorization")

	route.
		GroupContextualizer(instance.auth,
			"/api/v1/user",
			"/api/v1/action",
			"/api/v1/import/",
			"/api/v1/context",
			"/api/v1/historic",
			"/api/v1/request",
			"/api/v1/collection",
		).
		ErrorHandler(instance.error).
		Cors(cors)

	NewControllerLogin(route)
	NewControllerActions(route)
	NewControllerStorage(route, managerActions)
	NewControllerHistoric(route, managerActions, repositoryHisotric)
	NewControllerContext(route, managerContext)
	NewControllerCollection(route, managerCollection, managerActions, repositoryContext)

	return instance
}

func (c *Controller) auth(w http.ResponseWriter, r *http.Request, context router.Context) error {
	user := "anonymous"

	token, err := r.Cookie(COOKIE_NAME)
	if err != nil {
		context.Put(USER, user)
		return nil
	}

	claims, err := auth.ValidateJWT(token.Value)
	if err != nil {
		return err
	}

	user = claims.Username

	sessions := repository.InstanceManagerSession()

	_, exists := sessions.Find(user)
	if !exists {
		return errors.New("user not exists")
	}

	context.Put(USER, user)

	timeLeft := time.Until(claims.ExpiresAt.Time)
	if timeLeft < 10*time.Minute {
		defineSession(w, user)
	}

	return nil
}

func (c *Controller) error(w http.ResponseWriter, r *http.Request, context router.Context, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
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
