package configuration

import (
	"time"

	core_configuration "github.com/Rafael24595/go-api-core/src/commons/configuration"
	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/commons/utils"
)

const defaultPort = 8080

var instance *Configuration

type Configuration struct {
	core_configuration.Configuration
	Release *core_configuration.Release
	Front   FrontPackage
	debug   bool
	port    int
}

func Initialize(core *core_configuration.Configuration, kargs map[string]utils.Any, frontPackage *FrontPackage) Configuration {
	if instance != nil {
		log.Panics("The configuration is alredy initialized")
	}

	debug, ok := kargs["GO_API_DEBUG"].Bool()
	if !ok {
		debug = false
	}

	port, ok := kargs["GO_API_SERVER_PORT"].Int()
	if !ok {
		log.Messagef("Custom port flag is not defined; using default port %d", defaultPort)
		port = defaultPort
	}

	front, ok := kargs["GO_API_SERVER_FRONT"].Bool()
	if !ok {
		log.Message("Front flag is not defined; the frontend application will not be displayed")
		front = false
	}

	frontPackage.Enabled = front
	if !front {
		frontPackage.Name = ""
		frontPackage.Version = ""
	}

	instance = &Configuration{
		Configuration: *core,
		Front:         *frontPackage,
		debug:         debug,
		port:          port,
	}

	go instance.originLastVersion(kargs)

	return *instance
}

func (c *Configuration) originLastVersion(kargs map[string]utils.Any) {
	releaseTime, ok := kargs["GO_API_FETCH_RELEASE_TIME"].Int()
	if !ok || releaseTime < 1 {
		log.Message("Fetch release time is not defined; new releases will not be fetched")
		return
	}

	ticker := time.NewTicker(time.Duration(releaseTime) * time.Hour)
	defer ticker.Stop()

	fetchLastVersion(c)

	for {
		select {
		case <-c.Signal.Done():
			return
		case <-ticker.C:
			fetchLastVersion(c)
		}
	}
}

func fetchLastVersion(c *Configuration) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered from panic: %v", r)
		}
	}()

	release := core_configuration.OriginLastVersion("Rafael24595", "go-api-render")
	if release != nil {
		if release.TagName != c.Project.Version {
			log.Messagef("New release has been found %s", release.TagName)
		}
		c.Release = release
	}
}

func Instance() Configuration {
	if instance == nil {
		log.Panics("The configuration is not initialized yet")
	}
	return *instance
}

func (c Configuration) Debug() bool {
	return c.debug
}

func (c Configuration) Port() int {
	return c.port
}
