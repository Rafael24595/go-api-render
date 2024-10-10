package infrastructure

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Rafael24595/go-api-core/src/commons"
	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/domain/auth"
	"github.com/Rafael24595/go-api-core/src/domain/body"
	"github.com/Rafael24595/go-api-core/src/domain/header"
	"github.com/Rafael24595/go-api-core/src/domain/query"
	"github.com/Rafael24595/go-api-render/src/commons/configuration"
	"github.com/google/uuid"
)

var constants = configuration.GetConstants()

func proccessRequestAnonymous(r *http.Request) (*domain.Request, error) {
	name := fmt.Sprintf("temp_%s", uuid.NewString())
	return proccessRequest(r, name)
}

func proccessRequest(r *http.Request, name string) (*domain.Request, error) {
	method, err := proccessMethod(r)
	if err != nil {
		return nil, err
	}

	url, err := proccessUrl(r)
	if err != nil {
		return nil, err
	}

	request := domain.NewRequest(name, *method, url)

	request.Id = r.FormValue(constants.Client.Id)

	queries, err := proccessQueryParams(r)
	if err != nil {
		return nil, err
	}
	request.Queries = *queries

	headers, err := proccessHeaders(r)
	if err != nil {
		return nil, err
	}
	request.Headers = *headers

	body, err := proccessBody(r)
	if err != nil {
		return nil, err
	}
	request.Body = *body

	auth, err := proccessAuth(r)
	if err != nil {
		return nil, err
	}
	request.Auths = *auth

	return request, nil
}

func proccessMethod(r *http.Request) (*domain.HttpMethod, error) {
	sMethod := r.FormValue(constants.Client.Method)
	if sMethod == "" {
		return nil, commons.ApiErrorFrom(422, "Method is not defined")
	}
	return domain.HttpMethodFromString(sMethod)
}

func proccessUrl(r *http.Request) (string, error) {
	url := r.FormValue(constants.Client.Uri)
	if url == "" {
		return "", commons.ApiErrorFrom(422, "URL is not defined")
	}
	return url, nil
}

func proccessQueryParams(r *http.Request) (*query.Queries, error) {
	form := r.Form

	uuids := collection.FromMap(form).
		KeysCollection().
		Filter(func(k string) bool {
			return strings.Contains(k, constants.Format.FormatKey(constants.Query.Name, ""))
		}).
		MapSelf(func(k string) string {
			parts := strings.Split(k, constants.Format.KeySeparator)
			return parts[len(parts)-1]
		}).
		Collect()

	queries := query.NewQueries()
	for _, uuid := range uuids {
		status := form.Get(constants.Format.FormatKey(constants.Query.Status, uuid)) == "on"
		name := form.Get(constants.Format.FormatKey(constants.Query.Name, uuid))
		value := form.Get(constants.Format.FormatKey(constants.Query.Value, uuid))

		if name == "" {
			continue
		}

		queries.Add(query.NewQuery(status, name, value))
	}

	return queries, nil
}

func proccessHeaders(r *http.Request) (*header.Headers, error) {
	form := r.Form

	uuids := collection.FromMap(form).
	KeysCollection().
		Filter(func(k string) bool {
			return strings.Contains(k, constants.Format.FormatKey(constants.Header.Name, ""))
		}).
		MapSelf(func(k string) string {
			parts := strings.Split(k, constants.Format.KeySeparator)
			return parts[len(parts)-1]
		}).
		Collect()

	headers := header.NewHeaders()
	for _, uuid := range uuids {
		status := form.Get(constants.Format.FormatKey(constants.Header.Status, uuid)) == "on"
		name := form.Get(constants.Format.FormatKey(constants.Header.Name, uuid))
		value := form.Get(constants.Format.FormatKey(constants.Header.Value, uuid))

		if name == "" {
			continue
		}

		headers.Add(header.NewHeader(status, name, value))
	}

	return headers, nil
}

func proccessBody(r *http.Request) (*body.Body, error) {
	form := r.Form
	bodyType := form.Get(constants.Body.Type)
	contentType, _ := constants.Body.ContentTypeFromType(bodyType)
	parameterType, _ := constants.Body.BodyTypeFromType(bodyType)
	payload := form.Get(parameterType)
	return body.NewBodyString(contentType, payload), nil
}

func proccessAuth(r *http.Request) (*auth.Auths, error) {
	authStatus := r.Form.Get(constants.Auth.Enabled) == "on"

	auths := auth.NewAuths(authStatus)

	basic, err := proccessAuthBasic(r)
	if err != nil {
		return nil, err
	}
	if basic != nil {
		auths.PutAuth(*basic)
	}

	bearer, err := proccessAuthBearer(r)
	if err != nil {
		return nil, err
	}

	if bearer != nil {
		auths.PutAuth(*bearer)
	}

	return auths, nil
}

func proccessAuthBasic(r *http.Request) (*auth.Auth, error) {
	form := r.Form

	status := strings.ToLower(form.Get(constants.Auth.Basic.Status)) == "on"
	user := form.Get(constants.Auth.Basic.User)
	password := form.Get(constants.Auth.Basic.Password)

	return auth.NewAuthEmpty(status, auth.Basic).
		PutParam(auth.BASIC_PARAM_USER, user).
		PutParam(auth.BASIC_PARAM_PASSWORD, password), nil
}

func proccessAuthBearer(r *http.Request) (*auth.Auth, error) {
	form := r.Form

	status := strings.ToLower(form.Get(constants.Auth.Bearer.Status)) == "on"
	prefix := auth.DEFAULT_BEARER_PREFIX
	if form.Has(constants.Auth.Bearer.Prefix) {
		prefix = form.Get(constants.Auth.Bearer.Prefix)
	}
	token := form.Get(constants.Auth.Bearer.Token)

	return auth.NewAuthEmpty(status, auth.Bearer).
		PutParam(auth.BEARER_PARAM_PREFIX, prefix).
		PutParam(auth.BEARER_PARAM_TOKEN, token), nil
}
