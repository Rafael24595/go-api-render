package configuration

type Constants struct {
	Client Client
}

type Client struct {
	ClientType string
	TagQuery   string
	TagHeader  string
	TagAuth    string
	TagBody    string
}

var constants = getConstants()

func getConstants() Constants {
	return Constants{
		Client: Client{
			ClientType: "client-type",
			TagQuery:   "query",
			TagHeader:  "header",
			TagAuth:    "auth",
			TagBody:    "body",
		},
	}
}

func GetConstants() Constants {
	return constants
}