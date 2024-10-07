package configuration

import (
	"strings"

	"github.com/Rafael24595/go-api-core/src/domain/body"
)

type Constants struct {
	Format Format
	Client Client
	Query  Query
	Header Header
	Body   Body
	Auth   Auth
}

type Format struct {
	KeySeparator string
}

type Client struct {
	Id        string
	Type      string
	Method    string
	Uri       string
	TagQuery  string
	TagHeader string
	TagAuth   string
	TagBody   string
}

type Query struct {
	Status string
	Name   string
	Value  string
}

type Header struct {
	Status string
	Name   string
	Value  string
}

type Body struct {
	Type      string
	TagText   string
	TagJson   string
	ValueText string
	ValueJson string
}

type Auth struct {
	Type      string
	Enabled   string
	TagBasic  string
	TagBearer string
	Basic     Basic
	Bearer    Bearer
}

type Basic struct {
	Status   string
	User     string
	Password string
}

type Bearer struct {
	Status string
	Prefix string
	Token  string
}

var constants = getConstants()

func getConstants() Constants {
	return Constants{
		Format: Format{
			KeySeparator: "#",
		},
		Client: Client{
			Id:        "id",
			Type:      "client-type",
			Method:    "method",
			Uri:       "uri",
			TagQuery:  "query",
			TagHeader: "header",
			TagAuth:   "auth",
			TagBody:   "body",
		},
		Query: Query{
			Status: "query-status",
			Name:   "query-name",
			Value:  "query-value",
		},
		Header: Header{
			Status: "header-status",
			Name:   "header-name",
			Value:  "header-value",
		},
		Body: Body{
			Type:      "body-type",
			TagText:   "text",
			TagJson:   "json",
			ValueText: "body-text",
			ValueJson: "body-json",
		},
		Auth: Auth{
			Type:      "auth-type",
			Enabled:   "auth-enable",
			TagBasic:  "basic",
			TagBearer: "bearer",
			Basic: Basic{
				Status:   "auth-basic-status",
				User:     "auth-basic-user",
				Password: "auth-basic-password",
			},
			Bearer: Bearer{
				Status: "auth-bearer-status",
				Prefix: "auth-bearer-prefix",
				Token:  "auth-bearer-token",
			},
		},
	}
}

func GetConstants() Constants {
	return constants
}

func (b Body) BodyTypeFromType(typ string) (string, bool) {
	switch typ {
	case b.TagText:
		return b.ValueText, true
	case b.TagJson:
		return b.ValueJson, true
	}
	return b.ValueText, false
}

func (b Body) ContentTypeFromType(typ string) (body.ContentType, bool) {
	return body.ContentTypeFromString(typ)
}

func (f Format) FormatKey(keys ...string) string {
	return strings.Join(keys, f.KeySeparator)
}
