package configuration

import (
	"strings"

	"github.com/Rafael24595/go-api-core/src/domain/body"
)

type Constants struct {
	Format         Format
	Client         Client
	Collection     Collection
	Query          Query
	Header         Header
	Body           Body
	Auth           Auth
	SidebarRequest SidebarRequest
	Response       Response
}

type Format struct {
	KeySeparator string
}

type Client struct {
	Id        string
	Name      string
	Type      string
	Method    string
	Uri       string
	TagQuery  string
	TagHeader string
	TagAuth   string
	TagBody   string
	DoRequest string
}

type Collection struct {
	Variable Variable
}

type Variable struct {
	Status string
	Name   string
	Value  string
	Type   string
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
	TypeOf    TypeOf
}

type TypeOf struct {
	None body.ContentType
	Text body.ContentType
	Xml  body.ContentType
	Html body.ContentType
	Json body.ContentType
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

type SidebarRequest struct {
	Type        string
	TypeView    string
	TagHistoric string
	TagSaved    string
}

type Response struct {
	Type       string
	TagPayload string
	TagHeader  string
	TagCookie  string
}

var constants = getConstants()

func getConstants() Constants {
	return Constants{
		Format: Format{
			KeySeparator: "#",
		},
		Client: Client{
			Id:        "id",
			Name:      "name",
			Type:      "client-type",
			Method:    "method",
			Uri:       "uri",
			TagQuery:  "query",
			TagHeader: "header",
			TagAuth:   "auth",
			TagBody:   "body",
			DoRequest: "do-request",
		},
		Collection: Collection{
			Variable: Variable{
				Status: "collection-variable-status",
				Name:   "collection-variable-name",
				Value:  "collection-variable-value",
				Type:   "collection-variable-type",
			},
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
			TypeOf: TypeOf{
				None: body.None,
				Text: body.Text,
				Xml:  body.Xml,
				Html: body.Html,
				Json: body.Json,
			},
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
		SidebarRequest: SidebarRequest{
			Type:        "request-type",
			TypeView:    "request-view",
			TagHistoric: "historic",
			TagSaved:    "saved",
		},
		Response: Response{
			Type:       "request-container",
			TagPayload: "payload",
			TagHeader:  "header",
			TagCookie:  "cookie",
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
