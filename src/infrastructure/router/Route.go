package router

import "net/http"

type route struct {
	method  string
	route   string
	handler func(http.ResponseWriter, *Request)
}