package configuration

import (
	"os"
	"time"

	core_configuration "github.com/Rafael24595/go-api-core/src/commons/configuration"
	"github.com/Rafael24595/go-api-core/src/commons/log"
	"github.com/Rafael24595/go-api-core/src/commons/utils"
)

const defaultPort = 8080
const defaultCert = "./cert/cert.pem"
const defaultKey = "./cert/key.pem"

var instance *Configuration

type Configuration struct {
	core_configuration.Configuration
	Release *core_configuration.Release
	Front   FrontPackage
	debug   bool
	port    int
	portTLS int
	cert    string
	key     string
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

	portTLS, cert, key := findCertificate(kargs)

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
		portTLS:       portTLS,
		cert:          cert,
		key:           key,
	}

	go instance.originLastVersion(kargs)

	return *instance
}

func findCertificate(kargs map[string]utils.Any) (int, string, string) {
	if tls, _ := kargs["GO_API_SERVER_ENABLE_TLS"].Bool(); !tls {
		return 0, "", ""
	}

	portTLS, ok := kargs["GO_API_SERVER_PORT_TLS"].Int()
	if !ok {
		portTLS = 0
	}

	cert, ok := kargs["GO_API_SERVER_CERT"].String()
	if !ok || cert == "" {
		_, err := os.Stat(defaultCert)
		if !os.IsNotExist(err) {
			log.Message("Certificate flag is not defined; using default certificate")
			cert = defaultCert
		}
	}

	key, ok := kargs["GO_API_SERVER_KEY"].String()
	if !ok || key == "" {
		_, err := os.Stat(defaultKey)
		if !os.IsNotExist(err) {
			log.Message("Key flag is not defined; using default key")
			key = defaultKey
		}
	}

	_, err := os.Stat(cert)
	if os.IsNotExist(err) {
		log.Messagef("Certificate file '%s' does not exist", cert)
	}

	_, err = os.Stat(key)
	if os.IsNotExist(err) {
		log.Messagef("Key file '%s' does not exist", key)
	}

	return portTLS, cert, key
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

func (c Configuration) Secure() bool {
	return c.portTLS != 0 && c.cert != "" && c.key != ""
}

func (c Configuration) PortTLS() int {
	return c.portTLS
}

func (c Configuration) Cert() string {
	return c.cert
}

func (c Configuration) Key() string {
	return c.key
}
