package configuration

import "strings"

type Constants struct {
	Format Format
	Client Client
	Query  Query
}

type Format struct {
	KeySeparator string
}

type Client struct {
	ClientType string
	TagQuery   string
	TagHeader  string
	TagAuth    string
	TagBody    string
}

type Query struct {
	QueryStatus string
	QueryName   string
	QueryValue  string
}

var constants = getConstants()

func getConstants() Constants {
	return Constants{
		Format: Format{
			KeySeparator: "#",
		},
		Client: Client{
			ClientType: "client-type",
			TagQuery:   "query",
			TagHeader:  "header",
			TagAuth:    "auth",
			TagBody:    "body",
		},
		Query: Query{
			QueryStatus: "query-status",
			QueryName:   "query-name",
			QueryValue:  "query-value",
		},
	}
}

func GetConstants() Constants {
	return constants
}

func (f Format) FormatKey(keys ...string) string {
	return strings.Join(keys, f.KeySeparator)
}