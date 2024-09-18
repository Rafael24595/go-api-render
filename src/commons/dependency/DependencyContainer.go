package dependency

import (
	"github.com/Rafael24595/go-api-core/src/infrastructure/repository/request/memory"
)

var instance *DependencyContainer

type DependencyContainer struct {
	RequestQuery *memory.QueryMemory
	RequestCommand *memory.CommandMemory
}

func Initialize() *DependencyContainer {
	if instance != nil {
		panic("//TODO: Yet instanced")
	}

	requestQuery, err := memory.InitializeQueryMemory()
	if err != nil {
		panic(err)
	}

	requestCommand := memory.NewCommandMemory(requestQuery)

	container := &DependencyContainer{
		RequestQuery: requestQuery,
		RequestCommand: requestCommand,
	}

	instance = container
	
	return instance

}