package router

import (
	"fmt"
	"net/http"
	"regexp"
)

type Router struct {
	routeMap map[string]route
}

func NewRouter() *Router {
	router := &Router {
		routeMap: map[string]route{},
	}
	http.HandleFunc("/", router.routeHandler)
	return router
}

func (router *Router) ResourcesPath(path string) *Router {
	fs := http.FileServer(http.Dir(path))
	route := fmt.Sprintf("/%s/", path)
    http.Handle(route, http.StripPrefix(route, fs))
	return router
}

func (router *Router) Route(method string, uri string, handler func(http.ResponseWriter, *Request)) *Router {
	key := fmt.Sprintf("%s#%s", method, uri)
	route := route {
		method: method,
		route: uri,
		handler: handler,
	}
	router.routeMap[key] = route
	return router
}

func (router Router) routeHandler(w http.ResponseWriter, r *http.Request) {
	for _, v := range router.routeMap {
		matched, args := v.matches(r.URL.Path)
		if matched && r.Method == v.method  {
			v.handler(w, newCustomRequest(args, r))
			return
		}
	}
	http.Error(w, "Not found", http.StatusNotFound)
}

func (route route) matches(str string) (bool, map[string]string) {
	regexPattern := route.convertTemplateToRegex(route.route)

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return false, map[string]string{}
	}

	matched := re.MatchString(str)

	matchValues := re.FindStringSubmatch(str)
	if matchValues == nil {
		return matched, map[string]string{}
	}

	matchName := re.FindStringSubmatch(route.route)
	if matchName == nil {
		return matched, map[string]string{}
	}

	params := make(map[string]string)
	for i, name := range matchName {
		if i != 0 && i < len(matchName) {
			params[name[1:]] = matchValues[i]
		}
	}

	return matched, params
}

func (route route) convertTemplateToRegex(template string) string {
	escapedTemplate := regexp.QuoteMeta(template)

	re := regexp.MustCompile(`:[\w-]+`)
	escapedTemplate = re.ReplaceAllString(escapedTemplate, `([^/]+)`)

	escapedTemplate = `^` + escapedTemplate + `$`

	return escapedTemplate
}