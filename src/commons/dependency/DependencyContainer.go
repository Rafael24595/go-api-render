package dependency

import (
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/domain/context"
	core_infrastructure "github.com/Rafael24595/go-api-core/src/infrastructure"
	"github.com/Rafael24595/go-api-core/src/infrastructure/dto"
	repository "github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	repository_collection "github.com/Rafael24595/go-api-core/src/infrastructure/repository/collection"
	repository_context "github.com/Rafael24595/go-api-core/src/infrastructure/repository/context"
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
	RepositoryContext    repository.IRepositoryContext
	RepositoryActions    *repository.RequestManager
	RepositoryHistoric   repository.IRepositoryHistoric
	RepositoryCollection repository.IRepositoryCollection
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
	repositoryContext := loadRepositoryContext()
	repositoryCollection := loadRepositoryCollection(repositoryContext)

	container := &DependencyContainer{
		RepositoryContext:    repositoryContext,
		RepositoryActions:    repositoryActions,
		RepositoryHistoric:   repositoryHistoric,
		RepositoryCollection: repositoryCollection,
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
	file := repository.NewManagerCsvtFile(domain.NewHistoricDefault, repository.CSVT_FILE_PATH_HISTORIC)
	impl := collection.DictionarySyncEmpty[string, domain.Historic]()
	repository, err := historic.InitializeRepositoryMemory(impl, file)
	if err != nil {
		panic(err)
	}

	return repository
}

func loadRepositoryContext() repository.IRepositoryContext {
	file := repository.NewManagerCsvtFile(dto.NewDtoContextDefault, repository.CSVT_FILE_PATH_CONTEXT)
	impl := collection.DictionarySyncEmpty[string, context.Context]()
	repository, err := repository_context.InitializeRepositoryMemory(impl, file)
	if err != nil {
		panic(err)
	}

	return repository
}

func loadRepositoryCollection(context repository.IRepositoryContext) repository.IRepositoryCollection {
	file := repository.NewManagerCsvtFile(domain.NewCollectionDefault, repository.CSVT_FILE_PATH_COLLECTION)
	impl := collection.DictionarySyncEmpty[string, domain.Collection]()
	repository, err := repository_collection.InitializeRepositoryMemory(impl, file, context)
	if err != nil {
		panic(err)
	}

	return repository
}

func fixHistoricSize(request *domain.Request, repositoryRequest repository.IRepositoryRequest, repositoryResponse repository.IRepositoryResponse) error {
	//TODO: Implement
	return nil
}
