package controller

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/Rafael24595/go-api-core/src/domain"
)

func valideRequest(request *domain.Request) (int, error) {
	uri := strings.TrimSpace(request.Uri)

	if _, err := url.ParseRequestURI(uri); err != nil {
		return http.StatusUnprocessableEntity, err
	}

	return 0, nil
}
