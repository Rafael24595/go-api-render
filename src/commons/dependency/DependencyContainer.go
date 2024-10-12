package dependency

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/response"
	"github.com/Rafael24595/go-api-render/src/infrastructure/repository"
)

var instance *DependencyContainer

type DependencyContainer struct {
	RepositoryHisotric  *repository.RequestManager
	RepositoryPersisted *repository.RequestManager
}

func Initialize() *DependencyContainer {
	if instance != nil {
		panic("//TODO: Yet instanced")
	}
	
	requestQueryHistoric, err := request.InitializeMemoryQueryPath(request.HISTORIC_FILE_PATH)
	if err != nil {
		panic(err)
	}
	requestCommandHistoric := request.NewMemoryCommand(requestQueryHistoric)
	responseQueryHistoric, err := response.InitializeMemoryQueryPath(response.HISTORIC_FILE_PATH)
	if err != nil {
		panic(err)
	}
	responseCommandHistoric := response.NewMemoryCommand(responseQueryHistoric)

	requestQueryPersisted, err := request.InitializeMemoryQueryPath(request.DEFAULT_FILE_PATH)
	if err != nil {
		panic(err)
	}
	requestCommandPersisted := request.NewMemoryCommand(requestQueryPersisted)
	responseQueryPersisted, err := response.InitializeMemoryQueryPath(response.DEFAULT_FILE_PATH)
	if err != nil {
		panic(err)
	}
	responseCommandPersisted := response.NewMemoryCommand(responseQueryPersisted)


	repositoryHisotric := repository.NewRequestManagerLimited(10, requestQueryHistoric, requestCommandHistoric, responseQueryHistoric, responseCommandHistoric)
	repositoryPersisted := repository.NewRequestManager(requestQueryPersisted, requestCommandPersisted, responseQueryPersisted, responseCommandPersisted)

	container := &DependencyContainer{
		RepositoryHisotric: repositoryHisotric,
		RepositoryPersisted: repositoryPersisted,
	}

	instance = container
	
	return instance

}