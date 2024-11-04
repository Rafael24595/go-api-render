package dependency

import (
	"github.com/Rafael24595/go-api-core/src/domain"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	core_repository "github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/response"
	"github.com/Rafael24595/go-api-render/src/infrastructure/repository"
)

const (
	PRESIST_PREFIX = "sve"
	HISTORY_PREFIX = "hst"
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
	requestQuery, err := request.InitializeMemoryQuery(fileRequest)
	if err != nil {
		panic(err)
	}
	requestCommand := request.NewMemoryCommand(requestQuery)

	fileResponse := core_repository.NewManagerCsvtFile(domain.NewResponseDefault, response.CSVT_HISTORIC_FILE_PATH)
	responseQuery, err := response.InitializeMemoryQuery(fileResponse)
	if err != nil {
		panic(err)
	}
	responseCommand := response.NewMemoryCommand(responseQuery)

	return repository.NewRequestManagerLimited(10, requestQuery, requestCommand, responseQuery, responseCommand).
		SetPrefix(PRESIST_PREFIX)
}

func loadPersistedMemoryDependencies() *repository.RequestManager {
	fileRequest := core_repository.NewManagerCsvtFile(domain.NewRequestDefault, request.CSVT_PERSISTED_FILE_PATH)
	requestQuery, err := request.InitializeMemoryQuery(fileRequest)
	if err != nil {
		panic(err)
	}
	requestCommand := request.NewMemoryCommand(requestQuery)

	fileResponse := core_repository.NewManagerCsvtFile(domain.NewResponseDefault, response.CSVT_PERSISTED_FILE_PATH)
	responseQuery, err := response.InitializeMemoryQuery(fileResponse)
	if err != nil {
		panic(err)
	}
	responseCommand := response.NewMemoryCommand(responseQuery)

	return repository.NewRequestManager(requestQuery, requestCommand, responseQuery, responseCommand).
		SetPrefix(PRESIST_PREFIX)
}
