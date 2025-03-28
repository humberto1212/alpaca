package alpaca

import (
	"fmt"
	"net/http"
	"strings"
)

type Middleware func(Handler) Handler

func Chain(h Handler, middleware ...Middleware) Handler {

	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}

func (s *Server) ApplyMiddleware(h Handler) Handler {
	if len(s.middlewares) == 0 {
		return h
	}

	return Chain(h, s.middlewares...)
}

func (s *Server) Use(middleware ...Middleware) {
	s.middlewares = append(s.middlewares, middleware...)
}

// Implement Middlewares

func LoggingMiddleware(next Handler) Handler {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Apply Logging Middleware")
		next(w, r)
	}
}

func AuthMiddleware(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Apply Auth Middleware")
		next(w, r)
	}
}

func RecoveryMiddleware(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Apply Recovery Middleware")
		next(w, r)
	}
}

func EnableCors(router *Router) Middleware {
	return func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) {

			config := router.getCors()

			fmt.Println(strings.Join(config.allowAllOrigins, ", "))

			w.Header().Set("Access-Control-Allow-Origin", strings.Join(config.allowAllOrigins, ", "))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.allowMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.allowHeaders, ", "))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.exposeHeaders, ", "))
			w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%t", config.allowCredentials))
			w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.maxAge))

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	}
}
