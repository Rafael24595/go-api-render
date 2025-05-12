package datalist

import "sync"

type DataListLoader[T any] struct {
	mut sync.Mutex
	dataList *DataList
	loader   func(T) []string
}

func newDataListLoader[T any](DataList *DataList, loader func(T) []string) *DataListLoader[T] {
	return &DataListLoader[T]{
		dataList: DataList,
		loader:   loader,
	}
}

func (l *DataListLoader[T]) Update(item T) {
	l.mut.Lock()
	defer l.mut.Unlock()
	l.dataList.PushOption(l.loader(item)...)
}