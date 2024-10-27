package datalist

type DataList struct {
	Id      string
	Static  bool
	Options map[string]bool
}

func NewDataList(id string) *DataList {
	return newStaticDataList(id, false, make(map[string]bool))
}

func NewStaticDataList(id string) *DataList {
	return newStaticDataList(id, true, make(map[string]bool))
}

func newStaticDataList(id string, static bool, options map[string]bool) *DataList {
	return &DataList{
		Id:      id,
		Static:  static,
		Options: options,
	}
}

func (l *DataList) PushOption(option ...string) {
	for _, v := range option {
		l.Options[v] = true
	}
}

func (l *DataList) Clean() bool {
	if l.Static {
		return false
	}
	l.Options = make(map[string]bool)
	return true
}
