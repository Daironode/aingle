 
package restful

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

type paramsMap map[string]string

//http router
type Route struct {
	Method  string
	Path    *regexp.Regexp
	Params  []string
	Handler http.HandlerFunc
}
type Router struct {
	routes []*Route
}

func NewRouter() *Router {
	return &Router{}
}

func (this *Router) Try(path string, method string) (http.HandlerFunc, paramsMap, error) {

	for _, route := range this.routes {
		if route.Method == method {
			match := route.Path.MatchString(path)
			if !match {
				continue
			}
			params := paramsMap{}
			if len(route.Params) > 0 {
				params = parseParams(route, path)
			}
			return route.Handler, params, nil

		}
	}
	return nil, paramsMap{}, errors.New("Route not found")

}

func (this *Router) add(method string, path string, handler http.HandlerFunc) {
	route := &Route{}
	route.Method = method
	path = "^" + path + "$"
	route.Handler = handler

	if strings.Contains(path, ":") {
		matches := regexp.MustCompile(`:(\w+)`).FindAllStringSubmatch(path, -1)
		for _, v := range matches {
			route.Params = append(route.Params, v[1])
			path = strings.Replace(path, v[0], `(\w+)`, 1)
		}
	}
	compiledPath, err := regexp.Compile(path)
	if err != nil {
		panic(err)
	}
	route.Path = compiledPath
	this.routes = append(this.routes, route)
}

func (r *Router) Head(path string, handler http.HandlerFunc) {
	r.add("HEAD", path, handler)
}

func (r *Router) Connect(path string, handler http.HandlerFunc) {
	r.add("CONNECT", path, handler)
}

func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.add("GET", path, handler)
}

func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.add("POST", path, handler)
}

func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.add("PUT", path, handler)
}

func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.add("DELETE", path, handler)
}

func (r *Router) Options(path string, handler http.HandlerFunc) {
	r.add("OPTIONS", path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, params, err := r.Try(req.URL.Path, req.Method)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	ctx := context.WithValue(req.Context(), "params", params)
	handler(w, req.WithContext(ctx))
}

func parseParams(route *Route, path string) paramsMap {
	matches := route.Path.FindAllStringSubmatch(path, -1)
	params := paramsMap{}
	matchedParams := matches[0][1:]

	for k, v := range matchedParams {
		params[route.Params[k]] = v
	}
	return params
}

func getParam(r *http.Request, key string) string {
	ctx := r.Context()
	params := ctx.Value("params").(paramsMap)
	val := params[key]
	return val
}
