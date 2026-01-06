package dependency

import (
	"log"
	"sync"

	core_configuration "github.com/Rafael24595/go-api-core/src/commons/configuration"
	core_dependency "github.com/Rafael24595/go-api-core/src/commons/dependency"
	core_system "github.com/Rafael24595/go-api-core/src/commons/system"
	core_repository "github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/Rafael24595/go-api-render/src/commons/system"
	web_domain "github.com/Rafael24595/go-api-render/src/domain/web"
	"github.com/Rafael24595/go-api-render/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/infrastructure/repository/web"
	"github.com/Rafael24595/go-collections/collection"
)

var (
	instance *DependencyContainer
	once     sync.Once
)

type DependencyContainer struct {
	core_dependency.DependencyContainer
	ManagerWeb *repository.ManagerWeb
}

func Initialize(config configuration.Configuration, dependency core_dependency.DependencyContainer) *DependencyContainer {
	once.Do(func() {
		repositoryWeb := loadRepositoryWeb(config)

		managerWeb := loadManagerWeb(repositoryWeb)

		container := &DependencyContainer{
			DependencyContainer: dependency,
			ManagerWeb:          managerWeb,
		}

		instance = container
	})
	return instance
}

func loadRepositoryWeb(config configuration.Configuration) repository.IRepositoryWeb {
	var file core_repository.IFileManager[web_domain.WebData]
	file = core_repository.NewManagerCsvtFile[web_domain.WebData](repository.CSVT_FILE_PATH_WEB_DATA)

	snapshot := config.Snapshot()
	if snapshot.Enable {
		topic := system.SNAPSHOT_TOPIC_WEB_DATA
		file = loadManagerSnapshotFile(topic, snapshot, file)
	}

	impl := collection.DictionarySyncEmpty[string, web_domain.WebData]()
	repository, err := web.InitializeRepositoryMemory(impl, file)
	if err != nil {
		log.Panic(err)
	}

	return repository
}

func loadManagerSnapshotFile[T core_repository.IStructure](topic core_system.TopicSnapshot, snapshot core_configuration.Snapshot, file core_repository.IFileManager[T]) core_repository.IFileManager[T] {
	return core_repository.
		BuilderManagerSnapshotFile(topic, file).
		Limit(snapshot.Limit).
		Time(snapshot.Time).
		Make()
}

func loadManagerWeb(web repository.IRepositoryWeb) *repository.ManagerWeb {
	return repository.NewManagerWeb(web)
}
