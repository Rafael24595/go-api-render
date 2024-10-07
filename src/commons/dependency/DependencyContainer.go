package dependency

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request"
)

var instance *DependencyContainer

type DependencyContainer struct {
	RequestQueryHistoric request.RepositoryQuery
	RequestQueryPersisted request.RepositoryQuery
	RequestCommandManager *request.MemoryCommandManager
}

func Initialize() *DependencyContainer {
	if instance != nil {
		panic("//TODO: Yet instanced")
	}
	
	requestQueryHistoric, err := request.InitializeMemoryQueryPath(request.HISTORIC_FILE_PATH)
	if err != nil {
		panic(err)
	}

	requestQueryPersisted, err := request.InitializeMemoryQueryPath(request.DEFAULT_FILE_PATH)
	if err != nil {
		panic(err)
	}

	requestCommandHistoric := request.NewMemoryCommand(requestQueryHistoric)
	requestCommandPersisted := request.NewMemoryCommand(requestQueryPersisted)

	requestCommandManager := request.NewMemoryCommandManager(requestQueryHistoric, requestCommandHistoric, requestCommandPersisted)

	container := &DependencyContainer{
		RequestQueryHistoric: requestQueryHistoric,
		RequestQueryPersisted: requestQueryPersisted,
		RequestCommandManager: requestCommandManager,
	}

	instance = container
	
	return instance

}