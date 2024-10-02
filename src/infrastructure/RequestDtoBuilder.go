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
	"github.com/google/uuid"
)

const (
	AUTH_TYPE           = "auth-type"
	AUTH_BASIC_STATUS   = "auth-basic-status"
	AUTH_BASIC_USER     = "auth-basic-user"
	AUTH_BASIC_PASSWORD = "auth-basic-password"
	AUTH_BEARER_STATUS  = "auth-bearer-status"
	AUTH_BEARER_PREFIX  = "auth-bearer-prefix"
	AUTH_BEARER_TOKEN   = "auth-bearer-token"
)

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
	sMethod := r.FormValue("method")
	if sMethod == "" {
		return nil, commons.ApiErrorFrom(422, "Method is not defined")
	}
	return domain.HttpMethodFromString(sMethod)
}

func proccessUrl(r *http.Request) (string, error) {
	url := r.FormValue("url")
	if url == "" {
		return "", commons.ApiErrorFrom(422, "URL is not defined")
	}
	return url, nil
}

func proccessQueryParams(r *http.Request) (*query.Queries, error) {
	form := r.Form

	uuids := collection.FromMap(form).
		KeysList().
		Filter(func(k string) bool {
			return strings.Contains(k, "query-name#")
		}).
		MapSelf(func(k string) string {
			parts := strings.Split(k, "#")
			return parts[len(parts)-1]
		}).
		Collect()

	queries := query.NewQueries()
	for _, uuid := range uuids {
		sStatus := form.Get(fmt.Sprintf("query-status#%s", uuid))

		status := strings.ToLower(sStatus) == "on"
		name := form.Get(fmt.Sprintf("query-name#%s", uuid))
		value := form.Get(fmt.Sprintf("query-value#%s", uuid))

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
		KeysList().
		Filter(func(k string) bool {
			return strings.Contains(k, "header-name#")
		}).
		MapSelf(func(k string) string {
			parts := strings.Split(k, "#")
			return parts[len(parts)-1]
		}).
		Collect()

	headers := header.NewHeaders()
	for _, uuid := range uuids {
		sStatus := form.Get(fmt.Sprintf("header-status#%s", uuid))

		status := strings.ToLower(sStatus) == "on"
		name := form.Get(fmt.Sprintf("header-name#%s", uuid))
		value := form.Get(fmt.Sprintf("header-value#%s", uuid))

		if name == "" {
			continue
		}

		headers.Add(header.NewHeader(status, name, value))
	}

	return headers, nil
}

func proccessBody(r *http.Request) (*body.Body, error) {
	form := r.Form
	bodyType := form.Get("body-type")
	contentType, _ := body.ContentTypeFromString(bodyType)
	payload := form.Get(fmt.Sprintf("body-parameter-%s", bodyType))
	return body.NewBodyString(contentType, payload), nil
}

func proccessAuth(r *http.Request) (*auth.Auths, error) {
	auths := auth.NewAuths()

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

	status := strings.ToLower(form.Get(AUTH_BASIC_STATUS)) == "on"
	user := form.Get(AUTH_BASIC_USER)
	password := form.Get(AUTH_BASIC_PASSWORD)

	return auth.NewAuthEmpty(status, auth.Basic).
		PutParam(auth.BASIC_PARAM_USER, user).
		PutParam(auth.BASIC_PARAM_PASSWORD, password), nil
}

func proccessAuthBearer(r *http.Request) (*auth.Auth, error) {
	form := r.Form

	status := strings.ToLower(form.Get(AUTH_BEARER_STATUS)) == "on"
	prefix := auth.DEFAULT_BEARER_PREFIX
	if form.Has(AUTH_BEARER_PREFIX) {
		prefix = form.Get(AUTH_BEARER_PREFIX)
	}
	token := form.Get(AUTH_BEARER_TOKEN)

	return auth.NewAuthEmpty(status, auth.Bearer).
		PutParam(auth.BEARER_PARAM_PREFIX, prefix).
		PutParam(auth.BEARER_PARAM_TOKEN, token), nil
}
