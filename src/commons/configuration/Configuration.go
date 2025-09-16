package configuration

import (
	"os"
	"sync"
	"time"

	core_configuration "github.com/Rafael24595/go-api-core/src/commons/configuration"
	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/commons/utils"
)

const defaultPort = 8080
const defaultCert = "./cert/cert.pem"
const defaultKey = "./cert/key.pem"

var (
	instance *Configuration
	once     sync.Once
)

type Configuration struct {
	core_configuration.Configuration
	Release       *core_configuration.Release
	Front         FrontPackage
	debug         bool
	port          int
	onlyTLS       bool
	portTLS       int
	certTLS       string
	keyTLS        string
	enableSecrets bool
}

func Initialize(core *core_configuration.Configuration, kargs map[string]utils.Argument, frontPackage *FrontPackage) Configuration {
	once.Do(func() {
		debug := kargs["GO_API_DEBUG"].Boold(false)

		port, ok := kargs["GO_API_SERVER_PORT"].Int()
		if !ok {
			log.Messagef("Custom port flag is not defined; using default port %d", defaultPort)
			port = defaultPort
		}

		front := kargs["GO_API_SERVER_FRONT"].Boold(false)
		if !front {
			log.Message("Front flag is not defined; the frontend application will not be displayed")
			front = false
		}

		portTLS, certTLS, keyTLS, onlyTLS := tlsArgs(kargs)

		frontPackage.Enabled = front
		if !front {
			frontPackage.Name = ""
			frontPackage.Version = ""
		}

		enableSecrets := kargs["GO_API_ENABLE_SECRETS"].Boold(false)

		instance = &Configuration{
			Configuration: *core,
			Front:         *frontPackage,
			debug:         debug,
			port:          port,
			onlyTLS:       onlyTLS,
			portTLS:       portTLS,
			certTLS:       certTLS,
			keyTLS:        keyTLS,
			enableSecrets: enableSecrets,
		}

		go instance.originLastVersion(kargs)
	})

	if instance == nil {
		log.Panics("The configuration is not initialized properly")
	}

	return *instance
}

func tlsArgs(kargs map[string]utils.Argument) (int, string, string, bool) {
	if tls, _ := kargs["GO_API_SERVER_ENABLE_TLS"].Bool(); !tls {
		return 0, "", "", false
	}

	onlyTLS := kargs["GO_API_SERVER_ONLY_TLS"].Boold(false)
	portTLS := kargs["GO_API_SERVER_PORT_TLS"].Intd(0)

	certTLS := kargs["GO_API_SERVER_CERT"].String()
	if certTLS == "" {
		_, err := os.Stat(defaultCert)
		if !os.IsNotExist(err) {
			log.Message("Certificate flag is not defined; using default certificate")
			certTLS = defaultCert
		}
	}

	keyTLS := kargs["GO_API_SERVER_KEY"].String()
	if keyTLS == "" {
		_, err := os.Stat(defaultKey)
		if !os.IsNotExist(err) {
			log.Message("Key flag is not defined; using default key")
			keyTLS = defaultKey
		}
	}

	_, err := os.Stat(certTLS)
	if os.IsNotExist(err) {
		log.Warningf("Certificate file '%s' does not exist, TLS connection aborted", certTLS)
		return 0, "", "", false
	}

	_, err = os.Stat(keyTLS)
	if os.IsNotExist(err) {
		log.Warningf("Key file '%s' does not exist, TLS connection aborted", keyTLS)
		return 0, "", "", false
	}

	return portTLS, certTLS, keyTLS, onlyTLS
}

func (c *Configuration) originLastVersion(kargs map[string]utils.Argument) {
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

func (c Configuration) OnlyTLS() bool {
	return c.EnableTLS() && c.onlyTLS
}

func (c Configuration) EnableTLS() bool {
	return c.portTLS != 0 && c.certTLS != "" && c.keyTLS != ""
}

func (c Configuration) PortTLS() int {
	return c.portTLS
}

func (c Configuration) CertTLS() string {
	return c.certTLS
}

func (c Configuration) KeyTLS() string {
	return c.keyTLS
}

func (c Configuration) EnableSecrets() bool {
	return c.enableSecrets
}
