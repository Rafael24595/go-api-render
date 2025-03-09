package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router"
	"github.com/Rafael24595/go-api-render/src/infrastructure/router/templates"
	"github.com/Rafael24595/go-collections/collection"
)

type Controller struct {
	router  *router.Router
	manager templates.TemplateManager
}

func NewController(route *router.Router, repositoryManager *repository.RequestManager, repositoryHisotric repository.IRepositoryHistoric) Controller {
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

	NewControllerActions(route, repositoryManager)
	NewControllerStorage(route, repositoryManager, repositoryHisotric)

	return instance
}

func (c *Controller) contextualizer(w http.ResponseWriter, r *http.Request) (router.Context, error) {
	return collection.DictionaryEmpty[string, any](), nil
}

func (c *Controller) error(w http.ResponseWriter, r *http.Request, context router.Context, err error) {
	w.WriteHeader(http.StatusInternalServerError)
}

func jsonDeserialize[T any](r *http.Request) (*T, error) {
	var item T
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
