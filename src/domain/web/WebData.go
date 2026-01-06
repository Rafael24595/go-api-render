package web

type WebData struct {
	Id        string            `json:"id"`
	Timestamp int64             `json:"timestamp"`
	Data      map[string]string `json:"data"`
	Modified  int64             `json:"modified"`
	Owner     string            `json:"owner"`
}

func EmptyWebData(owner string) *WebData {
	return &WebData{
		Timestamp: 0,
		Data: make(map[string]string),
		Modified: 0,
		Owner: owner,
	}
}

func (r WebData) PersistenceId() string {
	return r.Id
}
