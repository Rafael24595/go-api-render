package infrastructure

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Rafael24595/go-api-core/src/commons"
	"github.com/Rafael24595/go-api-core/src/commons/collection"
	"github.com/Rafael24595/go-api-core/src/domain"
	"github.com/Rafael24595/go-api-core/src/domain/header"
	"github.com/Rafael24595/go-api-core/src/domain/query"
	"github.com/google/uuid"
)

func proccessRequest(r *http.Request) (*domain.Request, error) {
	name := fmt.Sprintf("temp_%s", uuid.NewString())

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
	form := r.URL.Query()

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