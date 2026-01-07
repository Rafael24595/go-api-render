package web

import (
	"sync"
	"time"

	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-render/src/commons/system"
	core_system "github.com/Rafael24595/go-api-core/src/commons/system"
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	web_domain "github.com/Rafael24595/go-api-render/src/domain/web"
	"github.com/Rafael24595/go-collections/collection"
	"github.com/google/uuid"
)

type RepositoryMemory struct {
	once       sync.Once
	muMemory   sync.RWMutex
	muFile     sync.RWMutex
	collection collection.IDictionary[string, web_domain.WebData]
	file       repository.IFileManager[web_domain.WebData]
	close      chan bool
}

func InitializeRepositoryMemory(impl collection.IDictionary[string, web_domain.WebData], file repository.IFileManager[web_domain.WebData]) (*RepositoryMemory, error) {
	requests, err := file.Read()
	if err != nil {
		return nil, err
	}

	instance := &RepositoryMemory{
		collection: impl.Merge(collection.DictionaryFromMap(requests)),
		file:       file,
	}

	go instance.watch()

	return instance, nil
}

func (r *RepositoryMemory) watch() {
	r.once.Do(func() {
		conf := configuration.Instance()
		if !conf.Snapshot().Enable {
			return
		}

		hub := make(chan core_system .SystemEvent, 1)
		defer close(hub)

		topics := []string{
			system.SNAPSHOT_TOPIC_WEB_DATA.TopicSnapshotApplyOutput(),
		}

		conf.EventHub.Subcribe(repository.RepositoryListener, hub, topics...)
		defer conf.EventHub.Unsubcribe(repository.RepositoryListener, topics...)

		for {
			select {
			case <-r.close:
				log.Customf(repository.SnapshotCategory, "Watcher stopped: local close signal received.")
				return
			case <-hub:
				if err := r.read(); err != nil {
					log.Custome(repository.SnapshotCategory, err)
					return
				}
			case <-conf.Signal.Done():
				log.Customf(repository.SnapshotCategory, "Watcher stopped: global shutdown signal received.")
				return
			}
		}
	})
}

func (r *RepositoryMemory) read() error {
	requests, err := r.file.Read()
	if err != nil {
		return err
	}

	r.collection = collection.DictionaryFromMap(requests)
	return nil
}

func (r *RepositoryMemory) Find(id string) (*web_domain.WebData, bool) {
	r.muMemory.RLock()
	defer r.muMemory.RUnlock()
	return r.collection.Get(id)
}

func (r *RepositoryMemory) FindByOwner(owner string) (*web_domain.WebData, bool) {
	r.muMemory.RLock()
	defer r.muMemory.RUnlock()
	return r.collection.FindOne(func(s string, w web_domain.WebData) bool {
		return w.Owner == owner
	})
}

func (r *RepositoryMemory) Resolve(owner string, webData *web_domain.WebData) *web_domain.WebData {
	r.muMemory.Lock()
	defer r.muMemory.Unlock()

	if webData.Id != "" {
		return r.insert(owner, webData)
	}

	key := uuid.New().String()
	if r.collection.Exists(key) {
		return r.Resolve(owner, webData)
	}

	webData.Id = key

	return r.insert(owner, webData)
}

func (r *RepositoryMemory) insert(owner string, webData *web_domain.WebData) *web_domain.WebData {
	webData.Owner = owner

	now := time.Now().UnixMilli()

	if webData.Timestamp == 0 {
		webData.Timestamp = now
	}

	webData.Modified = now

	r.collection.Put(webData.Id, *webData)

	go r.write(r.collection)

	return webData
}

func (r *RepositoryMemory) Delete(webData *web_domain.WebData) *web_domain.WebData {
	r.muMemory.Lock()
	defer r.muMemory.Unlock()

	cursor, _ := r.collection.Remove(webData.Id)
	go r.write(r.collection)

	return cursor
}

func (r *RepositoryMemory) write(snapshot collection.IDictionary[string, web_domain.WebData]) {
	r.muFile.Lock()
	defer r.muFile.Unlock()

	err := r.file.Write(snapshot.Values())
	if err != nil {
		log.Error(err)
	}
}
