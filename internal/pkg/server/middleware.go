package server

import (
	"net/http"
)

// MaxRequestsMiddleware limits the number of concurrent requests for specific handler
func MaxRequestsMiddleware(h http.Handler, max int) http.Handler {
	pool := make(chan struct{}, max)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pool <- struct{}{}
		defer func() {
			<-pool
		}()

		h.ServeHTTP(w, r)
	})
}
