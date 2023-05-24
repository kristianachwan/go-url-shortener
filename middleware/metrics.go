package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	apiCallsByRoute = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_calls_by_route",
			Help: "Number of API calls by route",
		},
		[]string{"route"},
	)
)

func init() {
	prometheus.MustRegister(apiCallsByRoute)

}

func TrackApiCalls(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCallsByRoute.With(prometheus.Labels{"route": r.URL.Path}).Inc()
		handler.ServeHTTP(w, r)
	})
}
