package datalist

import "sync"

type DataListManager struct {
	mu        sync.Mutex
	DataLists map[string]*DataList
	Loaders   map[string][]*DataListLoader[any]
}

func NewDataListManager() *DataListManager {
	return &DataListManager{
		DataLists: make(map[string]*DataList),
		Loaders:   make(map[string][]*DataListLoader[any]),
	}
}

func (m *DataListManager) Keys() []string {
	keys := make([]string, 0, len(m.DataLists))
	for k := range m.DataLists {
		keys = append(keys, k)
	}
	return keys
}

func (m *DataListManager) DataList(key string) (*DataList, bool) {
	dataList, ok := m.DataLists[key]
	return dataList, ok 
}

func (m *DataListManager) PutLoader(id, trigger string, loader func(any) []string) *DataListManager {
	dataList, ok := m.DataLists[id]
	if !ok {
		dataList = NewDataList(id)
		m.DataLists[id] = dataList
	}

	loaders, ok := m.Loaders[trigger]
	if !ok {
		loaders = make([]*DataListLoader[any], 0)
	}

	m.Loaders[trigger] = append(loaders, newDataListLoader(dataList, loader))

	return m
}

func (m *DataListManager) PutStatic(id string, list []string) *DataListManager {
	dataList := NewDataList(id)
	dataList.PushOption(list...)
	m.DataLists[id] = dataList
	return m
}

func (m *DataListManager) PushData(trigger string, item any) bool {
	loader, ok := m.Loaders[trigger]
	if ok {
		for _, v := range loader {
			v.Update(item)
		}
	}
	return ok
}

func (m *DataListManager) Clean() bool {
	for _, v := range m.DataLists {
		v.Clean()
	}
	return true
}