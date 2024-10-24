package dependency

import (
	"github.com/Rafael24595/go-api-core/src/domain"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	core_repository "github.com/Rafael24595/go-api-core/src/infrastructure/repository"
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

	_, err := core_infrastructure.WarmUp()
	if err != nil {
		println(err.Error())
	}

	repositoryHisotric := loadHistoricMemoryDependencies()
	repositoryPersisted := loadPersistedMemoryDependencies()

	container := &DependencyContainer{
		RepositoryHisotric:  repositoryHisotric,
		RepositoryPersisted: repositoryPersisted,
	}

	instance = container

	return instance
}

func loadHistoricMemoryDependencies() *repository.RequestManager {
	fileRequest := core_repository.NewManagerCsvtFile(domain.NewRequestDefault, request.CSVT_HISTORIC_FILE_PATH)
	requestQueryHistoric, err := request.InitializeMemoryQuery(fileRequest)
	if err != nil {
		panic(err)
	}
	requestCommandHistoric := request.NewMemoryCommand(requestQueryHistoric)

	fileResponse := core_repository.NewManagerCsvtFile(domain.NewResponseDefault, response.CSVT_HISTORIC_FILE_PATH)
	responseQueryHistoric, err := response.InitializeMemoryQuery(fileResponse)
	if err != nil {
		panic(err)
	}
	responseCommandHistoric := response.NewMemoryCommand(responseQueryHistoric)

	return repository.NewRequestManagerLimited(10, requestQueryHistoric, requestCommandHistoric, responseQueryHistoric, responseCommandHistoric)
}

func loadPersistedMemoryDependencies() *repository.RequestManager {
	fileRequest := core_repository.NewManagerCsvtFile(domain.NewRequestDefault, request.CSVT_PERSISTED_FILE_PATH)
	requestQueryPersisted, err := request.InitializeMemoryQuery(fileRequest)
	if err != nil {
		panic(err)
	}
	requestCommandPersisted := request.NewMemoryCommand(requestQueryPersisted)

	fileResponse := core_repository.NewManagerCsvtFile(domain.NewResponseDefault, response.CSVT_PERSISTED_FILE_PATH)
	responseQueryPersisted, err := response.InitializeMemoryQuery(fileResponse)
	if err != nil {
		panic(err)
	}
	responseCommandPersisted := response.NewMemoryCommand(responseQueryPersisted)

	return repository.NewRequestManager(requestQueryPersisted, requestCommandPersisted, responseQueryPersisted, responseCommandPersisted)
}
