package restapi

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

func RootPathMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)

			return
		}
		handler.ServeHTTP(w, r)
	})
}

func LogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		middleware.Logger.Debugf("got request: %v", r)
		handler.ServeHTTP(w, r)
	})
}
