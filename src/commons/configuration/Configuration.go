package configuration

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var instance *Configuration

type Configuration struct {
	debug bool
	kargs map[string]string
}

func Initialize(kargs map[string]string) Configuration {
	if instance != nil {
		panic("")
	}
	
	instance = &Configuration{
		debug : strings.ToLower(kargs["DEBUG"]) == "true",
		kargs: kargs,
	}

	return *instance
}

func Instance() Configuration {
	if instance == nil {
		panic("")
	}
	return *instance
}

func (i Configuration) Debug() bool {
	return i.debug;
}

func ReadEnv(file string) map[string]string {
	if len(file) > 0 {
		if err := godotenv.Load(".env"); err != nil {
			panic(fmt.Sprintf("Error loading %s file", file))
		}
	}

	envs := make(map[string]string)
	for _, env := range os.Environ() {
		pair := splitEnv(env)
		envs[pair[0]] = pair[1]
	}

	return envs
}

func splitEnv(env string) []string {
	var pair []string
	for i, char := range env {
		if char == '=' {
			pair = append(pair, env[:i], env[i+1:])
			break
		}
	}
	return pair
}
