package alpaca

import (
	"fmt"
	"net/http"
	"strings"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type Router struct {
	routes          map[string]map[string]Handler
	server          *Server
	notFoundHandler Handler
	cors            *Cors
}

func NewRouter(server *Server) *Router {

	return &Router{
		routes:          make(map[string]map[string]Handler),
		server:          server,
		notFoundHandler: defaultNotFoundHandler,
	}
}

func (r *Router) NotFound(handler Handler) {
	r.notFoundHandler = handler
}

func defaultNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func (r *Router) addRoute(method, path string, handler Handler) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]Handler)
	}

	r.routes[method][path] = handler
}

func (r *Router) GET(path string, handler Handler) {
	r.addRoute(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler Handler) {
	r.addRoute(http.MethodPost, path, handler)
}

func (r *Router) PUT(path string, handler Handler) {
	r.addRoute(http.MethodPut, path, handler)
}

func (r *Router) DELETE(path string, handler Handler) {
	r.addRoute(http.MethodDelete, path, handler)
}

func (r *Router) findHandler(method, path string) (Handler, error) {

	if methodRoutes, ok := r.routes[method]; ok {

		if handler, ok := methodRoutes[path]; ok {
			return handler, nil
		}

		for routePath, handler := range methodRoutes {
			if isWildCardMatch(routePath, path) {
				return handler, nil
			}
		}
	}

	return nil, fmt.Errorf("not Handler found for %s %s", method, path)
}

func isWildCardMatch(routePath, requestPath string) bool {

	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(routeParts) != len(requestParts) {
		return false
	}

	for i, part := range routeParts {

		if part == "*" {
			continue
		}

		if part == ":id" {
			continue
		}

		if part != requestParts[i] {
			return false
		}
	}

	return true
}

func (r *Router) ServerHTTP(w http.ResponseWriter, req *http.Request) {
	handler, err := r.findHandler(req.Method, req.URL.Path)

	if err != nil {
		r.notFoundHandler(w, req)
		return
	}

	if r.server != nil {
		handler = r.server.ApplyMiddleware(handler)
	}

	handler(w, req)
}

func (r *Router) SetCors(allowAllOrigins, allowMethods, allowHeaders, exposeHeaders []string, allowCredentials bool, maxAge int) {

	corsConfig := &Cors{
		allowAllOrigins,
		allowMethods,
		allowHeaders,
		exposeHeaders,
		allowCredentials,
		maxAge,
	}

	r.cors = newCors(corsConfig)

	if r.server != nil {
		r.server.Use(EnableCors(r.server.router))
	}
}

func (r *Router) getCors() *Cors {
	return r.cors
}
