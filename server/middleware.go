package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func TrackApiCalls(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCallsByRoute.With(prometheus.Labels{"route": r.URL.Path}).Inc()
		handler.ServeHTTP(w, r)
	})
}
