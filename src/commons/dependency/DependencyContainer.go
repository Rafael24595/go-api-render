package dependency

import (
	"github.com/Rafael24595/go-api-core/src/domain"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	repository "github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/historic"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/response"
	"github.com/Rafael24595/go-collections/collection"
)

const (
	PRESIST_PREFIX = "sve"
	HISTORY_PREFIX = "hst"
)

var instance *DependencyContainer

type DependencyContainer struct {
	RepositoryActions  *repository.RequestManager
	RepositoryHistoric repository.IRepositoryHistoric
}

func Initialize() *DependencyContainer {
	if instance != nil {
		panic("//TODO: Already instanced")
	}

	_, err := core_infrastructure.WarmUp()
	if err != nil {
		println(err.Error())
	}

	repositoryActions := loadRepositoryActions()
	repositoryHistoric := loadRepositoryHisotric()

	container := &DependencyContainer{
		RepositoryActions:  repositoryActions,
		RepositoryHistoric: repositoryHistoric,
	}

	instance = container

	return instance
}

func loadRepositoryActions() *repository.RequestManager {
	fileRequest := repository.NewManagerCsvtFile(domain.NewRequestDefault, repository.CSVT_FILE_PATH_REQUEST)
	implRequest := collection.DictionarySyncEmpty[string, domain.Request]()
	repositoryRequest, err := request.InitializeRepositoryMemory(implRequest, fileRequest)
	if err != nil {
		panic(err)
	}

	fileResponse := repository.NewManagerCsvtFile(domain.NewResponseDefault, repository.CSVT_FILE_PATH_RESPONSE)
	implResponse := collection.DictionarySyncEmpty[string, domain.Response]()
	repositoryResponse, err := response.InitializeRepositoryMemory(implResponse, fileResponse)
	if err != nil {
		panic(err)
	}

	return repository.NewRequestManager(repositoryRequest, repositoryResponse).
		SetInsertPolicy(fixHistoricSize)
}

func loadRepositoryHisotric() repository.IRepositoryHistoric {
	fileResponse := repository.NewManagerCsvtFile(domain.NewHistoricDefault, repository.CSVT_FILE_PATH_HISTORIC)
	implResponse := collection.DictionarySyncEmpty[string, domain.Historic]()
	repositoryResponse, err := historic.InitializeRepositoryMemory(implResponse, fileResponse)
	if err != nil {
		panic(err)
	}

	return repositoryResponse
}

func fixHistoricSize(request *domain.Request, repositoryRequest repository.IRepositoryRequest, repositoryResponse repository.IRepositoryResponse) error {
	//TODO: Implement
	return nil
}
