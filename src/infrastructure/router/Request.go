package router

import "net/http"

type Request struct {
	PathArgs map[string]string
	Request *http.Request
}

func newCustomRequest(args map[string]string, request *http.Request) *Request {
	return &Request{
		PathArgs: args,
		Request: request,
	}
}