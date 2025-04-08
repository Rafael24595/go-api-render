package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
	"github.com/Rafael24595/go-collections/collection"
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
		AllowedOrigins([]string{"*"}).
		AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}).
		AllowedHeaders([]string{"Content-Type", "Authorization"})

	route.
		Contextualizer(instance.contextualizer).
		ErrorHandler(instance.error).
		Cors(cors)

	NewControllerActions(route)
	NewControllerStorage(route, managerActions)
	NewControllerHistoric(route, managerActions, repositoryHisotric)
	NewControllerContext(route, managerContext)
	NewControllerCollection(route, managerCollection, managerActions, repositoryContext)

	return instance
}

func (c *Controller) contextualizer(w http.ResponseWriter, r *http.Request) (router.Context, error) {
	return collection.DictionaryFromMap(
		map[string]any{
			//TODO: Extract users from token.
			"user": "anonymous",
		},
	), nil
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
